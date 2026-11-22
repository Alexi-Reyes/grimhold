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
	fmt.Println(Colorize(Dim, "\nConsulting the bestiary..."))
	page, err := open5e.GetMonsterPage("")
	if err != nil {
		fmt.Printf("%s Failed to load bestiary: %v\n", Colorize(Red, "✖"), err)
		return nil
	}

	displayBestiaryPage(page)

	for {
		fmt.Print(Colorize(BoldCyan, "grimhold (bestiary) > "))
		line, err := bestiaryReader.ReadString('\n')
		if err != nil {
			break
		}

		cmd := strings.ToLower(strings.TrimSpace(line))
		switch cmd {
		case "next":
			if page.Next == nil || *page.Next == "" {
				fmt.Println(Colorize(Yellow, "You are already on the last page."))
				continue
			}
			fmt.Println(Colorize(Dim, "Turning page forward..."))
			nextPage, err := open5e.GetMonsterPage(*page.Next)
			if err != nil {
				fmt.Printf("%s Error loading page: %v\n", Colorize(Red, "✖"), err)
				continue
			}

			page = nextPage
			displayBestiaryPage(page)

		case "previous":
			if page.Previous == nil || *page.Previous == "" {
				fmt.Println(Colorize(Yellow, "You are already on the first page."))
				continue
			}
			fmt.Println(Colorize(Dim, "Turning page backward..."))
			prevPage, err := open5e.GetMonsterPage(*page.Previous)
			if err != nil {
				fmt.Printf("%s Error loading page: %v\n", Colorize(Red, "✖"), err)
				continue
			}

			page = prevPage
			displayBestiaryPage(page)

		case "back":
			fmt.Println("Returning to game.")
			return nil

		case "help", "h":
			fmt.Println(Colorize(BoldCyan, "Bestiary Controls:"))
			fmt.Println(Colorize(White, "  next (n) - Next page"))
			fmt.Println(Colorize(White, "  prev (p) - Previous page"))
			fmt.Println(Colorize(White, "  back (q) - Exit bestiary"))

		case "":
			continue

		default:
			fmt.Println(Colorize(Yellow, "Commands: next (n), prev (p), back (q)"))
		}
	}
	return nil
}

func displayBestiaryPage(page *open5e.MonsterList) {
	fmt.Println(Colorize(BoldCyan, "\n=== Bestiary of Foes ==="))
	if len(page.Results) == 0 {
		fmt.Println(Colorize(Yellow, "No monsters recorded."))
		return
	}

	for _, monster := range page.Results {
		var status string
		if game.IsSlain(monster.Name) {
			status = Colorize(BoldGreen, "[SLAIN]  ")
		} else {
			status = Colorize(Dim, "[UNSLAIN]")
		}

		var details []string
		if monster.ChallengeRating != "" {
			details = append(details, "CR "+monster.ChallengeRating)
		}

		if monster.Type != "" {
			details = append(details, monster.Type)
		}

		if monster.HitPoints > 0 {
			details = append(details, fmt.Sprintf("%d HP", monster.HitPoints))
		}

		if monster.ArmorClass > 0 {
			details = append(details, fmt.Sprintf("%d AC", monster.ArmorClass))
		}

		detailStr := ""
		if len(details) > 0 {
			detailStr = Colorize(Dim, " ("+strings.Join(details, ", ")+")")
		}

		fmt.Printf(" %s %s%s\n", status, Colorize(BoldWhite, monster.Name), detailStr)
	}

	var nav []string
	if page.Previous != nil && *page.Previous != "" {
		nav = append(nav, Colorize(Yellow, "prev"))
	}

	if page.Next != nil && *page.Next != "" {
		nav = append(nav, Colorize(Yellow, "next"))
	}
	nav = append(nav, Colorize(BoldCyan, "back"))

	fmt.Printf("\n[ %s ]\n\n", strings.Join(nav, " | "))
}
