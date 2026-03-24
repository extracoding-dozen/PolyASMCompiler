package ir

import "go.mod/pkg/ir/ir_constants_and_types"

type MacroExpander struct {
	ctx *IRContext
}

func NewMacroExpander(ctx *IRContext) *MacroExpander {
	return &MacroExpander{ctx: ctx}
}

func (e *MacroExpander) GetCtx() *IRContext {
	return e.ctx
}

func (e *MacroExpander) Expand(input []Instruction) []Instruction {
	var output []Instruction

	for _, inst := range input {
		switch inst.Op {

		case ir_constants_and_types.MACRO_COPY:
			expandedBlock := e.expandCopy(inst.Dst, inst.Args[0], inst.Args[1])
			output = append(output, expandedBlock...)

		case ir_constants_and_types.MACRO_USERADD:
			expandedBlock := e.expandUserAdd(inst.Dst, inst.Args[0], inst.Args[1])
			output = append(output, expandedBlock...)
		case ir_constants_and_types.MACRO_WRITE:
			output = append(output, e.expandWrite(inst.Dst, inst.Args[0], inst.Args[1], inst.Args[2])...)
		case ir_constants_and_types.MACRO_GET_FILE_SIZE:
			expandedBlock := e.expandGetFileSize(inst.Dst, inst.Args[0])
			output = append(output, expandedBlock...)
		case ir_constants_and_types.MACRO_SLEEP:
			expandedBlock := e.expandSleep(inst.Dst, inst.Args[0])
			output = append(output, expandedBlock...)
		default:
			output = append(output, inst)
		}
	}

	return output
}

func (e *MacroExpander) expandCopy(resReg, srcReg Value, dstReg Value) []Instruction {
	var block []Instruction

	fdIn := e.ctx.NextVReg()
	fdOut := e.ctx.NextVReg()
	errLabel := e.ctx.NextLabel()

	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdIn,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{srcReg, Imm{Value: 0}, Imm{Value: 0}},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.CMP,
		Dst:  fdIn,
		Src1: Imm{Value: 0},
	})

	block = append(block, Instruction{
		Op:  ir_constants_and_types.JL,
		Dst: errLabel,
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{dstReg, Imm{Value: 65}, Imm{Value: 511}},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.CMP,
		Dst:  fdOut,
		Src1: Imm{Value: 0},
	})

	block = append(block, Instruction{
		Op:  ir_constants_and_types.JL,
		Dst: errLabel,
	})

	resultReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  resultReg,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_SEND_FILE},
		Args: []Value{fdOut, fdIn, Imm{Value: 0}, Imm{Value: 2147483647}},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE},
		Args: []Value{fdIn},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE},
		Args: []Value{fdOut},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resultReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS},
	})
	block = append(block, Instruction{
		Op:  ir_constants_and_types.LABEL,
		Dst: errLabel,
	})
	return block
}

func (e *MacroExpander) expandWrite(resReg, pathReg, offsetReg, dataReg Value) []Instruction {
	var block []Instruction

	fdOut := e.ctx.NextVReg()
	errLabel := e.ctx.NextLabel()

	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{pathReg, Imm{Value: 65}, Imm{Value: 420}},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.CMP,
		Dst:  fdOut,
		Src1: Imm{Value: 0},
	})

	block = append(block, Instruction{
		Op:  ir_constants_and_types.JL,
		Dst: errLabel,
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_LSEEK},
		Args: []Value{fdOut, offsetReg, Imm{Value: 0}},
	})

	dataSizeReg := e.ctx.NextVReg() // Регистр, куда ляжет длина строки
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MACRO_STRLEN, // Тот самый макрос
		Dst:  dataSizeReg,                         // Сохраняем длину сюда
		Src1: dataReg,                             // Считаем длину этой строки
	})

	// 4. sys_write(fd, data, size)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE}, // sys_write
		Args: []Value{fdOut, dataReg, dataSizeReg},
	})

	// 5. sys_close(fd)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE}, // sys_close
		Args: []Value{fdOut},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS},
	})
	block = append(block, Instruction{
		Op:  ir_constants_and_types.LABEL,
		Dst: errLabel,
	})

	return block
}

func (e *MacroExpander) expandUserAdd(resReg, userReg, passReg Value) []Instruction {
	var block []Instruction

	// Нам понадобится строка "/etc/passwd". Мы создаем её прямо здесь!
	passwdPathReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.LOAD_STR,
		Dst:  passwdPathReg,
		Src1: Str{Value: "/etc/passwd"},
	})

	// 1. Открываем /etc/passwd на дозапись: fd = open("/etc/passwd", O_WRONLY | O_APPEND)
	// O_WRONLY(1) | O_APPEND(1024) = 1025
	fdOut := e.ctx.NextVReg()
	errorLabel := e.ctx.NextLabel()

	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_ERROR},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{passwdPathReg, Imm{Value: 1025}, Imm{Value: 420}},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.CMP,
		Dst:  fdOut,
		Src1: Imm{Value: 0},
	})

	block = append(block, Instruction{
		Op:  ir_constants_and_types.JL,
		Dst: errorLabel,
	})
	// 2. Пишем логин (userReg)
	// size = strlen(user)
	userSizeReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MACRO_STRLEN,
		Dst:  userSizeReg,
		Src1: userReg,
	})
	// write(fd, user, size)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE},
		Args: []Value{fdOut, userReg, userSizeReg},
	})

	// 3. Пишем хвост строки (":x:0:0::/root:/bin/bash\n")
	// Создаем эту строку:
	tailReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.LOAD_STR,
		Dst:  tailReg,
		Src1: Str{Value: ":x:0:0::/root:/bin/bash\n"},
	})

	// Хвост константный, мы знаем его длину (24 символа). Не нужно считать strlen.
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_WRITE},
		Args: []Value{fdOut, tailReg, Imm{Value: 24}},
	})

	// 4. Закрываем /etc/passwd
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE},
		Args: []Value{fdOut},
	})

	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resReg,
		Src1: Imm{Value: ir_constants_and_types.RETURN_SUCCESS},
	})

	block = append(block, Instruction{
		Op:  ir_constants_and_types.LABEL,
		Dst: errorLabel,
	})
	// (Для полного useradd нужно еще записать пароль в /etc/shadow,
	// но логика абсолютно такая же. Оставим как домашнее задание).

	return block
}

func (e *MacroExpander) expandGetFileSize(resultReg, filepathReg Value) []Instruction {
	var block []Instruction

	// Нам понадобится временный виртуальный регистр для хранения файлового дескриптора (fd)
	fdReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  resultReg,
		Src1: Imm{Value: -1},
	})
	// 1. Открываем файл: fd = sys_open(filepath, O_RDONLY)
	// Флаг O_RDONLY равен 0. Права доступа (третий аргумент) при чтении не важны, передаем 0.
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdReg,                                           // Куда сохранить fd
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, // Номер системного вызова
		Args: []Value{filepathReg, Imm{Value: 0}, Imm{Value: 0}},
	})
	block = append(block, Instruction{
		Op:   ir_constants_and_types.CMP,
		Dst:  fdReg,
		Src1: Imm{Value: 0},
	})
	errLaber := e.ctx.NextLabel()
	block = append(block, Instruction{
		Op:  ir_constants_and_types.JL,
		Dst: errLaber,
	})
	// 2. Узнаем размер: size = sys_lseek(fd, 0, SEEK_END)
	// sys_lseek(fd, offset, origin).
	// offset = 0 (сдвиг ноль байт).
	// origin = 2 (SEEK_END - считать от конца файла).
	// lseek вернет текущую позицию от начала файла, что равно его размеру.
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  resultReg, // ВАЖНО: пишем результат сразу в целевой регистр!
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_LSEEK},
		Args: []Value{fdReg, Imm{Value: 0}, Imm{Value: 2}},
	})

	// 3. Закрываем файл: sys_close(fd)
	// Закрывать файлы обязательно, иначе операционная система исчерпает лимит дескрипторов.
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(), // Выделяем мусорный регистр, результат закрытия нам не нужен
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_CLOSE},
		Args: []Value{fdReg},
	})
	block = append(block, Instruction{
		Op:  ir_constants_and_types.LABEL,
		Dst: errLaber,
	})
	return block
}

func (e *MacroExpander) expandSleep(resultReg, sleepTime Value) []Instruction {
	var block []Instruction

	// --- 1. ФОРМИРУЕМ СТРУКТУРУ TIMESPEC НА СТЕКЕ ---
	// Нам нужно 16 байт (2 виртуальных регистра подряд).
	// Запрашиваем первый регистр — он будет полем tv_sec (секунды).
	tvSecReg := e.ctx.NextVReg()

	// Запрашиваем второй регистр — он будет полем tv_nsec (наносекунды).
	tvNsecReg := e.ctx.NextVReg()

	// --- 2. ЗАПОЛНЯЕМ СТРУКТУРУ ---
	// Записываем переданное количество секунд в tv_sec
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  tvSecReg,
		Src1: Imm{Value: 0}, // Это может быть число (Imm{5}) или переменная (vReg)
	})

	// Записываем 0 в tv_nsec (нам нужны ровные секунды)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.MOV,
		Dst:  tvNsecReg,
		Src1: sleepTime,
	})

	// --- 3. ПОЛУЧАЕМ УКАЗАТЕЛЬ НА СТРУКТУРУ ---
	// Ядру нужен *указатель* (адрес) на эту структуру.
	// Так как tvSecReg — это просто смещение на стеке (например [rbp - 40]),
	// нам нужна новая псевдо-инструкция, чтобы получить реальный адрес в памяти (lea reg, [rbp - 40]).
	structPtrReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.LEA, // <-- НОВЫЙ ОПКОД (нужно добавить)
		Dst:  structPtrReg,
		Src1: tvSecReg, // Берем адрес начала структуры (tv_sec)
	})

	// --- 4. ВЫЗЫВАЕМ СИСКОЛЛ NANOSLEEP (35) ---
	// sys_nanosleep(*req, *rem)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  resultReg,                                            // Обычно nanosleep возвращает 0 при успехе
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_NANOSLEEP}, // 35
		// Передаем указатель на нашу структуру (req) и 0 (NULL) для rem (остаток времени нас не волнует)
		Args: []Value{structPtrReg, Imm{Value: 0}},
	})

	return block
}
