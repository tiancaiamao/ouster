package ouster

import (
	"io"
	"log"
	"net"
	"runtime"
	"encoding/binary"
)

func loop(client chan []byte) {
	aoi := make(chan interface{})
	for {
		select {
		case data, ok := (<-client):
			if ok {
				// 解析packet，决定自己处理或者向其它地方转发
			} else {
			}
		case <-aoi:
			// 来自aoi的消息
		default:
			runtime.Gosched()
		}
	}
}

func PlayerGoroutine(conn net.Conn) {
	ch := make(chan []byte)
	var header [4]byte
	go loop(ch)
	for {
		n, err := io.ReadFull(conn, header[:])
		if n == 0 && err == io.EOF {
			// 处理出错
			break
		} else if err != nil {
			log.Println("error receiving header:", err)
			break
		}

		// data
		size := binary.BigEndian.Uint16(header[:])
		data := make([]byte, size)
		n, err = io.ReadFull(conn, data)

		if err != nil {
			log.Println("error receiving msg:", err)
			break
		}
		ch <- data
	}
}
