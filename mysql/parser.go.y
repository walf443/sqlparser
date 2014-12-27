// vim: noet sw=8 sts=8
%{
package mysql

import (
    "fmt"
    "os"
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
    table_names []TableNameIdentifier
    table_name TableNameIdentifier
    database_name DatabaseNameIdentifier
    column_name ColumnNameIdentifier
    index_name IndexNameIdentifier
    alter_specifications []AlterSpecification
    alter_specification AlterSpecification
    tok       Token
}

%type<statements> statements
%type<statement> statement
%type<table_names> table_names
%type<table_name> table_name
%type<database_name> database_name
%type<column_name> column_name
%type<index_name> index_name
%type<alter_specifications> alter_specifications
%type<alter_specification> alter_specification

%token<tok> DROP CREATE ALTER
%token<tok> IDENT TABLE COLUMN DATABASE INDEX KEY

%%

statements
    :
    {
        $$ = nil
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }
    | statements statement
    {
        $$ = append([]Statement{$2}, $1...)
        if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
            l.statements = $$
        }
    }

statement
    : DROP TABLE table_names ';'
    {
        $$ = &DropTableStatement{TableNames: $3}
    }
    | DROP DATABASE database_name ';'
    {
        $$ = &DropDatabaseStatement{DatabaseName: $3}
    }
    | CREATE DATABASE database_name ';'
    {
        $$ = &CreateDatabaseStatement{DatabaseName: $3}
    }
    | ALTER TABLE table_name alter_specifications ';'
    {
        $$ = &AlterTableStatement{TableName: $3, AlterSpecifications: $4}
    }

table_names
    : table_name
    {
        $$ = []TableNameIdentifier{$1}
    }
    | table_names ',' table_name
    {
        $$ = append([]TableNameIdentifier{$3}, $1...)
    }

table_name
    : IDENT
    {
        $$ = TableNameIdentifier{Name: $1.lit}
    }
    | '`' IDENT '`'
    {
        $$ = TableNameIdentifier{Name: $2.lit}
    }
    | IDENT '.' IDENT
    {
        $$ = TableNameIdentifier{Database: $1.lit, Name: $3.lit}
    }

database_name
    : IDENT
    {
        $$ = DatabaseNameIdentifier{Name: $1.lit}
    }
    | '`' IDENT '`'
    {
        $$ = DatabaseNameIdentifier{Name: $2.lit}
    }

alter_specifications
    :
    {
        $$ = nil
    }
    | alter_specifications alter_specification
    {
        $$ = append([]AlterSpecification{$2}, $1...)
    }

alter_specification
    : DROP skipable_column column_name
    {
        $$ = &AlterSpecificationDropColumn{ColumnName: $3}
    }
    | DROP index_or_key index_name
    {
        $$ = &AlterSpecificationDropIndex{IndexName: $3}
    }

skipable_column
    :
    | COLUMN

index_or_key
    : INDEX
    | KEY

column_name
    : IDENT
    {
        $$ = ColumnNameIdentifier{Name: $1.lit}
    }
    | '`' IDENT '`'
    {
        $$ = ColumnNameIdentifier{Name: $2.lit}
    }

index_name
    : IDENT
    {
        $$ = IndexNameIdentifier{Name: $1.lit}
    }
    | '`' IDENT '`'
    {
        $$ = IndexNameIdentifier{Name: $2.lit}
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
    fmt.Printf("%s while processing near %q line %d, col: %d\n", e, l.recentLit, l.recentPos.Line, l.recentPos.Column)
    fmt.Printf("%s\n", l.scanner.CurrentLine())
    for i := 0; i < l.recentPos.Column-1; i++ {
        fmt.Printf(" ")
    }
    fmt.Printf("^\n")
    os.Exit(1)
}

func Parse(s *Scanner) []Statement {
    l := LexerWrapper{scanner: s}
    if yyParse(&l) != 0 {
        panic("Parse error")
    }
    return l.statements
}
