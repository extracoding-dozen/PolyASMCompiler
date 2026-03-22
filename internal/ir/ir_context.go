package ir

import "fmt"

type IRContext struct {
	vregCounter  int
	labelCounter int
}

func (ctx *IRContext) NextVReg() VReg {
	ctx.vregCounter++
	return VReg{ID: ctx.vregCounter}
}

func (ctx *IRContext) NextLabel() Lbl {
	ctx.labelCounter++
	return Lbl{Name: fmt.Sprintf(".L%d", ctx.labelCounter)}
}

func (ctx *IRContext) GetVRegCount() int {
	return ctx.vregCounter
}
