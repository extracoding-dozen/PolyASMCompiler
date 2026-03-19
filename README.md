# PolyASM Payload Builder

**PolyASM Builder** — это специализированный фреймворк (сборочный цех) для генерации высокозащищенных полиморфных shellcode-пейлоадов под Linux. 

Фреймворк компилирует исходный код, написанный на кастомном языке программирования, проводит множественные проходы обфускации на уровне промежуточного представления (IR), генерирует PIE Assembly, собирает его в бинарный файл и автоматически извлекает секцию `.text`, выдавая на выходе чистый raw-shellcode (`payload.bin`), готовый к инжекту.

## 🛠 Архитектура

Проект построен с использованием интерфейсов и композиции (Go-way OOP). Логика разделена на два основных конвейера: **Compiler Pipeline** (генерация и обфускация кода) и **Builder & Extractor** (сборка и извлечение байт).

### Диаграмма классов (Mermaid)

```mermaid
classDiagram
    %% Core Facade
    class PayloadBuilder {
        -compiler CompilerPipeline
        -toolchain ToolchainWrapper
        -extractor ElfExtractor
        +BuildPayload(sourceCode string, outputFile string) error
    }

    %% Subsystem 1: Compiler Pipeline
    class CompilerPipeline {
        -lexer Lexer
        -parser Parser
        -irGen IRGenerator
        -obfuscator ObfuscatorEngine
        -regAlloc RegisterAllocator
        -emitter AsmEmitter
        +CompileToAsm(sourceCode string) string
    }

    class Parser {
        +Parse(tokens) ASTNode
    }

    class IRGenerator {
        +GenerateIR(tree ASTNode) IR
    }

    class RegisterAllocator {
        +AllocatePhysicalRegisters(ir IR) IR
    }

    class AsmEmitter {
        +Emit(ir IR) string
    }

    %% Obfuscation Engine
    class ObfuscatorEngine {
        -passes []ObfuscationPass
        +RunObfuscation(ir IR) IR
        +AddPass(pass ObfuscationPass)
    }

    class ObfuscationPass {
        <<interface>>
        +Apply(ir IR) IR
    }

    class StringEncryptor {
        +Apply(ir IR) IR
    }
    class ControlFlowFlattener {
        +Apply(ir IR) IR
    }
    class AntiSandboxInjector {
        +Apply(ir IR) IR
    }
    class DeadCodeInjector {
        +Apply(ir IR) IR
    }

    %% Subsystem 2: Builder & Extractor
    class ToolchainWrapper {
        +CompileELF(asmCode string, tempElfPath string) error
    }

    class ElfExtractor {
        +ExtractTextSection(elfPath string) ([]byte, error)
    }

    %% Relationships
    PayloadBuilder --> CompilerPipeline : uses
    PayloadBuilder --> ToolchainWrapper : uses
    PayloadBuilder --> ElfExtractor : uses

    CompilerPipeline --> Parser : uses
    CompilerPipeline --> IRGenerator : uses
    CompilerPipeline --> ObfuscatorEngine : uses
    CompilerPipeline --> RegisterAllocator : uses
    CompilerPipeline --> AsmEmitter : uses

    ObfuscatorEngine *-- ObfuscationPass : contains
    ObfuscationPass <|.. StringEncryptor : implements
    ObfuscationPass <|.. ControlFlowFlattener : implements
    ObfuscationPass <|.. AntiSandboxInjector : implements
    ObfuscationPass <|.. DeadCodeInjector : implements
```

## 🧩 Описание компонентов

### 1. Compiler Pipeline (Компилятор)
Отвечает за трансляцию кастомного языка в полиморфный ассемблер.

* **`Parser`** *(ранее TreeGenerator)*: Читает исходный код языка и строит абстрактное синтаксическое дерево (AST).
* **`IRGenerator`** *(ранее AsmParser)*: Конвертирует AST в Intermediate Representation (IR) — промежуточный код, оперирующий абстрактными "виртуальными регистрами". Это позволяет применять обфускацию безопасно.
* **`ObfuscatorEngine`**: Ядро полиморфизма. Прогоняет IR через цепочку плагинов (интерфейс `ObfuscationPass`):
  * `StringEncryptor` — шифрует строковые литералы (например, XOR с динамическим ключом).
  * `ControlFlowFlattener` — ломает граф потока выполнения (защита от статического анализа).
  * `AntiSandboxInjector` — внедряет проверки среды (uptime, RAM, температура, cpuid).
  * `DeadCodeInjector` — зашумляет код мусорными инструкциями для изменения сигнатуры.
* **`RegisterAllocator`** *(ранее AsmFiller)*: Транслирует виртуальные регистры из IR в физические регистры архитектуры (RAX, RDI и т.д.), разрешая конфликты.
* **`AsmEmitter`**: Превращает финальный IR в чистый текст Assembly (`.s` файл).

### 2. Builder & Extractor (Сборщик)
Отвечает за взаимодействие с ОС и бинарными форматами.

* **`ToolchainWrapper`**: Вызывает системные утилиты (GCC/NASM/AS) с флагами для создания позиционно-независимого исполняемого файла (`-fPIE -pie -nostdlib`).
* **`ElfExtractor`**: Парсит сгенерированный ELF-файл, находит исполняемую секцию `.text` и извлекает из неё raw-байты.

### 3. Фасад
* **`PayloadBuilder`**: Главный оркестратор. Связывает воедино исходный код, компилятор, тулчейн и экстрактор.

## 🚀 Жизненный цикл сборки (Workflow)

1. Хакер пишет скрипт на кастомном языке (`source.pld`).
2. Фреймворк парсит код в AST, затем переводит в IR.
3. IR проходит многократные мутации (запутывание, шифрование, анти-песочница).
4. Виртуальные регистры заменяются на настоящие, генерируется файл `temp.s`.
5. Тулчейн компилирует `temp.s` -> `temp.elf`.
6. Экстрактор вырезает секцию `.text` из `temp.elf`.
7. На диск сохраняется итоговый **`payload.bin`**, готовый к загрузке в память целевого процесса.

## 💻 Использование

```bash
# Пример гипотетического вызова CLI
./polybuilder build -i source.pld -o payload.bin --obfuscation-level=max
```
