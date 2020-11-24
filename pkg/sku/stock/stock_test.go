package stock

import (
	"context"
	"fmt"
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/cache"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
)

func TestGetRealStock(t *testing.T) {
	spuId := uint64(100001)
	skuId := uint64(200001)
	stockSet := uint(100)
	stockKey := GetSkuStockKey(spuId, skuId)
	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	err = redisClient.Set(context.Background(), stockKey, stockSet, 0).Err()
	if !assert.Nil(t, err) {
		return
	}

	defer redisClient.Del(context.Background(), stockKey)

	stockGet, err := GetRealStock(spuId, skuId)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, stockSet, stockGet) {
		return
	}
}

func TestSetRealStock(t *testing.T) {
	spuId := uint64(100002)
	skuId := uint64(200002)
	stockSet := uint64(100)
	err := SetRealStock(spuId, skuId, stockSet)
	if !assert.Nil(t, err) {
		return
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	stockKey := GetSkuStockKey(spuId, skuId)
	defer redisClient.Del(context.Background(), stockKey)

	stockGet, err := redisClient.Get(context.Background(), stockKey).Int()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, stockSet, stockGet) {
		return
	}
}

func TestSetRealStocks(t *testing.T) {
	spuId := uint64(100003)
	skuStocks := map[uint64]uint64{
		200030: 100,
		200031: 200,
	}
	err := SetRealStocks(spuId, skuStocks)
	if !assert.Nil(t, err) {
		return
	}

	var stockKeys []string
	var skuIds []uint64
	for skuId, _ := range skuStocks {
		skuIds = append(skuIds, skuId)
		stockKeys = append(stockKeys, GetSkuStockKey(spuId, skuId))
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	defer func() {
		redisClient.Del(context.Background(), stockKeys...)
	}()

	stocksGet, err := redisClient.MGet(context.Background(), stockKeys...).Result()
	if !assert.Nil(t, err) {
		return
	}

	var stocksSet []uint64
	for _, skuId := range skuIds {
		stocksSet = append(stocksSet, skuStocks[skuId])
	}

	if !assert.EqualValues(t, len(stocksSet), len(stocksGet)) {
		return
	}

	for idx, stockSet := range stocksSet {
		stockGet := stocksGet[idx]
		if !assert.EqualValues(t, fmt.Sprintf("%v", stockSet), stockGet) {
			return
		}
	}
}

func TestIncrRealStock(t *testing.T) {
	spuId := uint64(100004)
	skuId := uint64(200004)
	stockSet := uint64(100)
	err := SetRealStock(spuId, skuId, stockSet)
	if !assert.Nil(t, err) {
		return
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	defer redisClient.Del(context.Background(), GetSkuStockKey(spuId, skuId))

	increment := uint64(2)
	newStock, err := IncrRealStock(spuId, skuId, increment)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, stockSet+increment, newStock) {
		return
	}
}

func TestDecrRealStock(t *testing.T) {
	spuId := uint64(100005)
	skuId := uint64(200005)
	stockSet := uint64(100)
	err := SetRealStock(spuId, skuId, stockSet)

	if !assert.Nil(t, err) {
		return
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	defer redisClient.Del(context.Background(), GetSkuStockKey(spuId, skuId))

	decrement := uint64(2)
	newStock, err := DecrRealStock(spuId, skuId, decrement)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, stockSet-decrement, newStock) {
		return
	}

	newStock, err = DecrRealStock(spuId, skuId, stockSet)
	if !assert.Equal(t, &OutOfStockError{SkuId: skuId}, err) {
		return
	}

	stockGet, err := GetRealStock(spuId, skuId)
	if !assert.Nil(t, err) {
		return
	}

	newStock, err = DecrRealStock(spuId, skuId, uint64(stockGet))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, 0, newStock) {
		return
	}
}
