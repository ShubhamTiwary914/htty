/*
	prefix completions taking some prefix substring to get K closest & shortest results from a file source
	ex: hel --> [hell, hello, ...]
*/

package htty

import (
	"bufio"
	"os"
	"slices"
	"sort"
	"strings"
)

func PrefixClosestSearch(substr string, filePath string) ([]string){
	return []string{}	
}

/*	Read from .txt file that have word per line
		Returns (string array of res[i]= ith line of file)
*/
func ReadTextLines_intoList(filePath string)([]string){
	var words[] string
	file, err := os.Open(filePath)
	if err != nil {
		Errorf("file %s doesn't exist!", filePath)	
		panic(err)
	}
	defer file.Close()

	ptr := bufio.NewScanner(file)
	for ptr.Scan(){
		words = append(words, ptr.Text())
	}
	return words
}

// Dump word into file at (filePath) with insertion sorted manner
func Insert_SortedFile(word string, filePath string){
	//sort line words string[] with case sensitivity (apple != Apple)
	var words []string = ReadTextLines_intoList(filePath)		
	if slices.Contains(words, word) {
		//skip if word already exists
		return;
	}
	words = append(words, word)
	sort.Slice(words, func(i, j int) bool {
		w1 := strings.ToLower(words[i])
		w2 := strings.ToLower(words[j])

		if w1 == w2 {
			return words[i] < words[j] 
		}
		return w1 < w2
	})
	//dump to file
	file, err := os.Create(filePath)
	if err != nil {
		Errorf("file %s unable to create!", filePath)	
		panic(err)
	}
	defer file.Close()

	//buffer writes
	wptr := bufio.NewWriter(file)
	for lineidx, word := range words {
		if _, err := wptr.WriteString(word + "\n"); err != nil {
			Errorf("unable to write word %s into file %s (line %d)", word, filePath, lineidx)	
			panic(err)
		}
	}
	//push saved buffer onto file
	if err := wptr.Flush(); err != nil {
		Errorf("unable to write onto file %s", filePath)
		panic(err)
	}
}

func CheckFileExists(filePath string) (bool) {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	file.Close()	
	return true
}

