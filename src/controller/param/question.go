package param

type ObjRspQuestion struct {
	ID        string   `json:"_id"`
	Title     string   `json:"title"`
	NextScene []string `json:"next_scene"`
	Status    int      `json:"status"`
}

type RspQuestionGetList struct {
	Question []ObjRspQuestion `json:"question"`
}

type ObjRspSubQuestion struct {
	Type       int      `json:"type"`
	Desc       string   `json:"desc"`
	Option     []string `json:"option"`
	FullPoint  float32  `json:"full_point"`
	PartPoint  float32  `json:"part_point"`
}

type ObjRspNextScene struct {
	Scene  string `json:"scene"`
	Option string `json:"option"`
}

type RspQuestionGetInfo struct {
	Title       string              `json:"title"`
	Desc        string              `json:"desc"`
	SubQuestion []ObjRspSubQuestion `json:"sub_question"`
	Author      string              `json:"author"`
	Audio       string              `json:"audio"`
	TimeLimit   int                 `json:"time_limit"`
	NextScene   []ObjRspNextScene   `json:"next_scene"`
	Status      int                 `json:"status"`
}
