package bloom

import "github.com/go-redis/redis/v7"

type RedisBitSetProvider struct {
	redisKey    string
	redisClient *redis.Client
}

func (r RedisBitSetProvider) Set(offset uint) error {
	_, err := r.redisClient.SetBit(r.redisKey, int64(offset), 1).Result()
	return err
}

func (r RedisBitSetProvider) Test(offset uint) (bool, error) {
	bitValue, err := r.redisClient.GetBit(r.redisKey, int64(offset)).Result()
	if err != nil {
		return false, err
	}

	return bitValue == 1, nil
}

func (r RedisBitSetProvider) New(m uint) {
}
