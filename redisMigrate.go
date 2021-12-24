package main

import (
	"flag"
	"fmt"
	"redis/util"
	"time"
)

//go run redisMigrate.go -SourcePrefix "*" -SourceIP 192.168.112.136 -SourcePort 6380 -SourceDB 1 -TargetIP 192.168.112.136 -TargetPort 6381 -TargetDB 1

func main() {
	inputvar := &util.InputVar{}

	flag.StringVar(&inputvar.SourcePrefix, "SourcePrefix", "spring*", "key Prefix")
	flag.StringVar(&inputvar.SourceIP, "SourceIP", "", "SourceIP")
	flag.StringVar(&inputvar.SourcePort, "SourcePort", "0", "SourcePort")
	flag.IntVar(&inputvar.SourceDB, "SourceDB", 0, "SourceDB")
	flag.StringVar(&inputvar.SourcePasswd, "SourcePasswd", "", "SourcePasswd")
	flag.StringVar(&inputvar.TargetIP, "TargetIP", "", "TargetIP")
	flag.StringVar(&inputvar.TargetPort, "TargetPort", "", "TargetPort")
	flag.IntVar(&inputvar.TargetDB, "TargetDB", 0, "TargetDB")
	flag.StringVar(&inputvar.TargetPasswd, "TargetPasswd", "", "TargetPasswd")
	flag.Parse()

	fmt.Println("start redis migrate:", time.Now())
	util.MigrateRedisDB(*inputvar)
	fmt.Println("stop redis migrate:", time.Now())
}
