package util

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

//设置key的过期时间
func SetKeyExpireTime(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	keyTTL := sourceCon.PTTL(ctx, key).Val()
	if keyTTL.Nanoseconds() != -1 {
		targetCon.Expire(ctx, key, keyTTL)
	}
}

//删除目标端list已经存在的key
func DelListKey(targetCon *redis.Client, ctx context.Context, key string) {
	var i int64
	isExist := targetCon.Exists(ctx, key).Val()
	if isExist == 1 {
		totalCount := targetCon.LLen(ctx, key).Val()
		for i = 0; i < totalCount; i++ {
			targetCon.LPop(ctx, key)
		}
		targetCon.Del(ctx, key)
	}
}

//删除目标端Set已经存在的key
func DelSetKey(targetCon *redis.Client, ctx context.Context, key string) {
	isExist := targetCon.Exists(ctx, key).Val()
	if isExist == 1 {
		iter := targetCon.SScan(ctx, key, 0, "*", ScanCount).Iterator()
		for iter.Next(ctx) {
			targetCon.SRem(ctx, key, iter.Val())
		}
		targetCon.Del(ctx, key)
	}

}

//删除目标端ZSet已经存在的key
func DelZSetKey(targetCon *redis.Client, ctx context.Context, key string) {
	var t_count int64 = 0
	isExist := targetCon.Exists(ctx, key).Val()
	if isExist == 1 {
		iter := targetCon.ZScan(ctx, key, 0, "*", ScanCount).Iterator()
		for iter.Next(ctx) {
			t_count++
			if t_count%2 == 1 {
				targetCon.ZRem(ctx, key, iter.Val())
			}
		}
		targetCon.Del(ctx, key)
	}
}

//删除目标端hash已经存在的key
func DelHashKey(targetCon *redis.Client, ctx context.Context, key string) {
	var t_count int64 = 0
	isExist := targetCon.Exists(ctx, key).Val()
	if isExist == 1 {
		iter := targetCon.HScan(ctx, key, 0, "*", ScanCount).Iterator()
		for iter.Next(ctx) {
			t_count++
			if t_count%2 == 1 {
				targetCon.HDel(ctx, key, iter.Val())
			}
		}
		targetCon.Del(ctx, key)
	}
}

//迁移string类型的key
func MigrateStringKey(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	keyValue := sourceCon.Get(ctx, key).Val()
	//获取剩余过期时间
	keyTTL := sourceCon.PTTL(ctx, key).Val()
	//如果为-1，则key不过期
	if keyTTL.Nanoseconds() == -1 {
		targetCon.Set(ctx, key, keyValue, 0)
	} else {
		targetCon.Set(ctx, key, keyValue, keyTTL)
	}
}

//迁移List类型的key
func MigrateListKey(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	var i int64
	//删除key
	DelListKey(targetCon, ctx, key)
	totalCount := sourceCon.LLen(ctx, key).Val()
	for i = totalCount - 1; i >= 0; i-- {
		value := sourceCon.LRange(ctx, key, i, i).Val()
		targetCon.LPush(ctx, key, value)
	}
	SetKeyExpireTime(sourceCon, targetCon, ctx, key)
}

//迁移set类型的key
func MigrateSetKey(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	//删除key
	DelSetKey(targetCon, ctx, key)
	iter := sourceCon.SScan(ctx, key, 0, "*", ScanCount).Iterator()
	for iter.Next(ctx) {
		targetCon.SAdd(ctx, key, iter.Val())
	}
	SetKeyExpireTime(sourceCon, targetCon, ctx, key)
}

//迁移Zset类型的key
func MigrateZsetKey(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	var t_count int64 = 0
	var value string
	//删除key
	DelZSetKey(targetCon, ctx, key)
	iter := sourceCon.ZScan(ctx, key, 0, "*", ScanCount).Iterator()
	for iter.Next(ctx) {
		t_count++
		if t_count%2 == 1 {
			value = iter.Val()
		} else {
			score, _ := strconv.ParseFloat(iter.Val(), 64)
			ranking := []*redis.Z{
				&redis.Z{Score: score, Member: value},
			}
			targetCon.ZAdd(ctx, key, ranking...)
		}
	}
	SetKeyExpireTime(sourceCon, targetCon, ctx, key)
}

//迁移Hash类型的key
func MigrateHashKey(sourceCon *redis.Client, targetCon *redis.Client, ctx context.Context, key string) {
	var t_count int64 = 0
	var keystr string
	//删除key
	DelHashKey(targetCon, ctx, key)
	iter := sourceCon.HScan(ctx, key, 0, "*", ScanCount).Iterator()
	for iter.Next(ctx) {
		t_count++
		if t_count%2 == 1 {
			keystr = iter.Val()
		} else {
			targetCon.HSet(ctx, key, keystr, iter.Val())
		}
	}
	SetKeyExpireTime(sourceCon, targetCon, ctx, key)
}
