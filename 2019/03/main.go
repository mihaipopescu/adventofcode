package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func BasenameWithoutExt(filename string) string {
	ext := path.Ext(filename)
	basename := path.Base(filename)
	return basename[0 : len(basename)-len(ext)]
}

func main() {
	filenamePtr := flag.String("input", "-", "Input file")
	debugPtr := flag.Bool("pgm", false, "Generates a pgm file (useful in debugging)")
	flag.Parse()

	filename := *filenamePtr

	var scanner *bufio.Scanner

	if filename == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		scanner = bufio.NewScanner(strings.NewReader(string(buf)))
	}

	onComma := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			if data[i] == ',' {
				return i + 1, data[:i], nil
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		return 0, data, bufio.ErrFinalToken
	}

	type move struct {
		dx int
		dy int
	}

	type wire struct {
		moves []move
	}

	var wires []wire

	type point struct {
		x int
		y int
	}

	var dimensions point
	var origin point

	for i := 0; scanner.Scan(); i++ {
		tokenizer := bufio.NewScanner(strings.NewReader(scanner.Text()))
		tokenizer.Split(onComma)

		// first two lines are dimensions and origin
		if i < 2 {
			var p point
			for k := 0; tokenizer.Scan(); k++ {
				v, err := strconv.Atoi(tokenizer.Text())
				if err != nil {
					log.Fatal(err)
				}
				if k == 0 {
					p.x = v
				} else if k == 1 {
					p.y = v
				}
			}
			if i == 0 {
				dimensions = p
			} else if i == 1 {
				origin = p
			}
			continue
		}

		// next, we read the wires instructions
		var w wire
		for step := 0; tokenizer.Scan(); step++ {
			re := regexp.MustCompile(`([LRUD])(\d+)`)
			groups := re.FindStringSubmatch(tokenizer.Text())
			normX := 0
			normY := 0
			if groups[1] == "L" {
				normX = -1
				normY = 0
			} else if groups[1] == "R" {
				normX = 1
				normY = 0
			} else if groups[1] == "U" {
				normX = 0
				normY = -1
			} else if groups[1] == "D" {
				normX = 0
				normY = 1
			}
			delta, err := strconv.Atoi(groups[2])
			if err != nil {
				log.Fatal(err)
			}

			w.moves = append(w.moves, move{dx: normX * delta, dy: normY * delta})

		}
		wires = append(wires, w)
	}

	// dynamically allocate matrix
	var matrix = make([][]int8, dimensions.y)
	for i := range matrix {
		matrix[i] = make([]int8, dimensions.x)
	}

	var cx = origin.x
	var cy = origin.y

	type intersection struct {
		point point
		steps int
	}

	var intersections []intersection

	// we run the wires instructions 3 times
	const passes = 3
	var wireIndexList = [3]byte{0, 1, 0}

	for pass := 0; pass < passes; pass++ {
		posX := cx
		posY := cy
		dist := 0
		i := wireIndexList[pass]
		var wireMask int8 = 1 << i
		for j := 0; j < len(wires[i].moves); j++ {
			m := wires[i].moves[j]
			x1 := posX
			x2 := posX + m.dx

			y1 := posY
			y2 := posY + m.dy

			fmt.Fprintf(os.Stderr, "move %d,%d to %d,%d\n", x1, y1, x2, y2)
			d := 0
			x := posX
			y := posY
			for k := 0; k < Abs(m.dx)+Abs(m.dy); k++ {
				matrix[y][x] = matrix[y][x] | wireMask // wires crossing themselves don't count = logical of power of two's
				if matrix[y][x] == 3 && dist > 0 {
					partialDist := dist + d
					fmt.Fprintf(os.Stderr, "Found intersection point at %d,%d with partial Dist=%d\n", y, x, partialDist)
					intersections = append(intersections, intersection{point: point{x: x, y: y}, steps: partialDist})
				}
				d++
				if m.dx < 0 {
					x--
				} else if m.dx > 0 {
					x++
				}
				if m.dy < 0 {
					y--
				} else if m.dy > 0 {
					y++
				}
			}

			posX = posX + m.dx
			posY = posY + m.dy
			dist += Abs(m.dx) + Abs(m.dy)
			fmt.Fprintf(os.Stderr, "cursor is at X=%d, Y=%d, dist=%d\n", posX, posY, dist)
		}
	}

	var uniqueIntersections []intersection
	// de-doubling of intersections
	for i := 0; i < len(intersections); i++ {
		for j := i + 1; j < len(intersections); j++ {
			if intersections[i].point == intersections[j].point {
				fmt.Fprintf(os.Stderr, "de-doubling %d,%d\n", i, j)
				uniqueIntersections = append(uniqueIntersections, intersection{
					point: intersections[i].point,
					steps: intersections[i].steps + intersections[j].steps})
			}
		}
	}

	minDist := math.MaxInt32
	minSteps := math.MaxInt32

	for _, p := range uniqueIntersections {
		d := Abs(p.point.x-cx) + Abs(p.point.y-cy)
		if d < minDist && d > 0 {
			minDist = d
		}
		if p.steps < minSteps {
			minSteps = p.steps
		}
	}

	fmt.Printf("minDist=%d\nminSteps=%d\n", minDist, minSteps)

	if *debugPtr {
		//create PGM file for debugging (see https://en.wikipedia.org/wiki/Netpbm_format#PGM_example)
		f, err := os.Create(fmt.Sprintf("%s.pgm", BasenameWithoutExt(filename)))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)

		w.WriteString(fmt.Sprintf("P2\n%d %d\n3\n", dimensions.x, dimensions.y))

		for i := 0; i < dimensions.y; i++ {
			var b strings.Builder
			for j := 0; j < dimensions.x; j++ {
				b.WriteString(fmt.Sprintf("%d ", matrix[i][j]))
			}
			w.WriteString(fmt.Sprintf("%s\n", b.String()))
		}

		w.Flush()
	}
}
