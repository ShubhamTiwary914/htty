//pure utility methods like calculations, rounding, array/object operations, etc..
package utils 

import (
	"os"
	"os/exec"
	"runtime"
	"bufio"
	"crypto/rand"
	"encoding/hex"
)

func GetPercent(percentage int, source int) int{
	return (percentage * source)/100
}

//write raw string contents into a file at filepath
func WriteFileContents(filePath string, contents string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	//write into buffer & flush push writes 
	wptr := bufio.NewWriter(file)
	wptr.WriteString(contents)
	if err := wptr.Flush(); err != nil {
		Errorf("unable to write onto file %s", filePath)
		return err;
	}
	return nil
}


func GenerateRandomUUID(size int) string{
	//1 byte = 2 hex characters
	numBytes := (size + 1) / 2
	bytes := make([]byte, numBytes)
	rand.Read(bytes)
	
	var uuid string = hex.EncodeToString(bytes)
	if len(uuid) > size {
		uuid = uuid[:size]
	}
	return uuid
}

//pops up dialog to save file(depending on OS may vary)
func SaveFileDialog(defaultName string) (string, error) {
	var cmd *exec.Cmd
	//TODO: adding support for darwin & windows 
	switch runtime.GOOS {
	case "linux": 
		cmd = exec.Command("zenity", "--file-selection", "--save", "--filename="+defaultName, "--confirm-overwrite")
	}
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// trim whitespace/newlines
	result := string(output)
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}
