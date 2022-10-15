package model

import (
	"strconv"
	"time"
)

func IntToInt32(cards []int) []int32 {
	var res []int32
	for _, v := range cards {
		res = append(res, int32(v))
	}
	return res
}

func AddNums(data map[int32]int64) (r int64) {
	for _, v := range data {
		r += v
	}
	return r
}
func CreatePerRoundId(g, r, d int) string {
	t := time.Now().Unix()
	tstr := strconv.FormatInt(t, 10)
	return strconv.Itoa(g) + strconv.Itoa(r) + strconv.Itoa(d) + tstr
}
