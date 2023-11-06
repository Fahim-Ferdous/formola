package main

import (
	"testing"
)

func Test_Crock(t *testing.T) {
	b := []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64}
	for x := range b {
		if b[x] >= 64 {
			break
		}
		c := (1 << b[1]) - 1
		i := uint64(0)
		if x != 0 {
			i = 1 << uint64(b[x-1])
		}
		for ; ; i += (1 << b[x]) {
			v := Crock(i)
			a, b := Uncrock(v)
			if a != i || b != nil {
				t.Log(i, v, a, b)
				t.FailNow()
			}
			if c == 0 {
				break
			}
			c--
		}
	}
}
