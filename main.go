package main

import (
	"fmt"

	"github.com/golem/commandlinegenerators"
	"github.com/golem/commandlineprocessors"
	"github.com/golem/globalstringsproviders"
)

func main() {
	commandlinegenerators.PrepareCommandLineProcessing()

	fmt.Println(globalstringsproviders.GetMenuPictureString())
	commandlinegenerators.ParseCommandLine()

	parameters := commandlinegenerators.GetParametersDict()
	processCommandLineProcessorError := commandlineprocessors.ProcessCommandLine(parameters)
	if processCommandLineProcessorError != nil {
		fmt.Println(fmt.Sprintf("%s: %s", "Golem->main->commandlineprocessors.ProcessCommandLine:", processCommandLineProcessorError.Error()))
	}
}
