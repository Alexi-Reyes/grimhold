package repl

import (
	"fmt"
	"slices"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/game"
	"github.com/Alexi-Reyes/grimhold/internal/json"
)

func showMap() error {
	regions, err := json.GetRegions()
	if err != nil {
		return fmt.Errorf("failed to load world map: %w", err)
	}

	fmt.Println(Colorize(BoldCyan, "\n=== Known Regions ==="))
	for _, r := range regions {
		marker := "  "
		nameColor := BoldWhite
		if r.Id == game.CurrentRegion {
			marker = Colorize(BoldGreen, "➤ ")
			nameColor = BoldGreen
		}

		fmt.Printf("%s%s (%s)\n", marker, Colorize(nameColor, r.Name), Colorize(Dim, r.Id))
		fmt.Printf("Danger: %s\n", FormatDanger(r.Danger))
		fmt.Printf("Exits:  %s\n", Colorize(Cyan, strings.Join(r.Exits, ", ")))
	}
	fmt.Println()
	return nil
}

func travelTo(dest string) error {
	dest = strings.ToLower(strings.TrimSpace(dest))
	current, err := json.GetRegion(game.CurrentRegion)
	if err != nil {
		return fmt.Errorf("failed to retrieve current location %q: %w", game.CurrentRegion, err)
	}

	if !slices.Contains(current.Exits, dest) {
		fmt.Println(Colorize(Red, fmt.Sprintf("✖ No such exit %q from %s. Type 'look' to see available routes.", dest, current.Name)))
		return nil
	}

	targetRegion, err := json.GetRegion(dest)
	if err != nil {
		return fmt.Errorf("destination region %q does not exist: %w", dest, err)
	}

	game.CurrentRegion = dest
	game.HasRestedInRegion = false

	fmt.Printf("\n%s Travelling to %s...\n", Colorize(Green, "✈"), Colorize(BoldCyan, targetRegion.Name))
	return look()
}

func look() error {
	region, err := json.GetRegion(game.CurrentRegion)
	if err != nil {
		return fmt.Errorf("failed to look at region %q: %w", game.CurrentRegion, err)
	}

	fmt.Println()
	fmt.Printf("%s [%s]\n", Colorize(BoldCyan, region.Name), Colorize(Dim, region.Id))
	fmt.Printf("%s\n", region.Description)
	fmt.Printf("Danger: %s\n", FormatDanger(region.Danger))
	fmt.Printf("Exits:  %s\n", Colorize(Cyan, strings.Join(region.Exits, ", ")))
	return nil
}
