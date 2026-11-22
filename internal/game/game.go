package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var CurrentRegion string
var HasRestedInRegion bool
var SlainBestiary = []string{}

const SaveFile = "save/save.json"
const FleeDC = 10

type Character struct {
	HP int `json:"hp"`
	MaxHP int `json:"max_hp"`
	Level int `json:"level"`
	XP int `json:"xp"`
	Gold int `json:"gold"`
	AttackBonus int `json:"attack_bonus"`
	Weapon string `json:"weapon"`
	DamageDice string `json:"damage_dice"`
	DexBonus int `json:"dex_bonus"`
}

var Player = Character{
	HP: 20,
	MaxHP: 20,
	Level: 1,
	XP: 0,
	Gold: 0,
	AttackBonus: 4,
	Weapon: "Longsword",
	DamageDice: "1d8+2",
	DexBonus: 2,
}

type SaveState struct {
	Location string `json:"location"`
	Character Character `json:"character"`
	SlainBestiary []string `json:"slain_bestiary"`
	Gold int `json:"gold"`
	RestState bool `json:"rest_state"`
}

func RecordSlain(name string) {
	SlainBestiary = append(SlainBestiary, name)
}

func IsSlain(name string) bool {
	for _, slain := range SlainBestiary {
		if strings.EqualFold(slain, name) {
			return true
		}
	}
	return false
}

func SaveGame(path string) error {
	if path == "" {
		path = SaveFile
	}

	slain := SlainBestiary
	if slain == nil {
		slain = []string{}
	}

	state := SaveState{
		Location: CurrentRegion,
		Character: Player,
		SlainBestiary: slain,
		Gold: Player.Gold,
		RestState: HasRestedInRegion,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode save state: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file %q: %w", path, err)
	}
	return nil
}

func LoadGame(path string) error {
	if path == "" {
		path = SaveFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return fmt.Errorf("failed to read save file %q: %w", path, err)
	}

	var state SaveState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to parse save file %q: %w", path, err)
	}

	if state.Location != "" {
		CurrentRegion = state.Location
	}

	Player = state.Character
	if state.Gold != 0 && Player.Gold == 0 {
		Player.Gold = state.Gold
	}

	if state.SlainBestiary != nil {
		SlainBestiary = state.SlainBestiary
	} else {
		SlainBestiary = []string{}
	}

	HasRestedInRegion = state.RestState
	return nil
}

func (character *Character) NextLevelXP() int {
	return character.Level * character.Level * 50
}

func (character *Character) AddXP(amount int) int {
	character.XP += amount
	levelsGained := 0
	for character.XP >= character.NextLevelXP() {
		character.Level++
		character.MaxHP += 5
		character.HP += 5
		character.AttackBonus += 1
		levelsGained++
	}
	return levelsGained
}

func D20() int {
	return rand.Intn(20) + 1
}

func RollDice(expr string) (int, error) {
	expr = strings.ToLower(strings.ReplaceAll(expr, " ", ""))
	mod := 0
	if i := strings.LastIndexAny(expr, "+-"); i > 0 {
		m, err := strconv.Atoi(expr[i:])
		if err != nil {
			return 0, fmt.Errorf("invalid dice modifier %q in expression %q", expr[i:], expr)
		}
		mod = m
		expr = expr[:i]
	}

	parts := strings.Split(expr, "d")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid dice format %q (expected XdY)", expr)
	}
	count := 1
	if parts[0] != "" {
		c, err := strconv.Atoi(parts[0])
		if err != nil || c < 1 {
			return 0, fmt.Errorf("invalid dice count %q in expression", parts[0])
		}
		count = c
	}
	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides < 1 {
		return 0, fmt.Errorf("invalid dice sides %q in expression", parts[1])
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
