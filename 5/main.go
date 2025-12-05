package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
	min int
	max int
}

func (r Range) Contains(id int) bool {
	return id >= r.min && id <= r.max
}

func (r Range) Size() int {
	return r.max - r.min + 1
}

func makeRanges(lines []string) []Range {
	ranges := make([]Range, 0, len(lines))
	for _, line := range lines {
		parts := strings.Split(line, "-")
		min, _ := strconv.Atoi(parts[0])
		max, _ := strconv.Atoi(parts[1])
		ranges = append(ranges, Range{min, max})
	}

	slices.SortFunc(ranges, func(a, b Range) int {
		return a.min - b.min
	})

	return ranges
}

func partOne(ranges []Range, ingredients []string) int {
	sum := 0

	for i := range ingredients {
		id, _ := strconv.Atoi(ingredients[i])

		for _, r := range ranges {
			if r.Contains(id) {
				sum++
				break
			}
		}
	}

	return sum
}

func partTwo(ranges []Range) int {
	currentRange := ranges[0]
	sum := 0

	for _, r := range ranges[1:] {
		if r.min <= currentRange.max {
			if r.max > currentRange.max {
				currentRange.max = r.max
			}
		} else {
			sum += currentRange.Size()
			currentRange = r
		}
	}
	sum += currentRange.Size()

	return sum
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	index := slices.Index(lines, "")

	ranges := makeRanges(lines[:index])
	ingredients := lines[index+1:]

	sum := partOne(ranges, ingredients)
	fmt.Println("Solution for part 1:", sum)

	sum = partTwo(ranges)
	fmt.Println("Solution for part 2:", sum)
}
