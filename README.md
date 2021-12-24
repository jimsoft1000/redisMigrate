# redisMigrate数据迁移工具

## 功能介绍
Redis数据库迁移工具，主要有以下特性  
1.支持string，list，set，zset，hash类型数据迁移，并设置key的剩余ttl  
2.list，set，zset，hash类型key，使用scan进行迁移，尽可能减少对源端数据库的影响  
3.支持源库和目标库的库不一致迁移，例如：源库为0号库，目标库为1号库  
4.支持迁移任务重跑，任务重跑，相同的key，会先删除目标库已存在的key，list，set，zset，hash类型key，使用scan进行删除 


## 使用说明

#### 迁移指定前缀的key
指定sping*之后，迁移工具只会找到spring开头的key 

`
./redisMigrate -SourcePrefix "spring*" -SourceIP "....rds.aliyuncs.com" -SourcePort 6379 -SourceDB 1 -SourcePasswd "xxxx" -TargetIP "......rds.aliyuncs.com" -TargetPort 6379 -TargetDB 1 -TargetPasswd "....."
`

#### 不同库迁移

`./redisMigrate -SourcePrefix "spring*" -SourceIP "....rds.aliyuncs.com" -SourcePort 6379 -SourceDB 1 -SourcePasswd "xxxx" -TargetIP "......rds.aliyuncs.com" -TargetPort 6379 -TargetDB 2 -TargetPasswd "....."`
