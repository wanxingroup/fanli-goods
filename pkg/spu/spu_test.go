package spu

import (
	"context"
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/cache"
	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	stock_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
)

func TestCreate(t *testing.T) {
	type args struct {
		spu *product_model.Spu
	}
	tests := []struct {
		name    string
		args    args
		wantSpu *product_model.Spu
		wantSku *product_model.Sku
		wantErr error
	}{
		{
			name: "test_cate_spu",
			args: args{
				spu: &product_model.Spu{
					ShopId:            100001,
					CategoryId:        200001,
					Name:              "title1",
					BannerImages:      "bannerImg1,bannerImg2",
					DetailText:        "detail text1",
					Price:             10000,
					Stock:             10,
					Unit:              "个",
					Barcode:           "abc",
					Support7DayReturn: true,
					Status:            product_model.SPUStatusOffline,
					Skus: []*product_model.Sku{
						{
							ShopId:        100001,
							CategoryId:    200001,
							Price:         10000,
							OriginalPrice: 10000,
							Stock:         10,
							Status:        product_model.Sku_Status_Normal,
						},
					},
				},
			},
			wantSpu: &product_model.Spu{
				ShopId:            100001,
				CategoryId:        200001,
				Name:              "title1",
				BannerImages:      "bannerImg1,bannerImg2",
				DetailText:        "detail text1",
				Price:             10000,
				Stock:             10,
				Unit:              "个",
				Barcode:           "abc",
				Support7DayReturn: true,
				Status:            product_model.SPUStatusOffline,
			},
			wantSku: &product_model.Sku{
				ShopId:        100001,
				CategoryId:    200001,
				Price:         10000,
				OriginalPrice: 10000,
				Stock:         10,
				Status:        product_model.Sku_Status_Normal,
			},
			wantErr: nil,
		},
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Create(tt.args.spu)
			if !assert.EqualValues(t, tt.wantErr, err) {
				return
			}

			actualSpu := new(product_model.Spu)
			err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSpu, tt.args.spu.ID).Error
			if !assert.Nil(t, err) {
				return
			}

			tt.wantSpu.ID = actualSpu.ID
			tt.wantSpu.CreatedAt = actualSpu.CreatedAt
			tt.wantSpu.UpdatedAt = actualSpu.UpdatedAt
			if !assert.EqualValues(t, tt.wantSpu, actualSpu) {
				return
			}

			actualSku := new(product_model.Sku)
			err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSku, tt.args.spu.Skus[0].ID).Error
			if !assert.Nil(t, err) {
				return
			}

			tt.wantSku.ID = actualSku.ID
			tt.wantSku.SpuId = actualSpu.ID
			tt.wantSku.CreatedAt = actualSku.CreatedAt
			tt.wantSku.UpdatedAt = actualSku.UpdatedAt
			if !assert.EqualValues(t, tt.wantSku, actualSku) {
				return
			}

			stock, err := stock_service.GetRealStock(actualSku.SpuId, actualSku.ID)
			if !assert.Nil(t, err) {
				return
			}
			if !assert.EqualValues(t, tt.args.spu.Stock, stock) {
				return
			}
			redisClient.Del(context.Background(), stock_service.GetSkuStockKey(actualSku.SpuId, actualSku.ID))
		})
	}
}

func TestSaveWithoutStock(t *testing.T) {
	existSpu := &product_model.Spu{
		ShopId:            100001,
		CategoryId:        200001,
		Name:              "title1",
		BannerImages:      "bannerImg1,bannerImg2",
		DetailText:        "detail text1",
		Price:             10000,
		Stock:             10,
		Unit:              "个",
		Barcode:           "abc",
		Support7DayReturn: true,
		Status:            product_model.SPUStatusOffline,
		Skus: []*product_model.Sku{
			{
				ShopId:        100001,
				CategoryId:    200001,
				Price:         10000,
				OriginalPrice: 10000,
				Stock:         10,
				Status:        product_model.Sku_Status_Normal,
			},
		},
	}
	err := Create(existSpu)
	if !assert.Nil(t, err) {
		return
	}

	existSpu.Name = "title changed"
	existSpu.Price = 12000
	existSpu.Stock = 20
	existSpu.Unit = "台"
	err = Save(existSpu)
	if !assert.Nil(t, err) {
		return
	}

	actualSpu := new(product_model.Spu)
	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSpu, existSpu.ID).Error
	if !assert.Nil(t, err) {
		return
	}

	wantSpu := &product_model.Spu{
		ShopId:            100001,
		CategoryId:        200001,
		Name:              "title changed",
		BannerImages:      "bannerImg1,bannerImg2",
		DetailText:        "detail text1",
		Price:             12000,
		Stock:             10,
		Unit:              "台",
		Barcode:           "abc",
		Support7DayReturn: true,
		Status:            product_model.SPUStatusOffline,
	}
	wantSpu.ID = existSpu.ID
	wantSpu.CreatedAt = actualSpu.CreatedAt
	wantSpu.UpdatedAt = actualSpu.UpdatedAt
	if !assert.EqualValues(t, wantSpu, actualSpu) {
		return
	}

	actualSku := new(product_model.Sku)
	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSku, existSpu.Skus[0].ID).Error
	if !assert.Nil(t, err) {
		return
	}

	wantSku := &product_model.Sku{
		ShopId:        100001,
		CategoryId:    200001,
		Price:         10000,
		OriginalPrice: 10000,
		Stock:         10,
		Status:        product_model.Sku_Status_Normal,
	}
	wantSku.ID = actualSku.ID
	wantSku.CreatedAt = actualSku.CreatedAt
	wantSku.UpdatedAt = actualSku.UpdatedAt
	if !assert.EqualValues(t, wantSpu, actualSpu) {
		return
	}

	stock, err := stock_service.GetRealStock(actualSku.SpuId, actualSku.ID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, wantSku.Stock, stock) {
		return
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	redisClient.Del(context.Background(), stock_service.GetSkuStockKey(actualSku.SpuId, actualSku.ID))
}

func TestSetStock(t *testing.T) {
	existSpu := &product_model.Spu{
		ShopId:            100001,
		CategoryId:        200001,
		Name:              "title1",
		BannerImages:      "bannerImg1,bannerImg2",
		DetailText:        "detail text1",
		Price:             10000,
		Stock:             10,
		Unit:              "个",
		Barcode:           "abc",
		Support7DayReturn: true,
		Status:            product_model.SPUStatusOffline,
		Skus: []*product_model.Sku{
			{
				ShopId:        100001,
				CategoryId:    200001,
				Price:         10000,
				OriginalPrice: 10000,
				Stock:         10,
				Status:        product_model.Sku_Status_Normal,
			},
		},
	}
	err := Create(existSpu)
	if !assert.Nil(t, err) {
		return
	}

	actualSpu := new(product_model.Spu)
	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSpu, existSpu.ID).Error
	if !assert.Nil(t, err) {
		return
	}

	if !assert.EqualValues(t, existSpu.Stock, actualSpu.Stock) {
		return
	}

	actualSku := new(product_model.Sku)
	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSku, existSpu.Skus[0].ID).Error
	if !assert.Nil(t, err) {
		return
	}

	if !assert.EqualValues(t, existSpu.Skus[0].Stock, actualSku.Stock) {
		return
	}

	stock, err := stock_service.GetRealStock(actualSku.SpuId, actualSku.ID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, existSpu.Skus[0].Stock, stock) {
		return
	}

	changedStock := uint64(30)
	err = SetStock(existSpu.ID, changedStock)
	if !assert.Nil(t, err) {
		return
	}

	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSpu, existSpu.ID).Error
	if !assert.Nil(t, err) {
		return
	}

	if !assert.EqualValues(t, changedStock, actualSpu.Stock) {
		return
	}

	err = database.GetDB(constants.DatabaseConfigKey).Find(&actualSku, existSpu.Skus[0].ID).Error
	if !assert.Nil(t, err) {
		return
	}

	if !assert.EqualValues(t, changedStock, actualSku.Stock) {
		return
	}

	stock, err = stock_service.GetRealStock(actualSku.SpuId, actualSku.ID)
	if !assert.Nil(t, err) {
		return
	}
	if !assert.EqualValues(t, changedStock, stock) {
		return
	}

	redisClient, err := cache.Get(constants.RedisConfigKey)
	assert.Nil(t, err, "can not get redis client")
	if err != nil {
		return
	}

	redisClient.Del(context.Background(), stock_service.GetSkuStockKey(actualSku.SpuId, actualSku.ID))
}
