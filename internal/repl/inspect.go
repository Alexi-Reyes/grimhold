package repl

import (
	"fmt"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/open5e"
)

func inspect(name string) error {
	name = strings.TrimSpace(name)
	result, err := open5e.GetMonster(name)
	if err != nil {
		return fmt.Errorf("failed to inspect monster %q: %w", name, err)
	}

	if result == nil {
		fmt.Println(Colorize(Red, fmt.Sprintf("✖ No monster found matching %q.", name)))
		return nil
	}

	fmt.Println()
	fmt.Printf("%s %s\n", Colorize(BoldRed, "⚔"), Colorize(BoldWhite, result.Name))
	if result.Size != "" || result.Type != "" {
		fmt.Printf("   %s\n", Colorize(Dim, fmt.Sprintf("%s %s", result.Size, result.Type)))
	}
	fmt.Printf("Armor Class: %s\n", Colorize(BoldCyan, fmt.Sprintf("%d", result.ArmorClass)))
	fmt.Printf("Hit Points: %s\n", Colorize(BoldGreen, fmt.Sprintf("%d", result.HitPoints)))
	if result.ChallengeRating != "" {
		fmt.Printf("Challenge Rating: %s\n", Colorize(BoldYellow, result.ChallengeRating))
	}

	fmt.Println(Colorize(Dim, "   ─────────────────────────────"))
	fmt.Printf("   %s %2d  %s %2d  %s %2d\n",
		Colorize(Cyan, "STR:"), result.Strength,
		Colorize(Cyan, "DEX:"), result.Dexterity,
		Colorize(Cyan, "CON:"), result.Constitution)
	fmt.Printf("   %s %2d  %s %2d  %s %2d\n",
		Colorize(Cyan, "INT:"), result.Intelligence,
		Colorize(Cyan, "WIS:"), result.Wisdom,
		Colorize(Cyan, "CHA:"), result.Charisma)

	if action, ok := result.AttackAction(); ok {
		fmt.Println(Colorize(Dim, "   ─────────────────────────────"))
		fmt.Printf("Action: %s (+%d to hit, %s dmg)\n",
			Colorize(BoldWhite, action.Name), action.AttackBonus, Colorize(Yellow, action.DamageDice))
	}
	fmt.Println()
	return nil
}
