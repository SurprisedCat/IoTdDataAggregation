package database

import (
	"crypto/sha256"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

/*
DataServerWriteAuth write the key into redis and set the expiration time
*/
func DataServerWriteAuth(key, value []byte, expireTime int64) bool {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return false
	}
	defer c.Close()
	_, err = c.Do("SET", key, value, "EX", expireTime)
	if err != nil {
		fmt.Printf("DataServerWriteAuth:%v\n", err)
		return false
	}
	return true
}

/*
DataServerEraseAuth write the key into redis and set the expiration time
*/
func DataServerEraseAuth(key []byte) bool {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return false
	}
	defer c.Close()
	_, err = c.Do("DEL", key)
	if err != nil {
		fmt.Printf("DataServerEraseAuth:%v\n", err)
		return false
	}
	return true
}

/*
DataServerGetKey server gets encrypted key from rediss
*/
func DataServerGetKey(cid []byte) ([]byte, bool) {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return nil, false
	}
	defer c.Close()
	key := sha256.Sum256(cid)
	encryptedKey, err := redis.String(c.Do("GET", key[:]))
	if err != nil {
		fmt.Printf(" DataServerGetKey:%v\n", err)
		return nil, false
	}
	return []byte(encryptedKey), true
}

/*
ConnTest Redis connecte test
*/
func ConnTest() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("conn redis failed,", err)
		return
	}

	defer c.Close()
	_, err = c.Do("MSet", "abc", 100, "efg", 300)
	if err != nil {
		fmt.Printf(" ConnTest:%v\n", err)
		return
	}
	_, err = c.Do("DEL", "efg")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := redis.Ints(c.Do("MGet", "abc", "efg", "hdf"))
	if err != nil {
		fmt.Println("get abc failed,", err)
		return
	}

	for _, v := range r {
		fmt.Println(v)
	}
	_, err = c.Do("expire", "abc", 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	key := sha256.Sum256([]byte("OAI1001"))
	r2, err := c.Do("GET", key[:])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(r2)
	r3, err := c.Do("GET", "edss")
	if err != nil {
		fmt.Println(err)
		return
	}
	if r3 == nil {
		fmt.Println(r3)
	} else {
		fmt.Println("nil is not nil")
	}
}
