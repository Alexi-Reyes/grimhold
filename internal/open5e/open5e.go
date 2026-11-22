package open5e

import (
    "encoding/json"
    "io"
    "net/http"
    "net/url"
)

const (
    apiUrl = "https://api.open5e.com/v1"
)

type Monster struct {
    Name string `json:"name"`
    Size string `json:"size"`
    Type string `json:"type"`
    ChallengeRating string `json:"challenge_rating"`
    Strength int `json:"strength"`
    Dexterity int `json:"dexterity"`
    Constitution int `json:"constitution"`
    Intelligence int `json:"intelligence"`
    Wisdom int `json:"wisdom"`
    Charisma int `json:"charisma"`
    ArmorClass int `json:"armor_class"`
    HitPoints int `json:"hit_points"`
    Actions json.RawMessage `json:"actions"`
}

type Action struct {
    Name string `json:"name"`
    AttackBonus int `json:"attack_bonus"`
    DamageDice string `json:"damage_dice"`
    DamageBonus int `json:"damage_bonus"`
}

func (m *Monster) AttackAction() (Action, bool) {
    var actions []Action
    if err := json.Unmarshal(m.Actions, &actions); err != nil {
        return Action{}, false
    }
    for _, a := range actions {
        if a.DamageDice != "" {
            return a, true
        }
    }
    return Action{}, false
}

type MonsterList struct {
    Count int `json:"count"`
    Next *string `json:"next"`
    Previous *string `json:"previous"`
    Results []Monster `json:"results"`
}

type MonsterResponse struct {
    Results []Monster `json:"results"`
}

func GetMonsterPage(urlStr string) (*MonsterList, error) {
    if urlStr == "" {
        urlStr = apiUrl + "/monsters/?format=json"
    }

    resp, err := http.Get(urlStr)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var list MonsterList
    if err := json.Unmarshal(body, &list); err != nil {
        return nil, err
    }

    return &list, nil
}

func GetMonsters() (*MonsterList, error) {
    return GetMonsterPage("")
}

func GetMonster(monster string) (*Monster, error) {
    resp, err := http.Get(apiUrl + "/monsters/?format=json&name__iexact=" + url.QueryEscape(monster))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var parsed MonsterResponse
    if err := json.Unmarshal(body, &parsed); err != nil {
        return nil, err
    }
    if len(parsed.Results) == 0 {
        return nil, nil
    }

    return &parsed.Results[0], nil
}
