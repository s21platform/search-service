package model

type UserInfo struct {
	Nickname   string   `json:"nickname"`
	Avatar     string   `json:"avatar"`
	Name       string   `json:"name"`
	Surname    string   `json:"surname"`
	Birthdate  string   `json:"birthdate"`
	Phone      string   `json:"phone"`
	City       string   `json:"city"`
	Telegram   string   `json:"telegram"`
	Git        string   `json:"git"`
	Os         GetOs    `json:"os"`
	Work       string   `json:"work"`
	University string   `json:"university"`
	Skills     []string `json:"skills"`
	Hobbies    []string `json:"hobbies"`
}

type GetOs struct {
	Id    int64  `json:"id"`
	Label string `json:"label"`
}
