package parser

import (
	"fmt"
	"strconv"

	"go.mod/pkg/ast"
	"go.mod/pkg/lexer"
)

// ParserImpl реализует логику парсера.
type ParserImpl struct {
	l              lexer.Lexer
	errors         []string
	curToken       lexer.Token
	peekToken      lexer.Token
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

// New создает и инициализирует новый экземпляр парсера.
func New(l lexer.Lexer) Parser {
	p := &ParserImpl{
		l:      l,
		errors: []string{},
	}
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifier)
	p.registerPrefix(lexer.NUMBER, p.parseIntegerLiteral)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.ASSIGN, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken продвигает парсер по потоку токенов.
func (p *ParserImpl) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors возвращает срез строк с описанием ошибок.
func (p *ParserImpl) Errors() []string { return p.errors }

// ParseProgram является входной точкой для парсинга программы.
func (p *ParserImpl) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

// parseStatement определяет тип текущего утверждения и вызывает соответствующий метод.
func (p *ParserImpl) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.QWORD, lexer.STR_KW:
		return p.parseLetStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.IMPORT:
		return p.parseImportStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement парсит объявление переменной (let-подобное).
func (p *ParserImpl) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(lexer.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	return stmt
}

// parseIfStatement парсит конструкцию условия if-else.
func (p *ParserImpl) parseIfStatement() ast.Statement {
	stmt := &ast.IfStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(lexer.ELSE) {
		p.nextToken()
		if !p.expectPeek(lexer.LBRACE) {
			return nil
		}
		stmt.Alternative = p.parseBlockStatement()
	}
	return stmt
}

// parseWhileStatement парсит цикл while.
func (p *ParserImpl) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStatement{Token: p.curToken}
	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

// parseBlockStatement парсит блок кода, заключенный в фигурные скобки.
func (p *ParserImpl) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(lexer.RBRACE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

// parseImportStatement парсит инструкцию импорта внешнего ресурса.
func (p *ParserImpl) parseImportStatement() ast.Statement {
	stmt := &ast.ImportStatement{Token: p.curToken}
	if !p.expectPeek(lexer.STRING) {
		return nil
	}
	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return stmt
}

// parseExpressionStatement парсит выражение, стоящее как отдельная инструкция.
func (p *ParserImpl) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

// parseExpression парсит выражение с учетом приоритетов операторов.
func (p *ParserImpl) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("Нет функции парсинга для токена %s", p.curToken.Type))
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(lexer.EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// parseIdentifier возвращает узел идентификатора.
func (p *ParserImpl) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral возвращает узел целочисленного литерала с поддержкой разных систем счисления.
func (p *ParserImpl) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Не могу распарсить %q как целое число", p.curToken.Literal))
		return nil
	}
	lit.Value = value
	return lit
}

// parseStringLiteral возвращает узел строкового литерала.
func (p *ParserImpl) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseInfixExpression парсит инфиксное выражение (оператор между операндами).
func (p *ParserImpl) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

// parseCallExpression парсит вызов функции с аргументами.
func (p *ParserImpl) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(lexer.RPAREN)
	return exp
}

// parseExpressionList парсит список выражений, разделенных запятыми.
func (p *ParserImpl) parseExpressionList(end lexer.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(end) {
		return nil
	}
	return list
}

// peekPrecedence возвращает приоритет следующего токена.
func (p *ParserImpl) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence возвращает приоритет текущего токена.
func (p *ParserImpl) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *ParserImpl) curTokenIs(t lexer.TokenType) bool  { return p.curToken.Type == t }
func (p *ParserImpl) peekTokenIs(t lexer.TokenType) bool { return p.peekToken.Type == t }

// expectPeek проверяет тип следующего токена и продвигает парсер, если он совпадает.
func (p *ParserImpl) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("Ожидался токен %s, но получен %s (Строка: %d)", t, p.peekToken.Type, p.peekToken.Line))
	return false
}

// registerPrefix регистрирует функцию парсинга для префиксных позиций.
func (p *ParserImpl) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix регистрирует функцию парсинга для инфиксных позиций.
func (p *ParserImpl) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
