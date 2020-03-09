package main

import (
	"fmt"

	"github.com/golem/commandlineparsers"
	"github.com/golem/commandlineprocessors"
	"github.com/golem/globalstringsproviders"
)

func main() {
	commandlineparsers.PrepareCommandLineProcessing()

	fmt.Println(globalstringsproviders.GetMenuPictureString())
	commandlineparsers.ParseCommandLine()

	parameters := commandlineparsers.GetParametersDict()
	processCommandLineProcessorError := commandlineprocessors.ProcessCommandLine(parameters)
	if processCommandLineProcessorError != nil {
		fmt.Println(fmt.Sprintf("%s: %s", "Golem->main->commandlineprocessors.ProcessCommandLine:", processCommandLineProcessorError.Error()))
	}
}
