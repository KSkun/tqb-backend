package param

type ObjRspScene struct {
	ID           string `json:"_id"`
	FromQuestion string `json:"from_question"`
	NextQuestion string `json:"next_question"`
	Title        string `json:"title"`
}

type RspSceneGetList struct {
	Scene []ObjRspScene `json:"scene"`
}

type RspSceneGetInfo struct {
	Title        string `json:"title"`
	Text         string `json:"text"`
	FromQuestion string `json:"from_question"`
	NextQuestion string `json:"next_question"`
	BGM          string `json:"bgm"`
}
