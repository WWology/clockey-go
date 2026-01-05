package main

import "time"

type Result struct {
	Result []ResultElement `json:"result"`
}

type ResultElement struct {
	Date            time.Time        `json:"date"`
	Dateexact       int64            `json:"dateexact"`
	Match2ID        string           `json:"match2id"`
	Pagename        string           `json:"pagename"`
	Namespace       int64            `json:"namespace"`
	Match2Opponents []Match2Opponent `json:"match2opponents"`
	Wiki            string           `json:"wiki"`
}

type Match2Opponent struct {
	ID            int64          `json:"id"`
	Type          string         `json:"type"`
	Name          string         `json:"name"`
	Template      string         `json:"template"`
	Icon          string         `json:"icon"`
	Score         int64          `json:"score"`
	Status        string         `json:"status"`
	Placement     int64          `json:"placement"`
	Match2Players []Match2Player `json:"match2players"`
	Extradata     []interface{}  `json:"extradata"`
	Teamtemplate  Teamtemplate   `json:"teamtemplate"`
}

type Match2Player struct {
	ID          int64         `json:"id"`
	Opid        int64         `json:"opid"`
	Name        string        `json:"name"`
	Displayname string        `json:"displayname"`
	Flag        Flag          `json:"flag"`
	Extradata   []interface{} `json:"extradata"`
}

type Teamtemplate struct {
	Template           string `json:"template"`
	Page               string `json:"page"`
	Name               string `json:"name"`
	Shortname          string `json:"shortname"`
	Bracketname        string `json:"bracketname"`
	Image              string `json:"image"`
	Imagedark          string `json:"imagedark"`
	Legacyimage        string `json:"legacyimage"`
	Legacyimagedark    string `json:"legacyimagedark"`
	Imageurl           string `json:"imageurl"`
	Imagedarkurl       string `json:"imagedarkurl"`
	Legacyimageurl     string `json:"legacyimageurl"`
	Legacyimagedarkurl string `json:"legacyimagedarkurl"`
}

type Flag string

const (
	Philippines Flag = "Philippines"
	Russia      Flag = "Russia"
	Ukraine     Flag = "Ukraine"
)
