package entity

type GamesConfig struct {
	ID                  int64   `gorm:"column:id" json:"id"`                                         //
	Name                string  `gorm:"column:name" json:"name"`                                     //
	Stock1              int64   `gorm:"column:stock_1" json:"stock_1"`                               //库存1基础库存
	Stock1WarnWater     int64   `gorm:"column:stock_1_warn_water" json:"stock_1_warn_water"`         //库存1警戒值
	DrawWater           int64   `gorm:"column:draw_water" json:"draw_water"`                         //库存1目标和报警之间的抽水概几率万分比
	PlayerServiceCharge float64 `gorm:"column:player_service_charge" json:"player_service_charge"`   //玩家赢的玩家服务费
	SystemServiceCharge float64 `gorm:"column:system_service_charge" json:"system_service_charge"`   //玩家赢的系统服务费(暂时没用到)
	Stock2ServiceCharge float64 `gorm:"column:stock_2_service_charge" json:"stock_2_service_charge"` //库存2奖励库存比例
	Stock2WarnWater     int64   `gorm:"column:stock_2_warn_water" json:"stock_2_warn_water"`         //库存2报警值
	Stock1State         int64   `gorm:"column:stock_1_state" json:"stock_1_state"`                   //库存1状态
	UpdateTime          int64   `gorm:"column:update_time" json:"update_time"`                       //变化时间
	ToStock1            float64 `gorm:"column:to_stock_1" json:"to_stock_1"`                         //玩家输了,钱进库存的比例
}

const TableGamesConfig = "games_config"
