package server

import (
	"fmt"
	"hhgame.bojiu.com/enum"
	"sync"
)

var HHGame *HHGameSession

// HHGameSession 场次信息   // 场次 1平民  2小资  3老板  4土豪
type HHGameSession struct {
	Sessions map[int32]*SessionData // int 场次
	SSLock   *sync.RWMutex
}

// SessionData 场次信息
type SessionData struct {
	MaxBet      int64
	MinBet      int64
	MaxBetLucky int64
	MinBetLucky int64
	Desks       []*DeskInfo
}

type DeskInfo struct {
	DeskId       int32                     // 桌号
	Num          string                    // 局数
	Status       int32                     // 当前牌桌的状态   1未下注  2下注中  3开牌中  4结算中
	NextTime     int32                     // 下次状态变化的时间
	PlayerList   map[string]*enum.UserInfo // 玩家信息
	RedAreaBet   int64                     // 红方投注区域的筹码
	BlackAreaBet int64                     // 黑方投注区域的筹码
	LuckyAreaBet int64                     // 幸运一击区域的筹码
	RedCards     []int                     // 红方牌组
	BlackCards   []int                     // 黑方牌组
	RPX          int32                     // 红方牌型
	BPX          int32                     // 黑方牌型
	RedOod       float32                   // 红方赔率
	BlackOod     float32                   // 黑方赔率
	LuckyOod     map[int32]float32         // 幸运一击赔率
	Set          string                    // 局数
	Record       map[string]string         // 记录
}

func NewHHGameSession() *HHGameSession {
	return &HHGameSession{
		Sessions: make(map[int32]*SessionData),
		SSLock:   new(sync.RWMutex),
	}
}

func NewDeskInfo() *DeskInfo {
	lucky := make(map[int32]float32)
	lucky[enum.PairBig] = enum.PairOod
	lucky[enum.Straight] = enum.StraightOod
	lucky[enum.Homogamy] = enum.HomogamyOod
	lucky[enum.Flush] = enum.FlushOod
	lucky[enum.Leopard] = enum.LeopardOod

	return &DeskInfo{
		DeskId:       0,
		Num:          "",
		Status:       0,
		NextTime:     0,
		PlayerList:   make(map[string]*enum.UserInfo),
		RedAreaBet:   0,
		BlackAreaBet: 0,
		LuckyAreaBet: 0,
		RedCards:     nil,
		BlackCards:   nil,
		RPX:          0,
		BPX:          0,
		RedOod:       enum.RedOod,
		BlackOod:     enum.BlackOod,
		LuckyOod:     lucky,
	}

}

func creatDesk() {
	deskInfo := NewDeskInfo()

	data := &SessionData{
		MaxBet:      0,
		MinBet:      0,
		MaxBetLucky: 0,
		MinBetLucky: 0,
		Desks:       []*DeskInfo{deskInfo},
	}

	HHGame.SSLock.Lock()
	HHGame.Sessions[enum.PM] = data
	//HHGame.Sessions[enum.XZ] = data
	//HHGame.Sessions[enum.LB] = data
	//HHGame.Sessions[enum.TH] = data
	HHGame.SSLock.Unlock()
	fmt.Println("创建的场次信息----", &HHGame.Sessions)
}

func StartSession() {
	HHGame = NewHHGameSession()

	// 创建桌子
	go creatDesk()
	// time.Sleep(time.Millisecond * 800)

	go waitTimer()

}
