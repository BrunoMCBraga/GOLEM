package util

import (
	"errors"
	"fmt"
	"os"
)

func WriteArrayOfStringsToFile(sliceOfStrings []string, filePath string) error {

	createResult, createError := os.Create(filePath)
	defer createResult.Close()

	var finalString string = ""

	if createError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "main->writeArrayOfStringsToFile->os.Create:", createError.Error()))
	}

	for _, str := range sliceOfStrings {
		finalString += fmt.Sprintf("%s\n", str)
	}
	createResult.WriteString(finalString)

	return nil
}
