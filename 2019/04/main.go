package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

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
		double := false
		valid := true

		prev := -1
		for n := i; n > 0; n = n / 10 {
			d := n % 10
			if prev >= 0 {
				if prev == d {
					double = true
				} else if prev < d {
					valid = false
					break
				}
			}
			prev = d
		}
		if valid && double {
			fmt.Fprintf(os.Stderr, "%d\n", i)
			count++
		}
	}
	fmt.Printf("%d\n", count)
}
