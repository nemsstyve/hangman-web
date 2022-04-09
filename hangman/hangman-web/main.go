package main

import (
	"fmt"
	"hangman/Utils"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func StartGame() Utils.HangManData {
	data := Utils.HangManData{Attempts: 10, Error: ""}
	data.Word = Utils.GetRandomWord(os.Args[1])
	data.ToFind = strings.Repeat("_", utf8.RuneCountInString(data.Word))
	data.HangmanPositions = Utils.ParseHangmanFile("./hangman.txt")

	Utils.RevealRandomLetter(&data)
	return data
}

func HandleHangman(w http.ResponseWriter, r *http.Request, data *Utils.HangManData) {

	l := r.FormValue("Letter")
	if l != "" {
		Utils.HangMan(data, l)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {

	if len(os.Args) != 2 {
		return
	}

	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	rand.Seed(time.Now().UnixNano())

	data := StartGame()

	fs := http.FileServer(http.Dir("./css"))

	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("./index.html"))
		if err := tmpl.Execute(rw, data); err != nil {
			fmt.Println(err)
		}
	})
	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		HandleHangman(w, r, &data)
	})
	http.HandleFunc("/replay", func(w http.ResponseWriter, r *http.Request) {
		data = StartGame()
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}
