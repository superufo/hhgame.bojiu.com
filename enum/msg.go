package enum

const (
	CMD_ENTER_2_GAME         uint16 = 2101 // 进入游戏,回复场次信息
	CMD_GAME_2_SESSION       uint16 = 2102 // 选择场次,回复该场次内信息
	CMD_Game_2_Eeter_Desk    uint16 = 2103 // 玩家进入桌子
	CMD_GAME_2_BETS          uint16 = 2104 // 玩家投注,返回投注的信息
	CMD_GAME_2_END_RESULT    uint16 = 2105 // 游戏结束的结果
	CMD_GAME_2_OUT           uint16 = 2106 // 推出游戏
	CMD_GAME_2_STATUS_CHANGE uint16 = 2107 // 状态改变
	CMD_GAME_2_GET_RECORD    uint16 = 2108 // 获取游戏记录

	CMDS string = "2101,2102,2103,2104,2105,2106,2107,2108"
)
