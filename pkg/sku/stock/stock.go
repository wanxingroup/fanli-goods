package stock

import (
	"context"
	"fmt"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/cache"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
)

type OutOfStockError struct {
	SkuId uint64
}

func (e *OutOfStockError) Error() string {
	return fmt.Sprintf("sku[%v] out of stock", e.SkuId)
}

func GetRealStock(spuId, skuId uint64) (int, error) {
	redisClient, err := cache.Get(constants.RedisConfigKey)
	if err != nil {
		return 0, err
	}
	return redisClient.Get(context.Background(), GetSkuStockKey(spuId, skuId)).Int()
}

func SetRealStock(spuId, skuId uint64, stock uint64) error {
	redisClient, err := cache.Get(constants.RedisConfigKey)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), GetSkuStockKey(spuId, skuId), stock, 0).Err()
}

func SetRealStocks(spuId uint64, skuStocks map[uint64]uint64) error {
	if len(skuStocks) == 0 {
		return nil
	}

	keyValues := make([]interface{}, 0, len(skuStocks))

	for skuId, stock := range skuStocks {
		keyValues = append(keyValues, GetSkuStockKey(spuId, skuId), stock)
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	if err != nil {
		return err
	}

	return redisClient.MSet(context.Background(), keyValues...).Err()
}

func DecrRealStock(spuId, skuId uint64, count uint64) (uint, error) {

	redisClient, err := cache.Get(constants.RedisConfigKey)
	if err != nil {
		return 0, err
	}
	newStock, err := CounterScript.Run(context.Background(), redisClient, []string{GetSkuStockKey(spuId, skuId)}, count).Int()

	if err != nil {
		return 0, err
	}

	if newStock < 0 {
		return 0, &OutOfStockError{SkuId: skuId}
	}

	return uint(newStock), nil
}

func IncrRealStock(spuId, skuId uint64, count uint64) (uint, error) {

	redisClient, err := cache.Get(constants.RedisConfigKey)
	if err != nil {
		return 0, err
	}

	newStock, err := redisClient.IncrBy(context.Background(), GetSkuStockKey(spuId, skuId), int64(count)).Result()
	if err != nil {
		return 0, err
	}
	return uint(newStock), nil
}

func GetSkuStockKey(spuId, skuId uint64) string {
	return fmt.Sprintf(constants.RedisKeySKUStock, spuId, skuId)
}
