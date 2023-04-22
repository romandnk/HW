package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type word struct {
	str   string
	count int
}

func Top10(s string) []string {
	strArr := strings.Fields(s)
	var (
		numWords = make([]word, 0, len(strArr))
		// a map with the number of repetitions of words
		countSameWords = make(map[string]int, len(strArr))
		result         []string
	)

	for i := range strArr {
		_, ok := countSameWords[strArr[i]]
		if ok {
			countSameWords[strArr[i]]++
		} else {
			countSameWords[strArr[i]] = 1
		}
	}
	for key, val := range countSameWords {
		numWords = append(numWords, word{key, val})
	}
	// sort slice numWords by the number of repetitions or if the number of repetitions is the same then lexicographically
	sort.Slice(numWords, func(i, j int) bool {
		descend := numWords[i].count > numWords[j].count
		descendAndLexicographically := numWords[i].count == numWords[j].count && numWords[i].str < numWords[j].str
		return descend || descendAndLexicographically
	})
	if len(numWords) < 10 {
		for i := 0; i < len(numWords); i++ {
			result = append(result, numWords[i].str)
		}
	} else {
		for i := 0; i < 10; i++ {
			result = append(result, numWords[i].str)
		}
	}
	return result
}
