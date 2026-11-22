package repl

import (
	"fmt"
	"math/rand"
	"github.com/Alexi-Reyes/grimhold/internal/game"
	"github.com/Alexi-Reyes/grimhold/internal/json"
	"github.com/Alexi-Reyes/grimhold/internal/open5e"
)

func explore() error {
	region, err := json.GetRegion(game.CurrentRegion)
	if err != nil {
		return fmt.Errorf("explore failed: %w", err)
	}

	if len(region.Foes) == 0 {
		fmt.Println(Colorize(Green, "The area is peaceful. No foes in sight."))
		return nil
	}

	rolls := region.Danger
	if rolls < 1 {
		rolls = 1
	}

	var chosenFoe string
	var worstHP = -1

	for i := 0; i < rolls; i++ {
		candidate := region.Foes[rand.Intn(len(region.Foes))]
		m, err := open5e.GetMonster(candidate)
		if err != nil || m == nil {
			if chosenFoe == "" {
				chosenFoe = candidate
			}
			continue
		}

		if m.HitPoints > worstHP {
			worstHP = m.HitPoints
			chosenFoe = candidate
		}
	}

	if chosenFoe == "" {
		chosenFoe = region.Foes[rand.Intn(len(region.Foes))]
	}

	return fight(chosenFoe)
}
