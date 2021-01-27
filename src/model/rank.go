package model

import "encoding/json"

const keyRankList = "rank_list"

type RankEntry struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Point    float64 `json:"point"`
}

type RankList []RankEntry

func (list RankList) Len() int {
	return len(list)
}

func (list RankList) Swap(i, j int) {
	tmp := list[i]
	list[i] = list[j]
	list[j] = tmp
}

// 排序关键字：得分、用户名、邮箱
func (list RankList) Less(i, j int) bool {
	if list[i].Point == list[j].Point {
		if list[i].Username == list[j].Username {
			return list[i].Email < list[j].Email
		}
		return list[i].Username < list[j].Username
	}
	return list[i].Point > list[j].Point
}

func (m *model) SaveRankList(list RankList) error {
	listStr, err := json.Marshal(list)
	if err != nil {
		return err
	}
	err = redisClient.Set(keyRankList, string(listStr), 0).Err()
	return err
}

func (m *model) GetRankList() (RankList, error) {
	result := redisClient.Get(keyRankList)
	if result.Err() != nil {
		return nil, result.Err()
	}
	list := make([]RankEntry, 0)
	err := json.Unmarshal([]byte(result.Val()), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
