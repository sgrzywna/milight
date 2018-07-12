package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sgrzywna/milight"
)

func main() {
	var host = flag.String("host", "", "Mi-Light network address")
	var port = flag.Int("port", 5987, "Mi-Light network port")

	flag.Parse()

	if *host == "" {
		fmt.Println("provide Mi-Light network address")
		os.Exit(1)
	}

	m, err := milight.NewMilight(*host, *port)
	if err != nil {
		fmt.Printf("milight error: %s\n", err)
		os.Exit(1)
	}
	defer m.Close()

	err = m.InitSession()
	if err != nil {
		fmt.Printf("milight session error: %s\n", err)
		os.Exit(1)
	}

	err = m.On()
	if err != nil {
		fmt.Printf("milight on error: %s\n", err)
		os.Exit(1)
	}

	err = m.White()
	if err != nil {
		fmt.Printf("milight white error: %s\n", err)
		os.Exit(1)
	}

	for b := 0; b < 0x32; b++ {
		err = m.Brightness(byte(b))
		if err != nil {
			fmt.Printf("milight brightness error: %s\n", err)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
	}

	for c := 0; c < 0xFF; c++ {
		err = m.Color(byte(c))
		if err != nil {
			fmt.Printf("milight color error: %s\n", err)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
	}

	colors := []byte{
		milight.Red,
		milight.Orange,
		milight.Yellow,
		milight.ChartreuseGreen,
		milight.Green,
		milight.SpringGreen,
		milight.Cyan,
		milight.Azure,
		milight.Blue,
		milight.Violet,
		milight.Magenta,
		milight.Rose,
	}

	for _, c := range colors {
		err = m.Color(c)
		if err != nil {
			fmt.Printf("milight color error: %s\n", err)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	err = m.Off()
	if err != nil {
		fmt.Printf("milight off error: %s\n", err)
		os.Exit(1)
	}
}
