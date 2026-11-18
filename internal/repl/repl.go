package repl

import (
    "fmt"
    "os"
    "slices"
    "github.com/Alexi-Reyes/grimhold/internal/json"
    "github.com/Alexi-Reyes/grimhold/internal/game"
    "github.com/Alexi-Reyes/grimhold/internal/open5e"
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
            Callback: func(args ...string) error { showMap(); return nil },
        },
        "travel": {
            Name: "travel",
            Desc: "Travel to a connected region (usage: travel <region>)",
            Callback: func(args ...string) error {
               if len(args) == 0 {
                    fmt.Println("Please enter a destination. (usage: travel <region>)")
                    return nil
               }

                if len(args) != 1 {
                    fmt.Println("More than one argument. (usage: travel <region>)")
                    return nil
                }

                return travelTo(args[0]) 
            },
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
        "bestiary": {
            Name: "bestiary",
            Desc: "List monsters",
            Callback: func(args ...string) error { bestiary(); return nil },
        },
        "inspect": {
            Name: "inspect",
            Desc: "View monster stats (usage: inspect <monster>).",
            Callback: func(args ...string) error {
               if len(args) == 0 {
                    fmt.Println("Please enter a monster. (usage: inspect <monster>)")
                    return nil
               }

                if len(args) != 1 {
                    fmt.Println("More than one argument. (usage: inspect <monster>)")
                    return nil
                }

                return inspect(args[0])
            },
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

func showMap() (error) {
    regions, err := json.GetRegions()
    if err != nil {
        return err
    }

    for _, region := range regions {
        fmt.Printf("Region: %v\n", region.Name)
        fmt.Printf("Danger level: %v\n", region.Danger)
        fmt.Printf("Exits: %v\n", region.Exits)
        fmt.Println("")
    }

    return nil
}

func travelTo(input string) error {
    region, err := json.GetRegion(game.CurrentRegion)
    if err != nil {
        return err
    }

    if slices.Contains(region.Exits, input) {
        fmt.Printf("Travelling to %v.\n", input)
        game.CurrentRegion = input
    } else {
        fmt.Println("No such exit for this region.")
    }
    return nil
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

func bestiary() {
    open5e.GetMonsters()
}

func inspect(input string) error {
    result, err := open5e.GetMonster(input)
    if err != nil {
        return err
    }

    if result == nil {
        fmt.Println("No such monster.")
        return nil
    }

    fmt.Println("Name:", result.Name)
    fmt.Println("Strength:", result.Strength)
    fmt.Println("Dexterity:", result.Dexterity)
    fmt.Println("Constitution:", result.Constitution)
    fmt.Println("Intelligence:", result.Intelligence)
    fmt.Println("Wisdom:", result.Wisdom)
    fmt.Println("Charisma:", result.Charisma)
    return nil
}