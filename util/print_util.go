package util

import "fmt"

func PrintStringsSlice(stringsSlice []string) {

	for _, str := range stringsSlice {
		fmt.Println(str)
	}
}
