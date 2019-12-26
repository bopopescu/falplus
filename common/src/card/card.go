package card

import (
	"fmt"
	"math/rand"
	"time"
)

const(
	GameTypeFAL = "fal"
)

const(
	CardTypeHeart = "heart"
	CardTypeSpade = "spade"
	CardTypeClub = "club"
	CardTypeDiamond = "diamond"
	CardTypeJoker = "joker"
)

var falType = map[int]int{
	1:12, 2:13, 3:1,
	4:2,  5:3,  6:4,
	7:5,  8:6,  9:7,
	10:8, 11:9, 12:10,
	13:11,14:14,15:15,
}

type Card struct {
	GameType   string
	CardSeq    int
	CardNumber int
	CardType   string
}

func (c Card) GetCardValue() int {
	switch c.GameType {
	case GameTypeFAL:
		return falType[c.CardNumber]
	default:
		return 0
	}
}

var Cards map[int]Card

func DistributeCards(num int) []int64 {
	rand.Seed(time.Now().UnixNano())
	dis := make([]int64, num+1)
	last := rand.Intn(num) + 1
	dis[last] = int64(num)
	for i:=num-1;i>0;i--{
		index := rand.Intn(i) + 1
		count := 0
		for j:=1;j<num+1;j++{
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
		Cards = make(map[int]Card, 54)
	}
	addCard := func(seq, value int, ctype string) {
		card := Card{
			GameType:   GameTypeFAL,
			CardSeq:    seq,
			CardNumber: value,
			CardType:   ctype,
		}
		Cards[seq] = card
	}
	var ctype string
	var value int
	for i:=0;i<52;i++{
		switch i/13 {
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

//获取重复数个数及其值
func GetRepeatNumAndValue(nums []int) (int, int) {
	max := 1
	repeat := 1
	value := nums[0]
	if len(nums) == 1 {
		return 1, nums[0]
	} else if len(nums) == 2 && nums[0] == 14 && nums[1] == 15 || len(nums) == 2 && nums[1] == 14 && nums[0] == 15 {
		return 4, nums[0]
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
	fmt.Println("重复次数", max, "重复值", value)
	return max, value
}

//获取出牌类型
func GetCardsType(repeatnum int, nums []int) string {
	if repeatnum == 1 {
		if len(nums) == 1 {
			str := "single"
			return str
		} else if len(nums) >= 5 && IsSequence(nums) {
			str := "sequence"
			return str
		} else {
			str := "err"
			return str
		}
	} else if repeatnum == 2 {
		if len(nums) == 2 {
			str := "double"
			return str
		} else if len(nums) >= 6 && len(nums)%2 == 0 && IsSequencePair(nums) {
			str := "sequencepair"
			return str
		} else {
			str := "err"
			return str
		}
	} else if repeatnum == 3 {
		if len(nums) == 3 {
			str := "three"
			return str
		} else if len(nums) == 4 {
			str := "threeandone"
			return str
		} else {
			str := "err"
			return str
		}
	} else if repeatnum == 4 {
		if len(nums) == 4 || len(nums) == 2 {
			str := "bomb"
			return str
		} else if len(nums) == 6 {
			str := "fourandtwo"
			return str
		} else {
			str := "err"
			return str
		}
	}
	str := "err"
	return str
}

//连牌
func IsSequence(nums []int) bool {
	for i := 1; i < len(nums); i++ {
		if nums[i]-nums[i-1] != 1 {
			return false
		}
	}
	return true
}

//连对
func IsSequencePair(nums []int) bool {
	for i := 2; i < len(nums); i++ {
		if nums[i]-nums[i-2] != 1 {
			return false
		}
		i = i + 1
	}
	return true
}