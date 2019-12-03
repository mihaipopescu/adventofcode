package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// https://adventofcode.com/2019/day/1
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var masses []int
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		mass, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalf("error parsing '%s' to int", line)
		}
		masses = append(masses, mass)
	}

	// part 2, integrate
	isIntegration := true

	var fuel int
	fuel = 0

	for _, mass := range masses {
		for fuel4mass := mass; isIntegration; {
			fuel4mass = (fuel4mass / 3) - 2
			if fuel4mass > 0 {
				fuel += fuel4mass
			} else {
				break
			}
		}
	}

	fmt.Printf("%d\n", fuel)
}
