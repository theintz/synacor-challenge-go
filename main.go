package main

import (
	"fmt"
	"strconv"
	"errors"
)

const BINFILE = "./assets/challenge.bin"

// all memory is words
type word uint16

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func executeNext(ctx *executionCtx) error {
	if ctx.ip > word(ctx.mem.Len()) {
		return errors.New("ip greater than program length")
	}

	// first extract 8 bytes from progam data, this will make up our Instruction
	instData, err := ctx.mem.GetRange(ctx.ip, ctx.ip + 4)
	if err != nil {
		return err
	}

	inst, err := NewInstruction(instData, ctx)
	if err != nil {
		return err
	}

	// increase ip first, since it's overwritten by the jmp instructions again
	ctx.ip += word(inst.Len())

	//fmt.Printf("%d %v %s %+v\n", ctx.ip, instData, reflect.TypeOf(inst), inst)
	inst.Handle(ctx)

	return nil
}

func main() {
	// the execution context, containing the memory
	ctx, err := setupExeCtx(BINFILE)
	check(err)

	fmt.Println("Successfully read binary data with length: " + strconv.Itoa(ctx.mem.Len()))

	// main loop
	for ctx.running {
		// fetch and execute Instruction
		err = executeNext(&ctx)

		if err != nil {
			fmt.Println(err.Error())
			ctx.running = false
		}
	}

	fmt.Println("Finished execution")
}
