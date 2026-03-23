package ir

import "go.mod/internal/ir/ir_constants_and_types"

type MacroExpander struct {
	ctx *IRContext
}

func NewMacroExpander(ctx *IRContext) *MacroExpander {
	return &MacroExpander{ctx: ctx}
}

// Главная функция. Принимает High-Level IR, возвращает Low-Level IR.
func (e *MacroExpander) Expand(input []Instruction) []Instruction {
	var output []Instruction

	for _, inst := range input {
		switch inst.Op {

		case ir_constants_and_types.MACRO_COPY:
			// Распутываем copy(src, dst)
			expandedBlock := e.expandCopy(inst.Args[0], inst.Args[1])
			// Вставляем результат распутывания вместо оригинальной инструкции
			output = append(output, expandedBlock...)

		case ir_constants_and_types.MACRO_USERADD:
			expandedBlock := e.expandUserAdd(inst.Args[0], inst.Args[1])
			output = append(output, expandedBlock...)
		case ir_constants_and_types.MACRO_WRITE:
			// Выдаем высокоуровневый макрос})
			output = append(output, e.expandWrite(inst.Args[0], inst.Args[1], inst.Args[2])...)
		case ir_constants_and_types.MACRO_GET_FILE_SIZE:
			expandedBlock := e.expandGetFileSize(inst.Dst, inst.Args[0])
			output = append(output, expandedBlock...)
		default:
			// Обычные инструкции (MOV, ADD, SYSCALL 59) пробрасываем без изменений
			output = append(output, inst)
		}
	}

	return output
}

func (e *MacroExpander) expandCopy(srcReg Value, dstReg Value) []Instruction {
	var block []Instruction

	// Запрашиваем новые виртуальные регистры для файловых дескрипторов
	fdIn := e.ctx.NextVReg()
	fdOut := e.ctx.NextVReg()

	// 1. fd_in = sys_open(src, 0 /* O_RDONLY */)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdIn,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, // номер sys_open
		Args: []Value{srcReg, Imm{Value: 0}, Imm{Value: 0}},
	})

	// 2. fd_out = sys_open(dst, 65 /* O_WRONLY|O_CREAT */, 0777)
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{dstReg, Imm{Value: 65}, Imm{Value: 511}}, // 511 = 0777 в восьмеричной
	})

	// 3. Вызываем sys_sendfile (номер 40 в x86_64).
	// Он копирует данные напрямую в ядре, не гоняя их в UserSpace!
	// sys_sendfile(out_fd, in_fd, offset, count)
	// count сделаем огромным числом, например 0x7fffffff, чтобы скопировать за 1 раз.
	resultReg := e.ctx.NextVReg()
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  resultReg,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_SEND_FILE}, // номер sys_sendfile
		Args: []Value{fdOut, fdIn, Imm{Value: 0}, Imm{Value: 2147483647}},
	})

	// 4. Закрываем дескрипторы
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

	return block
}

func (e *MacroExpander) expandWrite(pathReg, offsetReg, dataReg Value) []Instruction {
	var block []Instruction

	fdOut := e.ctx.NextVReg()

	// 1. fd = sys_open(path, O_WRONLY | O_CREAT (65), 0644 (420))
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN}, // sys_open
		Args: []Value{pathReg, Imm{Value: 65}, Imm{Value: 420}},
	})

	// 2. sys_lseek(fd, offset, SEEK_SET (0))
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  e.ctx.NextVReg(),
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_LSEEK}, // sys_lseek
		Args: []Value{fdOut, offsetReg, Imm{Value: 0}},
	})

	// 3. Вычисляем длину строки (sys_write требует указать размер в байтах)
	// В шеллкодах длину строки обычно считают в цикле, но для упрощения
	// мы пока передадим фиксированный размер (например 100 байт).
	// В идеале сюда нужно добавить макрос strlen.
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

	return block
}

func (e *MacroExpander) expandUserAdd(userReg, passReg Value) []Instruction {
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
	block = append(block, Instruction{
		Op:   ir_constants_and_types.SYSCALL,
		Dst:  fdOut,
		Src1: Imm{Value: ir_constants_and_types.SYSCALL_OPEN},
		Args: []Value{passwdPathReg, Imm{Value: 1025}, Imm{Value: 420}},
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
		Args: []Value{fdReg, Imm{Value: 0}},
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
