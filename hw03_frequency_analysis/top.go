package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`(?m)(\b|\s)[^a-zA-Zа-яА-Я0-9]+(\b|\s)|^[^a-zA-Zа-яА-Я0-9]+|[^a-zA-Zа-яА-Я0-9]+$`)

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

func Top10Additional(s string) []string {
	strArr := strings.Fields(strings.ToLower(s))
	numberOfRepetitions := make(map[string]int, len(strArr))
	numWords := make([]word, 0, len(strArr))
	var result []string
	for _, word := range strArr {
		switch {
		case !re.MatchString(word):
			_, ok := numberOfRepetitions[word]
			if ok {
				numberOfRepetitions[word]++
			} else {
				numberOfRepetitions[word] = 1
			}
		case word == "-":
			continue
		default:
			word = re.ReplaceAllString(word, "")
			_, ok := numberOfRepetitions[word]
			if ok {
				numberOfRepetitions[word]++
			} else {
				numberOfRepetitions[word] = 1
			}
		}
	}
	for key, val := range numberOfRepetitions {
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
