package main

import (
	"flag"
	"fmt"
	"pay.sc.weibo.com/accounts/finance"
	"pay.sc.weibo.com/accounts/logger"
)

func main() {
	flag.Parse()
	logger.Info("pay.sc.weibo.com start", logger.LogField{})
	fmt.Println("test")

	if Uri == "charge" {
		finance.WbpayCharge(Date)
		logger.Info("wbpaycharge done", logger.LogField{"uri":Uri,"date":Date})
	}
	//return

}

