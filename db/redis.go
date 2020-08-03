package db

import (
	//"github.com/go-redis/redis"

	"fmt"
	//"github.com/go-redis/redis"

	"github.com/bwmarrin/snowflake"
)

/*
func RedisNewClient() *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     Awsredisurl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong, err)
	return Rdb
	// Output: PONG <nil>
}

//返回一个shorturlval的自增长id 从900000000开始
func RedisKeyId() int64 {
	rdb := RedisNewClient()
	//err := rdb.Set(ctx, "shorturlkey", 900000000, 0).Err()
	shorturlval, err := rdb.Get(ctx, "shorturlkey").Result()
	if err == redis.Nil {
		fmt.Println("没有shorturlkey，开始初始化设置为900000000")
		err = rdb.Set(ctx, "shorturlkey", 900000000, 0).Err()
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("shorturlkey值目前存在为", shorturlval)

	}

	seqval, err := rdb.Incr(ctx, "shorturlkey").Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("resdis Incr自增ID值为", seqval)

	}
	//shorturlval, err = rdb.Get(ctx, "shorturlkey").Result()
	return seqval
}
*/

func snowflakeID() int64 {

	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)
	fmt.Printf("Base2  ID: %s\n", id.Base2())
	fmt.Printf("Base64 ID: %s\n", id.Base64())

	// Print out the ID's timestamp
	fmt.Printf("ID Time  : %d\n", id.Time())

	// Print out the ID's node number
	fmt.Printf("ID Node  : %d\n", id.Node())

	// Print out the ID's sequence number
	fmt.Printf("ID Step  : %d\n", id.Step())

	// Generate and print, all in one.
	fmt.Printf("ID       : %d\n", node.Generate().Int64())
	return node.Generate().Int64()
}

func RedisKeyId() int64 {
	return snowflakeID()
}
