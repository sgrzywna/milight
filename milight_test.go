package milight

import "testing"

func TestChecksum(t *testing.T) {
	packet := []byte{
		0x20, 0x00, 0x00, 0x00, 0x16, 0x02, 0x62, 0x3A,
		0xD5, 0xED, 0xA3, 0x01, 0xAE, 0x08, 0x2D, 0x46,
		0x61, 0x41, 0xA7, 0xF6, 0xDC, 0xAF, 0xD3, 0xE6,
		0x00, 0x00,
	}
	chksum := checksum(packet)
	if chksum != 0x1E {
		t.Errorf("checksum error: 0x%X\n", chksum)
	}
}
