package main

import (
	"bufio"
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
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: program [INPUT_FILE]\n")
		os.Exit(1)
	}

	buf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(buf)))

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

	var flag int8 = 1

	for i := 0; i < len(wires); i++ {
		posX := cx
		posY := cy
		for j := 0; j < len(wires[i].moves); j++ {
			m := wires[i].moves[j]
			x1 := posX
			x2 := posX + m.dx
			if m.dx < 0 {
				x2, x1 = x1, x2
			}
			y1 := posY
			y2 := posY + m.dy
			if m.dy < 0 {
				y2, y1 = y1, y2
			}

			fmt.Fprintf(os.Stderr, "move %d,%d to %d,%d\n", x1, y1, x2, y2)
			for y := y1; y <= y2; y++ {
				for x := x1; x <= x2; x++ {
					matrix[y][x] = matrix[y][x] | flag // wires crossing themselves don't count = logical of power of two's
				}
			}
			posX = posX + m.dx
			posY = posY + m.dy
			fmt.Fprintf(os.Stderr, "cursor is at X=%d, Y=%d\n", posX, posY)
		}
		flag = flag << 1
	}

	var intersections []point

	for i := 0; i < dimensions.y; i++ {
		for j := 0; j < dimensions.x; j++ {
			if matrix[i][j] == 3 {
				fmt.Fprintf(os.Stderr, "Found intersection point at %d,%d\n", j, i)
				intersections = append(intersections, point{x: j, y: i})
			}
		}
	}

	minDist := math.MaxInt32

	for _, p := range intersections {
		d := Abs(p.x-cx) + Abs(p.y-cy)
		if d < minDist && d > 0 {
			minDist = d
		}
	}

	fmt.Printf("%d\n", minDist)

	//create PGM file for debugging (see https://en.wikipedia.org/wiki/Netpbm_format#PGM_example)

	f, err := os.Create(fmt.Sprintf("%s.pgm", BasenameWithoutExt(os.Args[1])))
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
