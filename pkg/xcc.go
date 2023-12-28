package xcc

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Node struct {
	itm int
	clr string
	loc uint
}

type XCC struct {
	pCount int
	sCount int
	item   []uint // len(item) == pCount + sCount
	node   []Node
	set    []uint
	second int
	active int
}

func Builder(rd io.Reader) (*XCC, error) {
	var primary []string
	var secondary []string
	var options [][]string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line := scanner.Text()

		if primary == nil {
			primary = strings.Fields(line)
			continue
		}

		if secondary == nil {
			secondary = strings.Fields(line)
			continue
		}

		options = append(options, strings.Fields(line))
	}
	fmt.Println(primary)
	fmt.Println(secondary)
	fmt.Println(options)
	return &XCC{}, nil
}
