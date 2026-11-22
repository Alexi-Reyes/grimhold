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
		return err
	}
	if foe == nil {
		fmt.Println("No such monster.")
		return nil
	}

	action, ok := foe.AttackAction()
	if !ok {
		fmt.Printf("%v has no attack and wants no part of this.\n", foe.Name)
		return nil
	}

	foeHP := foe.HitPoints
	p := &game.Player

	fmt.Printf("\n⚔  %v attacks!  (AC %d, HP %d)\n", foe.Name, foe.ArmorClass, foeHP)
	fmt.Printf("   You: HP %d/%d, attack %+d, weapon %s (%s)\n", p.HP, p.MaxHP, p.AttackBonus, p.Weapon, p.DamageDice)

	for round := 1; ; round++ {
		fmt.Printf("\n── Round %d ── (attack / flee)\n", round)
		choice, err := askChoice()
		if err != nil {
			return err
		}

		if choice == "flee" {
			roll := game.D20()
			total := roll + p.DexBonus
			result := "FAILED"
			if total >= game.FleeDC {
				result = "SUCCESS"
			}
			fmt.Printf("Flee attempt:  d20 (%d) %+d (DEX) = %d vs DC %d → %s\n", roll, p.DexBonus, total, game.FleeDC, result)
			if total >= game.FleeDC {
				fmt.Println("You slip away.")
				return nil
			}
			fmt.Printf("%v gets a free strike!\n", foe.Name)
			if foeStrikes(foe.Name, action, p) {
				return playerFell(p)
			}
			continue
		}

		hit, dmg, err := resolveAttack("You attack: ", p.AttackBonus, foe.ArmorClass, p.DamageDice, 0)
		if err != nil {
			return err
		}

		if hit {
			foeHP -= dmg
			if foeHP < 0 {
				foeHP = 0
			}
			fmt.Printf("Damage: %d → %v takes %d (HP %d)\n", dmg, foe.Name, dmg, foeHP)
		}

		if foeHP <= 0 {
			fmt.Printf("%v is slain! You win.\n", foe.Name)
			game.RecordSlain(foe.Name)
			xpEarned := foe.HitPoints * 5

			if xpEarned < 25 {
				xpEarned = 25
			}

			p.AddXP(xpEarned)

			goldEarned := foe.HitPoints * 2
			if goldEarned < 10 {
				goldEarned = 10
			}
			p.Gold += goldEarned
			fmt.Printf("You found %d gold!\n", goldEarned)
			return nil
		}

		if foeStrikes(foe.Name, action, p) {
			return playerFell(p)
		}
	}
}

func askChoice() (string, error) {
	for {
		fmt.Print("grimhold (fight) > ")
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
			fmt.Println("Type attack or flee.")
		}
	}
}

func resolveAttack(label string, bonus, ac int, dice string, flatBonus int) (bool, int, error) {
	roll := game.D20()
	total := roll + bonus
	hit := roll == 20 || (roll != 1 && total >= ac)

	outcome := "MISS"
	if roll == 20 {
		outcome = "CRITICAL HIT"
	} else if hit {
		outcome = "HIT"
	}
	fmt.Printf("%s d20 (%d) %+d = %d vs AC %d → %s\n", label, roll, bonus, total, ac, outcome)
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
		fmt.Println(err)
		return false
	}
	if hit {
		p.HP -= dmg
		if p.HP < 0 {
			p.HP = 0
		}
		fmt.Printf("Damage: %d → You take %d (HP %d)\n", dmg, dmg, p.HP)
	}
	return p.HP <= 0
}

func playerFell(p *game.Character) error {
	fmt.Println("\n☠  You have fallen...")
	p.HP = p.MaxHP
	fmt.Println("You wake up, bruised but alive (HP restored).")
	return nil
}
