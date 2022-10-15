package msg

// 所有消息类型  消息为两位
const (
	CMD_ERROR uint16 = 1001

	CMD_NOTICE_MSG            uint16 = 1004
	CMD_USER                  uint16 = 1005
	CMD_USER_INFO             uint16 = 1006
	CMD_GAME_CONFIG           uint16 = 2001
	CMD_ENTER_GAME            uint16 = 2003
	CMD_LEAVE_GAME            uint16 = 2004
	CMD_GAME_1_DESK           uint16 = 3002
	CMD_GAME_1_BETS           uint16 = 3005
	CMD_GAME_1_CHANGE_STATE   uint16 = 3006
	CMD_GAME_1_END_RESULT     uint16 = 3007
	CMD_GAME_1_START_NEW_BETS uint16 = 3008
	CMD_GAME_1_END_AWARD      uint16 = 3009
	CMD_HEART_BIT             uint16 = 3001 // 心跳信息
	CMD_STEP                  uint16 = 3003 // 下发步长
	CMD_LOGIN                 uint16 = 3004
	CMD_SERVER_LIST           uint16 = 3016 // 下发服务信息

	CMDS string = "1001,1004,1005,1006,2001,2003,2004,3002,3005,3006,3007,3008,3009,3001,3003,3004,3007,3016"
)
