package repl

import (
	"fmt"

	"github.com/Alexi-Reyes/grimhold/internal/game"
)

const adventurerArt = `
      /\
     (o.o)
   <--|#|---==>
     / \
    _] [_`

func showCharacter() {
	player := &game.Player
	fmt.Println(Colorize(BoldCyan, adventurerArt))
	fmt.Println(Colorize(BoldCyan, "=== Character Sheet ==="))
	fmt.Printf("Level: %s\n", Colorize(BoldWhite, fmt.Sprintf("%d", player.Level)))
	fmt.Printf("XP: %s\n", Colorize(Yellow, fmt.Sprintf("%d / %d", player.XP, player.NextLevelXP())))
	fmt.Printf("Gold: %s\n", Colorize(BoldYellow, fmt.Sprintf("%d gp", player.Gold)))
	fmt.Printf("Health: %s\n", FormatHPBar(player.HP, player.MaxHP, 12))
	fmt.Printf("Attack Bonus: %s\n", Colorize(BoldWhite, fmt.Sprintf("%+d", player.AttackBonus)))
	fmt.Printf("Weapon: %s (%s)\n", Colorize(BoldWhite, player.Weapon), Colorize(Cyan, player.DamageDice))
	fmt.Println()
}

func rest() {
	player := &game.Player
	if game.HasRestedInRegion {
		fmt.Println(Colorize(Yellow, "You have already rested in this region. Travel elsewhere before resting again."))
		return
	}

	if player.HP >= player.MaxHP {
		fmt.Println(Colorize(Green, "You are already at full health."))
		return
	}

	player.HP = player.MaxHP
	game.HasRestedInRegion = true
	fmt.Printf("%s You set camp and take a long rest. HP fully restored: %s\n",
		Colorize(BoldGreen, "⛺"), FormatHPBar(player.HP, player.MaxHP, 12))
}
