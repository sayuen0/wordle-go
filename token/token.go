package token

import (
	"fmt"

	"github.com/sayuen0/wordle-go/util"
)

type Token struct {
	Literal string
	util.ColorScheme
}

func (t Token) String() string {
	return fmt.Sprintf(string(t.ColorScheme), t.Literal)
}

func (t Token) GoString() string {
	return fmt.Sprintf(string(t.ColorScheme), t.Literal)
}
