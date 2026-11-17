package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/repl"
)

const (
	redText = "\033[31m"
    resetText = "\033[0m"
)

func main() {
    fmt.Println("Welcome to grimhold, a text-based dungeon crawler that runs in your terminal. Type help to see what actions you can take.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("grimhold > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		input = strings.TrimSpace(input)
		if cmd, ok := repl.CliCommands[input]; ok {
		   	if err := cmd.Callback(os.Args[1:]...); err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(1)
            }
		} else if input == "" {
            fmt.Printf("%sPlease enter a command.%s\n", redText, resetText)
		} else {
            fmt.Printf("%sUnknown command. Type help to get a list of commands.%s\n", redText, resetText)
		}
	}
}
