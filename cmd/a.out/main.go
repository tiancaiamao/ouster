package main

import (
	"fmt"
	"github.com/tiancaiamao/ouster/data"
)

func main() {
	m := &data.Test
	l := m.Layers[2]
	fmt.Println(l.Type)
	fmt.Println("size of Layer[2] is", len(l.Data))
	return
}
