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

/*
esc,f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f11,
grave,1,2,3,4,5,6,7,8,9,0,minus,tab,
q,w,e,r,t,y,u,i,o,p,lbrace,caps,a,s,d,
f,g,h,j,k,l,colon,quote,lshift,bslash_iso,
z,x,c,v,b,n,m,comma,dot,slash,lctrl,lwin,
lalt,space,katahira,ralt,rwin,rmenu,light,
f12,prtscn,scroll,pause,rbrace,bslash,hash,
enter,ro,equal,yen,bspace,rshift,rctrl,up,
left,down,right,mute,muhenkan,henkan,fn,
topbar1,topbar2,topbar3,topbar4,topbar5,
topbar6,topbar7,topbar8,topbar9,topbar10,
topbar11,topbar12,topbar13,topbar14,topbar15,
topbar16,topbar17,topbar18,topbar19,logo,
side,ins,del,home,end,pgup,pgdn,lock,volup,voldn
*/

const (
	KeyEsc = iota
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyPrintScreen
	KeyScrollLock
	KeyPauseBreak

	KeyW
	KeyA
	KeyS
	KeyD
)

var Keys = []string{
	"esc",
	"f1",
	"f2",
	"f3",
	"f4",
	"f5",
	"f6",
	"f7",
	"f8",
	"f9",
	"f10",
	"f11",
	"f12",
	"prtscn",
	"scroll",
	"pause",

	"w",
	"a",
	"s",
	"d",
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
	r := bufio.NewReader(ckbNotify)

	w.WriteString("notify w a s d\n")
	if err := w.Flush(); err != nil {
		panic(err)
	}

	go func() {
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				panic(err)
			}

			fmt.Println(str)

			runtime.Gosched()
		}
	}()

	go func() {
		step := float32(360.0 / len(Keys))
		offset := float32(0.0)
		tick := time.Tick(17 * time.Millisecond)
		for range tick {
			for t := 0; t < len(Keys); t++ {
				h := float32(t)*step - offset
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
				cmd := fmt.Sprintf("rgb %s:%02x%02x%02x,", Keys[t], r, g, b)
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
