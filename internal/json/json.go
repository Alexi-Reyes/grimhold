package json

import (
	"encoding/json/v2"
	"fmt"
	"os"
    "github.com/Alexi-Reyes/grimhold/internal/game"
)

type Region struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Danger int `json:"danger"`
	Foes []string `json:"foes"`
	Exits []string `json:"exits"`
}

type World struct {
	Start string `json:"start"`
	Regions []Region `json:"regions"`
}

func LoadWorld() (World, error) {
	file, err := os.Open("data/world.json")
	if err != nil {
		return World{}, fmt.Errorf("failed to open world file : %w", err)
	}
	defer file.Close()

	var world World
	if err := json.UnmarshalRead(file, &world); err != nil {
		return World{}, fmt.Errorf("failed to decode world json : %w", err)
	}
	return world, nil
}

func GetRegion(input string) (Region, error) {
	world, err := LoadWorld()
	if err != nil {
		return Region{}, fmt.Errorf("cannot get region: %w", err)
	}

    for _, region := range world.Regions {
        if region.Id == input {
            return region, nil
        }
    }
    return Region{}, fmt.Errorf("region %q not found", input)
}

func GetRegions() ([]Region, error) {
	world, err := LoadWorld()
	if err != nil {
		return nil, fmt.Errorf("cannot get regions: %w", err)
	}

	return world.Regions, nil
}

func GetFoes() ([]string, error) {
	region, err := GetRegion(game.CurrentRegion)
	if err != nil {
		return nil, fmt.Errorf("cannot get foes for current region: %w", err)
	}

	return region.Foes, nil
}

func LoadStartData() error {
	world, err := LoadWorld()
	if err != nil {
		return fmt.Errorf("cannot load initial world data: %w", err)
	}

	game.CurrentRegion = world.Start
	return nil
}
