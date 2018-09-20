package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	errInvalidPacket = errors.New("invalid packet")
)

func main() {
	var port = flag.Int("port", 5987, "listening port")

	flag.Parse()

	log.Printf("Listening @ port %d...\n", *port)

	pc, err := net.ListenPacket("udp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer pc.Close()

	go rcvLoop(pc)

	select {}
}

func rcvLoop(pc net.PacketConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Printf("rcv err: %s", err)
			continue
		}
		// log.Printf("rcv from %s: %v\n", addr, buffer[:n])
		res, err := processPacket(buffer[:n])
		if err != nil {
			log.Printf("processPacket err: %s", err)
			continue
		}
		n, err = pc.WriteTo(res, addr)
		if err != nil {
			log.Printf("snd err: %s", err)
			continue
		}
		if n == 0 {
			log.Printf("snd err: wrote 0 bytes")
		}
	}
}

func processPacket(data []byte) ([]byte, error) {
	switch data[0] {
	case 0x20:
		return processCreateSession(data)
	case 0x80:
		return processCommand(data)
	case 0xD0:
		return processKeepAlive(data)
	default:
		return nil, errInvalidPacket
	}
}

func processCreateSession(data []byte) ([]byte, error) {
	log.Printf("processCreateSession\n")
	if len(data) != 27 {
		return nil, errInvalidPacket
	}
	res := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	return res, nil
}

func processCommand(data []byte) ([]byte, error) {
	// log.Printf("processCommand\n")
	switch data[14] {
	case 0x01:
		log.Printf("color")
	case 0x02:
		log.Printf("brightness")
	case 0x03:
		switch data[15] {
		case 0x03:
			log.Printf("on")
		case 0x04:
			log.Printf("off")
		case 0x05:
			log.Printf("white")
		default:
			return nil, errInvalidPacket
		}
	default:
		return nil, errInvalidPacket
	}
	res := []byte{0x88, 0x00, 0x00, 0x00, 0x03, 0x00, data[8], 0x00}
	return res, nil
}

func processKeepAlive(data []byte) ([]byte, error) {
	log.Printf("processKeepAlive\n")
	if len(data) != 8 {
		return nil, errInvalidPacket
	}
	res := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	return res, nil
}
