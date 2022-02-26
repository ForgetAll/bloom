package bloom

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"time"
)

var errorNotExists = errors.New("item not in bitset")

type RedisBitSetProvider struct {
	RedisKey    string
	RedisClient *redis.Client
	ExpireTime  time.Duration
}

func (r RedisBitSetProvider) Set(offset uint) error {
	_, err := r.RedisClient.SetBit(r.RedisKey, int64(offset), 1).Result()
	if err != nil {
		return err
	}

	_, err = r.RedisClient.Expire(r.RedisKey, r.ExpireTime).Result()
	return err
}

func (r RedisBitSetProvider) Test(offset uint) (bool, error) {
	bitValue, err := r.RedisClient.GetBit(r.RedisKey, int64(offset)).Result()
	if err != nil {
		return false, err
	}

	return bitValue == 1, nil
}

func (r RedisBitSetProvider) TestBatch(offset []uint) (bool, error) {
	_, err := r.RedisClient.Pipelined(func(pipeliner redis.Pipeliner) error {
		for i := range offset {
			pipeliner.GetBit(r.RedisKey, int64(offset[i]))
		}

		result, err := pipeliner.Exec()
		if err != nil {
			return err
		}

		for i := range result {
			if res, ok := result[i].(*redis.IntCmd); ok {
				bitValue, err := res.Result()
				if err != nil {
					return err
				}

				if bitValue != 1 {
					return errorNotExists
				}
			}
		}
		return err
	})
	if err != nil {
		if err == errorNotExists {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r RedisBitSetProvider) SetBatch(offset []uint) error {
	_, err := r.RedisClient.Pipelined(func(pipeliner redis.Pipeliner) error {
		for i := range offset {
			pipeliner.SetBit(r.RedisKey, int64(offset[i]), 1)
		}

		pipeliner.Expire(r.RedisKey, r.ExpireTime)
		_, err := pipeliner.Exec()
		return err
	})

	return err
}

func (r RedisBitSetProvider) New(_ uint) {
}
