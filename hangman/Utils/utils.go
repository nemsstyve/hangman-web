package Utils

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

type HangManData struct {
	Word             string
	ToFind           string
	Attempts         int
	LettersSuggested []string
	Error            string
	HangmanPositions [10]string
}

func ParseHangmanFile(filename string) [10]string {
	var hangmanPositions [10]string
	count := 0
	index := 0

	f, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if count%8 == 0 && count != 0 {
			index += 1
		}
		line := "\n" + scanner.Text()
		hangmanPositions[index] += line
		count += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return hangmanPositions
}

func GetRandomWord(fileName string) string {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	return txtlines[rand.Intn(len(txtlines)-1)]
}

func IndexToReveal(data *HangManData, letter string) []int {
	var results []int

	index := len(data.Word)
	tmp := data.Word
	for {
		match := strings.LastIndex(tmp[0:index], letter)
		if match == -1 {
			break
		} else {
			index = match
			results = append(results, match)
		}
	}

	return results
}

func RevealLetters(data *HangManData, indexes []int) {
	tmp := []rune(data.ToFind)
	for _, c := range indexes {
		tmp[c] = rune(data.Word[c])
	}
	data.ToFind = string(tmp)
}

func RevealRandomLetter(data *HangManData) {
	nbr := len(data.Word)/2 - 1
	var indexes []int = make([]int, nbr)
	for i := 0; i < nbr; i++ {
		indexes[i] = rand.Intn(len(data.Word) - 1)
	}
	RevealLetters(data, indexes)
}

func letterIsAlreadyPresent(data *HangManData, l string) bool {
	for _, letter := range data.LettersSuggested {
		if letter == l {
			return true
		}
	}
	return false
}

func HangMan(data *HangManData, l string) string {

	data.Error = ""
	if len(l) > 1 {
		if l == data.Word {
			data.ToFind = data.Word
		} else {
			data.Attempts -= 2
			data.Error = fmt.Sprintf("The word is not %s, %d attempts remaining\n", l, data.Attempts)
		}
	} else if !letterIsAlreadyPresent(data, l) {
		data.LettersSuggested = append(data.LettersSuggested, l)
		indexes := IndexToReveal(data, l)
		if len(indexes) > 0 {
			RevealLetters(data, indexes)
		} else {
			data.Attempts -= 1
			data.Error = fmt.Sprintf("Not present in the word, %d attempts remaining\n", data.Attempts)
		}
	} else {
		data.Error = fmt.Sprintf("Letter '%s' already suggested\n", l)
	}
	return data.Error
}
