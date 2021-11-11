package allocation

import (
	"math/rand"
	"time"
)

type Allocation struct {
	totalMoney int //总金额
	num         int //普通红包数
	lower       int //下界
	luck        []int
	luckNum    int //锦鲤数
}

func NewAllocation(totalMoney, num, lower, upper int) *Allocation {
	luckNum := num / 1000        //千分之一是锦鲤
	luck := make([]int, luckNum) //锦鲤数组
	avg := int(float32(totalMoney) / float32(num))
	rand.Seed(time.Now().Unix())
	var randNumber int
	for i := 0; i < luckNum; i++ {
		randNumber = rand.Intn(upper-20*avg+1) + 20*avg //锦鲤金额在范围[avg*20,upper]中随机
		luck[i] = randNumber
		totalMoney -= randNumber //总钱数扣去锦鲤部分
	}
	num -= luckNum //普通红包数
	return &Allocation{
		totalMoney: totalMoney,
		num:         num,
		lower:       lower,
		luck:        luck,
		luckNum:    luckNum,
	}
}

func (a Allocation) AllocateMoney(needNumber int) []int {
	res := make([]int, needNumber)
	for i := 0; i < needNumber; i++ {
		if a.luckNum > 0 && rand.Intn(a.num+a.luckNum) < a.luckNum { //概率判断是否为锦鲤
			a.luckNum--
			res[i] = a.luck[a.luckNum]
		} else {
			avg := int(float64(a.totalMoney) / float64(a.num)) //普通红包的均值
			upper := 2*avg - a.lower                            //框定范围[lower,upper]保证均值为avg
			res[i] = rand.Intn(upper-a.lower+1) + a.lower
			a.num--
			a.totalMoney -= res[i]
		}
	}
	return res
}

