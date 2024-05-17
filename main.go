package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type WordResult struct {
	YourWord                 string
	Definition               string
	IsPalindrore             bool
	DuplicateCharacter       []string
	WordPermutation          []string
	IsDefinitionDone         bool
	IsPalindromeDone         bool
	IsDuplicateCharacterDone bool
	IsWordPermutationDone    bool
}

func main() {
	fmt.Println("Welcome to word checker game")
	fmt.Print("What word do you want to check? (example: impostor) : ")
	var word string
	_, errScan := fmt.Scanln(&word)
	if errScan != nil {
		fmt.Println(errScan)
		return
	}

	var wg sync.WaitGroup
	wg.Add(4)

	resultChannel := make(chan WordResult)

	go func() {
		defer wg.Done()
		palindrome, _ := IsPalindrome(word)
		resultChannel <- WordResult{YourWord: word, IsPalindrore: palindrome, IsPalindromeDone: true}
	}()

	go func() {
		defer wg.Done()
		definition, errDefinition := GetDefinition(word)
		if errDefinition != nil {
			resultChannel <- WordResult{YourWord: word, Definition: "Word definition not found"}
		} else {
			resultChannel <- WordResult{YourWord: word, Definition: definition, IsDefinitionDone: true}
		}
	}()

	go func() {
		defer wg.Done()
		duplicateCharacter, _ := DuplicateCharacter(word)
		resultChannel <- WordResult{YourWord: word, DuplicateCharacter: duplicateCharacter, IsDuplicateCharacterDone: true}
	}()

	go func() {
		defer wg.Done()
		permutation, _ := WordPermutation(word)
		resultChannel <- WordResult{YourWord: word, WordPermutation: permutation, IsWordPermutationDone: true}
	}()

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var finalResult WordResult
	for result := range resultChannel {
		if result.IsDefinitionDone {
			finalResult.YourWord = result.YourWord
			finalResult.Definition = result.Definition
		}
		if result.IsPalindromeDone {
			finalResult.IsPalindrore = result.IsPalindrore
		}
		if result.IsDuplicateCharacterDone {
			finalResult.DuplicateCharacter = result.DuplicateCharacter
		}
		if result.IsWordPermutationDone {
			finalResult.WordPermutation = result.WordPermutation
		}
	}

	fmt.Println("----RESULT----")
	fmt.Println("Your word is", finalResult.YourWord)
	fmt.Println("Definition:", finalResult.Definition)
	fmt.Println("Is Palindrome:", finalResult.IsPalindrore)
	fmt.Println("Duplicate Character:", finalResult.DuplicateCharacter)
	fmt.Println("Word Permutation:", finalResult.WordPermutation)
}

func IsPalindrome(word string) (bool, error) {
	strBytes := []byte(word)
	totalLenStrBytes := len(strBytes)

	for i := 0; i < totalLenStrBytes/2; i++ {
		strBytes[i], strBytes[totalLenStrBytes-1-i] = strBytes[totalLenStrBytes-1-i], strBytes[i]
	}

	return word == string(strBytes), nil
}

func GetDefinition(word string) (string, error) {
	var apiResponse = []struct {
		Word      string `json:"word"`
		Phonetic  string `json:"phonetic"`
		Phonetics []struct {
			Text  string `json:"text"`
			Audio string `json:"audio,omitempty"`
		} `json:"phonetics"`
		Origin   string `json:"origin"`
		Meanings []struct {
			PartOfSpeech string `json:"partOfSpeech"`
			Definitions  []struct {
				Definition string        `json:"definition"`
				Example    string        `json:"example"`
				Synonyms   []interface{} `json:"synonyms"`
				Antonyms   []interface{} `json:"antonyms"`
			} `json:"definitions"`
		} `json:"meanings"`
	}{}

	apiUrl := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%v", word)
	client := &http.Client{}
	request, errRequest := http.NewRequest("GET", apiUrl, nil)
	if errRequest != nil {
		fmt.Println(errRequest)
		return "", errRequest
	}

	response, errResponse := client.Do(request)
	if errResponse != nil {
		fmt.Println(errResponse)
		return "", errRequest
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return "Word definition not found 😭", nil
	}

	errDecode := json.NewDecoder(response.Body).Decode(&apiResponse)
	if errDecode != nil {
		fmt.Println(errDecode)
		return "", errDecode
	}

	wordDefinition := apiResponse[0].Meanings[0].Definitions[0].Definition
	return wordDefinition, nil
}

func DuplicateCharacter(words string) ([]string, error) {
	charCount := make(map[rune]int)
	duplicates := []string{}

	for _, char := range words {
		charCount[char]++
	}

	for char, count := range charCount {
		if count > 1 {
			duplicates = append(duplicates, fmt.Sprintf("there are %d of word %s", count, string(char)))
		}
	}

	return duplicates, nil
}

func WordPermutation(word string) ([]string, error) {
	var results []string
	var generatePermutations func(prefix string, suffix string)

	generatePermutations = func(prefix string, suffix string) {
		if suffix == "" {
			results = append(results, prefix)
			return
		}

		for i := 0; i < len(suffix); i++ {
			newPrefix := prefix + string(suffix[i])
			newSuffix := suffix[:i] + suffix[i+1:]

			generatePermutations(newPrefix, newSuffix)
		}
	}

	generatePermutations("", word)

	return results, nil
}
