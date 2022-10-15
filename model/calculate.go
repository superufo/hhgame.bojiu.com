package model

import (
	"hhgame.bojiu.com/enum"
	"sort"
)

// 计算牌型
// 规则：豹子 > 同花顺 > 同花 > 顺子 > 对子 > 单牌

// CheckWhoWin 检测输赢
func CheckWhoWin(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	// 红方牌型
	red := Check(redCards)

	// 黑方牌型
	black := Check(blackCards)

	// 牌型不同时
	if red != black {
		if red > black {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 豹子
	if red == enum.Leopard {
		return CompareLeopard(tmpRedCards, tmpBlackCards)
	}
	// 同花顺
	if red == enum.Flush {
		return CompareFlush(tmpRedCards, tmpBlackCards)
	}
	// 同花
	if red == enum.Homogamy {
		return CompareHomogamy(tmpRedCards, tmpBlackCards)
	}
	// 顺子
	if red == enum.Straight {
		return CompareStraight(tmpRedCards, tmpBlackCards)
	}
	// 对子
	if red == enum.PairBig {
		return ComparePairBig(tmpRedCards, tmpBlackCards)
	}
	if red == enum.PairSmall {
		return ComparePairSmall(tmpRedCards, tmpBlackCards)
	}
	// 单张
	if red == enum.Leaflet {
		return CompareLeaflet(tmpRedCards, tmpBlackCards)
	}
	return 0
}

// Check 检测牌型
func Check(cards []int) int32 {
	tmpCards := make([]int, len(cards))
	copy(tmpCards, cards)

	// 牌值与花色分离
	valueCards, colorCards := CalculateValue(tmpCards)

	// 检测豹子
	if valueCards[0] == valueCards[2] {
		return enum.Leopard
	}

	// 检测同花顺/同花     特殊情况 A可以为1  可以为A
	if colorCards[0] == colorCards[2] {
		if (valueCards[0]+1 == valueCards[1] && valueCards[0]+2 == valueCards[2]) ||
			(valueCards[0]+11 == valueCards[1] && valueCards[0]+12 == valueCards[2]) {
			return enum.Flush // 同花顺
		}
		return enum.Homogamy // 同花
	}

	// 检测顺子  特殊情况 A可以为1  可以为A
	if (valueCards[0]+1 == valueCards[1] && valueCards[0]+2 == valueCards[2]) ||
		(valueCards[0]+11 == valueCards[1] && valueCards[0]+12 == valueCards[2]) {
		return enum.Straight
	}

	// 检测大对子
	if (valueCards[0] == valueCards[1] && valueCards[0] != valueCards[2] && valueCards[1] > 8) ||
		(valueCards[1] == valueCards[2] && valueCards[0] != valueCards[1] && valueCards[1] > 8) ||
		(valueCards[0] == valueCards[1] && valueCards[0] == 1) {
		return enum.PairBig
	}

	// 检测小对子
	if (valueCards[0] == valueCards[1] && valueCards[0] != valueCards[2] && valueCards[1] < 9) ||
		(valueCards[1] == valueCards[2] && valueCards[0] != valueCards[1] && valueCards[1] < 9) {
		return enum.PairSmall
	}

	// 单张
	return enum.Leaflet
}

// CalculateValue 计算牌值
func CalculateValue(cards []int) (value, color []int) {
	tmpCards := make([]int, len(cards))
	copy(tmpCards, cards)

	value = make([]int, 0)
	color = make([]int, 0)

	for _, v := range tmpCards {
		value = append(value, v/100)
	}

	for _, v := range tmpCards {
		color = append(color, v%100)
	}

	// 排序,从小到大
	sort.Ints(value)
	sort.Ints(color)

	return value, color
}

// CompareLeopard 豹子对比
func CompareLeopard(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	redCard, _ := CalculateValue(tmpRedCards)
	blackCard, _ := CalculateValue(tmpBlackCards)
	if redCard[0] != 1 && blackCard[0] != 1 {
		if redCard[0] > blackCard[0] {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	if redCard[0] == 1 {
		return enum.Red_Win
	}

	if blackCard[0] == 1 {
		return enum.Black_Win
	}

	return 0
}

// CompareFlush 同花顺对比
func CompareFlush(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	redCard, redColor := CalculateValue(tmpRedCards)
	blackCard, blackColor := CalculateValue(tmpBlackCards)

	// [1,12,13] [1,2,3] [2,3,4]
	if redCard[0]+11 == redCard[1] && blackCard[0]+11 == blackCard[1] {
		if redColor[0] < blackColor[0] {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	if redCard[0]+11 == redCard[1] && blackCard[0]+11 != blackCard[1] {
		return enum.Red_Win
	}

	if redCard[0]+11 != redCard[1] && blackCard[0]+11 == blackCard[1] {
		return enum.Black_Win
	}

	if redCard[0]+11 != redCard[1] && blackCard[0]+11 != blackCard[1] {
		if redCard[0] > blackCard[0] {
			return enum.Red_Win
		}

		if redCard[0] < blackCard[0] {
			return enum.Black_Win
		}

		if redCard[0] == blackCard[0] {
			if redColor[0] < blackColor[0] {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
	}

	return 0
}

// CompareHomogamy 对比同花
func CompareHomogamy(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	redCard, redColor := CalculateValue(tmpRedCards)
	blackCard, blackColor := CalculateValue(tmpBlackCards)

	// 都带A,且为同花
	if redCard[0] == 1 && redCard[0] == blackCard[0] {
		if redCard[1] == blackCard[1] { // 二位相同
			if redCard[2] == blackCard[2] { // 三位相同
				if redColor[0] < blackColor[0] {
					return enum.Red_Win
				} else {
					return enum.Black_Win
				}
			}
			if redCard[2] > blackCard[2] { // 三位不同
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
		if redCard[1] > blackCard[1] { // 二位不同
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 一方带A
	if redCard[0] == 1 && blackCard[0] != 1 {
		return enum.Red_Win
	}

	if redCard[0] != 1 && blackCard[0] == 1 {
		return enum.Black_Win
	}

	// 都不带A
	if redCard[0] != 1 && blackCard[0] != 1 {
		if redCard[0] == blackCard[0] { // 一位相同
			if redCard[1] == blackCard[1] { // 二位相同
				if redCard[2] == blackCard[2] { // 三位相同
					if redColor[0] < blackColor[0] {
						return enum.Red_Win
					} else {
						return enum.Black_Win
					}
				}
				if redCard[2] > blackCard[2] { // 三位不同
					return enum.Red_Win
				} else {
					return enum.Black_Win
				}
			}
			if redCard[1] > blackCard[1] { // 二位不同
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		if redCard[0] > blackCard[0] { // 一位不同
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	return 0
}

// CompareStraight 对比顺子
func CompareStraight(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)
	sort.Ints(tmpRedCards)
	sort.Ints(tmpBlackCards)
	// redCard, redColor := CalculateValue(tmpRedCards)
	// blackCard, blackColor := CalculateValue(tmpBlackCards)

	// [1,12,13] [1,2,3] [2,3,4]
	// 同为[1,12,13]
	if tmpRedCards[0]/100+11 == tmpRedCards[1]/100 && tmpBlackCards[0]+11 == tmpBlackCards[1]/100 {
		if tmpRedCards[0]%100 < tmpBlackCards[0]%100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	if tmpRedCards[0]/100+11 == tmpRedCards[1]/100 && tmpBlackCards[0]/100+11 != tmpBlackCards[1]/100 {
		return enum.Red_Win
	}

	if tmpRedCards[0]/100+11 != tmpRedCards[1]/100 && tmpBlackCards[0]/100+11 == tmpBlackCards[1]/100 {
		return enum.Black_Win
	}

	// 都不是[1,12,13]
	if tmpRedCards[0]/100+11 != tmpRedCards[1]/100 && tmpBlackCards[0]/100+11 != tmpBlackCards[1]/100 {
		if tmpRedCards[0]/100 > tmpBlackCards[0]/100 {
			return enum.Red_Win
		}

		if tmpRedCards[0]/100 < tmpBlackCards[0]/100 {
			return enum.Black_Win
		}

		// 数值相同但不为同花
		if tmpRedCards[0]/100 == tmpBlackCards[0]/100 {
			if tmpRedCards[2]%100 < tmpBlackCards[2]%100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
	}

	return 0
}

// ComparePairBig 对比大对子
func ComparePairBig(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	sort.Ints(tmpRedCards)
	sort.Ints(tmpBlackCards)
	// [1,1,2]
	// 牌型 [1,1,?] 都为对A
	if tmpRedCards[1]/100 == 1 && tmpBlackCards[1]/100 == 1 {
		// 剩余单张相同,比较花色
		if tmpRedCards[2]/100 == tmpBlackCards[2]/100 {
			if tmpRedCards[2]%100 < tmpBlackCards[2]%100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
		// 剩余单张不同,不存在A单独存在,直接比值
		if tmpRedCards[2]/100 > tmpBlackCards[2]/100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 一方存在对A时
	if tmpRedCards[1]/100 == 1 && tmpBlackCards[1]/100 != 1 {
		return enum.Red_Win
	}

	if tmpRedCards[1]/100 != 1 && tmpBlackCards[1]/100 == 1 {
		return enum.Red_Win
	}

	// 都不存在对A
	if tmpRedCards[1]/100 != 1 && tmpBlackCards[1]/100 != 1 {
		// 对子值相同
		// 同时占据1 2 位
		if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[0]/100 == tmpRedCards[1]/100 && tmpBlackCards[0]/100 == tmpBlackCards[1]/100 {
			// 剩余单张相同,比较花色
			if tmpRedCards[2]/100 == tmpBlackCards[2]/100 {
				if tmpRedCards[2]%100 < tmpBlackCards[2]%100 {
					return enum.Red_Win
				} else {
					return enum.Black_Win
				}
			}

			// 剩余单张不同,不存在A单独存在,直接比值
			if tmpRedCards[2]/100 > tmpBlackCards[2]/100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		// 对子值相同
		// 同时占据2 3 位
		if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[2]/100 == tmpRedCards[2]/100 && tmpBlackCards[2]/100 == tmpBlackCards[2]/100 {
			// 剩余单张相同,比较花色
			if tmpRedCards[0]/100 == tmpBlackCards[0]/100 {
				if tmpRedCards[0]%100 < tmpBlackCards[0]%100 {
					return enum.Red_Win
				} else {
					return enum.Black_Win
				}
			}
			// 剩余单张不同,A可以单独存在
			if tmpRedCards[0]/100 != tmpBlackCards[0]/100 {
				if tmpRedCards[0]/100 == 1 {
					return enum.Red_Win
				}
				if tmpBlackCards[0]/100 == 1 {
					return enum.Black_Win
				}
				if tmpRedCards[0]/100 > tmpBlackCards[0]/100 {
					return enum.Red_Win
				} else {
					return enum.Black_Win
				}
			}
		}

		// 对子值相同
		// 占据位置不同  exp:[445]  [?44]
		if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[0]/100 == tmpRedCards[1]/100 && tmpBlackCards[1]/100 == tmpBlackCards[2]/100 {
			if tmpBlackCards[0]/100 == 1 {
				return enum.Black_Win
			} else {
				return enum.Red_Win
			}
		}

		// 对子值相同
		// 占据位置不同    exp:  [?55]  [556]
		if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[1]/100 == tmpRedCards[2]/100 && tmpBlackCards[0]/100 == tmpBlackCards[1]/100 {
			if tmpRedCards[0]/100 == 1 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		// 对子值不同
		if tmpRedCards[1]/100 > tmpBlackCards[1]/100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	return 0
}

// ComparePairSmall 对比小对子
func ComparePairSmall(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	sort.Ints(tmpRedCards)
	sort.Ints(tmpBlackCards)

	// 对子值相同
	// 同时占据1 2 位
	if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[0]/100 == tmpRedCards[1]/100 && tmpBlackCards[0]/100 == tmpBlackCards[1]/100 {
		// 剩余单张相同,比较花色
		if tmpRedCards[2]/100 == tmpBlackCards[2]/100 {
			if tmpRedCards[2]%100 < tmpBlackCards[2]%100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		// 剩余单张不同,不存在A单独存在,直接比值
		if tmpRedCards[2]/100 > tmpBlackCards[2]/100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 对子值相同
	// 同时占据2 3 位
	if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[2]/100 == tmpRedCards[2]/100 && tmpBlackCards[2]/100 == tmpBlackCards[2]/100 {
		// 剩余单张相同,比较花色
		if tmpRedCards[0]/100 == tmpBlackCards[0]/100 {
			if tmpRedCards[0]%100 < tmpBlackCards[0]%100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
		// 剩余单张不同,A可以单独存在
		if tmpRedCards[0]/100 != tmpBlackCards[0]/100 {
			if tmpRedCards[0]/100 == 1 {
				return enum.Red_Win
			}
			if tmpBlackCards[0]/100 == 1 {
				return enum.Black_Win
			}
			if tmpRedCards[0]/100 > tmpBlackCards[0]/100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}
	}

	// 对子值相同
	// 占据位置不同  exp:[445]  [?44]
	if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[0]/100 == tmpRedCards[1]/100 && tmpBlackCards[1]/100 == tmpBlackCards[2]/100 {
		if tmpBlackCards[0]/100 == 1 {
			return enum.Black_Win
		} else {
			return enum.Red_Win
		}
	}

	// 对子值相同
	// 占据位置不同    exp:  [?55]  [556]
	if tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[1]/100 == tmpRedCards[2]/100 && tmpBlackCards[0]/100 == tmpBlackCards[1]/100 {
		if tmpRedCards[0]/100 == 1 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 对子值不同
	if tmpRedCards[1]/100 > tmpBlackCards[1]/100 {
		return enum.Red_Win
	} else {
		return enum.Black_Win
	}

}

// CompareLeaflet 对比单张
func CompareLeaflet(redCards, blackCards []int) int32 {
	tmpRedCards := make([]int, len(redCards))
	tmpBlackCards := make([]int, len(blackCards))
	copy(tmpRedCards, redCards)
	copy(tmpBlackCards, blackCards)

	sort.Ints(tmpRedCards)
	sort.Ints(tmpBlackCards)

	// 两组牌相同
	if tmpRedCards[0]/100 == tmpBlackCards[0]/100 && tmpRedCards[1]/100 == tmpBlackCards[1]/100 && tmpRedCards[2]/100 == tmpBlackCards[2]/100 {
		// 都存在A,直接判断A
		if tmpRedCards[0]/100 == 1 && tmpRedCards[0]/100 == tmpBlackCards[0] {
			if tmpRedCards[0]%100 < tmpBlackCards[0]%100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		// 不存在A,直接判断 三位花色
		if tmpRedCards[2]%100 < tmpBlackCards[2]%100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	// 两组牌不同 单A
	if tmpRedCards[0]/100 == 1 && tmpBlackCards[0]/100 != 1 {
		return enum.Red_Win
	}

	if tmpRedCards[0]/100 != 1 && tmpBlackCards[0]/100 == 1 {
		return enum.Black_Win
	}

	// 都不存在A
	if tmpRedCards[0]/100 != 1 && tmpBlackCards[0]/100 != 1 {
		if tmpRedCards[2]/100 == tmpBlackCards[2]/100 && tmpRedCards[1]/100 == tmpBlackCards[1]/100 { // 2  3  same
			if tmpRedCards[0]/100 > tmpBlackCards[0]/100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		if tmpRedCards[2]/100 == tmpBlackCards[2]/100 && tmpRedCards[1]/100 != tmpBlackCards[1]/100 { // 3 same
			if tmpRedCards[1]/100 > tmpBlackCards[1]/100 {
				return enum.Red_Win
			} else {
				return enum.Black_Win
			}
		}

		if tmpRedCards[2]/100 > tmpBlackCards[2]/100 {
			return enum.Red_Win
		} else {
			return enum.Black_Win
		}
	}

	return 0
}
