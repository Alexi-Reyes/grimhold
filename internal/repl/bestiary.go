package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/game"
	"github.com/Alexi-Reyes/grimhold/internal/open5e"
)

var bestiaryReader = bufio.NewReader(os.Stdin)

func runBestiary() error {
	page, err := open5e.GetMonsterPage("")
	if err != nil {
		fmt.Printf("Failed to load bestiary: %v\n", err)
		return nil
	}

	displayBestiaryPage(page)

	for {
		fmt.Print("grimhold (bestiary) > ")
		line, err := bestiaryReader.ReadString('\n')
		if err != nil {
			break
		}

		cmd := strings.ToLower(strings.TrimSpace(line))
		switch cmd {
		case "next":
			if page.Next == nil || *page.Next == "" {
				fmt.Println("You are on the last page.")
				continue
			}

			fmt.Println("Loading next page...")
			nextPage, err := open5e.GetMonsterPage(*page.Next)
			if err != nil {
				fmt.Printf("Error loading page: %v\n", err)
				continue
			}

			page = nextPage
			displayBestiaryPage(page)

		case "previous":
			if page.Previous == nil || *page.Previous == "" {
				fmt.Println("You are on the first page.")
				continue
			}

			fmt.Println("Loading previous page...")
			prevPage, err := open5e.GetMonsterPage(*page.Previous)
			if err != nil {
				fmt.Printf("Error loading page: %v\n", err)
				continue
			}

			page = prevPage
			displayBestiaryPage(page)

		case "back":
			fmt.Println("Returning to game.")
			return nil

		case "help", "h":
			fmt.Println("Bestiary commands:")
			fmt.Println("next (n) - View the next page")
			fmt.Println("prev (p) - View the previous page")
			fmt.Println("back (q) - Exit bestiary")

		case "":
			continue

		default:
			fmt.Println("Unknown bestiary command. Options: next, prev, back.")
		}
	}
	return nil
}

func displayBestiaryPage(page *open5e.MonsterList) {
	fmt.Println("\n=== Bestiary ===")
	if len(page.Results) == 0 {
		fmt.Println("No monsters found.")
		return
	}

	for _, monster := range page.Results {
		slainStatus := "[UNSLAIN]"
		if game.IsSlain(monster.Name) {
			slainStatus = "[SLAIN]  "
		}

		var details []string
		if monster.ChallengeRating != "" {
			details = append(details, "CR "+m.ChallengeRating)
		}

		if monster.Type != "" {
			details = append(details, monster.Type)
		}

		if monster.HitPoints > 0 {
			details = append(details, fmt.Sprintf("HP %d", monster.HitPoints))
		}

		if monster.ArmorClass > 0 {
			details = append(details, fmt.Sprintf("AC %d", monster.ArmorClass))
		}

		detailStr := ""
		if len(details) > 0 {
			detailStr = " (" + strings.Join(details, ", ") + ")"
		}

		fmt.Printf(" %s %s%s\n", slainStatus, monster.Name, detailStr)
	}

	nav := []string{}
	if page.Previous != nil && *page.Previous != "" {
		nav = append(nav, "prev")
	}

	if page.Next != nil && *page.Next != "" {
		nav = append(nav, "next")
	}

	nav = append(nav, "back")

	fmt.Printf("\n[Commands: %s]\n", strings.Join(nav, " | "))
}

