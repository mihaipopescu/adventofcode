package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func valid(number int) bool {
	valid := true

	var doubleCount [10]int

	prev := -1
	for n := number; n > 0; n = n / 10 {
		d := n % 10
		if prev >= 0 {
			if prev == d {
				doubleCount[d]++
			} else {
				if prev < d {
					valid = false
					break
				}
			}
		}
		prev = d
	}

	double := false
	for _, n := range doubleCount {
		if n == 1 {
			double = true
			break
		}
	}

	if valid && double {
		return true
	}

	return false
}

func main() {
	minS := os.Args[1]
	maxS := os.Args[2]
	rangeMin, err := strconv.Atoi(minS)
	if err != nil {
		log.Fatal(err)
	}
	rangeMax, err := strconv.Atoi(maxS)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	for i := rangeMin; i <= rangeMax; i++ {
		if valid(i) {
			fmt.Fprintf(os.Stderr, "%d\n", i)
			count++
		}
	}
	fmt.Printf("%d\n", count)
}
