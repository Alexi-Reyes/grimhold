package repl

import (
	"errors"
	"fmt"
	"os"
	"github.com/Alexi-Reyes/grimhold/internal/game"
)

func saveGameCommand() {
	if err := game.SaveGame(game.SaveFile); err != nil {
		fmt.Println(Colorize(Red, fmt.Sprintf("✖ Error saving game: %v", err)))
		return
	}

	fmt.Printf("%s Game saved successfully to %s.\n", Colorize(Green, "✔"), Colorize(Cyan, game.SaveFile))
}

func loadGameCommand() {
	if err := game.LoadGame(game.SaveFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(Colorize(Yellow, "No save file found."))
			return
		}
		fmt.Println(Colorize(Red, fmt.Sprintf("✖ Error loading save: %v", err)))
		return
	}
	fmt.Printf("%s Loaded save from %s. Resuming in %s.\n",
		Colorize(Green, "✔"), Colorize(Cyan, game.SaveFile), Colorize(BoldCyan, game.CurrentRegion))
	_ = look()
}
