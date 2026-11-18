package repl

import (
    "fmt"
    "os"
    "github.com/Alexi-Reyes/grimhold/internal/json"
    "github.com/Alexi-Reyes/grimhold/internal/game"
)

type CliCommand struct {
	Name     string
	Alias    string
	Desc     string
	Callback func(args ...string) error
}

var CliCommands map[string]CliCommand

func init() {
	CliCommands = map[string]CliCommand{
		"help": {
			Name: "help",
			Alias: "h",
			Desc: "Show this message",
			Callback: func(args ...string) error { listCommands(); return nil },
		},
        "exit": {
            Name: "exit",
            Desc: "Exit the program",
			Callback: func(args ...string) error { exitProgram(); return nil },
        },
        "map": {
            Name: "map",
            Desc: "List the known regions",
        },
        "travel": {
            Name: "travel",
            Desc: "Travel to a connected region (usage: travel <region>)",
        },
        "look": {
            Name: "look",
            Desc: "Describe your surroundings",
            Callback: func(args ...string) error { look(game.CurrentRegion); return nil },
        },
        "explore": {
            Name: "explore",
            Desc: "Wander the region (risk an encounter)",
        },
	}
}

func listCommands() {
    fmt.Println("Commands:")
	for _, command := range CliCommands {
		if command.Alias != "" {
			fmt.Printf("%s (%s): %s\n", command.Name, command.Alias, command.Desc)
		} else {
			fmt.Printf("%s: %s\n", command.Name, command.Desc)
		}
	}
}

func exitProgram() {
    fmt.Println("Thanks for playing.")
    os.Exit(0)
}

func showMap() {
}

func travelTo() {
}

func look(input string) (error) {
    region, err := json.GetRegion(input)
    if err != nil {
        return err
    }

    fmt.Println(region.Description)
    fmt.Printf("Danger level: %v\n", region.Danger)
    fmt.Printf("Exits: %v\n", region.Exits)
    return nil
}

func explore() {
}

func StartGame() {
    err := json.LoadStartData()
    if err != nil {
        fmt.Println(err)
        return
    }
}