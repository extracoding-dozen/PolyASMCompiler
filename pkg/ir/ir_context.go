// Package ir предоставляет инструменты для управления контекстом и генерации промежуточного представления кода.
package ir

import "fmt"

// IRContext хранит состояние генерации промежуточного представления, отслеживая счетчики регистров и меток.
type IRContext struct {
	vregCounter  int
	labelCounter int
}

// NextVReg генерирует и возвращает новый уникальный виртуальный регистр.
func (ctx *IRContext) NextVReg() VReg {
	ctx.vregCounter++
	return VReg{ID: ctx.vregCounter}
}

// NextLabel генерирует и возвращает новую уникальную текстовую метку для переходов.
func (ctx *IRContext) NextLabel() Lbl {
	ctx.labelCounter++
	return Lbl{Name: fmt.Sprintf(".L%d", ctx.labelCounter)}
}

// GetVRegCount возвращает общее количество выделенных виртуальных регистров в текущем контексте.
func (ctx *IRContext) GetVRegCount() int {
	return ctx.vregCounter
}
