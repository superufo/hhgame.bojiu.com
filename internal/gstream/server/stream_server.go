package server

import (
	"common.bojiu.com/def"
	Utils "common.bojiu.com/utils"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hhgame.bojiu.com/config"
	// c "hhgame.bojiu.com/create"
	"hhgame.bojiu.com/enum"
	"hhgame.bojiu.com/internal/gstream/pb"
	cproto "hhgame.bojiu.com/internal/gstream/proto"
	"hhgame.bojiu.com/model"
	"hhgame.bojiu.com/pkg/log"
	protoStructure "hhgame.bojiu.com/proto"
	"io"
	"net"
	"strings"
	"time"
)

// var center  = storage.StorageServerImpl
var stream *streamServer

type streamServer struct {
	GrpcRecvClientData chan *pb.StreamRequestData
	GrpcSendClientData chan *pb.StreamResponseData
}

func NewStreamServer() *streamServer {
	stream = &streamServer{
		make(chan *pb.StreamRequestData, 100),
		make(chan *pb.StreamResponseData, 100),
	}
	//
	return stream
}

//func init() {
//	GrpcRecvClientData = make(chan *pb.StreamRequestData, 100)
//	GrpcSendClientData = make(chan *pb.StreamResponseData, 100)
//}

// PPStream log.ZapLog.With(zap.Any("err", err)).Error("收到网关数据错误")
func (gs *streamServer) PPStream(stream pb.ForwardMsg_PPStreamServer) error {
	stop := make(chan struct{})
	defer func() {
		if e := recover(); e != nil {
			log.ZapLog.Info("PPStream recover", zap.Any("err", e.(error)))
		}
		close(stop)
	}()
	go gs.response(stream, stop)
	go gs.dispatch(stop)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.ZapLog.Info("PPStream recv io EOF", zap.Any("err", err))
			return nil
		}
		if err != nil {
			log.ZapLog.Info("PPStream recv error", zap.Any("err", err))
			return err
		}
		info := fmt.Sprintf("收到网关数据:协议号=%+v,加密字符=%+v,随机字符=%+v,protobuf=%+v", msg.GetMsg(), Utils.ToHexString(msg.GetSecret()), msg.GetSerialNum(), msg.GetData())
		log.ZapLog.Info(info)
		gs.GrpcRecvClientData <- msg
	}
}

func (gs *streamServer) response(stream pb.ForwardMsg_PPStreamServer, stop chan struct{}) {
	defer func() {
		if e := recover(); e != nil {
			log.ZapLog.Info("stream response", zap.Any("err", e.(error)))
		}
	}()
	for {
		select {
		case sd := <-gs.GrpcSendClientData:
			//业务代码
			if err := stream.Send(sd); err != nil {
				log.ZapLog.Info("", zap.Any("发给网关失败err", err))
			} else {
				log.ZapLog.Info("", zap.Any("发给网关成功msg", sd.String()))
			}
		case <-stop:
			return
		}
	}
}

func (gs *streamServer) dispatch(stop chan struct{}) {
	defer func() {
		if e := recover(); e != nil {
			log.ZapLog.Info("stream dispatch", zap.Any("err", e.(error)))
		}
	}()
	for {
		select {
		case cmsg := <-gs.GrpcRecvClientData:
			{
				log.ZapLog.Info("dispatch", zap.Any("Msg", cmsg.Msg))
				if !strings.Contains(enum.CMDS, fmt.Sprintf("%d", cmsg.Msg)) {
					log.ZapLog.Error("不存在的消息", zap.Any("msg", cmsg.Msg))
				}
				if err := gs.handlerMsg(cmsg); err != nil {
					log.ZapLog.Info("handlerMsg error ", zap.Any("Msg", cmsg.Msg), zap.Any("err", err))
				}
			}
		case <-stop:
			return
		}
	}
}

// 消息处理
func (gs *streamServer) handlerMsg(clientMsg *pb.StreamRequestData) error {
	// 进入游戏，处理
	if uint16(clientMsg.Msg) == enum.CMD_ENTER_2_GAME {
		gs.enterGame(clientMsg)
	}

	// 选择场次，回复场次信息
	if uint16(clientMsg.Msg) == enum.CMD_GAME_2_SESSION {
		gs.chooseSession(clientMsg)
	}

	// 选择桌子，回复桌子信息
	if uint16(clientMsg.Msg) == enum.CMD_Game_2_Eeter_Desk {
		gs.chooseDesk(clientMsg)
	}

	// 玩家下注，返回信息
	if uint16(clientMsg.Msg) == enum.CMD_GAME_2_BETS {
		gs.playerBet(clientMsg)
	}

	// 玩家退出游戏
	if uint16(clientMsg.Msg) == enum.CMD_GAME_2_OUT {
		gs.playerOut(clientMsg)
	}

	// 玩家请求游戏记录
	if uint16(clientMsg.Msg) == enum.CMD_GAME_2_GET_RECORD {
		gs.getGameRecord(clientMsg)
	}

	return nil
}

// sendToClientRes 告知游戏结果
func sendToClientRes(sessionId, deskId int32, playerList map[string]*enum.UserInfo) {
	// 获取牌组
	HHGame.SSLock.Lock()
	deskInfo := HHGame.Sessions[sessionId].Desks[deskId]
	pl := deskInfo.PlayerList
	HHGame.SSLock.Unlock()

	rCards := model.IntToInt32(deskInfo.RedCards)                    // 红方牌组
	bCards := model.IntToInt32(deskInfo.BlackCards)                  // 黑方牌组
	win := model.CheckWhoWin(deskInfo.RedCards, deskInfo.BlackCards) // 哪方赢
	var gameid uint32 = 2
	var w int32 = enum.Red_Win
	var l int32 = enum.Black_Win
	var b = true
	var sIds []string
	for _, v := range playerList {
		sIds = append(sIds, v.SId)
	}

	// 红色方赢
	if win == enum.Red_Win {
		var allRes []*protoStructure.PPlayerSettlementInfo
		var bet []*protoStructure.PBetInfo
		var sqlIns []*model.SqlInsOfUser
		var sqls []*model.SqlInsOfGame
		for _, v := range sIds {
			for _, playerInfo := range pl {
				if playerInfo.MyBet == nil {
					continue
				}
				tmpv := v

				if v == playerInfo.SId {
					for area, value := range playerInfo.MyBet {
						ta := area
						tv := value

						betInfo := protoStructure.PBetInfo{
							BetArea: &ta,
							Nums:    &tv,
						}

						bet = append(bet, &betInfo)
					}

					// 计算输赢
					nums, addnums := model.SettlementR(deskInfo.RPX, playerInfo.MyBet, deskInfo.RedOod, deskInfo.LuckyOod)
					fmt.Println("-------------------------------------------------->>>>>>>>red win", nums)
					// ------------------------------获取玩家数据-----------------------------------------
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
					client := GoClient()
					request := cproto.GameSettlementTos{GameType: 2, Score: nums}
					_, err := client.SendSettlement(ctx, &request)
					if err != nil {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result of SendSettlement")
					}
					t := uint32(deskId)
					requestAdd := cproto.ChangeBalanceReq{Uid: v, Gold: addnums, ChangeType: uint32(def.CHANGE_WINLOSE), PerRoundSid: &deskInfo.Set, GameId: &gameid, RoomId: &t, SerialNo: nil}
					res, err := client.AddBalance(ctx, &requestAdd)
					if err != nil {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result of ChangeBalanceReq")
					}

					if res.Code != 0 {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result err,res code not zero!")
					}

					// ---------------------------------------------------------------------------------
					HHGame.SSLock.Lock()
					HHGame.Sessions[sessionId].Desks[deskId].PlayerList[v].Coin += nums
					HHGame.SSLock.Unlock()

					allRes = append(allRes, &protoStructure.PPlayerSettlementInfo{
						SId:        &tmpv,
						Win:        &b,
						ChangeGold: &nums,
						BetInfo:    bet,
					})

					// 插入数据库
					sqlIns = append(sqlIns, &model.SqlInsOfUser{
						UserSid:             playerInfo.SId,
						PerRoundSid:         "",
						GameId:              int(gameid),
						RoomId:              int(deskId),
						Change:              nums,
						EndTime:             int(time.Now().Unix()),
						Bets:                model.Convert(playerInfo.MyBet),
						Result:              model.ResConvert(int32(0), int32(2), deskInfo.BlackCards, deskInfo.RedCards),
						PerRoundState:       0,
						Win:                 nums - model.AddNums(playerInfo.MyBet),
						BeforeMoney:         playerInfo.Coin,
						AfterMoney:          playerInfo.Coin + nums,
						Platform:            playerInfo.Platform,
						Agent:               playerInfo.Agent,
						PlayerServiceCharge: 0,
					})

					// 游戏数据
					sqls = append(sqls, &model.SqlInsOfGame{
						PerRoundSid: model.CreatePerRoundId(int(gameid), int(sessionId), int(deskId)),
						RoomId:      int(sessionId),
						DataTime:    int(time.Now().Unix()),
						UsersData:   "",
						Result:      model.ResConvert(int32(0), int32(2), deskInfo.BlackCards, deskInfo.RedCards),
						ResultState: enum.Red_Win,
					})
				}
			}
		}

		err := model.InsertUserLog(sqlIns)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("settlement insert err")
		}

		err = model.InsertGameLog(sqls)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("settlement insert err")
		}

		// 加减金币

		res := &protoStructure.MGameResultHongHei_2Toc{
			Win:            &w,
			RedPx:          &deskInfo.RPX,
			RedCards:       rCards,
			BlackPx:        &deskInfo.BPX,
			BlackCards:     bCards,
			SettlementInfo: allRes, // 结算信息
		}

		data, _ := proto.Marshal(res)
		sendCMsg := pb.StreamResponseData{
			ClientId: "",
			BAllUser: false,
			Uids:     sIds,
			Msg:      uint32(enum.CMD_GAME_2_END_RESULT),
			Data:     data,
		}
		stream.GrpcSendClientData <- &sendCMsg
	}

	if win == enum.Black_Win {
		var allRes []*protoStructure.PPlayerSettlementInfo
		// var res *protoStructure.PPlayerSettlementInfo
		var bet []*protoStructure.PBetInfo
		var sqlIns []*model.SqlInsOfUser
		var sqls []*model.SqlInsOfGame
		for _, v := range sIds {
			for _, playerInfo := range pl {
				if playerInfo.MyBet == nil {
					continue
				}
				tmpv := v

				if v == playerInfo.SId {
					for area, value := range playerInfo.MyBet {
						ta := area
						tv := value

						betInfo := protoStructure.PBetInfo{
							BetArea: &ta,
							Nums:    &tv,
						}

						bet = append(bet, &betInfo)
					}

					// 计算输赢
					nums, addnums := model.SettlementB(deskInfo.BPX, playerInfo.MyBet, deskInfo.BlackOod, deskInfo.LuckyOod)
					fmt.Println("-------------------------------------------------->>>>>>>>black win", nums, addnums)

					// ------------------------------获取玩家数据-----------------------------------------
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
					client := GoClient()
					request := cproto.GameSettlementTos{GameType: 2, Score: nums}
					_, err := client.SendSettlement(ctx, &request)

					if err != nil {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result")
					}
					t := uint32(deskId)
					requestAdd := cproto.ChangeBalanceReq{Uid: v, Gold: addnums, ChangeType: uint32(def.CHANGE_WINLOSE), PerRoundSid: &deskInfo.Set, GameId: &gameid, RoomId: &t, SerialNo: nil}
					res, err := client.AddBalance(ctx, &requestAdd)
					if err != nil {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result of ChangeBalanceReq")
					}

					if res.Code != 0 {
						log.ZapLog.With(zap.Error(err)).Error("grpc dial result err,res code not zero!")
					}
					// ---------------------------------------------------------------------------------
					HHGame.SSLock.Lock()
					HHGame.Sessions[sessionId].Desks[deskId].PlayerList[v].Coin += nums
					HHGame.SSLock.Unlock()

					allRes = append(allRes, &protoStructure.PPlayerSettlementInfo{
						SId:        &tmpv,
						Win:        &b,
						ChangeGold: &nums,
						BetInfo:    bet,
					})

					// 插入数据库
					sqlIns = append(sqlIns, &model.SqlInsOfUser{
						UserSid:             playerInfo.SId,
						PerRoundSid:         "",
						GameId:              int(gameid),
						RoomId:              int(deskId),
						Change:              nums,
						EndTime:             int(time.Now().Unix()),
						Bets:                model.Convert(playerInfo.MyBet),
						Result:              model.ResConvert(int32(0), int32(2), deskInfo.BlackCards, deskInfo.RedCards),
						PerRoundState:       0,
						Win:                 nums - model.AddNums(playerInfo.MyBet),
						BeforeMoney:         playerInfo.Coin,
						AfterMoney:          playerInfo.Coin + nums,
						Platform:            playerInfo.Platform,
						Agent:               playerInfo.Agent,
						PlayerServiceCharge: 0,
					})

					// 游戏数据
					sqls = append(sqls, &model.SqlInsOfGame{
						PerRoundSid: model.CreatePerRoundId(int(gameid), int(sessionId), int(deskId)),
						RoomId:      int(sessionId),
						DataTime:    int(time.Now().Unix()),
						UsersData:   "",
						Result:      model.ResConvert(int32(0), int32(2), deskInfo.BlackCards, deskInfo.RedCards),
						ResultState: enum.Red_Win,
					})
				}
			}
		}

		err := model.InsertUserLog(sqlIns)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("settlement insert err")
		}

		err = model.InsertGameLog(sqls)
		if err != nil {
			log.ZapLog.With(zap.Error(err)).Error("settlement insert err")
		}

		res := &protoStructure.MGameResultHongHei_2Toc{
			Win:            &l,
			BlackPx:        &deskInfo.BPX,
			BlackCards:     bCards,
			RedPx:          &deskInfo.RPX,
			RedCards:       rCards,
			SettlementInfo: allRes, // 结算信息
		}

		data, _ := proto.Marshal(res)
		sendCMsg := pb.StreamResponseData{
			ClientId: "",
			BAllUser: false,
			Uids:     sIds,
			Msg:      uint32(enum.CMD_GAME_2_END_RESULT),
			Data:     data,
		}
		stream.GrpcSendClientData <- &sendCMsg
	}
}

// sendToClientChangeStatus  告知客户端状态改变
func sendToClientChangeStatus(sessionId, nextStatus, time, deskId int32, playerList map[string]*enum.UserInfo) {
	var robotBets []*protoStructure.PRobotBet
	// 给机器人添加筹码
	if nextStatus == enum.Betting {
		var r int32 = 0
		var b int32 = 2
		var l int32 = 1
		var n int64 = 500000
		robotBets = append(robotBets, &protoStructure.PRobotBet{
			Area: &r,
			Nums: &n,
		}, &protoStructure.PRobotBet{
			Area: &b,
			Nums: &n,
		}, &protoStructure.PRobotBet{
			Area: &l,
			Nums: &n,
		})
		HHGame.SSLock.Lock()
		HHGame.Sessions[sessionId].Desks[deskId].RedAreaBet = n
		HHGame.Sessions[sessionId].Desks[deskId].BlackAreaBet = n
		HHGame.Sessions[sessionId].Desks[deskId].LuckyAreaBet = n
		HHGame.SSLock.Unlock()
	}

	statusChange := &protoStructure.MDeskStatusChangeToc{
		DeskId:     &deskId,
		NextStatus: &nextStatus,
		Time:       &time,
		Rbet:       robotBets,
	}

	var sIds []string
	for _, v := range playerList {
		sIds = append(sIds, v.SId)
	}

	fmt.Println("--------------------------------------------------------这里的uid是", sIds)
	data, _ := proto.Marshal(statusChange)
	sendCMsg := pb.StreamResponseData{
		ClientId: "",
		BAllUser: false,
		Uids:     sIds,
		Msg:      uint32(enum.CMD_GAME_2_STATUS_CHANGE),
		Data:     data,
	}

	stream.GrpcSendClientData <- &sendCMsg

}

// ENTER_GAME 进入游戏
func (gs *streamServer) enterGame(msg *pb.StreamRequestData) error {
	var l = protoStructure.MGame_2EnterGameTos{}

	if err := proto.Unmarshal(msg.Data, &l); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("SendGameSessionInfo")
		return errors.New("proto3解码错误")
	}

	fmt.Println("-------------------------------------", l.GetGameId())
	// check game id
	if l.GetGameId() != int32(2) {
		log.ZapLog.With(zap.Any("err", errors.New("GameId not match"))).Info("SendGameSessionInfo")
		return errors.New("GameId not match")
	}
	time.Now().Second()
	var sessions []*protoStructure.PGameHongHei_2Sessions
	var session *protoStructure.PGameHongHei_2Sessions
	for _, v := range HHGame.Sessions {
		session = &protoStructure.PGameHongHei_2Sessions{
			MinBet:      &v.MinBet,
			MaxBet:      &v.MaxBet,
			MinBetLucky: &v.MinBetLucky,
			MaxBetLucky: &v.MaxBetLucky,
		}
		sessions = append(sessions, session)
	}

	hhGameInfo := protoStructure.MGame_2EnterGameToc{
		GameId:   l.GameId,
		Room:     nil,
		Desk:     nil,
		Sessions: sessions, // 场次信息 场次 1 平民 2 小资 3 老板 4 土豪
	}

	data, _ := proto.Marshal(&hhGameInfo)
	sendCMsg := pb.StreamResponseData{
		ClientId: msg.GetClientId(),
		BAllUser: false,
		Uids:     nil,
		Msg:      uint32(enum.CMD_ENTER_2_GAME),
		Data:     data,
	}

	gs.GrpcSendClientData <- &sendCMsg
	return nil
}

// CHOOSE_SESSION 选择场次
func (gs *streamServer) chooseSession(msg *pb.StreamRequestData) error {
	var l = protoStructure.MPlayerIntoGameHongHei_2Tos{}

	if err := proto.Unmarshal(msg.Data, &l); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("SendGameSessionInfo")
		return errors.New("proto3解码错误")
	}

	var desks []*protoStructure.PAllDeskHongHei_2Info
	var desk *protoStructure.PAllDeskHongHei_2Info
	for _, v := range HHGame.Sessions[l.GetSessions()].Desks {
		desk = &protoStructure.PAllDeskHongHei_2Info{
			DeskId:   &v.DeskId,
			Status:   &v.Status,
			NextTime: &v.NextTime,
			Set:      nil,
		}

		desks = append(desks, desk)
	}

	allDeskInfo := protoStructure.MAllDeskHongHei_2InfoToc{
		SId:      l.SId,
		Nickname: nil,
		Coin:     nil,
		Desk:     desks,
	}

	data, _ := proto.Marshal(&allDeskInfo)
	sendCMsg := pb.StreamResponseData{
		ClientId: msg.GetClientId(),
		BAllUser: false,
		Uids:     nil,
		Msg:      uint32(enum.CMD_GAME_2_SESSION),
		Data:     data,
	}

	gs.GrpcSendClientData <- &sendCMsg

	return nil
}

// CHOOSE_DESK  选择桌子
func (gs *streamServer) chooseDesk(msg *pb.StreamRequestData) error {
	var l = protoStructure.MPlayerIntoDeskHongHei_2Tos{}
	if err := proto.Unmarshal(msg.Data, &l); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("SendGameSessionInfo")
		return errors.New("proto3解码错误")
	}

	// ----------------------------保存玩家数据-------------------
	PLock.Lock()
	PlayerAll[l.GetSId()] = &Data{
		SessionId: l.GetSessionId(),
		DeskId:    l.GetDeskId(),
	}
	PLock.Unlock()
	// ----------------------------------------------------------

	// ------------------------------获取玩家数据-----------------------------------------
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client := GoClient()
	request := cproto.UserRequest{Uid: *l.SId}
	res, err := client.GetUserInfo(ctx, &request)
	if err != nil {
		log.ZapLog.With(zap.Error(err)).Error("grpc dial result")
	}

	// ---------------------------------------------------------------------------------
	playerInfo := &enum.UserInfo{
		SId:      res.GetUser().SId,
		Name:     *res.GetUser().Name,
		Sex:      *res.GetUser().Sex,
		Nickname: *res.GetUser().Nickname,
		Platform: *res.GetUser().Platform,
		Agent:    *res.GetUser().Agent,
		Coin:     *res.GetUserInfo().Gold,
		MyBet:    nil,
	}

	fmt.Println("-----------------------------playerINfo", playerInfo)

	HHGame.SSLock.Lock()
	HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].PlayerList[l.GetSId()] = playerInfo
	HHGame.SSLock.Unlock()

	// 牌桌信息
	var deskInfo *protoStructure.MDeskInfoToc

	// 投注区域信息
	var allBetArea []*protoStructure.PDeskBetHongHei_2Area

	var r int64 = 0
	var b int64 = 2
	var ll int64 = 1

	var ro []*protoStructure.PDeskOodLuckyHongHei_2
	var bo []*protoStructure.PDeskOodLuckyHongHei_2
	var lo []*protoStructure.PDeskOodLuckyHongHei_2

	ro = append(ro, &protoStructure.PDeskOodLuckyHongHei_2{
		Px:  nil,
		Ood: &HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].RedOod,
	})
	bo = append(bo, &protoStructure.PDeskOodLuckyHongHei_2{
		Px:  nil,
		Ood: &HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].BlackOod,
	})

	for k, v := range HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].LuckyOod {
		tmpk := k
		tmpv := v
		tmpp := protoStructure.PDeskOodLuckyHongHei_2{
			Px:  &tmpk,
			Ood: &tmpv,
		}
		lo = append(lo, &tmpp)
	}

	allBetArea = append(allBetArea, &protoStructure.PDeskBetHongHei_2Area{
		Area:    &r,
		AreaBet: &HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].RedAreaBet,
		Ood:     ro,
	}, &protoStructure.PDeskBetHongHei_2Area{
		Area:    &b,
		AreaBet: &HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].BlackAreaBet,
		Ood:     bo,
	}, &protoStructure.PDeskBetHongHei_2Area{
		Area:    &ll,
		AreaBet: &HHGame.Sessions[l.GetSessionId()].Desks[l.GetDeskId()].LuckyAreaBet,
		Ood:     lo,
	})

	for _, v := range HHGame.Sessions[l.GetSessionId()].Desks {
		deskInfo = &protoStructure.MDeskInfoToc{
			SId:      l.SId,
			Coin:     res.GetUserInfo().Gold,
			Status:   &v.Status,
			NextTime: &v.NextTime,
			Set:      nil,
			History:  nil,
			BetArea:  allBetArea,
			MyBet:    nil,
		}
	}

	data, _ := proto.Marshal(deskInfo)
	sendCMsg := pb.StreamResponseData{
		ClientId: msg.GetClientId(),
		BAllUser: false,
		Uids:     nil,
		Msg:      uint32(enum.CMD_Game_2_Eeter_Desk),
		Data:     data,
	}

	gs.GrpcSendClientData <- &sendCMsg
	return nil
}

// PLAYER_BET   玩家下注
func (gs *streamServer) playerBet(msg *pb.StreamRequestData) error {
	var tmpG int64
	// 获取msg
	var l = protoStructure.MPlayerBetHongHei_2Tos{}
	if err := proto.Unmarshal(msg.Data, &l); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("SendGameSessionInfo")
		return errors.New("proto3解码错误")
	}

	var playerBets []*protoStructure.PPlayerBetHongHei_2
	PLock.Lock()
	u := PlayerAll[l.GetSId()]
	PLock.Unlock()

	// 保存内存玩家下注数据
	HHGame.SSLock.Lock()
	pl, ok := HHGame.Sessions[u.SessionId].Desks[u.DeskId].PlayerList[l.GetSId()]
	if ok {
		x := HHGame.Sessions[u.SessionId].Desks[u.DeskId].Set // 局数
		t := uint32(u.DeskId)                                 // 桌子Id
		g := uint32(2)                                        // gameId
		// ------------------------------获取玩家数据-----------------------------------------
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		client := GoClient()
		request := cproto.ChangeBalanceReq{Uid: *l.SId, Gold: l.GetNums(), ChangeType: uint32(def.CHANGE_BET), PerRoundSid: &x, GameId: &g, RoomId: &t, SerialNo: nil}
		_, err := client.ReduceBalance(ctx, &request)
		if err != nil {
			HHGame.SSLock.Unlock()
			log.ZapLog.With(zap.Error(err)).Error("grpc dial result")
			betInfo := protoStructure.MPlayerBetHongHei_2Toc{
				MyChip: &tmpG,
			}

			// 发送msg
			data, _ := proto.Marshal(&betInfo)
			sendCMsg := pb.StreamResponseData{
				ClientId: msg.GetClientId(),
				BAllUser: false,
				Uids:     nil,
				Msg:      uint32(enum.CMD_GAME_2_BETS),
				Data:     data,
			}
			gs.GrpcSendClientData <- &sendCMsg
			return nil
		}
		// ---------------------------------------------------------------------------------

		pl.Coin = pl.Coin - l.GetNums() // 减去玩家内存金币数量
		tmpG = pl.Coin

		if pl.MyBet != nil {
			pl.MyBet[l.GetArea()] = pl.MyBet[l.GetArea()] + l.GetNums()
		}

		if pl.MyBet == nil {
			pl.MyBet = make(map[int32]int64)
			pl.MyBet[l.GetArea()] = l.GetNums()

		}
	}

	for area, value := range pl.MyBet {
		ta := area
		tv := value

		v := protoStructure.PPlayerBetHongHei_2{
			Area:      &ta,
			MyAllChip: &tv,
		}

		playerBets = append(playerBets, &v)
	}

	HHGame.SSLock.Unlock()

	betInfo := protoStructure.MPlayerBetHongHei_2Toc{
		AreaNow:   l.Area,
		ChipNow:   l.Nums,
		MyChip:    &tmpG,
		PlayerBet: playerBets,
	}

	// 发送msg
	data, _ := proto.Marshal(&betInfo)
	sendCMsg := pb.StreamResponseData{
		ClientId: msg.GetClientId(),
		BAllUser: false,
		Uids:     nil,
		Msg:      uint32(enum.CMD_GAME_2_BETS),
		Data:     data,
	}

	gs.GrpcSendClientData <- &sendCMsg
	return nil
}

// 玩家退出游戏
func (gs *streamServer) playerOut(msg *pb.StreamRequestData) error {
	var l protoStructure.MPlayerOutHongHei_2Tos
	if err := proto.Unmarshal(msg.Data, &l); err != nil {
		log.ZapLog.With(zap.Any("err", err)).Info("SendGameSessionInfo")
		return errors.New("proto3解码错误")
	}
	PLock.Lock()
	player := PlayerAll[l.GetSId()]
	PLock.Unlock()

	HHGame.SSLock.Lock()
	playerList := HHGame.Sessions[player.SessionId].Desks[player.DeskId].PlayerList
	delete(playerList, l.GetSId())
	HHGame.Sessions[player.SessionId].Desks[player.DeskId].PlayerList = playerList
	HHGame.SSLock.Unlock()

	PLock.Lock()
	delete(PlayerAll, l.GetSId())
	PLock.Unlock()

	return nil
}

// 玩家请求游戏记录
func (gs *streamServer) getGameRecord(msg *pb.StreamRequestData) error {
	return nil
}

func Run() {
	//streamIp := viper.Vp.GetString("ser.stream.ip")
	//streamPort := viper.Vp.GetInt("ser.stream.port")
	var server pb.ForwardMsgServer
	sImpl := NewStreamServer()

	server = sImpl

	g := grpc.NewServer()

	// 2.注册逻辑到server中
	pb.RegisterForwardMsgServer(g, server)

	scfg := config.NewServerCfg()
	instance := fmt.Sprintf("%s:%d", scfg.GetIp(), scfg.GetPort())

	log.ZapLog.With(zap.Any("addr", instance)).Info("Run")
	// 3.启动server
	lis, err := net.Listen("tcp", instance)
	if err != nil {
		panic("监听错误:" + err.Error())
	}

	err = g.Serve(lis)
	if err != nil {
		panic("启动错误:" + err.Error())
	}

	//sImpl.dispatch()
}

//func Check() {
//	ticker := time.NewTicker(time.Second * 1)
//	go func() {
//		for {
//			<-ticker.C
//			HHGame.SSLock.RLock()
//			sessionDatas := HHGame.Sessions
//			HHGame.SSLock.RUnlock()
//			for sessionId, sessionData := range sessionDatas {
//				for deskId, deskInfo := range sessionData.Desks {
//					if deskInfo.NextTime == 1 { // TODO 优化---
//						if deskInfo.Status == enum.Bet_Before { //
//							sendToClientChangeStatus(sessionId, enum.Betting, 10, int32(deskId), deskInfo.PlayerList)
//						} else if deskInfo.Status == enum.Betting { // TODO 发送两条消息会导致时间不同步
//							sendToClientChangeStatus(sessionId, enum.Show_Cards, 8, int32(deskId), deskInfo.PlayerList)
//							sendToClientRes(sessionId, int32(deskId), deskInfo.PlayerList)
//						} else if deskInfo.Status == enum.Show_Cards { //
//							sendToClientChangeStatus(sessionId, enum.Bet_Before, 5, int32(deskId), deskInfo.PlayerList)
//						}
//					}
//				}
//			}
//		}
//	}()
//}
