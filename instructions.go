package main

import (
	"errors"
	"fmt"
	"strconv"
	"os"
	"bufio"
)

const OP_HALT = 0
const OP_SET = 1
const OP_PUSH = 2
const OP_POP = 3
const OP_EQ = 4
const OP_GT = 5
const OP_JMP = 6
const OP_JT = 7
const OP_JF = 8
const OP_ADD = 9
const OP_MULT = 10
const OP_MOD = 11
const OP_AND = 12
const OP_OR = 13
const OP_NOT = 14
const OP_RMEM = 15
const OP_WMEM = 16
const OP_CALL = 17
const OP_RET = 18
const OP_OUT = 19
const OP_IN = 20
const OP_NOOP = 21

type Instruction interface {
	Handle(ctx *executionCtx) error
	Len() int
}

// 0: HALT

type Inst_Halt struct {
	// does not have any operands
}

func (i Inst_Halt) Handle(ctx *executionCtx) error {
	ctx.running = false

	return nil
}

func (i Inst_Halt) Len() int {
	return 1
}

// 1: SET

type Inst_Set struct {
	reg word
	value word
}

func (i Inst_Set) Handle(ctx *executionCtx) error {
	return ctx.reg.Set(i.reg, i.value)
}

func (i Inst_Set) Len() int {
	return 3
}

// 2: PUSH

type Inst_Push struct {
	val word
}

func (i Inst_Push) Handle(ctx *executionCtx) error {
	ctx.stack.Push(i.val)

	return nil
}

func (i Inst_Push) Len() int {
	return 2
}

// 3: POP

type Inst_Pop struct {
	reg word
}

func (i Inst_Pop) Handle(ctx *executionCtx) error {
	val := ctx.stack.Pop()
	if val == nil {
		return errors.New("stack is empty")
	}

	return ctx.reg.Set(i.reg, val.(word))
}

func (i Inst_Pop) Len() int {
	return 2
}

// 4: EQ

type Inst_Eq struct {
	reg word
	op1, op2 word
}

func (i Inst_Eq) Handle(ctx *executionCtx) error {
	var val word
	if i.op1 == i.op2 {
		val = 1
	} else {
		val = 0
	}

	return ctx.reg.Set(i.reg, val)
}

func (i Inst_Eq) Len() int {
	return 4
}

// 5: GT

type Inst_Gt struct {
	reg word
	op1, op2 word
}

func (i Inst_Gt) Handle(ctx *executionCtx) error {
	var val word
	if i.op1 > i.op2 {
		val = 1
	} else {
		val = 0
	}

	return ctx.reg.Set(i.reg, val)
}

func (i Inst_Gt) Len() int {
	return 4
}

// 6: JMP

type Inst_Jmp struct {
	target word
}

func (i Inst_Jmp) Handle(ctx *executionCtx) error {
	ctx.ip = i.target

	return nil
}

func (i Inst_Jmp) Len() int {
	return 2
}

// 7: JT

type Inst_Jt struct {
	cond word
	target word
}

func (i Inst_Jt) Handle(ctx *executionCtx) error {
	if i.cond != 0 {
		ctx.ip = i.target
	}

	return nil
}

func (i Inst_Jt) Len() int {
	return 3
}

// 8: JF

type Inst_Jf struct {
	cond word
	target word
}

func (i Inst_Jf) Handle(ctx *executionCtx) error {
	if i.cond == 0 {
		ctx.ip = i.target
	}

	return nil
}

func (i Inst_Jf) Len() int {
	return 3
}

// 9: ADD

type Inst_Add struct {
	reg word
	op1, op2 word
}

func (i Inst_Add) Handle(ctx *executionCtx) error {
	sum := (i.op1 + i.op2) % 32768
	return ctx.reg.Set(i.reg, sum)
}

func (i Inst_Add) Len() int {
	return 4
}

// 10: MULT

type Inst_Mult struct {
	reg word
	op1, op2 word
}

func (i Inst_Mult) Handle(ctx *executionCtx) error {
	prod := (i.op1 * i.op2) % 32768
	return ctx.reg.Set(i.reg, prod)
}

func (i Inst_Mult) Len() int {
	return 4
}

// 11: MOD

type Inst_Mod struct {
	reg word
	op1, op2 word
}

func (i Inst_Mod) Handle(ctx *executionCtx) error {
	res := i.op1 % i.op2
	return ctx.reg.Set(i.reg, res)
}

func (i Inst_Mod) Len() int {
	return 4
}

// 12: AND

type Inst_And struct {
	reg word
	op1, op2 word
}

func (i Inst_And) Handle(ctx *executionCtx) error {
	res := i.op1 & i.op2
	return ctx.reg.Set(i.reg, res)
}

func (i Inst_And) Len() int {
	return 4
}

// 13: OR

type Inst_Or struct {
	reg word
	op1, op2 word
}

func (i Inst_Or) Handle(ctx *executionCtx) error {
	res := i.op1 | i.op2
	return ctx.reg.Set(i.reg, res)
}

func (i Inst_Or) Len() int {
	return 4
}

// 14: NOT

type Inst_Not struct {
	reg word
	op1 word
}

func (i Inst_Not) Handle(ctx *executionCtx) error {
	// need to apply a bitmask here since our numbers are only allowed to be 15 bit long
	res := ^i.op1 & 0x7FFF
	return ctx.reg.Set(i.reg, res)
}

func (i Inst_Not) Len() int {
	return 3
}

// 15: RMEM

type Inst_Rmem struct {
	reg word
	addr word
}

func (i Inst_Rmem) Handle(ctx *executionCtx) error {
	val, err := ctx.mem.Get(i.addr)
	if err != nil {
		return err
	}

	return ctx.reg.Set(i.reg, val)
}

func (i Inst_Rmem) Len() int {
	return 3
}

// 16: WMEM

type Inst_Wmem struct {
	addr word
	val word
}

func (i Inst_Wmem) Handle(ctx *executionCtx) error {
	return ctx.mem.Set(i.addr, i.val)
}

func (i Inst_Wmem) Len() int {
	return 3
}

// 17: CALL

type Inst_Call struct {
	addr word
}

func (i Inst_Call) Handle(ctx *executionCtx) error {
	ctx.stack.Push(ctx.ip)
	ctx.ip = i.addr

	return nil
}

func (i Inst_Call) Len() int {
	return 2
}

// 18: RET

type Inst_Ret struct {
	// does not have any operands
}

func (i Inst_Ret) Handle(ctx *executionCtx) error {
	val := ctx.stack.Pop()
	if val == nil {
		return errors.New("stack is empty")
	}

	ctx.ip = val.(word)

	return nil
}

func (i Inst_Ret) Len() int {
	return 1
}

// 19: OUT

type Inst_Out struct {
	char word
}

func (i Inst_Out) Handle(ctx *executionCtx) error {
	_, err := fmt.Print(string(i.char))

	return err
}

func (i Inst_Out) Len() int {
	return 2
}

// 20: IN

type Inst_In struct {
	reg word
}

func (i Inst_In) Handle(ctx *executionCtx) error {
	// check whether we need to read new input
	if len(ctx.inBuf) < 1 {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadBytes('\n')
		if err != nil {
			ctx.inBuf = nil
			return err
		}

		ctx.inBuf = input
	}

	// check again whether there is data buffered
	if len(ctx.inBuf) < 1 {
		return errors.New("could not read any input from buffer")
	}

	// pick the first character from the buffer
	char := ctx.inBuf[0]
	// and drop it from the buffer
	ctx.inBuf = ctx.inBuf[1:]

	return ctx.reg.Set(i.reg, bytesToWord(char, 0))
}

func (i Inst_In) Len() int {
	return 2
}

// 21: NOOP

type Inst_Noop struct {
	// does not have any operands
}

func (i Inst_Noop) Handle(ctx *executionCtx) error {
	// does nothing on purpose
	return nil
}

func (i Inst_Noop) Len() int {
	return 1
}

func NewInstruction(words []word, ctx *executionCtx) (Instruction, error) {
	if len(words) < 4 {
		return nil, errors.New("too few words")
	}

	opcode := words[0]
	op1Reg := reg(words[1])
	op1, err := value(words[1], ctx)
	check(err)
	op2, err := value(words[2], ctx)
	check(err)
	op3, err := value(words[3], ctx)
	check(err)

	// do lookup with opcode
	switch opcode {
	case OP_HALT: return Inst_Halt{}, nil
	case OP_PUSH: return Inst_Push{op1}, nil
	case OP_POP: return Inst_Pop{op1Reg}, nil
	case OP_SET: return Inst_Set{op1Reg,op2}, nil
	case OP_EQ: return Inst_Eq{op1Reg, op2, op3}, nil
	case OP_GT: return Inst_Gt{op1Reg, op2, op3}, nil
	case OP_JMP: return Inst_Jmp{op1}, nil
	case OP_JT: return Inst_Jt{op1, op2}, nil
	case OP_JF: return Inst_Jf{op1, op2}, nil
	case OP_ADD: return Inst_Add{op1Reg, op2, op3}, nil
	case OP_MULT: return Inst_Mult{op1Reg, op2, op3}, nil
	case OP_MOD: return Inst_Mod{op1Reg, op2, op3}, nil
	case OP_AND: return Inst_And{op1Reg, op2, op3}, nil
	case OP_OR: return Inst_Or{op1Reg, op2, op3}, nil
	case OP_NOT: return Inst_Not{op1Reg, op2}, nil
	case OP_RMEM: return Inst_Rmem{op1Reg, op2}, nil
	case OP_WMEM: return Inst_Wmem{op1, op2}, nil
	case OP_CALL: return Inst_Call{op1}, nil
	case OP_RET: return Inst_Ret{}, nil
	case OP_OUT: return Inst_Out{op1}, nil
	case OP_IN: return Inst_In{op1Reg}, nil
	case OP_NOOP: return Inst_Noop{}, nil
	}

	return nil, errors.New("No Instruction found for opcode: " + strconv.Itoa(int(opcode)))
}
