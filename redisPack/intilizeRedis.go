package redispack

import (
	"encoding/json"
	"project/internal/model"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
)

func NewRedisClient() *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	return redisDB
}

type Cache interface {
	CheckRedisKey(rdb *redis.Client, key string) (model.Job, error)
	SetRedisKey(rdb *redis.Client, key string, jobData model.Job)
}

type RadisDbConnection struct {
	rdc redis.Client
}

func NewRadisDbConnection(r redis.Client) RadisDbConnection {

	return RadisDbConnection{rdc: r}
}

func (r *RadisDbConnection) CheckRedisKey(key string) (model.Job, error) {

	val, err := r.rdc.Get(key).Result()
	if err == redis.Nil {
		return model.Job{}, err
	}
	var job model.Job
	err = json.Unmarshal([]byte(val), &job)
	if err != nil {
		log.Err(err)
	}
	return job, nil
}

func (r *RadisDbConnection) SetRedisKey(key string, jobData model.Job) {

	jobdata, err := json.Marshal(jobData)
	if err != nil {
		log.Err(err)
		return
	}
	data := string(jobdata)
	err = r.rdc.Set(key, data, 10*time.Minute).Err()
	if err != nil {
		log.Err(err)
		return
	}
}
