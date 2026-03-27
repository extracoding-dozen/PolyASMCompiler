// Package obfuscator реализует механизмы трансформации промежуточного представления (IR)
// для защиты кода от статического и динамического анализа.
package obfuscator

import (
	"math/rand"
	"time"

	"go.mod/pkg/ir"
	"go.mod/pkg/ir/ir_constants_and_types"
)

// ObfuscatorConfig определяет параметры конфигурации процесса обфускации.
type ObfuscatorConfig struct {
	// EnableSandboxNoise включает генерацию ложных системных вызовов.
	EnableSandboxNoise bool
	// NoiseFrequency задает вероятность вставки шумовых инструкций.
	NoiseFrequency int
	// EnableOpaquePreds включает создание непрозрачных предикатов.
	EnableOpaquePreds bool
	// OpaqueFrequency задает вероятность создания ложных ветвлений.
	OpaqueFrequency int
	// EnableStringCrypt включает XOR-шифрование строк.
	EnableStringCrypt bool
	// ObfuscateRepeat определяет количество проходов обфускации
	ObfuscateRepeat int
}

// DefaultObfuscatorConfig возвращает конфигурацию с активированными методами защиты по умолчанию.
func DefaultObfuscatorConfig() ObfuscatorConfig {
	return ObfuscatorConfig{
		EnableSandboxNoise: true,
		NoiseFrequency:     30,
		EnableOpaquePreds:  true,
		OpaqueFrequency:    20,
		EnableStringCrypt:  true,
		ObfuscateRepeat:    1,
	}
}

// Obfuscator выполняет полиморфные преобразования над списком инструкций IR.
type Obfuscator struct {
	ctx *ir.IRContext
	cfg ObfuscatorConfig
	rng *rand.Rand
}

// NewObfuscator создает новый экземпляр обфускатора с заданным контекстом и настройками.
func NewObfuscator(ctx *ir.IRContext, cfg ObfuscatorConfig) *Obfuscator {
	return &Obfuscator{
		ctx: ctx,
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (o *Obfuscator) SetNewConfig(newCfg ObfuscatorConfig) {
	o.cfg = newCfg
}

func (o *Obfuscator) Obfuscate(input []ir.Instruction) []ir.Instruction {
	o.cfg.ObfuscateRepeat %= 5
	if o.cfg.ObfuscateRepeat == 0 {
		return input
	}
	result := o.oneLayerObfuscate(input)
	o.cfg.ObfuscateRepeat -= 1
	for i := 0; i < o.cfg.ObfuscateRepeat; i++ {
		result = o.oneLayerObfuscate(result)
	}
	return result
}

// Obfuscate применяет последовательность трансформаций к входному набору инструкций.
func (o *Obfuscator) oneLayerObfuscate(input []ir.Instruction) []ir.Instruction {
	output := input
	if o.cfg.EnableStringCrypt {
		output = o.passStringCrypt(output)
	}
	if o.cfg.EnableSandboxNoise {
		output = o.passSandboxNoise(output)
	}
	if o.cfg.EnableOpaquePreds {
		output = o.passOpaquePredicates(output)
	}
	return output
}

// passSandboxNoise внедряет в поток команд безопасные системные вызовы для зашумления анализа strace.
func (o *Obfuscator) passSandboxNoise(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction
	safeSyscalls := []int64{39, 102, 104, 110}
	for _, inst := range input {
		output = append(output, inst)
		if o.rng.Intn(100) < o.cfg.NoiseFrequency {
			junkReg := o.ctx.NextVReg()
			randomSyscall := safeSyscalls[o.rng.Intn(len(safeSyscalls))]
			output = append(output, ir.Instruction{
				Op:   ir_constants_and_types.SYSCALL,
				Dst:  junkReg,
				Src1: ir.Imm{Value: randomSyscall},
				Args: []ir.Value{},
			})
		}
	}
	return output
}

// passOpaquePredicates вставляет логические условия, результат которых предопределен, для запутывания графа потока управления.
func (o *Obfuscator) passOpaquePredicates(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction
	for _, inst := range input {
		if inst.Op == ir_constants_and_types.LABEL || inst.Op == ir_constants_and_types.JMP || inst.Op == ir_constants_and_types.JE || inst.Op == ir_constants_and_types.JNE {
			output = append(output, inst)
			continue
		}
		if o.rng.Intn(100) < o.cfg.OpaqueFrequency {
			magicNum := int64(o.rng.Intn(9999) + 1)
			vReg := o.ctx.NextVReg()
			lblFake := o.ctx.NextLabel()
			lblEnd := o.ctx.NextLabel()
			output = append(output, ir.Instruction{Op: ir_constants_and_types.MOV, Dst: vReg, Src1: ir.Imm{Value: magicNum}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.CMP, Dst: vReg, Src1: ir.Imm{Value: magicNum}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.JNE, Dst: lblFake})
			output = append(output, inst)
			output = append(output, ir.Instruction{Op: ir_constants_and_types.JMP, Dst: lblEnd})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.LABEL, Dst: lblFake})
			junkReg1 := o.ctx.NextVReg()
			junkReg2 := o.ctx.NextVReg()
			output = append(output, ir.Instruction{Op: ir_constants_and_types.MOV, Dst: junkReg1, Src1: ir.Imm{Value: 0xDEADBEEF}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.ADD, Dst: junkReg2, Src1: junkReg1})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})
		} else {
			output = append(output, inst)
		}
	}
	return output
}

// passStringCrypt шифрует строковые константы XOR-ключом и добавляет инструкции для их динамической расшифровки.
func (o *Obfuscator) passStringCrypt(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction
	for _, inst := range input {
		if inst.Op == ir_constants_and_types.LOAD_STR {
			originalStr := inst.Src1.(ir.Str).Value
			dstReg := inst.Dst.(ir.VReg)
			strLen := len(originalStr)
			if strLen == 0 {
				output = append(output, inst)
				continue
			}
			xorKey := byte(o.rng.Intn(254) + 1)
			encryptedBytes := make([]byte, strLen)
			for i := 0; i < strLen; i++ {
				encryptedBytes[i] = originalStr[i] ^ xorKey
			}
			output = append(output, ir.Instruction{
				Op:  ir_constants_and_types.LOAD_STR,
				Dst: dstReg,
				Src1: ir.Str{
					Value: "",
					Bytes: encryptedBytes,
				},
			})
			output = append(output, ir.Instruction{
				Op:   ir_constants_and_types.MACRO_DECRYPT_STR,
				Dst:  dstReg,
				Src1: ir.Imm{Value: int64(xorKey)},
				Src2: ir.Imm{Value: int64(strLen)},
			})
		} else {
			output = append(output, inst)
		}
	}
	return output
}
