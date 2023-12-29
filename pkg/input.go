package xcc

import (
	"bufio"
	"io"
	"slices"
	"strings"
)

type SolverInput struct {
	primary   []string
	secondary []string
	options   [][]string

	optionsWithoutColors [][]string

	oCount int // options count
	iCount int // total items counts = sum of all terms in all options
}

func (si *SolverInput) Build(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if si.primary == nil {
			si.primary = fields
			continue
		}

		if si.secondary == nil {
			si.secondary = fields
			continue
		}

		si.options = append(si.options, fields)
		si.optionsWithoutColors = append(si.optionsWithoutColors, dropColors(fields))

		si.oCount += 1
		si.iCount += len(fields)
	}
	return nil
}

func (si SolverInput) GetPCount() int {
	return len(si.primary)
}

func (si SolverInput) GetSCount() int {
	return len(si.secondary)
}

func (si SolverInput) GetOCount() int {
	return si.oCount
}

func (si SolverInput) GetSetLen() int {
	// SET array has one element for item in options plus
	// two value for each (primary+secondary) item POS and SIZE
	return si.iCount + (si.GetPCount()+si.GetSCount())*2
}

func (si SolverInput) GetNodeLen() int {
	// Each item of all options needs a node plus
	// we need one extra for each option (to save its length) plus
	// element 0
	return 1 + si.iCount + si.oCount
}

func (si SolverInput) GetItemCount() []int {
	// return how many times each (primary or secondary) item appears in all options
	counts := make([]int, si.GetPCount()+si.GetSCount())

	for i, item := range append(si.primary, si.secondary...) {
		for _, options := range si.optionsWithoutColors {
			if slices.Contains(options, item) {
				counts[i] += 1
			}
		}
	}
	return counts
}

func dropColors(fields []string) []string {
	noColor := make([]string, len(fields))
	for i, item := range fields {
		noColor[i] = strings.Split(item, ":")[0]
	}
	return noColor
}
