package main

import (
	"fmt"
	"reflect"
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