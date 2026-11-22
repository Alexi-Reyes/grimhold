package game

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var CurrentRegion string
var HasRestedInRegion bool

type Character struct {
	HP int
	MaxHP int
	Level int
	XP int
	AttackBonus int
	Weapon string
	DamageDice string
	DexBonus int
}

var Player = Character{
	HP: 20,
	MaxHP: 20,
	Level: 1,
	XP: 0,
	AttackBonus: 4,
	Weapon: "Longsword",
	DamageDice: "1d8+2",
	DexBonus: 2,
}

const FleeDC = 10

func (character *Character) NextLevelXP() int {
	return character.Level * character.Level * 50
}

func (character *Character) AddXP(amount int) {
	character.XP += amount
	fmt.Printf("You gained %d XP!\n", amount)
	for character.XP >= character.NextLevelXP() {
		character.Level++
		character.MaxHP += 5
		character.HP += 5
		character.AttackBonus += 1
		fmt.Printf("\nLevel Up! You reached Level %d!\n", character.Level)
		fmt.Printf("Max HP increased to %d (+5)\n", character.MaxHP)
		fmt.Printf("Attack Bonus increased to %+d (+1)\n", character.AttackBonus)
	}
}

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
