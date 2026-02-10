package parser

import (
	"fmt"

	"github.com/impossibleclone/imposter/internal/ast"
	"github.com/impossibleclone/imposter/internal/lexer"
	"github.com/impossibleclone/imposter/internal/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l,
		errors: []string{}}

	//Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
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
		return nil
	}
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
