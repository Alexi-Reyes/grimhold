package repl

import (
	"fmt"
	"sort"
)

type CliCommand struct {
	Name string
	Alias string
	Desc string
	Callback func(args ...string) error
}

var CliCommands map[string]CliCommand

func init() {
	CliCommands = map[string]CliCommand{
		"help": {
			Name:  "help",
			Alias: "h",
			Desc:  "Show available commands",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "help takes no arguments."))
					return nil
				}

				listCommands()
				return nil
			},
		},
		"exit": {
			Name: "exit",
			Desc: "Exit the program",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "exit takes no arguments."))
					return nil
				}

				exitProgram()
				return nil
			},
		},
		"map": {
			Name: "map",
			Desc: "List all known regions and connections",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "map takes no arguments."))
					return nil
				}

				return showMap()
			},
		},
		"travel": {
			Name: "travel",
			Desc: "Travel to a connected region (usage: travel <region>)",
			Callback: func(args ...string) error {
				if len(args) == 0 {
					fmt.Println(Colorize(Yellow, "Please enter a destination. (usage: travel <region>)"))
					return nil
				}

				if len(args) != 1 {
					fmt.Println(Colorize(Yellow, "More than one argument. (usage: travel <region>)"))
					return nil
				}

				return travelTo(args[0])
			},
		},
		"look": {
			Name: "look",
			Desc: "Describe your surroundings and exits",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "look takes no arguments."))
					return nil
				}

				return look()
			},
		},
		"explore": {
			Name: "explore",
			Desc: "Wander the region (risk an encounter)",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "explore takes no arguments."))
					return nil
				}

				return explore()
			},
		},
		"bestiary": {
			Name: "bestiary",
			Desc: "Browse monsters in the bestiary",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "bestiary takes no arguments."))
					return nil
				}

				return runBestiary()
			},
		},
		"inspect": {
			Name: "inspect",
			Desc: "View monster stat block (usage: inspect <monster>)",
			Callback: func(args ...string) error {
				if len(args) == 0 {
					fmt.Println(Colorize(Yellow, "Please enter a monster name. (usage: inspect <monster>)"))
					return nil
				}

				if len(args) != 1 {
					fmt.Println(Colorize(Yellow, "More than one argument. (usage: inspect <monster>)"))
					return nil
				}

				return inspect(args[0])
			},
		},
		"character": {
			Name:  "character",
			Desc:  "Display your character sheet",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "character takes no arguments."))
					return nil
				}

				showCharacter()
				return nil
			},
		},
		"rest": {
			Name:  "rest",
			Desc:  "Rest to restore HP (once per region visit)",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "rest takes no arguments."))
					return nil
				}

				rest()
				return nil
			},
		},
		"save": {
			Name: "save",
			Desc: "Save the game state to save.json",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "save takes no arguments."))
					return nil
				}

				saveGameCommand()
				return nil
			},
		},
		"load": {
			Name: "load",
			Desc: "Load the game state from save.json",
			Callback: func(args ...string) error {
				if len(args) != 0 {
					fmt.Println(Colorize(Yellow, "load takes no arguments."))
					return nil
				}

				loadGameCommand()
				return nil
			},
		},
	}
}

func listCommands() {
	fmt.Println(Colorize(BoldCyan, "\nGrimhold — Commands:"))

	keys := make([]string, 0, len(CliCommands))
	for k := range CliCommands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		cmd := CliCommands[k]
		nameStr := cmd.Name
		if cmd.Alias != "" {
			nameStr = fmt.Sprintf("%s (%s)", cmd.Name, cmd.Alias)
		}
		fmt.Printf("  %-20s %s\n", Colorize(BoldWhite, nameStr), cmd.Desc)
	}
	fmt.Println()
}
