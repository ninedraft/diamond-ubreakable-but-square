package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func main() {
	var n, _ = strconv.Atoi(os.Args[1])
	var cc = makeCells(n)
	var r = newRand(time.Now().UnixNano())
	var size = len(cc)
	fmt.Println("size", size)
	cc.set(0, 0, r.Byte())
	cc.set(0, size-1, r.Byte())
	cc.set(size-1, 0, r.Byte())
	cc.set(size-1, size-1, r.Byte())

	var radius = size / 2

	for ; radius > 0; radius /= 2 {
		for x := radius; x < size; x += 2 * radius {
			for y := radius; y < size; y += 2 * radius {
				r.doSquare(cc, x, y, radius)
			}
		}
		for x := 0; x < size; x += 2 * radius {
			for y := radius; y < size; y += 2 * radius {
				r.doDiamond(cc, x, y, radius)
			}
		}
		for x := radius; x < size; x += 2 * radius {
			for y := 0; y < size; y += 2 * radius {
				r.doDiamond(cc, x, y, radius)
			}
		}
	}

	fmt.Println()
	_, _ = cc.WriteTo(os.Stdout)
}

func (r *randGen) doSquare(cc cells, x, y, radius int) {
	var a = cc.get(x-radius, y-radius)
	var b = cc.get(x-radius, y+radius)
	var c = cc.get(x+radius, y+radius)
	var d = cc.get(x+radius, y-radius)
	var v = (int(a) + int(b) + int(c) + int(d)) / 4
	cc.set(x, y, byte(v)+r.ByteN(32))
}

func (r *randGen) doDiamond(cc cells, x, y, radius int) {
	var size = len(cc)
	var sum int
	var counter int
	if x-radius > 0 {
		sum += int(cc.get(x-radius, y))
		counter++
	}
	if x+radius < size {
		sum += int(cc.get(x+radius, y))
		counter++
	}
	if y-radius > 0 {
		sum += int(cc.get(x, y-radius))
		counter++
	}
	if y+radius < size {
		sum += int(cc.get(x, y+radius))
		counter++
	}
	var v = sum / counter
	cc.set(x, y, byte(v)+r.ByteN(32))
}

type cells [][]uint8

func (c cells) set(x, y int, v uint8) {
	if x < 0 || y < 0 ||
		x >= len(c) || y >= len(c) {
		return
	}
	c[x][y] = v
}

func (c cells) get(x, y int) uint8 {
	if x < 0 || y < 0 ||
		x >= len(c) || y >= len(c) {
		return 0
	}
	return c[x][y]
}

func (c cells) WriteTo(w io.Writer) (int64, error) {
	var buf = make([]byte, 1)
	var written int64
	var err error
	var write = func(b byte) {
		if err != nil {
			return
		}
		buf[0] = b
		var n, errWrite = w.Write(buf)
		if errWrite != nil {
			err = errWrite
		}
		written += int64(n)
	}
	for _, row := range c {
		for _, cell := range row {
			var n, errWrite = w.Write(showCell(cell))
			err = errWrite
			written += int64(n)
			if err != nil {
				break
			}
		}
		write('\n')
	}
	return written, err
}

const (
	cod = iota

	bulb = iota + 80
	beetle
	dragonfly

	kea = iota + 140
)

func showCell(b byte) []byte {
	switch b {
	case bulb:
		return []byte("○°")
	case cod:
		return []byte("~α")
	case beetle:
		return []byte("·ⁿ")
	case dragonfly:
		return []byte("ⁿ∙")
	case kea:
		return []byte("k♣")
	}
	switch {
	case b < 64:
		return []byte("~~")
	case b < 128:
		return []byte(".ⁿ")
	case b < 192:
		return []byte("♠♣")
	default:
		return []byte("▲▲")
	}
}

func makeCells(n int) cells {
	var size = 1<<n + 1
	var buf = make([]uint8, size*size)
	var cc = make(cells, size)
	for i := 0; i < size; i++ {
		cc[i] = buf[i*size : (i+1)*size : (i+1)*size]
	}
	return cc
}
