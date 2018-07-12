// Package milight implements basic commands for control Mi-Light devices.
package milight

import (
	"fmt"
	"net"
	"time"
)

const (
	// Red color.
	Red byte = 0x00
	// Orange color.
	Orange byte = 0x15
	// Yellow color.
	Yellow byte = 0x2A
	// ChartreuseGreen color.
	ChartreuseGreen byte = 0x3F
	// Green color.
	Green byte = 0x55
	// SpringGreen color.
	SpringGreen byte = 0x6A
	// Cyan color.
	Cyan byte = 0x7F
	// Azure color.
	Azure byte = 0x94
	// Blue color.
	Blue byte = 0xAA
	// Violet color.
	Violet byte = 0xBF
	// Magenta color.
	Magenta byte = 0xD4
	// Rose color.
	Rose byte = 0xE9

	defaultZone byte = 0x01

	defaultKeepAlivePeriod time.Duration = 5 * time.Second

	maxBrightnessLevel byte = 0x64
)

var (
	// ErrResponseTooShort is returned when Mi-Light device responds with message too short.
	ErrResponseTooShort = fmt.Errorf("response too short")
)

// Milight represent Mi-Light controller.
type Milight struct {
	conn         net.Conn
	zone         byte
	quit         chan bool
	seqNum       byte
	sessionID    [2]byte
	lastActivity time.Time
}

// NewMilight returns initialized Mi-Light controller.
func NewMilight(addr string, port int) (*Milight, error) {
	d := net.Dialer{Timeout: 1 * time.Second}
	conn, err := d.Dial("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return nil, err
	}
	return &Milight{
		conn: conn,
		zone: defaultZone,
		quit: make(chan bool, 1),
	}, nil
}

// Close closes connection to Mi-Light device.
func (m *Milight) Close() error {
	close(m.quit)
	return m.conn.Close()
}

// InitSession creates session if needed.
func (m *Milight) InitSession() error {
	if m.sessionID[0] == 0 && m.sessionID[1] == 0 {
		err := m.createSession()
		if err != nil {
			return err
		}
		go m.keepAliveLoop()
	}
	return nil
}

// On turns light on.
func (m *Milight) On() error {
	cmd := []byte{0x31, 0x00, 0x00, 0x00, 0x03, 0x03, 0x00, 0x00, 0x00}
	return m.sendCommand(cmd)
}

// Off turns light off.
func (m *Milight) Off() error {
	cmd := []byte{0x31, 0x00, 0x00, 0x00, 0x03, 0x04, 0x00, 0x00, 0x00}
	return m.sendCommand(cmd)
}

// Color sets light color.
func (m *Milight) Color(color byte) error {
	cmd := []byte{0x31, 0x00, 0x00, 0x00, 0x01, color, color, color, color}
	return m.sendCommand(cmd)
}

// White sets white light.
func (m *Milight) White() error {
	cmd := []byte{0x31, 0x00, 0x00, 0x00, 0x03, 0x05, 0x00, 0x00, 0x00}
	return m.sendCommand(cmd)
}

// Brightness sets brightness level.
func (m *Milight) Brightness(brightness byte) error {
	if brightness > maxBrightnessLevel {
		brightness = maxBrightnessLevel
	}
	cmd := []byte{0x31, 0x00, 0x00, 0x00, 0x02, brightness, 0x00, 0x00, 0x00}
	return m.sendCommand(cmd)
}

// KeepAlive sustains session.
func (m *Milight) KeepAlive() error {
	packet := []byte{0xD0, 0x00, 0x00, 0x00, 0x02, m.sessionID[0], m.sessionID[1], 0x00}
	_, err := m.conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

// createSession creates Mi-Light communication session.
func (m *Milight) createSession() error {
	packet := []byte{
		0x20, 0x00, 0x00, 0x00, 0x16, 0x02, 0x62, 0x3A,
		0xD5, 0xED, 0xA3, 0x01, 0xAE, 0x08, 0x2D, 0x46,
		0x61, 0x41, 0xA7, 0xF6, 0xDC, 0xAF, 0xD3, 0xE6,
		0x00, 0x00, 0x1E,
	}
	_, err := m.conn.Write(packet)
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	m.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := m.conn.Read(buf)
	if err != nil {
		return err
	}
	if n != 22 {
		return ErrResponseTooShort
	}
	m.sessionID[0] = buf[19]
	m.sessionID[1] = buf[20]
	return nil
}

// keepAliveLoop periodically sends keep alive packets to sustain session.
func (m *Milight) keepAliveLoop() {
	m.lastActivity = time.Now()
	for {
		select {
		case <-m.quit:
			break
		case <-time.After(2 * time.Second):
			if time.Since(m.lastActivity) > defaultKeepAlivePeriod {
				m.KeepAlive()
				m.lastActivity = time.Now()
			}
		}
	}
}

// sendCommand send command to the Mi-Light device.
func (m *Milight) sendCommand(cmd []byte) error {
	m.lastActivity = time.Now()
	packet := []byte{0x80, 0x00, 0x00, 0x00, 0x11, m.sessionID[0], m.sessionID[1], 0x00, m.getSeqNum(), 0x00}
	packet = append(packet, cmd...)
	packet = append(packet, m.zone, 0x00)
	packet = append(packet, checksum(packet))
	_, err := m.conn.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

// getSeqNum returns next sequence number.
func (m *Milight) getSeqNum() byte {
	m.seqNum++
	return m.seqNum
}

// checksum calculates checksum for input data.
func checksum(data []byte) byte {
	var chksum byte
	if len(data) > 10 {
		for _, b := range data[len(data)-11:] {
			chksum += b
		}
	}
	return chksum
}
