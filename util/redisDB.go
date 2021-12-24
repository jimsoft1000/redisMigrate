package util

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
)

var ctx = context.Background()

//获取redis连接
func GetRedisDB(inputvar InputVar, flag string) *redis.Client {
	var redisCon *redis.Client
	if flag == SourceFlag {
		redisCon = redis.NewClient(&redis.Options{
			Addr:     inputvar.SourceIP + ":" + inputvar.SourcePort,
			Password: inputvar.SourcePasswd,
			DB:       inputvar.SourceDB,
		})
	} else if flag == TargetFlag {
		redisCon = redis.NewClient(&redis.Options{
			Addr:     inputvar.TargetIP + ":" + inputvar.TargetPort,
			Password: inputvar.TargetPasswd,
			DB:       inputvar.TargetDB,
		})
	}

	pong, err := redisCon.Ping(ctx).Result()
	if err != nil {
		fmt.Println("flag:", flag)
		fmt.Println("inputvar:", inputvar)
		fmt.Println(pong, err)
		os.Exit(0)
	}
	return redisCon
}

//迁移整库
func MigrateRedisDB(inputvar InputVar) {
	var stringCount int64 = 0
	var listCount int64 = 0
	var setCount int64 = 0
	var zsetCount int64 = 0
	var hashCount int64 = 0
	sourceCon := GetRedisDB(inputvar, SourceFlag)
	targetCon := GetRedisDB(inputvar, TargetFlag)
	keyinfo := sourceCon.Info(ctx, "Keyspace").Val()
	keyinfo = strings.TrimRight(keyinfo, "\n")
	info := strings.Split(keyinfo, "\n")
	for _, value := range info {
		fmt.Println("value:", value)
	}
	iter := sourceCon.Scan(ctx, 0, inputvar.SourcePrefix, ScanCount).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		//fmt.Println("key:", key)
		keytype := sourceCon.Type(ctx, key).Val()
		switch keytype {
		case "string":
			MigrateStringKey(sourceCon, targetCon, ctx, key)
			stringCount++
			if stringCount%ScanCount == 0 {
				fmt.Println("已同步String类型key:", stringCount)
			}
		case "list":
			MigrateListKey(sourceCon, targetCon, ctx, key)
			listCount++
			if listCount%ScanCount == 0 {
				fmt.Println("已同步List类型key:", listCount)
			}
		case "set":
			MigrateSetKey(sourceCon, targetCon, ctx, key)
			setCount++
			if setCount%ScanCount == 0 {
				fmt.Println("已同步set类型key:", setCount)
			}
		case "zset":
			MigrateZsetKey(sourceCon, targetCon, ctx, key)
			zsetCount++
			if zsetCount%ScanCount == 0 {
				fmt.Println("已同步zset类型key:", zsetCount)
			}
		case "hash":
			MigrateHashKey(sourceCon, targetCon, ctx, key)
			hashCount++
			if hashCount%ScanCount == 0 {
				fmt.Println("已同步hash类型key:", hashCount)
			}
		}
	}
	fmt.Println("已同步String类型key:", stringCount)
	fmt.Println("已同步List类型key:", listCount)
	fmt.Println("已同步set类型key:", setCount)
	fmt.Println("已同步zset类型key:", zsetCount)
	fmt.Println("已同步hash类型key:", hashCount)
}
