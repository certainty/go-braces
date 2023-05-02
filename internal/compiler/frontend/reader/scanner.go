package reader

import (
	"io"
)

type Position struct {
	Offset    uint64
	Line, Col uint64
}

type Scanner struct {
	buffer         *[]rune
	bufferLen      uint64
	pos, line, col uint64
	posStack       []uint64
}

func NewScanner(buffer *[]rune) *Scanner {
	return &Scanner{
		buffer:    buffer,
		bufferLen: uint64(len(*buffer)),
		pos:       0,
		line:      1,
		col:       1,
		posStack:  []uint64{},
	}
}

func (s *Scanner) IsEof() bool {
	return s.pos >= s.bufferLen
}

func (s *Scanner) Position() Position {
	return Position{
		Offset: s.pos,
		Line:   s.line,
		Col:    s.col,
	}
}

func (s *Scanner) SavePosition() error {
	s.posStack = append(s.posStack, s.pos)
	return nil
}

func (s *Scanner) RestorePosition() error {
	if len(s.posStack) == 0 {
		return io.ErrUnexpectedEOF
	}

	s.pos = s.posStack[len(s.posStack)-1]
	s.posStack = s.posStack[:len(s.posStack)-1]
	return nil
}

func (s *Scanner) Peek() (rune, error) {
	return s.PeekN(1)
}

func (s *Scanner) PeekN(n uint64) (rune, error) {
	if s.pos+n >= s.bufferLen {
		return 0, io.EOF
	}
	return (*s.buffer)[s.pos+n], nil
}

func (s *Scanner) Next() (rune, error) {
	s.pos++
	if s.IsEof() {
		return 0, io.EOF
	}

	ch := (*s.buffer)[s.pos]
	if ch == '\n' {
		s.line++
		s.col = 0
	} else {
		s.col++
	}

	return ch, nil
}

func (s *Scanner) Skip() error {
	_, err := s.Next()
	return err
}

func (s *Scanner) skipWhitespace() (bool, error) {
	ch, err := s.Peek()
	if err != nil {
		return false, err
	}

	if ch == ' ' || ch == '\t' {
		if err := s.Skip(); err != nil {
			return false, err
		}
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
		if err := s.Skip(); err != nil {
			return false, err
		}

		if ch == '\r' {
			nextCh, _ := s.Peek()
			if nextCh == '\n' {
				if err := s.Skip(); err != nil {
					return false, err
				}
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
		if err := s.SavePosition(); err != nil {
			return false, err
		}

		if err := s.Skip(); err != nil {
			return false, err
		}
		nextCh, _ := s.Peek()
		if nextCh == '|' {
			if err := s.Skip(); err != nil {
				return false, err
			}
			var commentNesting = 1

			for commentNesting > 0 {
				ch, err = s.Next()
				if err != nil {
					return false, err
				}

				if ch == '#' {
					nextCh, _ := s.Peek()
					if nextCh == '|' {
						if err := s.Skip(); err != nil {
							return false, err
						}
						commentNesting++
					}
				} else if ch == '|' {
					nextCh, _ := s.Peek()
					if nextCh == '#' {
						if err := s.Skip(); err != nil {
							return false, err
						}
						commentNesting--
					}
				}
			}
			return true, nil
		} else {
			if err := s.RestorePosition(); err != nil {
				return false, err
			}
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
		if err := s.Skip(); err != nil {
			return false, err
		}

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

func (s *Scanner) Attempt(expected string) bool {
	expectedLen := uint64(len(expected))

	if s.pos+expectedLen > s.bufferLen {
		return false
	}
	actualRunes := (*s.buffer)[s.pos : s.pos+expectedLen]

	if expected != string(actualRunes) {
		return false
	}

	s.pos += expectedLen
	return true
}
