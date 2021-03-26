package main

import "math/rand"

type randGen struct {
	r *rand.Rand
}

func newRand(seed int64) *randGen {
	var source = rand.NewSource(seed)
	return &randGen{
		r: rand.New(source),
	}
}

func (r *randGen) Byte() byte {
	return byte(r.r.Intn(256))
}

func (r *randGen) ByteN(n byte) byte {
	var max = int(n)
	return byte(r.r.Intn(max) - max/2)
}
