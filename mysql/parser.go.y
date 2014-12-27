// vim: noet sw=8 sts=8
%{
package mysql

import (
    "fmt"
    "os"
    "strconv"
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
    column_definition ColumnDefinition
    alter_specifications []AlterSpecification
    alter_specification AlterSpecification
    data_type DataTypeDefinition
    bool bool
    data_type_type DataType
    uint uint
    fraction_option [2]uint
    tok       Token
}

%type<statements> statements
%type<statement> statement
%type<table_names> table_names
%type<table_name> table_name
%type<database_name> database_name
%type<column_name> column_name
%type<index_name> index_name
%type<column_definition> column_definition
%type<alter_specifications> alter_specifications
%type<alter_specification> alter_specification
%type<data_type> data_type
%type<data_type_type> data_type_number data_type_fraction data_type_decimal
%type<bool> unsigned_option zerofill_option
%type<uint> length_option
%type<fraction_option> fraction_option decimal_option

%token<tok> IDENT NUMBER
%token<tok> DROP CREATE ALTER ADD
%token<tok> TABLE COLUMN DATABASE INDEX KEY
%token<tok> BIT TINYINT SMALLINT MEDIUMINT INT INTEGER BIGINT REAL DOUBLE FLOAT DECIMAL NUMERIC DATE TIME TIMESTAMP DATETIME YEAR CHAR VARCHAR BINARY VARBINARY TINYBLOB BLOB MEDIUMBLOB LONGBLOB TINYTEXT TEXT MEDIUMTEXT LONGTEXT UNSIGNED ZEROFILL

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
    : ADD skipable_column column_name column_definition
    {
        $$ = &AlterSpecificationAddColumn{ColumnName: $3, ColumnDefinition: $4}
    }
    | DROP index_or_key index_name
    {
        $$ = &AlterSpecificationDropIndex{IndexName: $3}
    }
    | DROP skipable_column column_name
    {
        $$ = &AlterSpecificationDropColumn{ColumnName: $3}
    }

skipable_column
    :
    | COLUMN

column_definition
    : data_type nullable default autoincrement key_options column_comment
    {
        $$ = ColumnDefinition{$1}
    }

nullable
    :

default
    :

autoincrement
    :

key_options
    :

column_comment
    :

data_type
    : BIT
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BIT }
    }
    | data_type_number length_option unsigned_option zerofill_option
    {
        $$ = &DataTypeDefinitionNumber{Type: $1, Length: $2, Unsigned: $3, Zerofill: $4 }
    }
    | data_type_fraction fraction_option unsigned_option zerofill_option
    {
        fraction := $2
        $$ = &DataTypeDefinitionFraction{Type: $1, Length: fraction[0], Decimals: fraction[1], Unsigned: $3, Zerofill: $4 }
    }
    | data_type_decimal decimal_option unsigned_option zerofill_option
    {
        fraction := $2
        $$ = &DataTypeDefinitionFraction{Type: $1, Length: fraction[0], Decimals: fraction[1], Unsigned: $3, Zerofill: $4 }
    }
    | DATE
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_DATE }
    }
    | TIME
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TIME }
    }
    | TIMESTAMP
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TIMESTAMP }
    }
    | DATETIME
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_DATETIME }
    }
    | YEAR
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_YEAR }
    }
    | CHAR length_option
    {
        $$ = &DataTypeDefinitionString{Type: DATATYPE_CHAR, Length: $2 }
    }
    | VARCHAR length_option
    {
        $$ = &DataTypeDefinitionString{Type: DATATYPE_VARCHAR, Length: $2 }
    }
    | BINARY
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BINARY }
    }
    | VARBINARY
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_VARBINARY }
    }
    | TINYBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_TINYBLOB }
    }
    | BLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_BLOB }
    }
    | MEDIUMBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_MEDIUMBLOB }
    }
    | LONGBLOB
    {
        $$ = &DataTypeDefinitionSimple{Type: DATATYPE_LONGBLOB }
    }
    | TINYTEXT
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_TINYTEXT }
    }
    | TEXT
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_TEXT }
    }
    | MEDIUMTEXT
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_MEDIUMTEXT }
    }
    | LONGTEXT
    {
        $$ = &DataTypeDefinitionTextBlob{Type: DATATYPE_LONGTEXT }
    }

data_type_number
    : TINYINT
    {
        $$ = DATATYPE_TINYINT
    }
    | SMALLINT
    {
        $$ = DATATYPE_SMALLINT
    }
    | MEDIUMINT
    {
        $$ = DATATYPE_MEDIUMINT
    }
    | INT
    {
        $$ = DATATYPE_INT
    }
    | INTEGER
    {
        $$ = DATATYPE_INT
    }
    | BIGINT
    {
        $$ = DATATYPE_BIGINT
    }

data_type_fraction
    : REAL
    {
        $$ = DATATYPE_REAL
    }
    | DOUBLE
    {
        $$ = DATATYPE_DOUBLE
    }
    | FLOAT
    {
        $$ = DATATYPE_FLOAT
    }

data_type_decimal
    : DECIMAL
    {
        $$ = DATATYPE_DECIMAL
    }
    | NUMERIC
    {
        $$ = DATATYPE_NUMERIC
    }

length_option
    :
    {
        $$ = 0
    }
    | '(' NUMBER ')'
    {
        num, err := strconv.Atoi($2.lit)
        if err != nil {
            num = 0
        }
        $$ = uint(num)
    }

fraction_option
    :
    {
        $$ = [2]uint{0, 0}
    }
    | '(' NUMBER ',' NUMBER ')'
    {
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        num2, err := strconv.Atoi($4.lit)
        if err != nil {
            num2 = 0
        }
        result := [2]uint{0, 0}
        result[0] = uint(num1)
        result[1] = uint(num2)
        $$ = result
    }

decimal_option
    :
    {
        $$ = [2]uint{0, 0}
    }
    | '(' NUMBER ')'
    {
        result := [2]uint{0, 0}
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        result[0] = uint(num1)
        $$ = result
    }
    | '(' NUMBER ',' NUMBER ')'
    {
        num1, err := strconv.Atoi($2.lit)
        if err != nil {
            num1 = 0
        }
        num2, err := strconv.Atoi($4.lit)
        if err != nil {
            num2 = 0
        }
        result := [2]uint{0, 0}
        result[0] = uint(num1)
        result[1] = uint(num2)
        $$ = result
    }

unsigned_option
    :
    {
        $$ = false
    }
    | UNSIGNED
    {
        $$ = true
    }

zerofill_option
    :
    {
        $$ = false
    }
    | ZEROFILL
    {
        $$ = true
    }


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
