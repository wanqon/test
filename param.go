package main

import (
	"flag"
	"fmt"
)

var (
	Uri		string
	Date	string
)

func init()  {
	flag.StringVar(&Uri, "uri", "", "请求地址")
	flag.StringVar(&Date, "date", "", "日期")
	fmt.Println("init success")
}
