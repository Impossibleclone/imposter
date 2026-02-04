package ast

import "github.com/impossibleclone/imposter/internal/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type VarStatement struct {
	Token token.Token //The token for the "var" keyword
	Name  *Identifier //The identifier for the variable
	Value Expression  //The value of the variable
}

func (ls *VarStatement) statementNode()       {}
func (ls *VarStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type FnStatement struct {
	Token  token.Token
	Name   *Identifier
	Params []*Identifier
	Body   any
}

func (fs *FnStatement) statementNode()       {}
func (fs *FnStatement) TokenLiteral() string { return fs.Token.Literal }
