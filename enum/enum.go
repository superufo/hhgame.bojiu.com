package enum

const (
	Red_Win   = 2 // 红方胜
	Black_Win = 0 // 黑方胜
)

const (
	Leopard   int32 = 100 // 豹子
	Flush     int32 = 99  // 同花顺
	Homogamy  int32 = 98  // 同花
	Straight  int32 = 97  // 顺子
	PairBig   int32 = 96  // 大对  --> 大于等于对9
	PairSmall int32 = 95  // 小对  --> 小于对9
	Leaflet   int32 = 94  // 单张
)

const (
	Bet_Before = 1 // 下注前
	Betting    = 2 // 下注中
	Show_Cards = 3 // 开牌中
)

const (
	PM = 1 // 平民场
	XZ = 2 // 小资场
	LB = 3 // 老板场
	TH = 4 // 土豪场
)

const (
	Only_Red_Win   = 1 //  只让红方胜利的桌子号
	Only_Red_Black = 2 // 只让黑方胜利的桌子号
)

// 赔率
const (
	RedOod      float32 = 0.95
	BlackOod    float32 = 0.95
	PairOod     float32 = 1
	StraightOod float32 = 2
	HomogamyOod float32 = 2
	FlushOod    float32 = 11
	LeopardOod  float32 = 24
)

// 限红
const ()
