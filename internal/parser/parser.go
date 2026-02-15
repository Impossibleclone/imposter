package parser

import (
	"fmt"
	"strconv"

	"github.com/impossibleclone/imposter/internal/ast"
	"github.com/impossibleclone/imposter/internal/lexer"
	"github.com/impossibleclone/imposter/internal/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFn   map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // !x or -x
	CALL        // ()
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []string{}}

	//Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the next token is of the given type, and if so, it
// advances the parser's position to the next token.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// parseStatement parses a statement
// A statement is either a var statement or a function statement
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
		// case token.FUNCTION:
		// 	return p.parseFunctionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseVarStatement parses a var statement
func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: We're skipping the expression until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	//TODO: We're skipping the expression until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// func (p *Parser) parseFunctionStatement() *ast.FnStatement {
// 	stmt := &ast.FnStatement{Token: p.curToken}
//
// 	if !p.expectPeek(token.IDENT) {
// 		return nil
// 	}
// 	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
//
// 	if !p.expectPeek(token.LPAREN) {
// 		return nil
// 	}
//
// 	stmt.Params = p.parseFunctionParams()
//
// 	if !p.expectPeek(token.LBRACE) {
// 		return nil
// 	}
//
// 	// stmt.Body = p.parseBlockStatement()
//
// 	for !p.curTokenIs(token.SEMICOLON) {
// 		p.nextToken()
// 	}
// 	return stmt
// }
//
// func (p *Parser) parseFunctionParams() []*ast.Identifier {
// 	identifiers := []*ast.Identifier{}
//
// 	if p.peekTokenIs(token.RPAREN) {
// 		p.nextToken()
// 		return identifiers
// 	}
// 	p.nextToken()
//
// 	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
// 	identifiers = append(identifiers, ident)
//
// 	for p.peekTokenIs(token.COMMA) {
// 		p.nextToken()
// 		p.nextToken()
// 		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
// 		identifiers = append(identifiers, ident)
// 	}
//
// 	if !p.expectPeek(token.RPAREN) {
// 		return nil
// 	}
//
// 	return identifiers
// }

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) registerPrefix(tokentype token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokentype] = fn
}

func (p *Parser) registerInfix(tokentype token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokentype] = fn
}
