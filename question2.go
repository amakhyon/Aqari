package main

import (
	"fmt"
	"sort"
)

func rearrangeString(inputStr string) string {
	charCount := make(map[rune]int)

	for _, char := range inputStr { //count how many times each letter appeared
		charCount[char]++
	}

	// Create a slice of characters
	sortedChars := make([]rune, 0, len(charCount))
	for char := range charCount {
		sortedChars = append(sortedChars, char)
	}
	sort.Slice(sortedChars, func(i, j int) bool { //sort them by how many times they appeared
		return charCount[sortedChars[i]] > charCount[sortedChars[j]]
	})

	maxCount := charCount[sortedChars[0]] //check if it's possible to rearrange the string without having to deleete any character
	if maxCount > (len(inputStr)+1)/2 {
		return ""
	}

	// Rearrange the characters
	result := make([]rune, len(inputStr))
	index := 0
	for _, char := range sortedChars {
		for charCount[char] > 0 {
			result[index] = char
			index += 2
			if index >= len(inputStr) {
				index = 1
			}
			charCount[char]--
		}
	}

	return string(result)
}

func main() {
	fmt.Println(rearrangeString("aab"))
	fmt.Println(rearrangeString("aaab"))
	fmt.Println(rearrangeString("yyyxxxzzz"))
	fmt.Println(rearrangeString("yytyye"))
}
