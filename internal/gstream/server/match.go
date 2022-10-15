package server

// Limits coin限制
var Limits *map[int32]NumLimit

type NumLimit struct {
	MinBet      int32
	MaxBet      int32
	MinBetLucky int32
	MaxBetLucky int32
}

// 2  服务器先创建桌子 ---> 玩家选择场次过后直接进入桌子
// TODO V2.0   每个场次开四张桌子
