package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var CurrentRegion string

type Character struct {
	HP int
	MaxHP int
	AttackBonus int
	DamageDice string
	DexBonus int
}

var Player = Character{
	HP: 20,
	MaxHP: 20,
	AttackBonus: 4,
	DamageDice: "1d8+2",
	DexBonus: 2,
}

const FleeDC = 10

func D20() int {
	return rand.Intn(20) + 1
}

func RollDice(expr string) (int, error) {
	mod := 0
	if i := strings.LastIndexAny(expr, "+-"); i > 0 {
		m, err := strconv.Atoi(expr[i:])
		if err != nil {
			return 0, fmt.Errorf("invalid dice %q", expr)
		}
		mod = m
		expr = expr[:i]
	}

	parts := strings.Split(expr, "d")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid dice %q", expr)
	}
	count := 1
	if parts[0] != "" {
		c, err := strconv.Atoi(parts[0])
		if err != nil || c < 1 {
			return 0, fmt.Errorf("invalid dice %q", expr)
		}
		count = c
	}
	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides < 1 {
		return 0, fmt.Errorf("invalid dice %q", expr)
	}

	total := mod
	for i := 0; i < count; i++ {
		total += rand.Intn(sides) + 1
	}
	if total < 0 {
		total = 0
	}
	return total, nil
}
