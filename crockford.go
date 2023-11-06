package main

import (
	"errors"
	"strconv"
)

var aliases = [][]byte{
	{'0', 'O', 'o'},
	{'1', 'I', 'L', 'i', 'l'},
	{'2'},
	{'3'},
	{'4'},
	{'5'},
	{'6'},
	{'7'},
	{'8'},
	{'9'},
	{'A', 'a'},
	{'B', 'b'},
	{'C', 'c'},
	{'D', 'd'},
	{'E', 'e'},
	{'F', 'f'},
	{'G', 'g'},
	{'H', 'h'},
	{'J', 'j'},
	{'K', 'k'},
	{'M', 'm'},
	{'N', 'n'},
	{'P', 'p'},
	{'Q', 'q'},
	{'R', 'r'},
	{'S', 's'},
	{'T', 't'},
	{'V', 'v'},
	{'W', 'w'},
	{'X', 'x'},
	{'Y', 'y'},
	{'Z', 'z'},
	{'*'},
	{'~'},
	{'$'},
	{'='},
	{'U', 'u'},
}

var u = []byte{
	// Encoding chars
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Y', 'Z',
	// Checksum chars
	'*', '~', '$', '=', 'U',
}

var d [256]byte // For decoding
var s [256]byte // For checksum

func init() {
	if len(u) != len(aliases) {
		panic("sth wrong with decode map")
	}

	for i := range d {
		d[i] = 0xff
		s[i] = 0xff
	}

	for a := range aliases[:32] {
		for _, c := range aliases[a] {
			d[c] = byte(a)
		}
	}

	for a := range aliases {
		for _, c := range aliases[a] {
			s[c] = byte(a)
		}
	}
}

var (
	CrockSumErr = errors.New("checksum mismatch")
	CrockChrErr = errors.New("contains illegal character")
)

func Crock(x uint64) string {
	if x == 0 {
		return "00"
	}

	var (
		b  [14]byte
		i  int
		x0 = x
		p  = 12
	)

	// Count the leading zeros in base 32.
	for ; x&(0x1f<<(p*5)) == 0 && p > 0; p-- {
	}

	for i = (12 - p); i < 13; i++ {
		b[i] = u[byte((x&(0x1f<<((12-i)*5)))>>((12-i)*5))]
	}
	b[i] = u[byte(x0%37)]
	i++

	return string(b[12-p : i])
}

func Uncrock(v string) (uint64, error) {
	if v == "" {
		return 0, strconv.ErrSyntax
	}
	if len(v) > 14 || len(v) == 14 && d[v[0]] > 15 {
		return 0, strconv.ErrRange
	}

	var (
		x uint64
	)

	v, k := v[:len(v)-1], s[v[len(v)-1]]
	for _, c := range v {
		if d[c] == 0xff {
			return 0, CrockChrErr
		}
		x <<= 5
		x |= uint64(d[c])
	}

	if byte(x%37) != k {
		return 0, CrockSumErr
	}

	return x, nil
}
