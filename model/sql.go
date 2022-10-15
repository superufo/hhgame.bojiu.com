package model

import (
	"common.bojiu.com/models/bj_log"
	"encoding/json"
	"go.uber.org/zap"
	"hhgame.bojiu.com/pkg/log"
	"hhgame.bojiu.com/pkg/mysql"
)

type SqlInsOfUser bj_log.LogUserPerRound
type SqlInsOfGame bj_log.Log2PerRoundHonghei

var tableNameOfUser = "log_user_per_round_"
var tableNameOfGame = "log_2_per_round_honghei"

func InsertGameLog(sqls []*SqlInsOfGame) (errorCode error) {
	for _, v := range sqls {
		in := &SqlInsOfGame{
			PerRoundSid: v.PerRoundSid, //string `xorm:"not null pk VARCHAR(32)"`
			DataTime:    v.DataTime,    //int    `xorm:"default 0 comment('记录时间') INT"`
			UsersData:   v.UsersData,   //string `xorm:"comment('参与的用户:[{sid,输赢值},....]0平') MEDIUMTEXT"`
			Result:      v.Result,      //string `xorm:"default '' comment('开奖结果') VARCHAR(64)"`
			ResultState: v.ResultState, //int    `xorm:"default 0 comment('平台输赢') TINYINT"`
			Amount:      v.Amount,      // int64  `xorm:"default 0 comment('平台输赢值') BIGINT"`
		}
		_, err := mysql.S1().Table(tableNameOfGame).Insert(in)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("sql insert err")
			continue
		}
	}

	return nil
}

// InsertUserLog 插入玩家记录表
func InsertUserLog(sqls []*SqlInsOfUser) (errorCode error) {
	for _, v := range sqls {
		in := &SqlInsOfUser{
			UserSid:             v.UserSid,
			PerRoundSid:         v.PerRoundSid,
			GameId:              v.GameId,
			RoomId:              v.RoomId,
			Change:              v.Change,
			EndTime:             v.EndTime,
			Bets:                v.Bets,
			Result:              v.Result,
			PerRoundState:       v.PerRoundState,
			Win:                 v.Win,
			BeforeMoney:         v.BeforeMoney,
			AfterMoney:          v.AfterMoney,
			Platform:            v.Platform,
			Agent:               v.Agent,
			PlayerServiceCharge: v.PlayerServiceCharge,
		}

		println("-----------------------------------crtb", crtb(v.UserSid))
		_, err := mysql.S1().Table(crtb(v.UserSid)).Insert(in)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("sql insert err")
			continue
		}
	}

	return nil
}

// create 表名
func crtb(sid string) string {
	// sid = "58488346"
	return tableNameOfUser + sid[len(sid)-1:]
}

// Convert ----> string
func Convert(arg interface{}) string {
	data, err := json.Marshal(arg)
	if err != nil {
		log.ZapLog.With(zap.Error(err)).Error("json err")
	}

	return string(data)
}

func ResConvert(black, red int32, bCards, rCards []int) string {
	var tmp = make(map[int32]interface{})
	tmp[black] = bCards
	tmp[red] = rCards
	data, err := json.Marshal(tmp)
	if err != nil {
		log.ZapLog.With(zap.Error(err)).Error("json err")
	}

	return string(data)
}
