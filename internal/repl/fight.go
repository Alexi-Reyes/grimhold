package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/Alexi-Reyes/grimhold/internal/game"
	"github.com/Alexi-Reyes/grimhold/internal/open5e"
)

var fightReader = bufio.NewReader(os.Stdin)

func fight(name string) error {
	foe, err := open5e.GetMonster(name)
	if err != nil {
		return fmt.Errorf("failed to fetch foe %q: %w", name, err)
	}
	if foe == nil {
		fmt.Printf("%s Monster %q not found.\n", Colorize(Red, "✖"), name)
		return nil
	}

	action, ok := foe.AttackAction()
	if !ok {
		fmt.Printf("%s %s has no attacks and flees into the shadows.\n", Colorize(Yellow, "⚠"), foe.Name)
		return nil
	}

	foeHP := foe.HitPoints
	p := &game.Player

	fmt.Printf("\n%s %s emerges from the dark! (AC %d, HP %d)\n",
		Colorize(BoldRed, "⚔"), Colorize(BoldWhite, foe.Name), foe.ArmorClass, foeHP)
	fmt.Printf("   You: %s | Attack: %+d | Weapon: %s (%s)\n",
		FormatHPBar(p.HP, p.MaxHP, 10), p.AttackBonus, p.Weapon, p.DamageDice)

	for round := 1; ; round++ {
		fmt.Printf("\n%s (type %s or %s)\n",
			Colorize(BoldCyan, fmt.Sprintf("── Round %d ──", round)),
			Colorize(BoldGreen, "attack"),
			Colorize(Yellow, "flee"))

		choice, err := askChoice()
		if err != nil {
			return fmt.Errorf("failed reading fight choice: %w", err)
		}

		if choice == "flee" {
			roll := game.D20()
			total := roll + p.DexBonus
			if total >= game.FleeDC {
				fmt.Printf("Flee attempt: d20 (%d) %+d = %d vs DC %d → %s\n",
					roll, p.DexBonus, total, game.FleeDC, Colorize(BoldGreen, "SUCCESS"))
				fmt.Println(Colorize(Green, "You successfully slip away back into the area."))
				return nil
			}

			fmt.Printf("Flee attempt: d20 (%d) %+d = %d vs DC %d → %s\n",
				roll, p.DexBonus, total, game.FleeDC, Colorize(BoldRed, "FAILED"))
			fmt.Printf("%s %s gets a free attack as you stumble!\n", Colorize(Yellow, "⚠"), foe.Name)
			if foeStrikes(foe.Name, action, p) {
				return playerFell(p)
			}
			continue
		}

		hit, dmg, err := resolveAttack("You attack:   ", p.AttackBonus, foe.ArmorClass, p.DamageDice, 0)
		if err != nil {
			return fmt.Errorf("attack calculation failed: %w", err)
		}

		if hit {
			foeHP -= dmg
			if foeHP < 0 {
				foeHP = 0
			}
			fmt.Printf("Damage:        %s takes %s damage (HP %d)\n",
				foe.Name, Colorize(BoldGreen, fmt.Sprintf("%d", dmg)), foeHP)
		}

		if foeHP <= 0 {
			fmt.Printf("\n%s %s is slain! %s\n",
				Colorize(BoldGreen, "★"), Colorize(BoldWhite, foe.Name), Colorize(BoldGreen, "You win!"))
			game.RecordSlain(foe.Name)

			xpEarned := foe.HitPoints * 5
			if xpEarned < 25 {
				xpEarned = 25
			}
			levelsGained := p.AddXP(xpEarned)
			fmt.Printf("You gained %s!\n", Colorize(BoldYellow, fmt.Sprintf("%d XP", xpEarned)))
			if levelsGained > 0 {
				fmt.Printf("\n%s You reached %s!\n",
					Colorize(BoldGreen, "✨ LEVEL UP!"),
					Colorize(BoldWhite, fmt.Sprintf("Level %d", p.Level)))
				fmt.Printf("   Max HP increased to %s (+5)\n", Colorize(BoldGreen, fmt.Sprintf("%d", p.MaxHP)))
				fmt.Printf("   Attack Bonus increased to %s (+1)\n", Colorize(BoldCyan, fmt.Sprintf("%+d", p.AttackBonus)))
			}

			goldEarned := foe.HitPoints * 2
			if goldEarned < 10 {
				goldEarned = 10
			}
			p.Gold += goldEarned
			fmt.Printf("You scavenged %s from the corpse.\n", Colorize(BoldYellow, fmt.Sprintf("%d gold", goldEarned)))
			return nil
		}

		if foeStrikes(foe.Name, action, p) {
			return playerFell(p)
		}
	}
}

func askChoice() (string, error) {
	for {
		fmt.Print(Colorize(BoldCyan, "grimhold (fight) > "))
		line, err := fightReader.ReadString('\n')
		if err != nil {
			return "", err
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "attack", "a":
			return "attack", nil
		case "flee", "f":
			return "flee", nil
		default:
			fmt.Println(Colorize(Yellow, "Choose: attack (a) or flee (f)"))
		}
	}
}

func resolveAttack(label string, bonus, ac int, dice string, flatBonus int) (bool, int, error) {
	roll := game.D20()
	total := roll + bonus
	hit := roll == 20 || (roll != 1 && total >= ac)

	outcome := Colorize(Dim, "MISS")
	if roll == 20 {
		outcome = Colorize(BoldGreen, "CRITICAL HIT!")
	} else if hit {
		outcome = Colorize(Green, "HIT")
	}
	fmt.Printf("%s d20 (%d) %+d = %d vs AC %d → %s\n",
		Colorize(BoldWhite, label), roll, bonus, total, ac, outcome)
	if !hit {
		return false, 0, nil
	}

	dmg, err := game.RollDice(dice)
	if err != nil {
		return false, 0, err
	}
	dmg += flatBonus
	if roll == 20 {
		dmg *= 2
	}
	return true, dmg, nil
}

func foeStrikes(name string, a open5e.Action, p *game.Character) bool {
	const playerAC = 12
	hit, dmg, err := resolveAttack(name+":", a.AttackBonus, playerAC, a.DamageDice, a.DamageBonus)
	if err != nil {
		fmt.Println(Colorize(Red, fmt.Sprintf("Attack error: %v", err)))
		return false
	}
	if hit {
		p.HP -= dmg
		if p.HP < 0 {
			p.HP = 0
		}
		fmt.Printf("Damage: You take %s damage! %s\n",
			Colorize(BoldRed, fmt.Sprintf("%d", dmg)), FormatHPBar(p.HP, p.MaxHP, 10))
	}
	return p.HP <= 0
}

func playerFell(p *game.Character) error {
	fmt.Println(Colorize(BoldRed, "\n☠  You have fallen in combat..."))
	p.HP = 1
	game.CurrentRegion = "crossroads"
	game.HasRestedInRegion = false
	fmt.Println(Colorize(Yellow, "You wake up at the crossroads with 1 HP."))
	return nil
}
