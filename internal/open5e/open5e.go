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
    Strength int `json:"strength"`
    Dexterity int `json:"dexterity"`
    Constitution int `json:"constitution"`
    Intelligence int `json:"intelligence"`
    Wisdom int `json:"wisdom"`
    Charisma int `json:"charisma"`
}

type MonsterResponse struct {
    Results []Monster `json:"results"`
}

func GetMonsters() ([]byte, error) {
    resp, err := http.Get(apiUrl + "/monsters/?format=json")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
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
