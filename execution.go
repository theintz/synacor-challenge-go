package main

import (
	"io/ioutil"
	"github.com/golang-collections/collections/stack"
)

type executionCtx struct {
	// the Instruction pointer into main memory
	ip word
	// a heap, contains the program at start
	mem memory
	// our stack
	stack stack.Stack
	// the registers
	reg registers
	// whether or not we are running
	running bool
	// our buffered input
	inBuf []byte
}

func setupExeCtx(fileName string) (executionCtx, error) {
	progData, err := ioutil.ReadFile(fileName)
	check(err)

	mainMem, err := initMemory(progData)
	check(err)

	regMem := initRegisters()

	goStack := stack.New()

	return executionCtx{ip: 0, mem: mainMem, stack: *goStack, reg: regMem, running: true, inBuf: nil}, err
}