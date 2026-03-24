package obfuscator

import (
	"math/rand"
	"time"

	"go.mod/pkg/ir"
	"go.mod/pkg/ir/ir_constants_and_types"
)

type ObfuscatorConfig struct {
	EnableSandboxNoise bool
	NoiseFrequency     int

	EnableOpaquePreds bool
	OpaqueFrequency   int

	EnableStringCrypt bool
}

func DefaultObfuscatorConfig() ObfuscatorConfig {
	return ObfuscatorConfig{
		EnableSandboxNoise: true,
		NoiseFrequency:     30,
		EnableOpaquePreds:  true,
		OpaqueFrequency:    20,
		EnableStringCrypt:  true,
	}
}

type Obfuscator struct {
	ctx *ir.IRContext
	cfg ObfuscatorConfig
	rng *rand.Rand
}

func NewObfuscator(ctx *ir.IRContext, cfg ObfuscatorConfig) *Obfuscator {
	return &Obfuscator{
		ctx: ctx,
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Obfuscate — ЕДИНСТВЕННЫЙ ПУБЛИЧНЫЙ МЕТОД (Оркестратор)
func (o *Obfuscator) Obfuscate(input []ir.Instruction) []ir.Instruction {
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

// passSandboxNoise вставляет безопасные сисколлы (getpid, getuid), чтобы заспамить strace песочницы
func (o *Obfuscator) passSandboxNoise(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction

	safeSyscalls := []int64{
		39,
		102,
		104,
		110,
	}

	for _, inst := range input {
		output = append(output, inst) // Добавляем оригинальную инструкцию

		// С вероятностью NoiseFrequency вставляем мусорный сисколл
		if o.rng.Intn(100) < o.cfg.NoiseFrequency {
			junkReg := o.ctx.NextVReg()
			randomSyscall := safeSyscalls[o.rng.Intn(len(safeSyscalls))]

			output = append(output, ir.Instruction{
				Op:   ir_constants_and_types.SYSCALL,
				Dst:  junkReg,
				Src1: ir.Imm{Value: randomSyscall},
				Args: []ir.Value{}, // Аргументы не нужны
			})
		}
	}
	return output
}

// passOpaquePredicates создает "Непрозрачные предикаты".
// Это условия, результат которых мы (компилятор) знаем заранее, а дизассемблер/анализатор — нет.
func (o *Obfuscator) passOpaquePredicates(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction

	for _, inst := range input {
		// Не ломаем структуру уже существующих меток и прыжков
		if inst.Op == ir_constants_and_types.LABEL || inst.Op == ir_constants_and_types.JMP || inst.Op == ir_constants_and_types.JE || inst.Op == ir_constants_and_types.JNE {
			output = append(output, inst)
			continue
		}

		if o.rng.Intn(100) < o.cfg.OpaqueFrequency {
			// Генерируем ложный блок кода.
			// Логика:
			//   vX = СЛУЧАЙНОЕ_ЧИСЛО
			//   cmp vX, СЛУЧАЙНОЕ_ЧИСЛО
			//   jne .L_FAKE  (Никогда не выполнится, так как они равны!)
			//   ... оригинальная инструкция ...
			//   jmp .L_END
			// .L_FAKE:
			//   ... мусорный код ...
			// .L_END:

			magicNum := int64(o.rng.Intn(9999) + 1)
			vReg := o.ctx.NextVReg()
			lblFake := o.ctx.NextLabel()
			lblEnd := o.ctx.NextLabel()

			// Устанавливаем предикат (всегда TRUE)
			output = append(output, ir.Instruction{Op: ir_constants_and_types.MOV, Dst: vReg, Src1: ir.Imm{Value: magicNum}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.CMP, Dst: vReg, Src1: ir.Imm{Value: magicNum}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.JNE, Dst: lblFake}) // Прыжок на фейк, если НЕ равно

			// Реальный код
			output = append(output, inst)
			output = append(output, ir.Instruction{Op: ir_constants_and_types.JMP, Dst: lblEnd}) // Пропускаем фейк

			// Ложный код (Junk Block - Дизассемблер IDA Pro сойдет с ума, пытаясь это проанализировать)
			output = append(output, ir.Instruction{Op: ir_constants_and_types.LABEL, Dst: lblFake})
			junkReg1 := o.ctx.NextVReg()
			junkReg2 := o.ctx.NextVReg()
			output = append(output, ir.Instruction{Op: ir_constants_and_types.MOV, Dst: junkReg1, Src1: ir.Imm{Value: 0xDEADBEEF}})
			output = append(output, ir.Instruction{Op: ir_constants_and_types.ADD, Dst: junkReg2, Src1: junkReg1})

			// Конец конструкции
			output = append(output, ir.Instruction{Op: ir_constants_and_types.LABEL, Dst: lblEnd})
		} else {
			output = append(output, inst)
		}
	}
	return output
}

// passStringCrypt шифрует строки XOR-ом во время компиляции
// и вставляет новую IR-инструкцию MACRO_DECRYPT_STR для расшифровки в рантайме.
func (o *Obfuscator) passStringCrypt(input []ir.Instruction) []ir.Instruction {
	var output []ir.Instruction

	for _, inst := range input {
		if inst.Op == ir_constants_and_types.LOAD_STR {
			// Оригинальная строка (взятая из исходника)
			originalStr := inst.Src1.(ir.Str).Value
			dstReg := inst.Dst.(ir.VReg)
			strLen := len(originalStr)

			if strLen == 0 {
				output = append(output, inst)
				continue
			}

			// 1. Генерируем случайный XOR-ключ (от 1 до 255)
			xorKey := byte(o.rng.Intn(254) + 1)

			// 2. Шифруем строку в сырой массив байт
			encryptedBytes := make([]byte, strLen)
			for i := 0; i < strLen; i++ {
				encryptedBytes[i] = originalStr[i] ^ xorKey
			}

			// 3. Выдаем LOAD_STR, но передаем сырые байты, а не строку!
			output = append(output, ir.Instruction{
				Op:  ir_constants_and_types.LOAD_STR,
				Dst: dstReg,
				Src1: ir.Str{
					Value: "",             // Текст больше не нужен
					Bytes: encryptedBytes, // Передаем сырой зашифрованный массив
				},
			})

			// 4. Сразу выдаем инструкцию на расшифровку
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
