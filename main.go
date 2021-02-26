package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/warthog618/gpiod"
)

/*
I connected GPIO pins 0 and 1 to my 2PH63091A 2-relay board to switch a big AC fan on and off.

			Raspberry Pi1 B.
Line 17: 	WiringPi pin 0
Line 18: 	WiringPi pin 1 (PCM_CLK)
*/

func main() {
	port := "8080"
	fmt.Printf("Fan Controller starting.. listening on port %v\n", port)

	// gpiod.WithConsumer("softwire")
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		fmt.Printf("Failed to connect to /dev/gpiochip0: %v\n", err)
		os.Exit(1)
	}
	/*
		for i := 10; i < 20; i++ {
			if i == 4 || i == 5 || i == 6 || i == 9 {
				continue
			}
			info, err := c.LineInfo(i)
			if err == nil {
				fmt.Printf("LineInfo %v: %v\n", i, info.Name)
			}
			out, err := c.RequestLine(i, gpiod.AsOutput())
			fmt.Printf("RequestLine %v: %v\n", i, err)
			if err != nil {
				continue
			}
			out.SetValue(1)
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(500 * time.Millisecond)
	*/

	// 1 = off
	// 2 = on
	pin0, err1 := c.RequestLine(17, gpiod.AsOutput(1))
	pin1, err2 := c.RequestLine(18, gpiod.AsOutput(1))
	if err1 != nil {
		fmt.Printf("Failed to connect to pin 17: %v\n", err)
		os.Exit(1)
	}
	if err2 != nil {
		fmt.Printf("Failed to connect to pin 18: %v\n", err)
		os.Exit(1)
	}
	//time.Sleep(500 * time.Millisecond)
	//pin0.SetValue(1)
	//pin1.SetValue(1)
	switches := []*gpiod.Line{
		pin0,
		pin1,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		//fmt.Printf("Parts: %v\n", parts)
		if len(parts) == 3 {
			switchNum, onOff := parts[1], parts[2]
			iswitch, _ := strconv.ParseInt(switchNum, 10, 64)
			if (iswitch == 1 || iswitch == 2) && (onOff == "on" || onOff == "off") {
				val := 1
				if parts[2] == "on" {
					val = 0
				}
				switches[iswitch-1].SetValue(val)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))

	c.Close()
}
