package milight

import (
	"bytes"
	"net"
	"strconv"
	"testing"
)

const (
	WB1 byte = 0x98
	WB2 byte = 0x65
)

func TestChecksum(t *testing.T) {
	packet := []byte{
		0x31, 0x00, 0x00, 0x08, 0x04, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00,
	}
	chksum := checksum(packet)
	if chksum != 0x3F {
		t.Errorf("checksum error: 0x%X\n", chksum)
	}
}

func TestInitSession(t *testing.T) {
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer pc.Close()

	go func() {
		handleInitSession(t, pc)
	}()

	host, port, err := net.SplitHostPort(pc.LocalAddr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	iPort, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	m, err := NewMilight(host, iPort)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer m.Close()

	if m.sessionID[0] != WB1 || m.sessionID[1] != WB2 {
		t.Fatalf("expected [0x%X, 0x%X], got [0x%X, 0x%X]", WB1, WB2, m.sessionID[0], m.sessionID[1])
	}
}

func TestSendCommand(t *testing.T) {
	var cmd = []byte{0x31, 0x00, 0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x00}

	var packet = []byte{
		0x80, 0x00, 0x00, 0x00, 0x11, WB1, WB2, 0x00,
		0x01, 0x00, 0x31, 0x00, 0x00, 0x00, 0x03, 0x03,
		0x00, 0x00, 0x00, defaultZone, 0x00, 0x38,
	}

	var res = []byte{0x88, 0x00, 0x00, 0x00, 0x03, 0x00, 0x01, 0x00}

	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer pc.Close()

	go func() {
		handleInitSession(t, pc)
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if !bytes.Equal(packet, buffer[:n]) {
			t.Fatalf("expected %v, got %v", packet, buffer[:n])
		}
		n, err = pc.WriteTo(res, addr)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if n == 0 {
			t.Fatalf("write 0 bytes")
		}
	}()

	host, port, err := net.SplitHostPort(pc.LocalAddr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	iPort, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	m, err := NewMilight(host, iPort)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer m.Close()

	err = m.sendCommand(cmd)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestKeepAlive(t *testing.T) {
	var cmd = []byte{0xD0, 0x00, 0x00, 0x00, 0x02, WB1, WB2, 0x00}

	var res = []byte{0xD8, 0x00, 0x00, 0x00, 0x07, 0xF0, 0xFE, 0x6B, 0xC2, 0x7F, 0x8A, 0x01}

	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer pc.Close()

	go func() {
		handleInitSession(t, pc)
		buffer := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if !bytes.Equal(cmd, buffer[:n]) {
			t.Fatalf("expected %v, got %v", cmd, buffer[:n])
		}
		n, err = pc.WriteTo(res, addr)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if n == 0 {
			t.Fatalf("write 0 bytes")
		}
	}()

	host, port, err := net.SplitHostPort(pc.LocalAddr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	iPort, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	m, err := NewMilight(host, iPort)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer m.Close()

	err = m.KeepAlive()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func handleInitSession(t *testing.T, pc net.PacketConn) {
	var cmd = []byte{
		0x20, 0x00, 0x00, 0x00, 0x16, 0x02, 0x62, 0x3A,
		0xD5, 0xED, 0xA3, 0x01, 0xAE, 0x08, 0x2D, 0x46,
		0x61, 0x41, 0xA7, 0xF6, 0xDC, 0xAF, 0xD3, 0xE6,
		0x00, 0x00, 0x1E,
	}

	var res = []byte{
		0x28, 0x00, 0x00, 0x00, 0x11, 0x00, 0x02, 0xAC,
		0xCF, 0x23, 0xF5, 0x7A, 0xD4, 0x69, 0xF0, 0x3C,
		0x23, 0x00, 0x01, WB1, WB2, 0x00,
	}

	buffer := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buffer)
	if err != nil {
		t.Fatalf("handleInitSession: err: %s", err)
	}
	if !bytes.Equal(cmd, buffer[:n]) {
		t.Fatalf("handleInitSession: expected %v, got %v", cmd, buffer[:n])
	}
	n, err = pc.WriteTo(res, addr)
	if err != nil {
		t.Fatalf("handleInitSession: err: %s", err)
	}
	if n == 0 {
		t.Fatalf("handleInitSession: write 0 bytes")
	}
}
