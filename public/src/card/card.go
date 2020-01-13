package card

import (
	"math/rand"
	"sort"
	"time"
)

const (
	GameTypeFAL = "fal"
)

const (
	CardTypeHeart   = "heart"
	CardTypeSpade   = "spade"
	CardTypeClub    = "club"
	CardTypeDiamond = "diamond"
	CardTypeJoker   = "joker"
)

const (
	ValueTypeSingle   = "single"
	ValueTypeSequence = "sequence"
	ValueTypeDouble   = "double"
	ValueTypePairs    = "pairs"
	ValueType3take1   = "3take1"
	ValueType4take2   = "4take2"
	ValueTypeBomb     = "bomb"
	ValueTypeThree    = "three"
	ValueTypeUnknown  = "unknown"
)

var falType = map[int64]int64{
	1: 12, 2: 13, 3: 1,
	4: 2, 5: 3, 6: 4,
	7: 5, 8: 6, 9: 7,
	10: 8, 11: 9, 12: 10,
	13: 11, 14: 14, 15: 15,
}

type Card struct {
	GameType   string
	CardSeq    int64
	CardNumber int64
	CardType   string
}

func (c Card) GetCardValue() int64 {
	switch c.GameType {
	case GameTypeFAL:
		return falType[c.CardNumber]
	default:
		return 0
	}
}

var Cards map[int64]Card

func DistributeCards(num int64) []int64 {
	rand.Seed(time.Now().UnixNano())
	dis := make([]int64, num+1)
	last := rand.Int63n(num) + 1
	dis[last] = num
	for i := num - 1; i > 0; i-- {
		index := rand.Int63n(i) + 1
		count := int64(0)
		for j := int64(1); j < num+1; j++ {
			if dis[j] == 0 {
				count++
			}
			if count == index {
				dis[j] = int64(i)
				break
			}
		}
	}
	dis[last] = 0
	return dis[1:]
}

func init() {
	if Cards == nil {
		Cards = make(map[int64]Card, 54)
	}
	addCard := func(seq, value int64, ctype string) {
		card := Card{
			GameType:   GameTypeFAL,
			CardSeq:    seq,
			CardNumber: value,
			CardType:   ctype,
		}
		Cards[seq] = card
	}
	var ctype string
	var value int64
	for i := int64(0); i < 52; i++ {
		switch i / 13 {
		case 0:
			ctype = CardTypeSpade
		case 1:
			ctype = CardTypeHeart
		case 2:
			ctype = CardTypeClub
		case 3:
			ctype = CardTypeDiamond
		}
		value++
		if i%13 == 0 {
			value = 1
		}
		addCard(i, value, ctype)
	}
	addCard(52, 14, CardTypeJoker)
	addCard(53, 15, CardTypeJoker)
}

//获取重复数个数及其值和长度
func GetRepeatNumAndValue(cards []int64) (int64, int64, int64) {
	if len(cards) == 0 {
		return 0, 0, 0
	}
	var nums []int64
	for _, seq := range cards {
		nums = append(nums, Cards[seq].GetCardValue())
	}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})
	max := int64(1)
	repeat := int64(1)
	value := nums[0]
	if len(nums) == 1 {
		return 1, nums[0], 1
	} else if len(nums) == 2 && nums[0] == 14 && nums[1] == 15 || len(nums) == 2 && nums[1] == 14 && nums[0] == 15 {
		return 4, nums[0], 4
	}
	for i := 1; i < len(nums); i++ {
		if nums[i] == nums[i-1] {
			repeat = repeat + 1
			if repeat > max {
				max = repeat
				value = nums[i]
			}
		} else {
			repeat = 1
		}
	}
	//fmt.Println("重复次数", max, "重复值", value)
	return max, value, int64(len(nums))
}

//获取出牌类型
func GetCardsType(repeatnum int64, nums []int64) string {
	if repeatnum == 1 {
		if len(nums) == 1 {
			return ValueTypeSingle
		} else if len(nums) >= 5 && IsSequence(nums) {
			return ValueTypeSequence
		}
	} else if repeatnum == 2 {
		if len(nums) == 2 {
			return ValueTypeDouble
		} else if len(nums) >= 6 && len(nums)%2 == 0 && IsSequencePair(nums) {
			return ValueTypePairs
		}
	} else if repeatnum == 3 {
		if len(nums) == 3 {
			return ValueTypeThree
		} else if len(nums) == 4 {
			return ValueType3take1
		}
	} else if repeatnum == 4 {
		if len(nums) == 4 || len(nums) == 2 {
			return ValueTypeBomb
		} else if len(nums) == 6 {
			return ValueType4take2
		}
	}
	return ValueTypeUnknown
}

//连牌
func IsSequence(nums []int64) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i]-nums[i-1] != 1 {
			return false
		}
	}
	return true
}

//连对
func IsSequencePair(nums []int64) bool {
	for i := 2; i < len(nums); i++ {
		if nums[i]-nums[i-2] != 1 {
			return false
		}
		i = i + 1
	}
	return true
}
