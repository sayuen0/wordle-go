package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sayuen0/wordle-go/wordle"
)

func main() {
	filepath := "./wordle.txt"
	if len(os.Args) > 1 {
		filepath = os.Args[1]
	}

	wordle := wordle.New(nil)
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := wordle.OpenQuestion(f, r); err != nil {
		panic(err)
	}
	fmt.Println(wordle.Rule())
	fmt.Println("Game start!!")
	for {
		var s string
		fmt.Print(">")
		fmt.Scan(&s)
		if s == "quit" || s == "exit" {
			fmt.Println("You quit the game...")
			return
		}
		// check valid word
		if !wordle.ValidInput(s) {
			fmt.Printf("Invalid input: %s \n", s)
			fmt.Println(wordle.Table())
			fmt.Println(wordle.ProcessMessage())
			continue
		}
		// search word
		wordle.SetInput(s)
		fmt.Println(wordle.Table())
		if wordle.ExactMatchWord(s) {
			fmt.Println(wordle.WinMessage())
			return
		}
		if wordle.TimeUp() {
			fmt.Printf("Game over... answer: %s\n", wordle.Answer())
			return
		}
		fmt.Println(wordle.ProcessMessage())
	}
}
