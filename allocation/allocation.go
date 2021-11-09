package allocation

import (
	"math/rand"
	"time"
)

type Allocation struct {
	total_money int //总金额
	num         int //普通红包数
	lower       int //下界
	luck        []int
	luck_num    int //锦鲤数
}

func NewAllocation(total_money, num, lower, upper int) *Allocation {
	luck_num := num / 1000        //千分之一是锦鲤
	luck := make([]int, luck_num) //锦鲤数组
	avg := int(float32(total_money) / float32(num))
	rand.Seed(time.Now().Unix())
	var rand_number int
	for i := 0; i < luck_num; i++ {
		rand_number = rand.Intn(upper-20*avg+1) + 20*avg //锦鲤金额在范围[avg*20,upper]中随机
		luck[i] = rand_number
		total_money -= rand_number //总钱数扣去锦鲤部分
	}
	num -= luck_num //普通红包数
	return &Allocation{
		total_money: total_money,
		num:         num,
		lower:       lower,
		luck:        luck,
		luck_num:    luck_num,
	}
}
func (a Allocation) AllocateMoney(need_number int) []int {
	res := make([]int, need_number)
	for i := 0; i < need_number; i++ {
		if a.luck_num > 0 && rand.Intn(a.num+a.luck_num) < a.luck_num { //概率判断是否为锦鲤
			a.luck_num--
			res[i] = a.luck[a.luck_num]
		} else {
			avg := int(float64(a.total_money) / float64(a.num)) //普通红包的均值
			upper := 2*avg - a.lower                            //框定范围[lower,upper]保证均值为avg
			res[i] = rand.Intn(upper-a.lower+1) + a.lower
			a.num--
			a.total_money -= res[i]
		}
	}
	return res
}

