// vim: noet sw=8 sts=8
%{
package mysql

import (
    "log"
)

type Token struct {
    tok int
    lit string
    pos Position
}

%}

%union{
    statements []Statement
    statement Statement
    tok       Token
}

%type<statements> statements
%type<statement> statement

%token<tok> IDENT DROP TABLE

%%

statements
    :
    {
        $$ = nil
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }
    | statement statements
    {
        $$ = append([]Statement{$1}, $2...)
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }

statement
    : DROP TABLE table_name ';'
    {
    }


table_name
    : IDENT
    {
    }

%%

type LexerWrapper struct {
    scanner *Scanner
    recentLit   string
    recentPos   Position
    statements []Statement
}

func (l *LexerWrapper) Lex(lval *yySymType) int {
    tok, lit, pos := l.scanner.Scan()
    if tok == EOF {
        return 0
    }
    lval.tok = Token{tok: tok, lit: lit, pos: pos}
    l.recentLit = lit
    l.recentPos = pos
    return tok
}

func (l *LexerWrapper) Error(e string) {
    log.Fatalf("Line %d, Column %d: %q %s", l.recentPos.Line, l.recentPos.Column, l.recentLit, e)
}

func Parse(s *Scanner) []Statement {
    l := LexerWrapper{scanner: s}
    if yyParse(&l) != 0 {
        panic("Parse error")
    }
    return l.statements
}
