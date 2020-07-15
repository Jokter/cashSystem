package main

import (
	"fmt"
	"net/http"
	"strconv"
)

/*
题目
	使用 Go 实现一个 商场收银系统，要求:
		1、输入商品的单价，数量，输出总价
		2、商场会举行活动，给用户优惠，活动有两种:
			a.打折活动(打5，6，7，8，9折扣)
			b.满减活动(按照总价参与例如满700减200，满300减50) 此时这个收银系统也需要支持
要求
	用工厂or策略模式实现
 */

//创建接口
type strategy interface {
	doOperation(cash int) float32
}

type discount struct {}
type reduce struct {}

//打折活动(打5，6，7，8，9折扣)
func (*discount) doOperation(cash int) float32 {
	if cash>1000{
		return float32(cash)*0.9
	}else if cash>900{
		return float32(cash)*0.8
	}else if cash>800{
		return float32(cash)*0.7
	}else if cash>700{
		return float32(cash)*0.6
	}else if cash>600{
		return float32(cash)*0.5
	}else{
		return float32(cash)
	}
}

//满减活动(按照总价参与例如满700减200，满300减50)
func (*reduce) doOperation(cash int) float32 {
	if cash>700{
		return float32(cash-200)
	}else if cash>300{
		return float32(cash-50)
	}else{
		return float32(cash)
	}
}

//创建Context
type context struct {
	strategy strategy
}

func (context *context) executeStrategy(strategy strategy,cash int) float32 {
	context.strategy = strategy
	return strategy.doOperation(cash)
}



func handler(w http.ResponseWriter,r *http.Request)  {
	//获取money和number
	r.ParseForm()
	money,_ := strconv.Atoi(r.PostFormValue("money"))
	number,_ := strconv.Atoi(r.PostFormValue("number"))
	cash := money*number
	fmt.Fprintln(w,"总价",cash)
	digit,_ := strconv.Atoi(r.PostFormValue("digit"))

	var context context

	//创建折扣
	var discount discount
	var strategy1 strategy = &discount

	//创建满减
	var reduce reduce
	var strategy2 strategy = &reduce

	switch digit {
	case 1:
		fmt.Fprintf(w,"打折活动价格 %.2f",context.executeStrategy(strategy1,cash))
		break
	case 2:
		fmt.Fprintf(w,"满减活动价格 %.2f",context.executeStrategy(strategy2,cash))
		break
	default:
		fmt.Fprintln(w,"价格",cash)
		break
	}
}

func main()  {
	http.HandleFunc("/cash",handler)
	//创建路由
	http.ListenAndServe(":8803",nil)
}