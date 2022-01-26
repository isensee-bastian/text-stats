package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Structs are good for bundling data in a composite type.
type wordStat struct {
	word  string
	count int
}

func main() {
	words := readWordsFromFile("input.txt")
	words = cleanupWords(words)
	fmt.Printf("Total count of words: %v\n", len(words))

	// Start measuring.
	start := time.Now()
	mapResult := countWordOccurrencesFast(words)
	// Stop measuring.
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Duration for counting with map implementation: %v\n", duration)
	fmt.Printf("Counting with map result: %v\n", mapResult)

	wordStats := countWordOccurrencesSlow(words)

	// Pass inline function for comparing elements. Needed to perform actual
	// sorting since we need to specify our sort criteria. Note that we
	// have access to enclosing variables like wordStats.
	sort.Slice(wordStats, func(leftIndex, rightIndex int) bool {
		return wordStats[leftIndex].count < wordStats[rightIndex].count
	})
	fmt.Printf("Sort result wordStats: %v\n", wordStats)

	wordToCount := statSliceToMap(wordStats)
	input := readLineFromTerminal()

	count, found := wordToCount[input]
	if found {
		fmt.Printf("Found it! Word count: %d\n", count)
	} else {
		fmt.Printf("Word is not present.\n")
	}
}

// This is a simple and efficient approach for organizing elements that have
// a common characteristic. Our key is the word and our value the count.
// Checking our map for existing words is easy and efficient, no searching required.
func countWordOccurrencesFast(words []string) map[string]int {
	wordToCount := map[string]int{}

	for _, word := range words {
		count, found := wordToCount[word]

		if found {
			wordToCount[word] = count + 1
		} else {
			wordToCount[word] = 1
		}
	}

	return wordToCount
}

// Build a map from a slice of wordStat elements. The word becomes the key
// for simple and efficient lookup. The count is our value.
func statSliceToMap(stats []wordStat) map[string]int {
	var wordToCount map[string]int = make(map[string]int)

	for _, stat := range stats {
		wordToCount[stat.word] = stat.count
	}

	return wordToCount
}

// Careful: This is a naive approach that is inefficient for large
// input slices. E.g.:
// 1 mio words * len(wordStats) <-- grows fast and becomes slow to process
func countWordOccurrencesSlow(words []string) []wordStat {
	var wordStats []wordStat = make([]wordStat, 0)

	// For every word: Look for other elements containing the same word and increment.
	for _, word := range words {
		found := false

		// Keep in mind: range returns a copy of the element. Since we want to modify
		// the original element we need to use access by index.
		for index, stat := range wordStats {
			if word == stat.word {
				// Found it: increment count.
				wordStats[index].count++
				found = true
			}
		}

		if !found {
			// Create a fresh struct instance with count 1.
			wordStats = append(wordStats, wordStat{word, 1})
		}
	}

	return wordStats
}

func cleanupWords(words []string) []string {
	cleanedWords := make([]string, 0)

	// Quick and dirty way is to use string replace:
	// removableStrings := []string{",", ".", ";", "[", "]", "{", "}", "(", ")", "%", "\""}

	removablePattern := regexp.MustCompile("[^a-z]")

	for _, word := range words {
		word = strings.ToLower(word)

		/* See regexp version for more efficient and robust solution.
		for _, removable := range removableStrings {
			word = strings.ReplaceAll(word, removable, "")
		}
		*/

		word = removablePattern.ReplaceAllString(word, "")

		if len(word) > 1 || word == "a" {
			cleanedWords = append(cleanedWords, word)
		}
	}

	return cleanedWords
}

func readLineFromTerminal() string {
	scanner := bufio.NewScanner(os.Stdin)
	// By default scanner splits after a new line.
	fmt.Println("Enter word to search count for (confirm with enter):")

	if scanner.Scan() {
		// Normalized input for counting independently of case.
		return strings.ToLower(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error on reading input: %v\n", err)
	}

	return ""
}

func readWordsFromFile(name string) []string {
	// Go supports multiple return values.
	file, err := os.Open(name)
	// Defer ensures that your cleanup call is performed before exiting the function,
	// even if an error occurrs. Use it to release system resources after usage.
	defer file.Close()

	// Check for potential errors.
	if err != nil {
		fmt.Printf("Error on reading file: %v\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error on reading words: %v\n", err)
	}

	return words
}
