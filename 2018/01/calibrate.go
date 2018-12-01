package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	b, err := ioutil.ReadFile("input.txt")
	if err != nil {
		fmt.Print(err)
	}

	str := string(b)
	deltaArr := strings.Split(str, "\n")

	partTwo := true

	known := make(map[int]bool)
	known[0] = true
	
	sum := 0
	i := 0
	for {
		if i == len(deltaArr) {
			i = 0
		}
		elem := deltaArr[i]
		i++
		if len(elem) > 0 {
			delta, err := strconv.Atoi(elem)
			if err == nil {
				prevSum := sum
				sum = sum + delta

				if partTwo {
					value, ok := known[sum]
					if ok {
						if value {
							fmt.Fprintf(os.Stderr, "first reaches %d twice.\n", sum)
							fmt.Println(sum)
							break
						}
					} else {
						known[sum] = true
					}
				} else {
					fmt.Fprintf(os.Stderr, "Current frequency  %d, change of %d; resulting frequency  %d.\n", prevSum, delta, sum)
					if i == len(deltaArr)-1 {
						fmt.Println(sum)
						break
					}
				}
			}
		}
	}
}
