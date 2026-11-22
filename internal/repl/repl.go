package repl

import (
	"fmt"
	"github.com/Alexi-Reyes/grimhold/internal/json"
	"os"
)

func StartGame() error {
	if err := json.LoadStartData(); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	return nil
}

func exitProgram() {
	fmt.Println(Colorize(BoldYellow, "Thanks for playing."))
	os.Exit(0)
}
