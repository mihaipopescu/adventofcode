package main

import (
	"testing"
)

func Test1(t *testing.T) {
	data := []byte("123456789012")
	sx := 3
	sy := 2

	e := 1
	c := checksum(data, sx, sy)

	if c != e {
		t.Errorf("got %v expected %v", c, e)
	}

}
