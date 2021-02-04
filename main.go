package main

import (
	"bufio"
	"flag"
	"mtSecKill/global"
	"mtSecKill/logs"
	"mtSecKill/secKill"
	"os"
	"strings"
	"time"
)

var skuId = flag.String("sku", "100012043978", "茅台商品ID")
var num = flag.Int("num", 1, "茅台商品数量")
var works = flag.Int("works", 1, "并发标签数")
var start = flag.String("time", "09:59:59", "开始时间---不带日期")
var browserPath = flag.String("execPath", "", "浏览器执行路径，路径不能有空格")

func init() {
	flag.Parse()
}

func main() {
	var err error
	execPath := ""
	if *browserPath != "" {
		execPath = *browserPath
	}
RE:
	logs.PrintlnInfo("浏览器路径为：", execPath)
	jdSecKill := secKill.NewJdSecKill(execPath, *skuId, *num, *works)
	jdSecKill.StartTime, err = global.Hour2Unix(*start)
	if err != nil {
		logs.Fatal("开始时间初始化失败", err)
	}

	// 判断时间是否是第二天
	if jdSecKill.StartTime.Unix() < time.Now().Unix() {
		jdSecKill.StartTime = jdSecKill.StartTime.AddDate(0, 0, 1)
	}
	jdSecKill.SyncJdTime()
	logs.PrintlnInfo("开始执行时间为：", jdSecKill.StartTime.Format(global.DateTimeFormatStr))

	// 启动秒杀任务
	err = jdSecKill.Run()
	if err != nil {
		if strings.Contains(err.Error(), "exec") {
			logs.PrintErr("默认浏览器执行路径未找到，" + execPath + "  请重新输入：")
			// 获取控制台输入的新浏览器命令路径
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				execPath = scanner.Text()
				if execPath != "" {
					break
				}
			}
			// 重载秒杀服务
			goto RE
		}
		logs.Fatal(err)
	}
}
