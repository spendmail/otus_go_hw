package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
)

const sliceLen = 10

var ErrTextShorterThenNeeded = fmt.Errorf(fmt.Sprintf("given text is shorter than %d unique words", sliceLen))

var regexpSplitter = regexp.MustCompile(`\s+`)

type wordCounterStruct struct {
	Word    string
	Counter int
}

func Top10(sourceText string) ([]string, error) {
	// If text is empty, return empty slice
	if sourceText == "" {
		return []string{}, nil
	}

	resultSlice := make([]string, sliceLen)
	var wordsSlice []string
	wordsMap := make(map[string]int)
	wordCounterSlice := make([]wordCounterStruct, 100)

	// Split text to a separate words
	wordsSlice = regexpSplitter.Split(sourceText, -1)

	// Counting the number of words
	for _, value := range wordsSlice {
		wordsMap[value]++
	}

	// If given text is shorter than 10 unique words, return an error
	if len(wordsMap) < sliceLen {
		return []string{}, ErrTextShorterThenNeeded
	}

	// Write words to a slice of structs for sorting
	for word, counter := range wordsMap {
		wordCounterSlice = append(wordCounterSlice, wordCounterStruct{word, counter})
	}

	// Sorting the slice "lexicographically"
	sort.Slice(wordCounterSlice, func(i, j int) bool {
		if wordCounterSlice[i].Counter == wordCounterSlice[j].Counter {
			return wordCounterSlice[i].Word < wordCounterSlice[j].Word
		}
		return wordCounterSlice[i].Counter > wordCounterSlice[j].Counter
	})

	// Taking first ${sliceLen} items
	for i := 0; i < sliceLen; i++ {
		resultSlice[i] = wordCounterSlice[i].Word
	}

	return resultSlice, nil
}
