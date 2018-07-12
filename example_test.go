package milight_test

import (
	"flag"
	"fmt"
	"os"
	"sgrzywna/milight"
	"strconv"
	"time"
)

// POST /api/v1/color/{color}

// POST /api/v1/brightness/{level}

// POST /api/v1/light
// DELETE /api/v1/light

// GET /api/v1/sequence/
// GET /api/v1/sequence/{seqID}
// POST /api/v1/sequence/ -> seqID
// POST /api/v1/sequence/{seqID}
// DELETE /api/v1/sequence/{seqID}

// /light:on/color:white/sleep:0.3/color:green/sleep:0.3/light:off
// sequence: {
//   name: "alert",
//     steps: [
//       "light:on",
//       "color:white",
//       "sleep:0.3",
//       "color:green",
//       "sleep:0.3",
//       "light:off"
//    ]
// }

func main() {
	var miHost = flag.String("mihost", "", "Mi-Light network address")
	var miPort = flag.Int("miport", 5987, "Mi-Light network port")
	var port = flag.Int("port", 8080, "listening port")

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
		fmt.Printf("milight on error: %s\n", err)
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)

	b, err := strconv.ParseInt(flag.Arg(0), 0, 16)
	if err != nil {
		fmt.Printf("brightness error: %s\n", err)
		os.Exit(1)
	}

	err = m.Brightness(byte(b))
	if err != nil {
		fmt.Printf("brightness error: %s\n", err)
		os.Exit(1)
	}

	c, err := strconv.ParseInt(flag.Arg(1), 0, 16)
	if err != nil {
		fmt.Printf("color error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("color: %d\n", c)
	err = m.Color(byte(c))
	if err != nil {
		fmt.Printf("milight color error: %s\n", err)
		os.Exit(1)
	}

	time.Sleep(3 * time.Second)

	err = m.Off()
	if err != nil {
		fmt.Printf("milight off error: %s\n", err)
		os.Exit(1)
	}

	// var i byte
	// for i = 0; i < 0x64; i++ {
	// 	err = brightness(conn, session, i)
	// 	if err != nil {
	// 		fmt.Printf("brightness error: %s\n", err)
	// 		break
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// }

	// for i = 0x64; i > 0; i-- {
	// 	err = brightness(conn, session, i)
	// 	if err != nil {
	// 		fmt.Printf("brightness error: %s\n", err)
	// 		break
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// }

	// err = m.Brightness(0)
	// if err != nil {
	// 	fmt.Printf("milight brightness error: %s\n", err)
	// 	os.Exit(1)
	// }

	// var i byte
	// for {
	// 	fmt.Printf("hue: %d\n", i)
	// 	err = m.Color(i)
	// 	if err != nil {
	// 		fmt.Printf("milight color error: %s\n", err)
	// 		break
	// 	}
	// 	time.Sleep(500 * time.Millisecond)
	// 	i += 15
	// 	if i > 0x64 {
	// 		break
	// 	}
	// }

	// colors := []byte{
	// 	milight.Red,
	// 	milight.Lavender,
	// 	milight.Blue,
	// 	milight.Aqua,
	// 	milight.Green,
	// 	milight.Lime,
	// 	milight.Yellow,
	// 	milight.Orange,
	// }

	// for i := 0; i < 100; i++ {
	// 	for _, c := range colors {
	// 		err = m.Color(c)
	// 		if err != nil {
	// 			fmt.Printf("milight color error: %s\n", err)
	// 			break
	// 		}
	// 		time.Sleep(100 * time.Millisecond)
	// 	}
	// }

	// err = m.On()
	// if err != nil {
	// 	fmt.Printf("milight off error: %s\n", err)
	// 	os.Exit(1)
	// }

	// err = m.Color(milight.Blue)
	// if err != nil {
	// 	fmt.Printf("milight color error: %s\n", err)
	// 	os.Exit(1)
	// }

	// err = m.Brightness(0)
	// if err != nil {
	// 	fmt.Printf("milight color error: %s\n", err)
	// 	os.Exit(1)
	// }

	// for i := 0; i < 10; i++ {
	// 	err = m.Off()
	// 	if err != nil {
	// 		fmt.Printf("milight off error: %s\n", err)
	// 		os.Exit(1)
	// 	}

	// 	time.Sleep(500 * time.Millisecond)

	// 	err = m.On()
	// 	if err != nil {
	// 		fmt.Printf("milight off error: %s\n", err)
	// 		os.Exit(1)
	// 	}

	// 	time.Sleep(500 * time.Millisecond)
	// }

	// err = m.Off()
	// if err != nil {
	// 	fmt.Printf("milight off error: %s\n", err)
	// 	os.Exit(1)
	// }

	// for {
	// 	err = on(conn, session)
	// 	if err != nil {
	// 		fmt.Printf("on error: %s\n", err)
	// 		break
	// 	}

	// 	time.Sleep(500 * time.Millisecond)

	// 	err = off(conn, session)
	// 	if err != nil {
	// 		fmt.Printf("off error: %s\n", err)
	// 		break
	// 	}

	// 	time.Sleep(500 * time.Millisecond)
	// }
}
