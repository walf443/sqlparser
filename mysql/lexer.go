package mysql

import (
	"strings"
)

const (
	EOF     = -1
	UNKNOWN = 0
)

type DataType uint

const (
	DATATYPE_BIT DataType = iota
	DATATYPE_TINYINT
	DATATYPE_SMALLINT
	DATATYPE_MEDIUMINT
	DATATYPE_INT
	DATATYPE_BIGINT
	DATATYPE_REAL
	DATATYPE_DOUBLE
	DATATYPE_FLOAT
	DATATYPE_DECIMAL
	DATATYPE_NUMERIC
	DATATYPE_DATE
	DATATYPE_TIME
	DATATYPE_TIMESTAMP
	DATATYPE_DATETIME
	DATATYPE_YEAR
	DATATYPE_CHAR
	DATATYPE_VARCHAR
	DATATYPE_BINARY
	DATATYPE_VARBINARY
	DATATYPE_TINYBLOB
	DATATYPE_BLOB
	DATATYPE_MEDIUMBLOB
	DATATYPE_LONGBLOB
	DATATYPE_TINYTEXT
	DATATYPE_TEXT
	DATATYPE_MEDIUMTEXT
	DATATYPE_LONGTEXT
)

var keywords = map[string]int{
	"ADD": ADD,
	"DROP": DROP,
	"CREATE": CREATE,
	"ALTER": ALTER,
	"COLUMN": COLUMN,
	"TABLE": TABLE,
	"INDEX": INDEX,
	"KEY": KEY,
	"DATABASE": DATABASE,
	"NULL": NULL,
	"NOT": NOT,
	"AUTO_INCREMENT": AUTO_INCREMENT,

	// datatypes
	"BIT": BIT,
	"TINYINT": TINYINT,
	"SMALLINT": SMALLINT,
	"MEDIUMINT": MEDIUMINT,
	"INT": INT,
	"INTEGER": INTEGER,
	"BIGINT": BIGINT,
	"REAL": REAL,
	"DOUBLE": DOUBLE,
	"FLOAT": FLOAT,
	"DECIMAL": DECIMAL,
	"NUMERIC": NUMERIC,
	"DATE": DATE,
	"TIME": TIME,
	"TIMESTAMP": TIMESTAMP,
	"DATETIME": DATETIME,
	"YEAR": YEAR,
	"CHAR": CHAR,
	"VARCHAR": VARCHAR,
	"BINARY": BINARY,
	"VARBINARY": VARBINARY,
	"TINYBLOB": TINYBLOB,
	"BLOB": BLOB,
	"MEDIUMBLOB": MEDIUMBLOB,
	"LONGBLOB": LONGBLOB,
	"TINYTEXT": TINYTEXT,
	"TEXT": TEXT,
	"MEDIUMTEXT": MEDIUMTEXT,
	"LONGTEXT": LONGTEXT,

	// datatype options
	"UNSIGNED": UNSIGNED,
	"ZEROFILL": ZEROFILL,
}

type Position struct {
	Line   int
	Column int
}

type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *Scanner) Init(src string) {
	s.src = []rune(src)
}

func (s *Scanner) Scan() (tok int, lit string, pos Position) {
	s.skipWhiteSpace()
	pos = s.position()
	switch ch := s.peek(); {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if keyword, ok := keywords[strings.ToUpper(lit)]; ok {
			tok = keyword
		} else {
			tok = IDENT
		}
	case isNumber(ch):
		lit = s.scanNumber()
		tok = NUMBER
	default:
		switch ch {
		case -1:
			tok = EOF
		case ';', ',', '`', '.', '(', ')':
			tok = int(ch)
			lit = string(ch)
		}
		s.next()
	}
	return
}

func (s *Scanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *Scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *Scanner) CurrentLine() string {
	cursor := s.lineHead
	var bytes []rune
	for {
		ch :=  s.src[cursor]

		if ch == '\n' {
			break
		}
		bytes = append(bytes, ch)
		cursor++
		if len(s.src) <= cursor {
			break
		}
	}
	return string(bytes)
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_';
}

func isNumber(ch rune) bool {
	return '0' <= ch && ch <= '9';
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *Scanner) position() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *Scanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *Scanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isNumber(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}

	return string(ret)
}

func (s *Scanner) scanNumber() string {
	var ret []rune
	for isNumber(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}
