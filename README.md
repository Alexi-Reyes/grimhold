# grimhold

## Description

A text-based dungeon crawler that runs in your terminal. You explore cursed regions, stumble into random encounters, fight foes in dice-based combat, and record your kills in a bestiary. 

## World map

Generated from `data/world.json`. `<-->` is a two-way exit, `-->` a one-way exit (start region in bold).

```mermaid
flowchart LR
    crossroads["Crossroads<br/>danger 0"]
    ashen_marsh["Ashen Marsh<br/>danger 1"]
    drowned_chapel["Drowned Chapel<br/>danger 2"]
    bleak_forest["Bleak Forest<br/>danger 2"]
    iron_keep["Iron Keep<br/>danger 3"]
    bone_catacombs["Bone Catacombs<br/>danger 4"]
    rotting_village["Rotting Village<br/>danger 1"]
    hag_hollow["Hag's Hollow<br/>danger 2"]
    frozen_pass["Frozen Pass<br/>danger 3"]
    obsidian_spire["Obsidian Spire<br/>danger 5"]
    crossroads <--> ashen_marsh
    crossroads <--> bleak_forest
    crossroads <--> iron_keep
    crossroads <--> rotting_village
    ashen_marsh <--> drowned_chapel
    drowned_chapel <--> bone_catacombs
    bleak_forest <--> iron_keep
    bleak_forest <--> hag_hollow
    iron_keep <--> bone_catacombs
    iron_keep <--> frozen_pass
    bone_catacombs <--> obsidian_spire
    rotting_village <--> hag_hollow
    frozen_pass <--> obsidian_spire
    style crossroads stroke-width:4px
```
