package reader

import (
	"bufio"
	"io"
)

type Position struct {
	Offset    int
	Line, Col int
}

type Scanner struct {
	reader         io.Reader
	runeReader     *bufio.Reader
	pos, line, col int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader:     r,
		runeReader: bufio.NewReader(r),
		line:       1,
	}
}

func (s *Scanner) Peek() (rune, error) {
	ch, _, err := s.runeReader.ReadRune()
	if err != nil {
		return 0, err
	}
	err = s.runeReader.UnreadRune()
	if err != nil {
		return 0, err
	}
	return ch, nil
}

func (s *Scanner) Next() (rune, error) {
	ch, _, err := s.runeReader.ReadRune()
	if err != nil {
		return 0, err
	}
	s.pos++
	if ch == '\n' {
		s.line++
		s.col = 0
	} else {
		s.col++
	}
	return ch, nil
}

func (s *Scanner) Position() Position {
	return Position{
		Offset: s.pos,
		Line:   s.line,
		Col:    s.col,
	}
}

func (s *Scanner) SkipIrrelevant() {
	// skip inter-token-space and comments
}
