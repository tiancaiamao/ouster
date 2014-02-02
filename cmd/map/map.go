package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tiancaiamao/ouster/tool/config"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

type Layer struct {
	Type string
	Data string // as a tool...it don't need to know that Data is []uint16 actually
}

func (l *Layer) String() string {
	return fmt.Sprintf("Layer{\nType:%s,\nData:[]uint16{%s},\n}", strings.ToUpper(l.Type), l.Data)
}

type EnemyGroup struct {
	Type     string
	Location string
	Level    string
	Number   string
}

func (eg *EnemyGroup) String() string {
	level := strings.Split(eg.Level, ",")
	number := strings.Split(eg.Number, ",")
	return fmt.Sprintf("EnemyGroup{\nType:%s,\nLocation:Rect{%s},\nLevelMin:%s,\nLevelMax:%s,\nNumberMin:%s,\nNumberMax:%s,\n}",
		strings.ToUpper(eg.Type),
		eg.Location,
		level[0], level[1], number[0], number[1])
}

type Enemy struct {
	Type      string
	Location  string
	Direction string
}

func (e *Enemy) String() string {
	return fmt.Sprintf("Enemy{\nType:%s,\nWanderArea:Rect{%s},\nDirection:%s,\n}",
		strings.ToUpper(e.Type),
		e.Location,
		e.Direction)
}

type Map struct {
	Title  string
	Width  string
	Height string

	Layers      []Layer
	EnemyGroups []EnemyGroup
	Enemies     []Enemy
}

func (m *Map) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("package data\nimport (\n. \"github.com/tiancaiamao/ouster\"\n)\n\nvar test = Map{\nTitle:\"%s\",\nWidth:%s,\nHeight:%s,\nLayers:[]Layer{\n", m.Title, m.Width, m.Height))
	for _, l := range m.Layers {
		b.WriteString(l.String())
		b.WriteString(",\n")
	}
	b.WriteString("\n},\nEnemyGroups:[]EnemyGroup{\n")
	for _, e := range m.EnemyGroups {
		b.WriteString(e.String())
		b.WriteString(",\n")
	}
	b.WriteString("\n},\nEnemies:[]Enemy{\n")
	for _, e := range m.Enemies {
		b.WriteString(e.String())
		b.WriteString(",\n")
	}
	b.WriteString("\n},\n}")

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", b.String(), 0)
	if err != nil {
		fmt.Println("not pretty print!!!")
		return b.String()
	}
	fmt.Println("Success!")
	var pretty bytes.Buffer
	printer.Fprint(&pretty, fset, f)
	return pretty.String()
}

var (
	input string
	output string
)

func init() {
	const (
		msgInput = "input file"
		msgOutput = "output file"
	)
	flag.StringVar(&input, "input", "", msgInput)
	flag.StringVar(&input, "i", "", msgInput)

	flag.StringVar(&output, "output", "", msgOutput)
	flag.StringVar(&output, "o", "", msgOutput)
}

func main() {
	flag.Parse()
	fmt.Println("input is:", input, "and output is:", output)
	m, err := loadMap(input)
	if err != nil {
		fmt.Println("error occured:", err)
		os.Exit(-1)
	}
	f, err := os.Create(output)
	if err != nil {
		panic(fmt.Sprintf("os.Create file error:", err))
	}
	fmt.Fprintf(f, "%v", m)
	// fmt.Printf("%v", m)
	return
}

func loadMap(fileName string) (*Map, error) {
	c, err := config.ReadDefault(fileName)
	if err != nil {
		return nil, err
	}

	m := new(Map)
	m.Width, err = c.String("header", "width")
	if err != nil {
		return nil, err
	}
	m.Height, err = c.String("header", "height")
	if err != nil {
		return nil, err
	}
	m.Title, err = c.String("header", "title")
	if err != nil {
		return nil, err
	}

	secs, err := c.Sections("layer")
	if err != nil {
		return nil, err
	}
	m.Layers = make([]Layer, len(secs))
	for i := 0; i < len(secs); i++ {
		sec := secs[i]
		m.Layers[i].Type, err = sec.String("type")
		if err != nil {
			fmt.Println("read layer's type error:", err)
		}
		m.Layers[i].Data, err = sec.String("data")
		if err != nil {
			fmt.Println("read layer's data error:", err)
		}
	}

	secs, err = c.Sections("enemygroup")
	if err != nil {
		return nil, err
	}
	m.EnemyGroups = make([]EnemyGroup, len(secs))
	for i := 0; i < len(secs); i++ {
		sec := secs[i]
		m.EnemyGroups[i].Type, err = sec.String("type")
		if err != nil {
			fmt.Println("read enemygroup's type error:", err)
			return nil, err
		}
		m.EnemyGroups[i].Location, err = sec.String("location")
		if err != nil {
			fmt.Println("read enemygroup's location error:", err)
		}
		m.EnemyGroups[i].Level, _ = sec.String("level")
		m.EnemyGroups[i].Number, _ = sec.String("number")
	}

	secs, err = c.Sections("enemy")
	if err != nil {
		return nil, err
	}
	m.Enemies = make([]Enemy, len(secs))
	for i := 0; i < len(secs); i++ {
		sec := secs[i]
		m.Enemies[i].Type, _ = sec.String("type")
		m.Enemies[i].Location, _ = sec.String("location")
		m.Enemies[i].Direction, _ = sec.String("direction")
	}
	return m, nil
}
