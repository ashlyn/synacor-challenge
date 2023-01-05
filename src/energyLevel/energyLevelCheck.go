package energylevel

import (
	"fmt"
	"math"
)

var max, mod int = 32767, 32768

// GetMaxEnergyLevel gets the energy level to set r7 to
func GetMaxEnergyLevel(r0 int, r1 int) int {
	r7 := max
	for r7 > 0 {
		if ackermannBetter([3]int{ r0, r1, r7 }) == 6 {
			return r7
		}
		r7--
	}
	
	return -1
}

/* original recursive version decompiled
var r0, r1, r7 int = 4, 1, _
var stack []int = []int{}
func energyLevelCheck() {
	if r0 != 0 {
		if r1 != 0 {
			stack = append([]int{r0}, stack...)
			r1 = (r1 + max) % mod
			energyLevelCheck()
			r1 = r0
			r0 = stack[0]
			stack = stack[1:]
			r0 = (r0 + max) % mod
			energyLevelCheck()
			return
		}
		r0 = (r0 + max) % mod
		r1 = r7
		energyLevelCheck()
		return
	}
	r0 = (r1 + 1) % mod
	return
}
*/

// https://rosettacode.org/wiki/Ackermann_function#Go
func ackermann(original [3]int) int {
	r7 := original[2]
	var ackermanMemoized func(regs [2]int) int
	ackermanMemoized = memoized(func(regs [2]int) int {
		if regs[0] == 0 {
			return (regs[1] + 1) % mod
		}
		if regs[1] == 0 {
			return ackermanMemoized([2]int { regs[0] - 1, r7 })
		}
		return ackermanMemoized([2]int { regs[0] - 1, ackermanMemoized([2]int { regs[0], regs[1] - 1 }) })
	})
	return ackermanMemoized([2]int{ original[0], original[1]})
}

// riff on the Go "Expanded version" from the link above
// couldn't figure out the form for the r0 == 3 case so it's really not much better
func ackermannBetter(original [3]int) int {
	r7 := original[2]
	var ackermanMemoized func(regs [2]int) int
	ackermanMemoized = memoized(func(regs [2]int) int {
		if regs[0] == 0 {
			return (regs[1] + 1) % mod
		}
		if regs[0] == 1 {
			return (regs[1] + r7 + 1) % mod
		}
		if regs[0] == 2 {
			return ((regs[1] + 2) * r7 + regs[1] + 1) % mod
		}
		if regs[1] == 0 {
			return ackermanMemoized([2]int { regs[0] - 1, r7 })
		}
		return ackermanMemoized([2]int { regs[0] - 1, ackermanMemoized([2]int { regs[0], regs[1] - 1 }) })
	})
	return ackermanMemoized([2]int{ original[0], original[1]})
}

func memoized(fn func(original [2]int) int) func([2]int) int {
	cache := make(map[string]int)
	return func(regs [2]int) int {
		if val, found := cache[fmt.Sprint(regs)]; found {
			return val
		}
		result := fn(regs)
		cache[fmt.Sprint(regs)] = result
		return result
	}
}

func pow(base int, power int) int {
	return int(math.Pow(float64(base), float64(power)))
}
