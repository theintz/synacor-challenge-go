package main

import (
	"errors"
	"strconv"
)

const OFFSET_REG = 32768

type memory interface {
	Set(addr word, value word) error
	Get(addr word) (word, error)
	GetRange(start word, end word) ([]word, error)
	Len() int
}

type registers interface {
	Set(idx word, value word) error
	Get(idx word) (word, error)
}

type mainMemory struct {
	data []word
}

func (m mainMemory) Set(addr word, value word) error {
	if addr >= word(len(m.data)) || addr < 0 {
		return errors.New("Address out of bounds: " + string(addr))
	}

	m.data[addr] = value

	return nil
}

func (m mainMemory) Get(addr word) (word, error) {
	if addr >= word(len(m.data)) || addr < 0 {
		return 0, errors.New("Address out of bounds: " + string(addr))
	}

	return m.data[addr], nil
}

func (m mainMemory) GetRange(start word, end word) ([]word, error) {
	if start > end || start < 0 || start >= word(len(m.data)) ||
		end < 0 || end >= word(len(m.data)) {
		return nil, errors.New("Addresses out of bounds: " + string(start) + ":" + string(end))
	}

	return m.data[start:end], nil
}

func (m mainMemory) Len() int {
	return len(m.data)
}

func initMemory(data []byte) (memory, error) {
	num := len(data) / 2
	words := make([]word, num)

	for i := 0; i < num; i++ {
		words[i] = bytesToWord(data[i * 2], data[i * 2 + 1])
	}

	return mainMemory{words}, nil
}

func initRegisters() registers {
	return mainMemory{make([]word, 8)}
}

func reg(n word) word {
	return n - OFFSET_REG
}

func value(n word, ctx *executionCtx) (word, error) {
	if n < OFFSET_REG {
		return n, nil
	} else if n < OFFSET_REG + 8 {
		return ctx.reg.Get(reg(n))
	} else {
		return 0, errors.New("invalid number: " + strconv.Itoa(int(n)))
	}
}

func bytesToWord(a, b byte) word {
	aUint := uint16(a)
	bUint := uint16(b)

	return word(aUint | bUint << 8)
}