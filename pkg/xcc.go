package xcc

import (
	"errors"
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
	// Following fields are precisely based on Knuth's Dancing Cells presentation
	item   []int // len(item) == pCount + sCount
	node   []Node
	set    []int
	second int
	active int

	// Following fields were added to cache important values for easier processing
	pCount       int
	sCount       int
	optionsIndex []int // node[] index to point to beginning of each option
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
	xcc.optionsIndex = make([]int, si.GetOCount())

	xcc.Initialize(si.GetItemCount())
	xcc.SetItmClr(append(si.primary, si.secondary...), si.options)
	xcc.SetLoc(append(si.primary, si.secondary...), si.options)
	return &xcc, nil
}

// These 3 methods are used to initialize a XCC struct based on the solver input
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
	for i, option := range options {
		xcc.optionsIndex[i] = index
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
			// LOC and SET are Sparse and Dense arrays that points to one another
			xcc.node[index].loc = latestItemIndex[itemIndex]
			xcc.set[latestItemIndex[itemIndex]] = index

			latestItemIndex[itemIndex] += 1
			index += 1
		}
	}
}

// The following 2 methods implement Dancing Cells idea as presented by Knuth
func (xcc *XCC) ProcessOption(oIndex int) error {
	if oIndex < 0 || oIndex >= len(xcc.optionsIndex) {
		return errors.New("OptionIndex is out of range")
	}

	index := xcc.optionsIndex[oIndex]
	n := xcc.node[index]
	for n.itm >= 0 {
		i, j, err := denseDelete(xcc.set, n.itm-1, n.loc)
		if err != nil {
			// Something has gone seriously wrong!
			return err
		}
		swapLoc(xcc.node, i, j)

		index += 1
		n = xcc.node[index]
	}
	return nil
}

func (xcc *XCC) UndoOption(oIndex int) error {
	if oIndex < 0 || oIndex >= len(xcc.optionsIndex) {
		return errors.New("OptionIndex is out of range")
	}

	index := xcc.optionsIndex[oIndex]
	n := xcc.node[index]
	for n.itm >= 0 {
		iSize := n.itm - 
		// As Knuth said in his talk, the main benefit of dancing cell is that unDelete is super simple
		// We just need to increase the size!
		xcc.set[iSize] += 1

		index += 1
		n = xcc.node[index]
	}
	return nil
}

func (xcc *XCC) IsConsistent() bool {
	// This function checks the validity of sparse set:
	// all elements of dense and sparse must point together

	for i := 1; i < len(xcc.node); i++ {
		n := xcc.node[i]
		if n.itm < 0 {
			continue
		}
		if n.loc < 0 || n.loc > len(xcc.set) || xcc.set[n.loc] != i {
			fmt.Println("Found inconsistency at index ", i, n, xcc.set[n.loc])
			return false
		}
	}
	fmt.Println("No inconsistency was found.")
	return true
}

// Simple helper functions
func denseDelete(d []int, sizeI int, delI int) (int, int, error) {
	// Please refer to delete operation on the dense array of sparse set data structure.
	// d is XCC.set array that hosts multiple dense arrays as well as size value.
	if d[sizeI] <= 0 {
		return -1, -1, errors.New("Cannot delete from a dense array with size 0")
	}
	endI := sizeI + d[sizeI]
	if delI > endI || delI <= sizeI {
		return -1, -1, errors.New(fmt.Sprintf("denseDelete is called with wrong inputs: ", d, sizeI, delI))
	}
	swap(d, delI, endI)
	d[sizeI] -= 1
	return d[endI], d[delI], nil
}

func swap(s []int, i int, j int) error {
	if i < 0 || i >= len(s) || j < 0 || j >= len(s) {
		return errors.New("index out of range")
	}

	temp := s[i]
	s[i] = s[j]
	s[j] = temp
	return nil
}

func swapLoc(n []Node, i int, j int) error {
	if i < 0 || i >= len(n) || j < 0 || j >= len(n) {
		return errors.New("index out of range")
	}

	temp := n[i].loc
	n[i].loc = n[j].loc
	n[j].loc = temp
	return nil
}
