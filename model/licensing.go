package model

import (
	"hhgame.bojiu.com/enum"
	"math/rand"
	"time"
)

// LicensingFunc 发牌：根据概率来
func LicensingFunc() (redCards, blackCards []int) {
	redCards, blackCards = RandMath()
	return

}

// Cards 牌组
var (
	Cards = []int{101, 201, 301, 401, 501, 601, 701, 801, 901, 1001, 1101, 1201, 1301, // 黑桃
		102, 202, 302, 402, 502, 602, 702, 802, 902, 1002, 1102, 1202, 1302, // 红桃
		103, 203, 303, 403, 503, 603, 703, 803, 903, 1003, 1103, 1203, 1303, // 梅花
		104, 204, 304, 404, 504, 604, 704, 804, 904, 1004, 1104, 1204, 1304, // 方块
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandMath 正常发牌，随机取六个数
func RandMath() ([]int, []int) {
	// 拷贝
	AllPai := make([]int, 0, 0)
	tmpCards := make([]int, len(Cards))
	copy(tmpCards, Cards)
	tmp := Cards
	FaPaiCiShu := 6               // 初始化牌型
	SuiJiMap := make(map[int]int) // 记录随机数
	for i := 0; i < FaPaiCiShu; i++ {
		WeiZhi := r.Intn(52)
		_, ok := SuiJiMap[WeiZhi]
		if ok {
			FaPaiCiShu++
			continue
		}

		SuiJiMap[WeiZhi] = WeiZhi
		AllPai = append(AllPai, tmp[WeiZhi])

	}

	redResult := AllPai[0:3]
	blackResult := AllPai[3:6]

	return redResult, blackResult
}

// ChooseWhoWin 某方必胜的发牌
func ChooseWhoWin(str string) (redCards, blackCards []int) {

	// 发牌！红方必胜
	if str == "red" {
		for {
			redCards, blackCards = RandMath()
			if CheckWhoWin(redCards, blackCards) == enum.Red_Win {
				return
			}
		}
	}
	// 发牌！黑方必胜
	if str == "black" {
		for {
			redCards, blackCards = RandMath()
			if CheckWhoWin(redCards, blackCards) == enum.Black_Win {
				return
			}
		}
	}

	return nil, nil
}
