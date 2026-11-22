package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/repl"
)

const (
	redText = "\033[31m"
    resetText = "\033[0m"
)

func main() {
    fmt.Println("Welcome to grimhold, a text-based dungeon crawler that runs in your terminal. Type help to see what actions you can take.")
	repl.StartGame()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("grimhold > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fields := strings.Fields(input)
		if len(fields) == 0 {
            fmt.Printf("%sPlease enter a command.%s\n", redText, resetText)
		} else if cmd, ok := repl.CliCommands[fields[0]]; ok {
		   	if err := cmd.Callback(fields[1:]...); err != nil {
                fmt.Printf("%s%v%s\n", redText, err, resetText)
            }
		} else {
            fmt.Printf("%sUnknown command. Type help to get a list of commands.%s\n", redText, resetText)
		}
	}
}
