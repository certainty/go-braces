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

func (s *Scanner) skipWhitespace() (bool, error) {
	ch, err := s.Peek()
	if err != nil {
		return false, err
	}

	if ch == ' ' || ch == '\t' {
		s.Next()
		return true, nil
	}
	return false, nil
}

func (s *Scanner) skipEOL() (bool, error) {
	ch, err := s.Peek()
	if err != nil {
		return false, err
	}

	if ch == '\n' || ch == '\r' {
		s.Next()
		if ch == '\r' {
			nextCh, _ := s.Peek()
			if nextCh == '\n' {
				s.Next()
			}
		}
		return true, nil
	}

	return false, nil
}

func (s *Scanner) SkipIrrelevant() error {
	for {
		accepted, err := s.skipWhitespace()
		if err != nil {
			return err
		}
		if accepted {
			continue
		}

		accepted, err = s.skipEOL()
		if err != nil {
			return err
		}
		if accepted {
			continue
		}

		accepted, err = s.skipSkipLineComment()
		if err != nil {
			return err
		}

		if accepted {
			continue
		}

		accepted, err = s.skipMultiLineComment()
		if err != nil {
			return err
		}
		if accepted {
			continue
		}

		break
	}
	return nil
}

func (s *Scanner) skipMultiLineComment() (bool, error) {
	ch, err := s.Peek()
	if err != nil {
		return false, err
	}

	if ch == '#' {
		s.Next()
		nextCh, _ := s.Peek()
		if nextCh == '|' {
			s.Next()
			var commentNesting = 1

			for commentNesting > 0 {
				ch, err = s.Next()
				if err != nil {
					return false, err
				}

				if ch == '#' {
					nextCh, _ := s.Peek()
					if nextCh == '|' {
						s.Next()
						commentNesting++
					}
				} else if ch == '|' {
					nextCh, _ := s.Peek()
					if nextCh == '#' {
						s.Next()
						commentNesting--
					}
				}
			}
			return true, nil
		}
	}
	return false, nil
}

// consumes the input stream skipping a line comment
func (s *Scanner) skipSkipLineComment() (bool, error) {
	ch, err := s.Peek()
	if err != nil {
		return false, err
	}

	if ch == ';' {
		s.Next()
		for {
			ch, err = s.Next()
			if err != nil {
				return false, err
			}
			if ch == '\n' || ch == '\r' {
				return true, nil
			}
		}
	}
	return false, nil
}
