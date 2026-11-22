package main

import (
	"bufio"
	"fmt"
	"github.com/Alexi-Reyes/grimhold/internal/repl"
	"os"
	"strings"
)

func main() {
	fmt.Println(repl.Colorize(repl.BoldCyan, "----- Welcome to Grimhold! -----"))
	fmt.Println(repl.Colorize(repl.BoldCyan, "a text-based dungeon crawler that runs in your terminal"))
	fmt.Println(repl.Colorize(repl.Dim, "Type 'help' to see available commands, or 'look' to check surroundings.\n"))

	if err := repl.StartGame(); err != nil {
		fmt.Printf("%s\n", repl.Colorize(repl.BoldRed, fmt.Sprintf("Error initializing game: %v", err)))
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(repl.Colorize(repl.BoldCyan, "grimhold > "))
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		fields := strings.Fields(input)
		if len(fields) == 0 {
			fmt.Println(repl.Colorize(repl.Yellow, "Please enter a command (type 'help' for a list)."))
			continue
		}

		cmd, ok := repl.CliCommands[fields[0]]
		if !ok {
			fmt.Println(repl.Colorize(repl.Red, "Unknown command. Type 'help' to get a list of commands."))
			continue
		}

		if err := cmd.Callback(fields[1:]...); err != nil {
			fmt.Printf("%s\n", repl.Colorize(repl.BoldRed, fmt.Sprintf("Error: %v", err)))
		}
	}
}
