package model

import (
	"github.com/shopspring/decimal"
)

/*
** 赔率
**
 */

// SettlementR 结算
func SettlementR(rPx int32, bet map[int32]int64, rOod float32, lOod map[int32]float32) (int64, int64) { // 加金币的总额

	// 只下注了红色方
	if len(bet) == 1 && bet[2] > 0 {
		// de := decimal.NewFromFloat(float64(rOod)) // 0.95 --> 0.94
		de, _ := addRet(rOod, 0)
		be := decimal.New(bet[2], 0) //  30000^1

		return deleteRet(de.Mul(be).IntPart()+1000*bet[2], de.Mul(be).IntPart())
		// return de.Mul(be).IntPart() + 1*bet[0], de.Mul(be).IntPart()
		// return r1, r2
	}

	// 只下注了幸运一击
	if len(bet) == 1 && bet[1] > 0 {
		_, ok := lOod[rPx]
		if ok { // 特殊牌型
			// de := decimal.NewFromFloat(float64(lOod[rPx])) //
			de, _ := addRet(lOod[rPx], 0)
			be := decimal.New(bet[1], 0)
			return deleteRet(de.Mul(be).IntPart()+1000*bet[1], de.Mul(be).IntPart())
			// return r1, r2
			// return de.Mul(be).IntPart() + bet[1], de.Mul(be).IntPart()
			// return (i + 1) * bet[1], bet[0] // 赢
		} else {
			return 0, 0 // 输
		}
	}

	// 只下注了黑方
	if len(bet) == 1 && bet[0] > 0 {
		return 0, 0
	}

	// 下注了红色方和幸运一击
	if len(bet) == 2 && bet[2] > 0 && bet[1] > 0 {
		_, ok := lOod[rPx]
		if ok { // 特殊牌型
			// de0 := decimal.NewFromFloat(float64(rOod))
			// de1 := decimal.NewFromFloat(float64(lOod[rPx]))
			de0, de1 := addRet(rOod, lOod[rPx])

			be0 := decimal.New(bet[2], 0)
			be1 := decimal.New(bet[1], 0)

			return deleteRet(de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart()+1000*bet[2]+1000*bet[1], de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart())
			//return r1, r2
			//return de0.Mul(be0).IntPart() + de1.Mul(be1).IntPart() + bet[0] + bet[1], de0.Mul(be0).IntPart() + de1.Mul(be1).IntPart()
		} else {
			// de := decimal.NewFromFloat(float64(rOod))
			de, _ := addRet(rOod, 0)

			be := decimal.New(bet[2], 0)
			// deleteRet(de.Mul(be).IntPart() + 1*bet[0], de.Mul(be).IntPart())
			return deleteRet(de.Mul(be).IntPart()+1000*bet[2], de.Mul(be).IntPart()) // 赢
		}
	}

	// 下注了黑色方和幸运一击
	if len(bet) == 2 && bet[0] > 0 && bet[1] > 0 {
		_, ok := lOod[rPx]
		if ok {
			_, de1 := addRet(0, lOod[rPx])
			be1 := decimal.New(bet[1], 0)
			return deleteRet(de1.Mul(be1).IntPart()+1000*bet[1], 0)
		} else {
			return 0, 0
		}

	}

	// 下注了红黑区域
	if len(bet) == 2 && bet[0] > 0 && bet[2] > 0 {
		// de := decimal.NewFromFloat(float64(rOod))

		de, _ := addRet(rOod, 0)
		be := decimal.New(bet[2], 0)
		return deleteRet(de.Mul(be).IntPart()+1000*bet[2], de.Mul(be).IntPart())
	}

	// 下注了红黑幸运一击
	if len(bet) == 3 {
		_, ok := lOod[rPx]
		if ok { // 特殊牌型
			//de1 := decimal.NewFromFloat(float64(lOod[rPx]))
			//de0 := decimal.NewFromFloat(float64(rOod))
			de0, de1 := addRet(rOod, lOod[rPx])
			be0 := decimal.New(bet[2], 0)
			be1 := decimal.New(bet[1], 0)
			return deleteRet(de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart()+1000*bet[2]+1000*bet[1], de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart())
		} else { // 不存在特殊牌型
			// de := decimal.NewFromFloat(float64(rOod))
			de, _ := addRet(rOod, 0)
			be := decimal.New(bet[2], 0)
			return deleteRet(de.Mul(be).IntPart()+1000*bet[2], de.Mul(be).IntPart())
		}
	}

	return 0, 0
}

// SettlementB 结算
func SettlementB(bPx int32, bet map[int32]int64, bOod float32, lOod map[int32]float32) (int64, int64) { // 加金币的总额
	// 只下注了黑色方
	if len(bet) == 1 && bet[0] > 0 {
		// de := decimal.NewFromFloat(float64(rOod)) // 0.95 --> 0.94
		de, _ := addRet(bOod, 0)
		be := decimal.New(bet[0], 0) //  30000^1

		return deleteRet(de.Mul(be).IntPart()+1000*bet[0], de.Mul(be).IntPart())
		// return de.Mul(be).IntPart() + 1*bet[0], de.Mul(be).IntPart()
		// return r1, r2
	}

	// 只下注了幸运一击
	if len(bet) == 1 && bet[1] > 0 {
		_, ok := lOod[bPx]
		if ok { // 特殊牌型
			// de := decimal.NewFromFloat(float64(lOod[rPx])) //
			de, _ := addRet(lOod[bPx], 0)
			be := decimal.New(bet[1], 0)
			return deleteRet(de.Mul(be).IntPart()+1000*bet[1], de.Mul(be).IntPart())
			// return r1, r2
			// return de.Mul(be).IntPart() + bet[1], de.Mul(be).IntPart()
			// return (i + 1) * bet[1], bet[0] // 赢
		} else {
			return 0, 0 // 输
		}
	}

	// 只下注了红方
	if len(bet) == 1 && bet[2] > 0 {
		return 0, 0
	}

	// 下注了黑色方和幸运一击
	if len(bet) == 2 && bet[0] > 0 && bet[1] > 0 {
		_, ok := lOod[bPx]
		if ok { // 特殊牌型
			// de0 := decimal.NewFromFloat(float64(rOod))
			// de1 := decimal.NewFromFloat(float64(lOod[rPx]))
			de0, de1 := addRet(bOod, lOod[bPx])

			be2 := decimal.New(bet[0], 0)
			be1 := decimal.New(bet[1], 0)

			return deleteRet(de0.Mul(be2).IntPart()+de1.Mul(be1).IntPart()+1000*bet[0]+1000*bet[1], de0.Mul(be2).IntPart()+de1.Mul(be1).IntPart())
			//return r1, r2
			//return de0.Mul(be0).IntPart() + de1.Mul(be1).IntPart() + bet[0] + bet[1], de0.Mul(be0).IntPart() + de1.Mul(be1).IntPart()
		} else {
			// de := decimal.NewFromFloat(float64(rOod))
			de, _ := addRet(bOod, 0)

			be := decimal.New(bet[0], 0)
			// deleteRet(de.Mul(be).IntPart() + 1*bet[0], de.Mul(be).IntPart())
			return deleteRet(de.Mul(be).IntPart()+1000*bet[0], de.Mul(be).IntPart()) // 赢
		}
	}

	// 下注了红色方和幸运一击
	if len(bet) == 2 && bet[2] > 0 && bet[1] > 0 {
		_, ok := lOod[bPx]
		if ok {
			_, de1 := addRet(0, lOod[bPx])
			be1 := decimal.New(bet[1], 0)
			return deleteRet(de1.Mul(be1).IntPart()+1000*bet[1], 0)
		} else {
			return 0, 0
		}

	}

	// 下注了红黑区域
	if len(bet) == 2 && bet[0] > 0 && bet[2] > 0 {
		// de := decimal.NewFromFloat(float64(rOod))

		de, _ := addRet(bOod, 0)
		be := decimal.New(bet[0], 0)
		return deleteRet(de.Mul(be).IntPart()+1000*bet[0], de.Mul(be).IntPart())
	}

	// 下注了红黑幸运一击
	if len(bet) == 3 {
		_, ok := lOod[bPx]
		if ok { // 特殊牌型
			//de1 := decimal.NewFromFloat(float64(lOod[rPx]))
			//de0 := decimal.NewFromFloat(float64(rOod))
			de0, de1 := addRet(bOod, lOod[bPx])
			be0 := decimal.New(bet[0], 0)
			be1 := decimal.New(bet[1], 0)
			return deleteRet(de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart()+bet[0]+1000*bet[1], de0.Mul(be0).IntPart()+de1.Mul(be1).IntPart())
		} else { // 不存在特殊牌型
			// de := decimal.NewFromFloat(float64(rOod))
			de, _ := addRet(bOod, 0)
			be := decimal.New(bet[0], 0)
			return deleteRet(de.Mul(be).IntPart()+1000*bet[0], de.Mul(be).IntPart())
		}
	}

	return 0, 0
}

// 增加1000倍
func addRet(f1, f2 float32) (decimal.Decimal, decimal.Decimal) {
	var ret float32 = 1000
	return decimal.NewFromFloat(float64(f1 * ret)), decimal.NewFromFloat(float64(f2 * ret))

}

// 去掉1000倍
func deleteRet(i1, i2 int64) (int64, int64) {
	var ret int64 = 1000
	return i1 / ret, i2 / ret
}
