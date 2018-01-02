package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(2)
}

func HSLToRGB(h float32, s, v float32) (uint8, uint8, uint8) {
	var r uint8
	var g uint8
	var b uint8

	H := h
	if H >= 360 {
		H = 0.0
	} else {
		H /= 60.0
	}
	fract := H - float32(math.Floor(float64(H)))
	S := s
	V := v

	P := V * (1.0 - S)
	Q := V * (1.0 - S*fract)
	T := V * (1.0 - S*(1.0-fract))

	if 0.0 <= H && H < 1.0 {
		r = uint8(V * 255)
		g = uint8(T * 255)
		b = uint8(P * 255)
	} else if 1.0 <= H && H < 2.0 {
		r = uint8(Q * 255)
		g = uint8(V * 255)
		b = uint8(P * 255)
	} else if 2.0 <= H && H < 3.0 {
		r = uint8(P * 255)
		g = uint8(V * 255)
		b = uint8(T * 255)
	} else if 3.0 <= H && H < 4.0 {
		r = uint8(P * 255)
		g = uint8(Q * 255)
		b = uint8(V * 255)
	} else if 4.0 <= H && H < 5.0 {
		r = uint8(T * 255)
		g = uint8(P * 255)
		b = uint8(V * 255)
	} else if 5.0 <= H && H < 6.0 {
		r = uint8(V * 255)
		g = uint8(P * 255)
		b = uint8(Q * 255)
	} else {
		r = uint8(0)
		g = uint8(0)
		b = uint8(0)
	}

	return r, g, b
}

func main() {
	ckbCmd, err := os.OpenFile("/dev/input/ckb1/cmd", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ckbCmd.Close(); err != nil {
			panic(err)
		}
	}()
	ckbNotify, err := os.Open("/dev/input/ckb1/notify0")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := ckbNotify.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(ckbCmd)
	//r := bufio.NewReader(ckbNotify)

	go func() {
		offset := float32(0.0)
		tick := time.Tick(17 * time.Millisecond)
		for range tick {
			for t := 0; t < 12; t++ {
				h := float32(t)*(360.0/12.0) - offset
				for h < 0 {
					h += 360
				}
				for h > 360 {
					h -= 360
				}

				offset += 0.5
				if offset > 360 {
					offset = 0
				}

				r, g, b := HSLToRGB(h, 1.0, 1.0)
				cmd := fmt.Sprintf("rgb f%d:%02x%02x%02x,", t+1, r, g, b)
				w.WriteString(cmd)
			}

			// Send cmd
			w.Write([]byte{'\n'})
			if err := w.Flush(); err != nil {
				panic(err)
			}

			runtime.Gosched()
		}
	}()

	fmt.Println("ENTER to exit")
	fmt.Scanln()
}
