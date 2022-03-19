package wordle

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/sayuen0/wordle-go/token"
	"github.com/sayuen0/wordle-go/util"
)

const testQuestion = `water
melon
apple
hello
world
`

type dummyRandom struct{}

func (d dummyRandom) Intn(n int) int {
	return 1
}

var testWords = []string{"water", "melon", "apple", "hello", "world"}
var wrongTestWords = []string{"water", "halloween"}

func TestOpenQuestion(t *testing.T) {
	test := []struct {
		question string
		wantErr  bool
		err      error
	}{
		{question: testQuestion},
		{question: "hello\nhalloween\neaster", wantErr: true, err: errors.New("invalid source, word: halloween")},
		{question: "", wantErr: true, err: errors.New("no word found for question in source")},
	}
	for _, tt := range test {
		w := New(nil)
		if err := w.OpenQuestion(strings.NewReader(tt.question), dummyRandom{}); (err != nil) != tt.wantErr {
			t.Errorf("error opening question: %v", err)
		}
	}
}

func TestRule(t *testing.T) {
	test := []struct {
		count  int
		length int
		want   string
	}{
		{count: 10, length: 5, want: "rule\n" +
			"\t- \033[1;104mBlue\033[0m letter means the letter is not used in answer.\n" +
			"\t- \033[1;102mGreen\033[0m letter means the letter is used in the answer but in a different position.\n" +
			"\t- \033[1;101mRed\033[0m letter means the letter is used in the answer and in the same position."},
	}
	for _, tt := range test {
		w := New(&WordleOption{GameCount: 10})
		if got := w.Rule(); got != tt.want {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}

func TestToken(t *testing.T) {
	type tokenFunc func(string) token.Token
	test := []struct {
		s         string
		tokenFunc tokenFunc
		want      token.Token
	}{
		{s: "water", tokenFunc: wordle{}.ExactMatchToken, want: token.Token{Literal: "water", ColorScheme: exactColor}},
		{s: "final", tokenFunc: wordle{}.SomeMatchToken, want: token.Token{Literal: "final", ColorScheme: someMatchColor}},
		{s: "match", tokenFunc: wordle{}.NoMatchToken, want: token.Token{Literal: "match", ColorScheme: noMatchColor}},
	}
	for _, tt := range test {
		if got := tt.tokenFunc(tt.s); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}

func assertEqual(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSetInput(t *testing.T) {
	w := &wordle{
		validateFunc: util.ValidLowerAlphabet(5),
		gameCount:    6,
		inputs:       make([][]token.Token, 0),
	}
	if err := w.OpenQuestion(strings.NewReader(testQuestion), dummyRandom{}); err != nil {
		t.Fatal(err)
	}
	w.SetInput("write")
	assertEqual(t, len(w.inputs), 1)
}

func Test_answer(t *testing.T) {
	test := []struct {
		words []string
		ans   int
		want  string
	}{
		{words: testWords, ans: 2, want: "apple"},
		{words: testWords, ans: 3, want: "hello"},
		{words: testWords, ans: 4, want: "world"},
	}
	for _, tt := range test {
		w := &wordle{}
		w.setQuestion(tt.words, tt.ans)
		if got := w.Answer(); got != tt.want {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}

func Test_createToken(t *testing.T) {
	test := []struct {
		input string
		want  []token.Token
	}{
		{input: "water", want: []token.Token{
			{Literal: "w", ColorScheme: noMatchColor},
			{Literal: "a", ColorScheme: someMatchColor},
			{Literal: "t", ColorScheme: noMatchColor},
			{Literal: "e", ColorScheme: someMatchColor},
			{Literal: "r", ColorScheme: noMatchColor},
		}},
		{input: "anger", want: []token.Token{
			{Literal: "a", ColorScheme: exactColor},
			{Literal: "n", ColorScheme: noMatchColor},
			{Literal: "g", ColorScheme: noMatchColor},
			{Literal: "e", ColorScheme: someMatchColor},
			{Literal: "r", ColorScheme: noMatchColor},
		}},
	}

	w := &wordle{}
	w.setQuestion(testWords, 2)
	for _, tt := range test {
		if got := w.createTokens(tt.input); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("want %v, got %v", tt.want, got)
		}
	}
}

func TestExactMatchWord(t *testing.T) {
	w := New(nil)
	if err := w.OpenQuestion(strings.NewReader(testQuestion), dummyRandom{}); err != nil {
		t.Fatal(err)
	}
	test := []struct {
		input string
		want  bool
	}{
		{input: "water", want: false},
		{input: "world", want: false},
		{input: "melon", want: true},
	}
	for _, tt := range test {
		assertEqual(t, w.ExactMatchWord(tt.input), tt.want)
	}
}

func TestWinMessage(t *testing.T) {
	w := New(&WordleOption{GameCount: 3})
	if err := w.OpenQuestion(strings.NewReader(testQuestion), dummyRandom{}); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 0/3 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 1/3 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 2/3 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 3/3 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 4/3 hands!\n")

	w = New(&WordleOption{GameCount: 5})
	if err := w.OpenQuestion(strings.NewReader(testQuestion), dummyRandom{}); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 0/5 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 1/5 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 2/5 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 3/5 hands!\n")
	w.SetInput("aaaaa")
	assertEqual(t, w.WinMessage(), "Congrats! You find the word 'melon' in 4/5 hands!\n")
}

func TestTimeUp(t *testing.T) {
	w := New(&WordleOption{GameCount: 3})
	if err := w.OpenQuestion(strings.NewReader(testQuestion), dummyRandom{}); err != nil {
		t.Fatal(err)
	}
	w.SetInput("aaaaa")
	assertEqual(t, w.TimeUp(), false)
	w.SetInput("aaaaa")
	assertEqual(t, w.TimeUp(), false)
	w.SetInput("aaaaa")
	assertEqual(t, w.TimeUp(), true)

}
