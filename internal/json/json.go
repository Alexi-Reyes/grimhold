package json

import (
	"encoding/json/v2"
	"os"
	"fmt"
    "github.com/Alexi-Reyes/grimhold/internal/game"
)

type Region struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Danger      int      `json:"danger"`
	Foes        []string `json:"foes"`
	Exits       []string `json:"exits"`
}

type World struct {
	Start   string   `json:"start"`
	Regions []Region `json:"regions"`
}

func LoadWorld(path string) (World, error) {
	file, err := os.Open(path)
	if err != nil {
		return World{}, err
	}
	defer file.Close()

	var world World
	if err := json.UnmarshalRead(file, &world); err != nil {
		return World{}, err
	}
	return world, nil
}

const worldPath = "data/world.json"

func GetRegion(input string) (Region, error) {
    world, err := LoadWorld(worldPath)
    if err != nil {
        return Region{}, err
    }

    for _, region := range world.Regions {
        if region.Id == input {
            return region, nil
        }
    }
    return Region{}, fmt.Errorf("region %q not found", input)
}

func GetRegions() ([]Region, error) {
    world, err := LoadWorld(worldPath)
    if err != nil {
        return nil, err
    }

    return world.Regions, nil
}

func LoadStartData() error {
	world, err := LoadWorld(worldPath)
    if err != nil {
        return err
    }

	game.CurrentRegion = world.Start
	return nil
}