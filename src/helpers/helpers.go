package main

import (
	"fmt"
	"reflect"

	energylevel "github.com/ashlyn/synacor-challenge/src/energyLevel"
)

func main() {
	coins := []int{2,3,7,9,5}
	permutations := generatePermutations(coins)
	for _, p := range permutations {
		if checkCoins(p) {
			fmt.Println(p)
			break
		}
	}
	e := energylevel.GetMaxEnergyLevel(4, 1)
	fmt.Println(e)
	orb()
}

func checkCoins(coins []int) bool {
	if len(coins) != 5 {
		return false
	}
	return coins[0] + coins[1] * (coins[2] * coins[2]) + (coins[3] * coins[3] * coins[3]) - coins[4] == 399
}

func generatePermutations(items []int) [][]int {
	length := len(items)

	initial, itemsCopy := make([]int, length), make([]int, length)
	copy(initial, items)
	copy(itemsCopy, items)

	permutations := [][]int{initial}

	indexes := make([]int, length)

	i := 0
	for i < length {
		if indexes[i] < i {
			if i%2 == 0 {
				swap(itemsCopy, 0, i)
			} else {
				swap(itemsCopy, indexes[i], i)
			}
			permutation := make([]int, length)
			copy(permutation, itemsCopy)
			permutations = append(permutations, permutation)
			indexes[i] = indexes[i] + 1
			i = 0
		} else {
			indexes[i] = 0
			i++
		}
	}

	return permutations
}

func swap(slice interface{}, i int, j int) {
	if reflect.TypeOf(slice).Kind() == reflect.Slice {
		reflect.Swapper(slice)(i, j)
	}
}

func orb() {
	minus, plus, times := int('-'), int('+'), int('*')
	target := 30
	grid := [][]int{
		{times, 8, minus, 1},
		{4, times, 11, times},
		{plus, 4, minus, 18},
		{22, minus, 9, times},
	}
	queue := [][][2]int{{{0, 3}}}
	shortest := make([][2]int, 100)
	for len(queue) > 0 {
		path, loc := queue[0], queue[0][len(queue[0])-1]
		queue = queue[1:]
		if len(path) > len(shortest) {
			continue
		}
		if loc[0] == 3 && loc[1] == 0 {
			if getTotal(path, grid) == target {
				shortest = path
				pathStr, dirStr := pathToString(path, grid)
				println(pathStr)
				println(dirStr)
			}
			continue
		}
		neighbors := getValidNeighbors(loc, grid)
		for _, n := range neighbors {
			value := getValueFromGrid(n, grid)
			if value != 22 {
				newPath := make([][2]int, len(path))
				copy(newPath, path)
				newPath = append(newPath, n)
				queue = append(queue, newPath)
			}
		}
	}
}

func getTotal(path [][2]int, grid [][]int) int {
	total := getValueFromGrid(path[0], grid)
	for i := 1; i < len(path) - 1; i++ {
		current, next := getValueFromGrid(path[i], grid), getValueFromGrid(path[i+1], grid)
		if current == '-' {
			total -= next
			i++
		} else if current == '+' {
			total += next
			i++
		} else if current == '*' {
			total *= next
			i++
		}
	}

	return total
}

func getValueFromGrid(coords [2]int, grid [][]int) int {
	return grid[coords[1]][coords[0]]
}

func isInGrid(coords [2]int, grid [][]int) bool {
	return coords[0] >= 0 && coords[0] < len(grid[0]) &&
		coords[1] >= 0 && coords[1] < len(grid)
}

func getValidNeighbors(coords [2]int, grid [][]int) [][2]int {
	validNeighbors := [][2]int{}
	x, y := coords[0], coords[1]
	neighbors := [][2]int{
		{ x + 1, y },
		{ x - 1, y },
		{ x, y + 1 },
		{ x, y - 1 },
	}
	for _, n := range neighbors {
		if isInGrid(n, grid) {
			validNeighbors = append(validNeighbors, n)
		}
	}
	return validNeighbors
}

func pathToString(path [][2]int, grid [][]int) (string, string) {
	str, dirStr := "", ""
	for i, p := range path {
		value := getValueFromGrid(p, grid)
		if value == '-' {
			str += "- "
		} else if value == '+' {
			str += "+ "
		} else if value == '*' {
			str += "* "
		} else {
			str += fmt.Sprintf("%d ", value)
		}
		if i != 0 {
			dirStr += fmt.Sprintf("%s ", getDirection([2]int{ p[0] - path[i-1][0], p[1] - path[i-1][1] }))
		}
	}
	return str, dirStr
}

func getDirection(vec [2]int) string {
	x, y := vec[0], vec[1]
	if x == 0 && y == 1 {
		return "south"
	}
	if x == 0 && y == -1 {
		return "north"
	}
	if x == 1 && y == 0 {
		return "east"
	}
	if x == -1 && y == 0 {
		return "west"
	}
	return "unknown"
}
