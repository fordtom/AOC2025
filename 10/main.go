package main

import (
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	desired           uint16
	button_bitfields  []uint16
	button_increments [][]int
	desired_joltages  []int
	joltage_sum       int
}

func extractDesired(line string) (uint16, int) {
	d_start := strings.IndexByte(line, '[')
	d_end := strings.IndexByte(line, ']')
	pattern := line[d_start+1 : d_end]
	width := len(pattern)

	var desired uint16 = 0
	for i, c := range pattern {
		if c == '#' {
			desired |= 1 << i
		}
	}
	return desired, width
}

func extractButtons(line string, width int) ([]uint16, [][]int, error) {
	var raw []string

	for i := 0; i < len(line); i++ {
		if line[i] == '(' {
			end := strings.IndexByte(line[i:], ')')
			if end != -1 {
				raw = append(raw, line[i+1:i+end])
				i += end
			}
		}
	}

	bitfields := make([]uint16, 0, len(raw))
	increments := make([][]int, 0, len(raw))

	for _, b := range raw {
		parts := strings.Split(b, ",")

		var bitfield uint16 = 0
		var increment []int = make([]int, width)

		for _, p := range parts {
			num, err := strconv.ParseUint(p, 10, 16)
			if err != nil {
				return nil, nil, err
			}
			bitfield |= 1 << num
			increment[num] = 1
		}

		bitfields = append(bitfields, bitfield)
		increments = append(increments, increment)
	}

	return bitfields, increments, nil
}

func extractJoltages(line string) ([]int, error) {
	j_start := strings.IndexByte(line, '{')
	j_end := strings.IndexByte(line, '}')

	parts := strings.Split(line[j_start+1:j_end], ",")
	joltages := make([]int, 0, len(parts))

	for _, j := range parts {
		jolt, err := strconv.ParseInt(j, 10, 32)
		if err != nil {
			return nil, err
		}
		joltages = append(joltages, int(jolt))
	}
	return joltages, nil
}

func extractMachines(lines []string) []Machine {
	machines := make([]Machine, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		desired, width := extractDesired(line)

		buttons, increments, err := extractButtons(line, width)
		if err != nil {
			panic(err)
		}

		joltages, err := extractJoltages(line)
		if err != nil {
			panic(err)
		}

		joltageSums := 0
		for _, j := range joltages {
			joltageSums += j
		}

		machines = append(machines, Machine{desired, buttons, increments, joltages, joltageSums})
	}

	return machines
}

func (m Machine) MinTurnOn() int {
	min := math.MaxInt

	for i := 0; i < (1 << len(m.button_bitfields)); i++ {
		var result uint16 = 0
		for j := 0; j < len(m.button_bitfields); j++ {
			if i&(1<<j) != 0 {
				result ^= m.button_bitfields[j]
			}
		}

		if result == m.desired {
			if bits.OnesCount(uint(i)) < min {
				min = bits.OnesCount(uint(i))
			}
		}
	}

	return min
}

func partOne(machines []Machine) int {
	sum := 0

	for _, machine := range machines {
		sum += machine.MinTurnOn()
	}

	return sum
}

func partTwo(machines []Machine) int {
	sum := 0

	// to be implemented

	return sum
}

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	machines := extractMachines(lines)

	sum := partOne(machines)
	fmt.Println("Solution for part 1:", sum)

	sum = partTwo(machines)
	fmt.Println("Solution for part 2:", sum)
}
