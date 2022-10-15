package server

import (
	"hhgame.bojiu.com/enum"
	"hhgame.bojiu.com/model"
	"runtime"
	"time"
)

// waitTimer 未下注  5S
func waitTimer() {
	HHGame.SSLock.Lock()
	sessionData := HHGame.Sessions
	for sessionId, session := range sessionData {
		for deskId, _ := range session.Desks {
			HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Bet_Before
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = 5
		}
	}

	HHGame.SSLock.Unlock()
	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			ss := HHGame.Sessions
			for sessionId, session := range ss {
				for deskId, _ := range session.Desks {
					HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(5 - count)
				}
			}
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 5 {
				HHGame.SSLock.Lock()
				ss := HHGame.Sessions
				for sessionId, session := range ss {
					for deskId, deskInfo := range session.Desks {
						go sendToClientChangeStatus(sessionId, enum.Betting, 10, int32(deskId), deskInfo.PlayerList)
					}
				}
				HHGame.SSLock.Unlock()

				ticker.Stop()
				go bettingTimer()
				runtime.Goexit()
			}
		}
	}()
}

// bettingTimer 下注中  10S
func bettingTimer() {
	HHGame.SSLock.Lock()
	sessionData := HHGame.Sessions
	for sessionId, session := range sessionData {
		for deskId, _ := range session.Desks {
			rC, bC := model.LicensingFunc()
			Rpx := model.Check(rC)
			Bpx := model.Check(bC)
			HHGame.Sessions[sessionId].Desks[deskId].RedCards = rC
			HHGame.Sessions[sessionId].Desks[deskId].BlackCards = bC
			HHGame.Sessions[sessionId].Desks[deskId].RPX = Rpx
			HHGame.Sessions[sessionId].Desks[deskId].BPX = Bpx
			HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Betting
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = 10
		}
	}

	HHGame.SSLock.Unlock()

	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			ss := HHGame.Sessions
			for sessionId, session := range ss {
				for deskId, _ := range session.Desks {
					HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(10 - count)
				}
			}
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 10 {
				HHGame.SSLock.Lock()
				ss := HHGame.Sessions
				for sessionId, session := range ss {
					for deskId, deskInfo := range session.Desks {
						go sendToClientRes(sessionId, int32(deskId), deskInfo.PlayerList)
						time.Sleep(time.Millisecond * 100)
						go sendToClientChangeStatus(sessionId, enum.Show_Cards, 8, int32(deskId), deskInfo.PlayerList)
					}
				}
				HHGame.SSLock.Unlock()
				ticker.Stop()
				go openResTimer()
				runtime.Goexit()
			}
		}
	}()
}

// openResTimer 开牌中  5S
func openResTimer() {
	HHGame.SSLock.Lock()
	sessionData := HHGame.Sessions
	for sessionId, session := range sessionData {
		for deskId, _ := range session.Desks {
			HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Show_Cards
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = 8
		}
	}

	HHGame.SSLock.Unlock()
	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			ss := HHGame.Sessions
			for sessionId, session := range ss {
				for deskId, _ := range session.Desks {
					HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(8 - count)
				}
			}
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 8 {
				// 清除桌子数据
				HHGame.SSLock.Lock()
				sss := HHGame.Sessions
				for sessionId, session := range sss {
					for deskId, deskInfo := range session.Desks {
						// 先发消息
						go sendToClientChangeStatus(sessionId, enum.Bet_Before, 5, int32(deskId), deskInfo.PlayerList)
						HHGame.Sessions[sessionId].Desks[deskId].RedAreaBet = 0
						HHGame.Sessions[sessionId].Desks[deskId].BlackAreaBet = 0
						HHGame.Sessions[sessionId].Desks[deskId].LuckyAreaBet = 0
						for _, v := range HHGame.Sessions[sessionId].Desks[deskId].PlayerList {
							v.MyBet = nil
						}
					}
				}
				HHGame.SSLock.Unlock()

				ticker.Stop()
				go waitTimer()
				runtime.Goexit()
			}
		}
	}()
}

//// --------------------------------------------------------other version----------------------------------------------------
// waitTimer 未下注  5S
/*
func waitTimer(sessionId, deskId int32) {
	HHGame.SSLock.Lock()
	HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Bet_Before
	HHGame.Sessions[sessionId].Desks[deskId].NextTime = 5
	HHGame.SSLock.Unlock()

	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(5 - count)
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 5 {
				go bettingTimer(sessionId, deskId)
				ticker.Stop()
				runtime.Goexit()
			}
		}
	}()
}

// bettingTimer 下注中  10S
func bettingTimer(sessionId, deskId int32) {
	HHGame.SSLock.Lock()
	HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Betting
	HHGame.Sessions[sessionId].Desks[deskId].NextTime = 10
	HHGame.SSLock.Unlock()

	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(10 - count)
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 10 {
				// push.SendToClientChangeStatus(sessionId, enum.Bet_Before, enum.Betting, deskId)
				go openResTimer(sessionId, deskId)
				ticker.Stop()
				runtime.Goexit()
			}
		}
	}()
}

// openResTimer 开牌中  3S
func openResTimer(sessionId, deskId int32) {
	HHGame.SSLock.Lock()
	// rC, bC := model.LicensingFunc()
	// Rpx := model.Check(rC)
	// Bpx := model.Check(bC)
	// HHGame.Sessions[sessionId].Desks[deskId].RedCards = rC
	// HHGame.Sessions[sessionId].Desks[deskId].BlackCards = bC
	// HHGame.Sessions[sessionId].Desks[deskId].RPX = int32(Rpx)
	// HHGame.Sessions[sessionId].Desks[deskId].BPX = int32(Bpx)
	HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Show_Cards
	HHGame.Sessions[sessionId].Desks[deskId].NextTime = 3
	HHGame.SSLock.Unlock()

	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(3 - count)
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 3 {
				go settleTimer(sessionId, deskId)
				ticker.Stop()
				runtime.Goexit()
			}
		}
	}()
}

// settleTimer 结算中   5S
func settleTimer(sessionId, deskId int32) {
	HHGame.SSLock.Lock()
	HHGame.Sessions[sessionId].Desks[deskId].Status = enum.Settlement
	HHGame.Sessions[sessionId].Desks[deskId].NextTime = 5
	HHGame.SSLock.Unlock()

	ticker := time.NewTicker(time.Second * 1)
	count := 0
	go func() {
		for {
			<-ticker.C
			count++
			HHGame.SSLock.Lock()
			HHGame.Sessions[sessionId].Desks[deskId].NextTime = int32(5 - count)
			HHGame.SSLock.Unlock()

			// 触发条件
			if count == 5 {
				// 清除桌子数据
				HHGame.SSLock.Lock()
				HHGame.Sessions[sessionId].Desks[deskId].RedAreaBet = 0
				HHGame.Sessions[sessionId].Desks[deskId].BlackAreaBet = 0
				HHGame.Sessions[sessionId].Desks[deskId].LuckyAreaBet = 0
				HHGame.SSLock.Unlock()

				go waitTimer(sessionId, deskId)
				ticker.Stop()
				runtime.Goexit()
			}
		}
	}()
}
*/
