// Go Predictions by Tim Lucca
// This program outputs a list of words based on what the user types

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/marcusolsson/tui-go"
)

var (
	root        *node
	once        sync.Once
	predictions [20]string
)

type node struct {
	letter   string
	word     bool
	children [26]*node
}

// ensures the root is only initialized once
func init() {
	once.Do(initTree)
}

// initializes the root of the tree
func initTree() {
	p := node{
		letter: "",
		word:   false,
	}
	root = &p
}

// creates branches based on a given word
// takes in a node pointer and a string
func buildTree(current *node, s string) {
	i := []rune(s)[0] - 97
	if i < 0 || i > 25 {
		return
	}
	if current.children[i] == nil {
		p := node{
			letter: s[0:1],
			word:   false,
		}
		if len(s) == 1 {
			p.word = true
		}
		current.children[i] = &p
	}
	if len(s) > 1 {
		buildTree(current.children[i], s[1:])
	}

}

// reads all of the words from the input file and builds the tree
func readAndBuild() {
	file, e := os.Open("word.txt")
	if e != nil {
		log.Fatal(e)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		buildTree(root, strings.ToLower(scanner.Text()))
	}
}

// used for resetting the list of predictions
func resetPredictions() {
	for i := 0; i < len(predictions); i++ {
		predictions[i] = ""
	}
}

// sets predictions based on input string
// takes in a string
func setPredictions(s string) {
	resetPredictions()
	if s == "" {
		return
	}
	current := getCurrent(root, s)
	if current == nil {
		return
	}
	_ = traverse(current, s, 0)
}

// finds the deepest node based on the input string
// takes in a node pointer and a string; returns a node pointer
func getCurrent(current *node, s string) *node {
	i := []rune(s)[0] - 97
	if i < 0 || i > 25 {
		return nil
	}
	if len(s) > 1 {
		if current.children[i] == nil {
			return nil
		} else {
			return getCurrent(current.children[i], s[1:])
		}
	} else {
		return current.children[i]
	}
}

// traverses the tree based on some root node, searches for nodes that are endpoints for words
// takes in a node pointer, a string, and an int; returns an int
func traverse(current *node, t string, n int) int {
	if current.word && n < len(predictions) {
		predictions[n] = t
		n++
	}
	if n < len(predictions) {
		for x, _ := range current.children {
			if current.children[x] != nil {
				n = traverse(current.children[x], t+current.children[x].letter, n)
			}
		}
	}
	return n
}

// main function, calls the readandbuild, creates the gui
func main() {

	readAndBuild()

	words := tui.NewVBox()

	wordsScroll := tui.NewScrollArea(words)
	wordsScroll.SetAutoscrollToBottom(true)

	wordsBox := tui.NewVBox(wordsScroll)
	wordsBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	prompt := tui.NewLabel(" Input | ")

	inputBox := tui.NewHBox(prompt, input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	// updates the output box whenever a new character is typed
	input.OnChanged(func(e *tui.Entry) {
		setPredictions(strings.ToLower(e.Text()))

		next := tui.NewHBox(
			tui.NewLabel("Showing predictions for: "),
			tui.NewLabel(strings.ToLower(e.Text())),
			tui.NewSpacer(),
		)
		next.SetBorder(true)

		words.Append(next)

		for _, w := range predictions {
			words.Append(tui.NewHBox(
				tui.NewLabel(w),
				tui.NewSpacer(),
			))
		}
	})

	app := tui.NewVBox(wordsBox, inputBox)
	app.SetSizePolicy(tui.Expanding, tui.Expanding)

	ui, err := tui.New(app)
	if err != nil {
		panic(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
