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

// rational number for exact linear algebra
type rat struct{ num, den int64 }

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}

func newrat(num, den int64) rat {
	if num == 0 {
		return rat{0, 1}
	}
	if den < 0 {
		num, den = -num, -den
	}
	g := gcd(num, den)
	return rat{num / g, den / g}
}

func (r rat) sub(o rat) rat { return newrat(r.num*o.den-o.num*r.den, r.den*o.den) }
func (r rat) mul(o rat) rat { return newrat(r.num*o.num, r.den*o.den) }
func (r rat) div(o rat) rat { return newrat(r.num*o.den, r.den*o.num) }

// gaussian elimination to reduced row echelon form
// returns pivot column indices and whether system is consistent
func rref(aug [][]rat, rows, cols int) ([]int, bool) {
	pivot_row := 0
	pivot_cols := []int{}

	for col := 0; col < cols && pivot_row < rows; col++ {
		// find pivot
		pivot := -1
		for r := pivot_row; r < rows; r++ {
			if aug[r][col].num != 0 {
				pivot = r
				break
			}
		}
		if pivot == -1 {
			continue
		}
		aug[pivot_row], aug[pivot] = aug[pivot], aug[pivot_row]

		// scale pivot row to 1
		pv := aug[pivot_row][col]
		for c := col; c <= cols; c++ {
			aug[pivot_row][c] = aug[pivot_row][c].div(pv)
		}

		// eliminate column in other rows
		for r := range rows {
			if r == pivot_row || aug[r][col].num == 0 {
				continue
			}
			factor := aug[r][col]
			for c := col; c <= cols; c++ {
				aug[r][c] = aug[r][c].sub(aug[pivot_row][c].mul(factor))
			}
		}

		pivot_cols = append(pivot_cols, col)
		pivot_row++
	}

	// check consistency: row of all zeros with nonzero RHS means no solution
	for r := range rows {
		all_zero := true
		for c := range cols {
			if aug[r][c].num != 0 {
				all_zero = false
				break
			}
		}
		if all_zero && aug[r][cols].num != 0 {
			return nil, false
		}
	}

	return pivot_cols, true
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

func (m Machine) MinCorrectJoltages() int {
	num_counters := len(m.desired_joltages)
	num_buttons := len(m.button_increments)

	// build augmented matrix [A|t] for system A*x = t
	aug := make([][]rat, num_counters)
	for i := 0; i < num_counters; i++ {
		row := make([]rat, num_buttons+1)
		for j := 0; j < num_buttons; j++ {
			row[j] = newrat(int64(m.button_increments[j][i]), 1)
		}
		row[num_buttons] = newrat(int64(m.desired_joltages[i]), 1)
		aug[i] = row
	}

	pivot_cols, ok := rref(aug, num_counters, num_buttons)
	if !ok {
		return math.MaxInt
	}

	// identify free variables (columns without pivots)
	is_pivot := make([]bool, num_buttons)
	for _, c := range pivot_cols {
		is_pivot[c] = true
	}
	free_cols := []int{}
	for c := 0; c < num_buttons; c++ {
		if !is_pivot[c] {
			free_cols = append(free_cols, c)
		}
	}

	// upper bound: button can't be pressed more than smallest target it affects
	upper_bound := func(btn int) int {
		ub := math.MaxInt
		for k := 0; k < num_counters; k++ {
			if m.button_increments[btn][k] > 0 && m.desired_joltages[k] < ub {
				ub = m.desired_joltages[k]
			}
		}
		if ub == math.MaxInt {
			return 0
		}
		return ub
	}

	// extract coefficients for free variables from RREF result
	num_pivots := len(pivot_cols)
	rhs := make([]rat, num_pivots)
	coeff := make([][]rat, num_pivots)
	for i := 0; i < num_pivots; i++ {
		rhs[i] = aug[i][num_buttons]
		coeff[i] = make([]rat, len(free_cols))
		for f, col := range free_cols {
			coeff[i][f] = aug[i][col]
		}
	}

	best := math.MaxInt

	// evaluate a candidate assignment of free variables
	eval := func(free []int) {
		total := 0
		for _, v := range free {
			total += v
		}
		if total >= best {
			return
		}

		// compute pivot variables: x_pivot = rhs - sum(coeff * free)
		for i := range num_pivots {
			val := rhs[i]
			for f := 0; f < len(free_cols); f++ {
				if free[f] != 0 {
					val = val.sub(newrat(coeff[i][f].num*int64(free[f]), coeff[i][f].den))
				}
			}
			if val.den != 1 || val.num < 0 {
				return // not a valid non-negative integer solution
			}
			total += int(val.num)
			if total >= best {
				return
			}
		}
		best = total
	}

	// enumerate all non-negative integer values for free variables
	switch len(free_cols) {
	case 0:
		eval(nil)
	case 1:
		for a := 0; a <= upper_bound(free_cols[0]); a++ {
			eval([]int{a})
		}
	case 2:
		for a := 0; a <= upper_bound(free_cols[0]); a++ {
			for b := 0; b <= upper_bound(free_cols[1]); b++ {
				eval([]int{a, b})
			}
		}
	case 3:
		for a := 0; a <= upper_bound(free_cols[0]); a++ {
			for b := 0; b <= upper_bound(free_cols[1]); b++ {
				for c := 0; c <= upper_bound(free_cols[2]); c++ {
					eval([]int{a, b, c})
				}
			}
		}
	default:
		panic(fmt.Sprintf("too many free variables: %d", len(free_cols)))
	}

	return best
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

	for _, machine := range machines {
		sum += machine.MinCorrectJoltages()
	}

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
