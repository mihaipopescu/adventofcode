package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func checksum(data []byte, sx int, sy int) int {
	var layers []string

	for i := 0; i < len(data); i += sx * sy {
		var b bytes.Buffer
		for j := 0; j < sy; j++ {
			b.WriteString(string(data[i+j*sx:i+(j+1)*sx]) + "\n")
		}
		layers = append(layers, b.String())
	}

	var counters [][]int
	counters = make([][]int, len(layers))
	for i := range counters {
		counters[i] = make([]int, 3)
	}

	for i, l := range layers {
		counters[i][0] = strings.Count(l, "0")
		counters[i][1] = strings.Count(l, "1")
		counters[i][2] = strings.Count(l, "2")
	}

	fewestZeros := math.MaxInt32
	cs := 0
	for i, c := range counters {
		if c[0] < fewestZeros {
			fmt.Fprintf(os.Stderr, "found layer %d having fewest zeros %d\n", i, c[0])
			fewestZeros = c[0]
			cs = c[1] * c[2]
		}
	}

	return cs
}

func main() {
	var data []byte
	// if len(os.Args) != 4 {
	// 	log.Fatal("usage: program FILE SIZE_X SIZE_Y\n")
	// }
	// if os.Args[1] == "-" {
	// 	var err error
	// 	data, err = ioutil.ReadAll(bufio.NewReader(os.Stdin))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	var err error
	// 	data, err = ioutil.ReadFile(os.Args[1])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	data, _ = ioutil.ReadFile("in.txt")

	sx, err := strconv.Atoi("25")
	if err != nil {
		log.Fatal(err)
	}
	sy, err := strconv.Atoi("6")
	if err != nil {
		log.Fatal(err)
	}

	cs := checksum(data[0:len(data)-1], sx, sy)
	fmt.Printf("%d\n", cs)
}
