package xcc

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Node struct {
	itm int
	clr string
	loc int
}

type XCC struct {
	pCount int
	sCount int
	item   []int // len(item) == pCount + sCount
	node   []Node
	set    []int
	second int
	active int
}

func Builder(rd io.Reader) (*XCC, error) {
	si := SolverInput{}
	err := si.Build(rd)
	if err != nil {
		return nil, err
	}
	fmt.Println(si)

	xcc := XCC{}
	xcc.pCount = si.GetPCount()
	xcc.sCount = si.GetSCount()
	xcc.item = make([]int, si.GetPCount()+si.GetSCount())
	xcc.set = make([]int, si.GetSetLen())
	xcc.node = make([]Node, si.GetNodeLen())

	xcc.Initialize(si.GetItemCount())
	fmt.Println(xcc)
	fmt.Println("=====")

	xcc.SetItmClr(append(si.primary, si.secondary...), si.options)
	fmt.Println(xcc)
	fmt.Println("=====")

	xcc.SetLoc(append(si.primary, si.secondary...), si.options)
	fmt.Println(xcc)
	fmt.Println("=====")

	return &xcc, nil
}

func (xcc *XCC) Initialize(itemCount []int) {
	index := 0
	for i, count := range itemCount {
		xcc.set[index] = i // POS value
		index += 1
		xcc.set[index] = count // SIZE value
		index += 1
		xcc.item[i] = index // item: value (p:, q:, r:, x:, y:)
		index += count
	}
}

func (xcc *XCC) SetItmClr(items []string, options [][]string) {
	index := 1 // Itm[0] is left to be 0
	for _, option := range options {
		for _, item := range option {
			if len(strings.Split(item, ":")) == 2 {
				xcc.node[index].clr = strings.Split(item, ":")[1]
				item = strings.Split(item, ":")[0]
			}
			xcc.node[index].itm = xcc.item[slices.Index(items, item)]

			index += 1
		}
		xcc.node[index].itm = -1 * len(option)
		index += 1
	}
}

func (xcc *XCC) SetLoc(items []string, options [][]string) {
	index := 0
	latestItemIndex := make([]int, len(xcc.item))
	copy(latestItemIndex, xcc.item)
	for _, option := range options {
		xcc.node[index].loc = len(option)
		index += 1
		for _, item := range option {
			item = strings.Split(item, ":")[0] // drop color
			itemIndex := slices.Index(items, item)
			// loc and set are Sparse and Dense that points to one another
			xcc.node[index].loc = latestItemIndex[itemIndex]
			xcc.set[latestItemIndex[itemIndex]] = index

			latestItemIndex[itemIndex] += 1
			index += 1
		}
	}
}
