package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(in string) []string {
	lengthOutput := 10

	words := strings.Fields(in)

	mapWords := make(map[string]int)

	for _, word := range words {
		mapWords[word]++
	}

	topWords := make([]string, 0, len(mapWords))

	for word := range mapWords {
		topWords = append(topWords, word)
	}

	sort.Slice(topWords, func(i, j int) bool {
		if mapWords[topWords[i]] == mapWords[topWords[j]] {
			return (topWords[i] < topWords[j])
		}
		return (mapWords[topWords[i]] > mapWords[topWords[j]])
	})

	if len(mapWords) < lengthOutput {
		lengthOutput = len(mapWords)
	}

	return topWords[0:lengthOutput]
}
