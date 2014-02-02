package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster/packet"
)

func main() {
	for k, v := range packet.PacketMap {
		fmt.Println(k, "----", v)
	}
	return
}
