package RedisConnector

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.com/capslock-ltd/reconciler/backend-golang/shared/Constants"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr:     Constants.REDIS_CONNECTION_STRING,
	Password: Constants.REDIS_PASSWORD,
	DB:       0, // use default DB
})

//SaveItemInRedis saves an Item in the redis DataStore
func SaveItemInRedis(key string, jsonValue string) error {

	err := rdb.Set(ctx, key, jsonValue, 0).Err()

	return err
}

//GetItemFromRedis retrieves the item with the key specified from Redis
func GetItemFromRedis(key string) (string, error) {

	rawValue, err := rdb.Get(ctx, key).Result()

	//value is there but is null
	if err == redis.Nil {
		return "", nil
	}

	//error on getting the value
	if err != nil {
		return "", err
	}

	//item found

	return rawValue, nil
}

//RemoveItemFromRedis removes items with the specified key from redis
func RemoveItemFromRedis(key string) error {

	_, err := rdb.Del(ctx, key).Result()

	return err
}
