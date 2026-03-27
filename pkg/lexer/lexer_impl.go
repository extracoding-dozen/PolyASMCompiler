package lexer

// LexerImpl — конкретная реализация интерфейса Lexer.
// Она хранит состояние сканирования входной строки.
type LexerImpl struct {
	input        string // Исходный код
	position     int    // Текущая позиция в input (указывает на текущий символ)
	readPosition int    // Позиция для следующего чтения (после текущего символа)
	ch           byte   // Текущий исследуемый символ
	line         int    // Текущая строка (для отладки и ошибок)
	column       int    // Текущая колонка (для отладки и ошибок)
}

// New создает и инициализирует новый экземпляр LexerImpl с входной строкой.
func New() Lexer {
	l := &LexerImpl{input: "", line: 1, column: 0}
	//l.readChar() // Чтение первого символа
	return l
}

// SetInput задает строку кода.
func (l *LexerImpl) SetInput(input string) {
	l.input = input
	l.readChar()
}

// readChar считывает следующий символ из input и переводит указатели позиции.
// Если достигнут конец ввода, устанавливает ch в 0 (NUL).
func (l *LexerImpl) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

// NextToken — основная функция лексера. Распознает текущий символ (или последовательность)
// и возвращает соответствующую структуру Token.
func (l *LexerImpl) NextToken() Token {
	var tok Token

	l.skipWhitespaceAndComments()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: EQ, Literal: "==", Line: l.line, Column: l.column}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: "!=", Line: l.line, Column: l.column}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '>':
		tok = Token{Type: GT, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '<':
		tok = Token{Type: LT, Literal: string(l.ch), Line: l.line, Column: l.column}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch), Line: l.line, Column: l.column}
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = l.line
		tok.Column = l.column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line
		tok.Column = l.column
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			tok.Line = l.line
			tok.Column = l.column
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column}
		}
	}

	l.readChar()
	return tok
}

// readIdentifier считывает последовательность букв, цифр или подчеркиваний, формируя идентификатор.
func (l *LexerImpl) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber считывает числовой литерал. Поддерживает десятичные числа и шестнадцатеричные (0x...).
func (l *LexerImpl) readNumber() string {
	position := l.position

	// Обработка Hex-формата (0x...)
	if l.ch == '0' && l.peekChar() == 'x' {
		l.readChar()
		l.readChar()

		for isDigit(l.ch) || (l.ch >= 'a' && l.ch <= 'f') || (l.ch >= 'A' && l.ch <= 'F') {
			l.readChar()
		}
		return l.input[position:l.position]
	}

	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString считывает содержимое строкового литерала, заключенного в двойные кавычки.
func (l *LexerImpl) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// skipWhitespaceAndComments пропускает пробельные символы и однострочные комментарии (//).
// Увеличивает счетчик строк при обнаружении символа новой строки.
func (l *LexerImpl) skipWhitespaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
			l.readChar()
		} else if l.ch == '\n' {
			l.line++
			l.column = 0
			l.readChar()
		} else if l.ch == '/' && l.peekChar() == '/' {
			// Пропуск комментария до конца строки
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else {
			break
		}
	}
}

// isLetter проверяет, является ли символ буквой латинского алфавита или знаком подчеркивания.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit проверяет, является ли символ цифрой.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// lookupIdent проверяет, является ли считанный идентификатор зарезервированным ключевым словом.
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// peekChar возвращает следующий символ во входном потоке, не перемещая текущую позицию.
func (l *LexerImpl) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
