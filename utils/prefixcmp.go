package htty

import (
	"bufio"
	"os"
	"slices"
	"sort"
	"strings"
)

/*
	Prefix completions taking some prefix substring to get closest & shortest results from a file source
	ex: hel --> [hell, hello, Hello...] (NOTE: it is case insensitive)
*/
func PrefixClosestSearch(substr string, filePath string) ([]string, error){
	var words []string
	var err error
	words, err = ReadTextLines_intoList(filePath)
	if err != nil { return nil, err; }
	return PrefixClosestSearch_withOptions(substr, words);
}


/*
	Prefix completions taking some prefix substring to get closest & shortest results from a source string aray
	ex: hel --> [hell, hello, Hello...] (NOTE: it is case insensitive)
*/

func PrefixClosestSearch_withOptions(substr string, words []string) ([]string, error){	//first index where prefix matches
	prefix := strings.ToLower(substr)
	start := sort.Search(len(words), func(i int) bool {
		return strings.ToLower(words[i]) >= prefix
	})
	// no possible matches
	if start >= len(words) || !strings.HasPrefix(strings.ToLower(words[start]), prefix) {
		return []string{}, nil
	}
	// scan forward until prefix doesn't match
	end := start
	for end < len(words) && strings.HasPrefix(strings.ToLower(words[end]), prefix) {
		end++
	}
	return words[start:end], nil
	
}

/*	Read from .txt file that have word per line
		Returns (string array of res[i]= ith line of file)
*/
func ReadTextLines_intoList(filePath string) ([]string, error){
	var words[] string
	file, err := os.Open(filePath)
	if err != nil {
		Errorf("file %s doesn't exist!", filePath)	
		return nil, err
	}
	defer file.Close()

	ptr := bufio.NewScanner(file)
	for ptr.Scan(){
		words = append(words, ptr.Text())
	}
	return words, nil 
}

// Dump word into file at (filePath) with insertion sorted manner
func Insert_SortedFile(word string, filePath string) error {
	//sort line words string[] with case sensitivity (apple != Apple)
	var words []string;
	var err error;
	words, err = ReadTextLines_intoList(filePath)

	if slices.Contains(words, word) {
		//skip if word already exists
		return nil;
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
		return err
	}
	defer file.Close()

	//buffer writes
	wptr := bufio.NewWriter(file)
	for lineidx, word := range words {
		if _, err := wptr.WriteString(word + "\n"); err != nil {
			Errorf("unable to write word %s into file %s (line %d)", word, filePath, lineidx)	
			return err;	
		}
	}
	//push saved buffer onto file
	if err := wptr.Flush(); err != nil {
		Errorf("unable to write onto file %s", filePath)
		return err;
	}
	return nil
}

func CheckFileExists(filePath string) (bool) {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	file.Close()	
	return true
}
