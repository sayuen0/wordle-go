package wordle

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sayuen0/wordle-go/token"
	"github.com/sayuen0/wordle-go/util"
)

const (
	defaultGameCount  int = 6
	defaultCharLength int = 5
)
const (
	exactColor     = util.BgBrightRed
	someMatchColor = util.BgBrightGreen
	noMatchColor   = util.BgBrightBlue
)

type RandomIntGenerator interface {
	Intn(n int) int
}

type Wordle interface {
	Rule() string
	OpenQuestion(r io.Reader, rand RandomIntGenerator) error
	Pick(r RandomIntGenerator) int
	ValidInput(s string) bool
	SetInput(s string)
	ExactMatchToken(s string) token.Token
	NoMatchToken(s string) token.Token
	SomeMatchToken(s string) token.Token
	Table() string
	ExactMatchWord(s string) bool
	WinMessage() string
	TimeUp() bool
	ProcessMessage() string
	Answer() string
}

type wordle struct {
	answerIndex  int
	words        []string
	inputs       [][]token.Token
	validateFunc util.ValidateFunc
	gameCount    int
	curGame      int
}

type WordleOption struct {
	GameCount    int
	ValidateFunc util.ValidateFunc
}

var DefaultWordleOption = WordleOption{
	GameCount:    defaultGameCount,
	ValidateFunc: util.ValidLowerAlphabet(defaultCharLength),
}

func New(option *WordleOption) Wordle {
	if option == nil {
		option = &DefaultWordleOption
	}
	if option.ValidateFunc == nil {
		option.ValidateFunc = DefaultWordleOption.ValidateFunc
	}
	if option.GameCount <= 0 {
		option.GameCount = DefaultWordleOption.GameCount
	}

	gameCount := option.GameCount
	validateFunc := option.ValidateFunc

	tokens := make([][]token.Token, 0)
	return &wordle{
		inputs:       tokens,
		gameCount:    gameCount,
		curGame:      0,
		validateFunc: validateFunc,
	}
}

func (w *wordle) Rule() string {
	rule := fmt.Sprintf(`rule
	- %s letter means the letter is not used in answer.
	- %s letter means the letter is used in the answer but in a different position.
	- %s letter means the letter is used in the answer and in the same position.`,
		fmt.Sprintf(string(noMatchColor), "Blue"),
		fmt.Sprintf(string(someMatchColor), "Green"),
		fmt.Sprintf(string(exactColor), "Red"),
	)
	return rule
}

func (w *wordle) OpenQuestion(r io.Reader,
	rand RandomIntGenerator) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := scanner.Text()
		if !w.ValidInput(s) {
			return fmt.Errorf("invalid source, word: %s", s)
		}
		w.words = append(w.words, s)
	}
	if len(w.words) == 0 {
		return errors.New("no word found for question in source")
	}
	w.answerIndex = w.Pick(rand)
	return nil
}

func (w *wordle) ValidInput(s string) bool {
	return w.validateFunc(s)
}

func (w *wordle) Pick(r RandomIntGenerator) int {
	return r.Intn(len(w.words))
}

func (w *wordle) SetInput(s string) {
	// create 5 token with matched string
	tokens := w.createTokens(s)
	w.inputs = append(w.inputs, tokens)
	w.curGame++
}

func (w *wordle) createTokens(s string) []token.Token {
	tokens := make([]token.Token, 0)
	for i, r := range s {
		if strings.ContainsRune(w.Answer(), r) {
			if s[i] == w.Answer()[i] {
				// same place, same letter
				tokens = append(tokens,
					w.ExactMatchToken(string(r)))
			} else {
				tokens = append(tokens,
					w.SomeMatchToken(string(r)))
			}
		} else {
			tokens = append(tokens,
				w.NoMatchToken(string(r)))
		}
	}
	return tokens
}

func (w *wordle) Answer() string {
	return w.words[w.answerIndex]
}

// debug method
func (w *wordle) setQuestion(words []string, index int) {
	w.words = words
	w.answerIndex = index
}

func (w wordle) ExactMatchToken(s string) token.Token {
	return w.token(s, exactColor)
}

func (w wordle) SomeMatchToken(s string) token.Token {
	return w.token(s, someMatchColor)
}

func (w wordle) NoMatchToken(s string) token.Token {
	return w.token(s, noMatchColor)
}

func (w wordle) token(s string, c util.ColorScheme) token.Token {
	return token.Token{Literal: s, ColorScheme: c}
}

func (w *wordle) Table() string {
	s := ""
	for i := 0; i < w.curGame; i++ {
		tokenArr := w.inputs[i]
		for _, t := range tokenArr {
			s += fmt.Sprintf("%v", t)
		}
		s += "\n"
	}
	return s
}

func (w *wordle) ExactMatchWord(s string) bool {
	return w.Answer() == s
}

func (w *wordle) WinMessage() string {
	return fmt.Sprintf(
		"Congrats! You find the word '%s' in %d/%d hands!\n",
		w.Answer(), w.curGame, w.gameCount)
}

func (w *wordle) TimeUp() bool {
	return w.curGame >= w.gameCount
}

func (w *wordle) ProcessMessage() string {
	return fmt.Sprintf("Used %d/%d hands.\n",
		w.curGame, w.gameCount)
}
