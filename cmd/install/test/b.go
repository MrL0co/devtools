package main

import (
	"log"
	"os"
)

func foo(x byte) byte { return x + 1 }
func bar(y byte) byte { return y * 2 }

func ReadByte() byte {
	b1 := make([]byte, 1)
	for {
		n, _ := os.Stdin.Read(b1)
		if n == 1 {
			return b1[0]
		}
	}
}
func WriteByte(b byte) {
	b1 := []byte{b}
	for {
		n, _ := os.Stdout.Write(b1)
		if n == 1 {
			return
		}
	}
}
func main() {
	var res byte
	for {
		fn := ReadByte()
		log.Println("fn=", fn)
		arg := ReadByte()
		log.Println("arg=", arg)
		if fn == 1 {
			res = foo(arg)
		} else if fn == 2 {
			res = bar(arg)
		} else if fn == 0 {
			return //exit
		} else {
			res = fn //echo
		}
		WriteByte(1)
		WriteByte(res)
	}
}
