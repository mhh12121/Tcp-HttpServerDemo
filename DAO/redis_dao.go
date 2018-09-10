package dao

import (
	"fmt"
	"log"
	"time"

	"../Util"
	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     Util.RedisAddr,
		Password: "",
		DB:       0,
		PoolSize: 300,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		//do nothing
		// panic(err)
	}

	fmt.Println("initialize redis:", pong)
}

// InvaildCache
func InvalidCache(username string, token string) error {
	//todo
	// dont know the return value
	errinfo := client.HSet(username, "valid", "0").Err()
	if errinfo != nil {

		return errinfo
	}
	client.Del(tokenFormat(username))
	return nil
}

// SetToken
func SetToken(username string, token string, expiration int64) error {
	err := client.Set(tokenFormat(username), token, time.Duration(expiration)).Err()
	if err != nil {
		return err
	}
	return nil
}

// CheckToken
func CheckToken(username string, token string) (bool, error) {
	val, err := client.Get(tokenFormat(username)).Result()
	if err != nil {
		return false, err
	}

	return token == val, nil
}

//not used
func SaveCacheInfo(username string, nickname string, avatar string) bool {
	tmp := map[string]interface{}{
		"valid":    "1",
		"nickname": nickname,
		"avatar":   avatar,
	}

	err := client.HMSet(username, tmp).Err()
	if err != nil {
		fmt.Println("redis save cache fail:", err)
		return false
	}
	// client.Save()
	return true
}

// CacheInfo
//not used
func GetCacheInfo(username string) (*Util.RealUser, bool, error) {
	val, err := client.HGetAll(username).Result()
	log.Println("val", val)
	log.Println("redis val", val)
	if err != nil {
		return nil, false, err
	}
	if val["valid"] != "0" && len(val) != 0 {
		tmpuser := &Util.RealUser{Username: username, Avatar: val["avatar"], Nickname: val["nickname"]}
		return tmpuser, true, nil
	}
	return nil, false, err
}

//not used
func UpdateCacheNickname(username string, nickname string) error {
	row := map[string]interface{}{
		"valid":    "1",
		"nickname": nickname,
	}
	err := client.HMSet(username, row).Err()
	if err != nil {
		return err
	}
	return nil
}

// update avatar
//not used
func UpdateCacheAvatar(username string, avatar string) error {
	row := map[string]interface{}{
		"valid":  "1",
		"avatar": avatar,
	}
	err := client.HMSet(username, row).Err()
	if err != nil {
		return err
	}
	return nil
}

func tokenFormat(username string) string {
	return "token_" + username
}
