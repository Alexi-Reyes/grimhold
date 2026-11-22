package repl

import (
    "fmt"
    "os"
    "slices"
    "math/rand"
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
            Callback: func(args ...string) error { explore(); return nil },
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
        "fight": {
            Name: "fight",
            Desc: "Fight a monster in rounds (usage: fight <monster>)",
            Callback: func(args ...string) error {
                if len(args) == 0 {
                    fmt.Println("Please enter a monster. (usage: fight <monster>)")
                    return nil
                }

                if len(args) != 1 {
                    fmt.Println("More than one argument. (usage: fight <monster>)")
                    return nil
                }

                return fight(args[0])
            },
        },
        "character": {
            Name: "character",
            Desc: "Display your character sheet",
            Callback: func(args ...string) error { showCharacter(); return nil },
        },
        "rest": {
            Name: "rest",
            Desc: "Rest to restore HP (once per region visit)",
            Callback: func(args ...string) error { rest(); return nil },
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
        game.HasRestedInRegion = false
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

func explore() error {
    foes, err := json.GetFoes()
    if err != nil {
        return err
    }

    foe := foes[rand.Intn(len(foes))]

    fmt.Printf("%v appears\n", foe)
    fmt.Println("")
    return inspect(foe)
}

const adventurerArt = `
      /\
     (o.o)
   <--|#|---==>
     / \
    _] [_`

func showCharacter() {
    player := &game.Player
    fmt.Println(adventurerArt)
    fmt.Println("--- Character Sheet ---")
    fmt.Printf("Level: %d\n", player.Level)
    fmt.Printf("XP: %d / %d\n", player.XP, player.NextLevelXP())
    fmt.Printf("HP: %d / %d\n", player.HP, player.MaxHP)
    fmt.Printf("Attack Bonus: %+d\n", player.AttackBonus)
    fmt.Printf("Weapon: %s (%s)\n", player.Weapon, player.DamageDice)
}

func rest() {
    if game.HasRestedInRegion {
        fmt.Println("You have already rested in this region. You must travel elsewhere before resting again.")
        return
    }

    if game.Player.HP >= game.Player.MaxHP {
        fmt.Println("You are already at full health.")
        return
    }

    game.Player.HP = game.Player.MaxHP
    game.HasRestedInRegion = true
    fmt.Printf("You take a long rest and recover your strength. HP restored to %d/%d.\n", game.Player.HP, game.Player.MaxHP)
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