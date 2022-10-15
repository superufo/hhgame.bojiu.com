package entity

type GamesInfo struct {
	ID            int64 `gorm:"column:id" json:"id"`
	CurrentStock1 int64 `gorm:"column:current_stock_1" json:"current_stock_1"` // 当前水位
	CurrentStock2 int64 `gorm:"column:current_stock_2" json:"current_stock_2"` // 当前库存2水位
	ChangeTime    int64 `gorm:"column:change_time" json:"change_time"`         // 当前水位变化的时间
	UpdataTime    int64 `gorm:"column:updata_time" json:"updata_time"`         // 写入数据库的时间
}

const TableGamesInfo = "games_info"
