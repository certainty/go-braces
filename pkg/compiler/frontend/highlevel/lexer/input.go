package lexer

import (
	"bufio"
	"fmt"
	"github.com/certainty/go-braces/pkg/compiler/frontend/highlevel/token"
	"io"
	"os"
	"unicode/utf8"
)

type Input struct {
	Origin token.Origin
	Buffer *[]rune
}

func (i Input) Source() string {
	return string(*i.Buffer)
}

func NewStringInput(label string, s string) *Input {
	buffer := []rune(s)
	return &Input{
		Origin: token.StringOrigin{Label: label},
		Buffer: &buffer,
	}
}

func NewReplInput(inputCount uint64, s string) *Input {
	buffer := []rune(s)
	return &Input{
		Origin: token.ReplOrigin{InputCount: inputCount},
		Buffer: &buffer,
	}
}

func NewFileInput(path string) (*Input, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buffer, err := readRunesFromReader(file)
	if err != nil {
		return nil, err
	}

	return &Input{
		Origin: token.FileOrigin{Path: path},
		Buffer: &buffer,
	}, nil
}

func readRunesFromReader(r io.Reader) ([]rune, error) {
	reader := bufio.NewReader(r)
	var runes []rune
	buf := make([]byte, utf8.UTFMax)

	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if r == utf8.RuneError && size == 1 {
			// Handle invalid UTF-8 bytes
			_, _ = reader.Read(buf[1:size])
			continue
		}

		runes = append(runes, r)
	}

	return runes, nil
}

func (i Input) Description() string {
	return fmt.Sprintf("%s  excerpt: %s", i.Origin.Description(), excerpt(i.Buffer, 20))
}

func excerpt(s *[]rune, max int) string {
	if len(*s) < max {
		return string(*s)
	} else {
		return string((*s)[0:max]) + "..."
	}
}
