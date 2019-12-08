package main

import (
	"fmt"
	"testing"
)

var testData = []struct {
	number int
	valid  bool
}{
	{112233, true},
	{123444, false},
	{111122, true},
	{588999, true},
}

func TestNumbersAreValid(t *testing.T) {
	for _, tt := range testData {
		t.Run(fmt.Sprintf("Testing %d", tt.number), func(t *testing.T) {
			v := valid(tt.number)
			if tt.valid != v {
				t.Errorf("got %v, want %v", v, tt.valid)
			}
		})
	}
}
