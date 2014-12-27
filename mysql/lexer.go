package mysql

const (
	EOF     = -1
	UNKNOWN = 0
)

var keywords = map[string]int{
	"DROP": DROP,
	"drop": DROP,
	"TABLE": TABLE,
	"table": TABLE,
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
		if keyword, ok := keywords[lit]; ok {
			tok = keyword
		} else {
			tok = IDENT
		}
	default:
		switch ch {
		case -1:
			tok = EOF
		case ';', ',', '`', '.':
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
	for isLetter(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}

	return string(ret)
}
