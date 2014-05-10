package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/tiancaiamao/ouster/packet/msgpack"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8782")
	if err != nil {
		panic(err)
	}

	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadBytes('\n')
		if err != nil {
			fmt.Println("read line error...")
			continue
		}

		if bytes.HasPrefix(line, []byte("move")) {
			off := 4
			X, err := readInt(line, &off)
			if err != nil {
				fmt.Println("read move error")
				continue
			}
			Y, err := readInt(line, &off)
			if err != nil {
				fmt.Println("read move error")
				continue
			}

			fmt.Println(X, Y, "-------")
			msgpack.Write(conn, msgpack.PCMove, msgpack.CMovePacket{
				X: float32(X),
				Y: float32(Y),
			})
		}
	}
}

func readInt(buf []byte, offset *int) (ret int, err error) {

	ofst := *offset
	// skip whitespace
	for ofst < len(buf) && buf[ofst] == ' ' {
		ofst++
	}

	read := false
	for ofst < len(buf) {
		if buf[ofst] <= '9' && buf[ofst] >= '0' {
			ret = ret*10 + int(buf[ofst]-'0')
			read = true
			ofst++
		} else {
			break
		}
	}

	if !read {
		err = errors.New("ReadInt")
	}
	*offset = ofst
	return
}
