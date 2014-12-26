package mysql

import (
	"testing"
)

func testScanner(t *testing.T, src string, expectTok int) {
	s := new(Scanner)
	s.Init(src)
	tok, lit, _ := s.Scan()
	if tok != expectTok {
		t.Errorf("Expect Scanner{%q}.Scan() expected %#v, but got %#v", src, expectTok, tok)
	}
	if lit != src {
		t.Errorf("Expect Scanner{%q}.Scan(): lit: %#v src: %q", src, lit, src)
	}
	tok, lit, _ = s.Scan()
	if tok != EOF {
		t.Errorf("Expect Scanner{%q}.Scan() expected EOF but got %#v", src, tok)
	}
}

func TestScanner(t *testing.T) {
	testScanner(t, "(", '(')
}
