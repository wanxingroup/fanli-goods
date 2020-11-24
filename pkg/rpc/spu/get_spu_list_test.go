package spu

import (
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
	productModel "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
)

func TestController_GetSPUList(t *testing.T) {

	tests := []struct {
		initData  func() error
		input     *protos.GetSPUListRequest
		want      []*protos.SPUInformation
		wantCount uint64
		wantError *protos.Error
	}{
		{
			initData: func() (err error) {
				err = database.GetDB(constants.DatabaseConfigKey).Create(&productModel.Spu{
					BaseModel: databases.BaseModel{
						ID: 1000,
					},
					ShopId:            10001,
					CategoryId:        100000,
					Name:              "打印纸",
					BannerImages:      "/images/a/b/c.jpg,/images/b/c/d.jpg",
					Price:             1000,
					Stock:             5,
					Unit:              "张",
					Barcode:           "123-456-789",
					Support7DayReturn: true,
					Sort:              0,
					Status:            productModel.SPUStatusOnline,
					Skus: []*productModel.Sku{
						{
							BaseModel: databases.BaseModel{
								ID: 100001,
							},
							SpuId:         1000,
							ShopId:        10001,
							CategoryId:    100000,
							Price:         1000,
							Name:          "0.01mm",
							OriginalPrice: 1000,
							Stock:         5,
							Status:        productModel.Sku_Status_Normal,
						},
						{
							BaseModel: databases.BaseModel{
								ID: 100002,
							},
							SpuId:         1000,
							ShopId:        10001,
							CategoryId:    100000,
							Price:         1100,
							Name:          "0.02mm",
							OriginalPrice: 1100,
							Stock:         5,
							Status:        productModel.Sku_Status_Normal,
						},
					},
				}).Error

				if err != nil {
					return
				}

				err = database.GetDB(constants.DatabaseConfigKey).Create(&productModel.Spu{
					BaseModel: databases.BaseModel{
						ID: 1001,
					},
					ShopId:            10001,
					CategoryId:        100000,
					Name:              "照片纸",
					BannerImages:      "/images/a/b/c.jpg,/images/b/c/d.jpg",
					Price:             1100,
					Stock:             5,
					Unit:              "张",
					Barcode:           "123-456-780",
					Support7DayReturn: true,
					Sort:              0,
					Status:            productModel.SPUStatusOnline,
					Skus: []*productModel.Sku{
						{
							BaseModel: databases.BaseModel{
								ID: 100101,
							},
							SpuId:         1001,
							ShopId:        10001,
							CategoryId:    100000,
							Price:         1100,
							Name:          "0.01mm",
							OriginalPrice: 1100,
							Stock:         5,
							Status:        productModel.Sku_Status_Normal,
						},
						{
							BaseModel: databases.BaseModel{
								ID: 100102,
							},
							SpuId:         1001,
							ShopId:        10001,
							CategoryId:    100000,
							Price:         1200,
							Name:          "0.02mm",
							OriginalPrice: 1200,
							Stock:         5,
							Status:        productModel.Sku_Status_Normal,
						},
						{
							BaseModel: databases.BaseModel{
								ID: 100103,
							},
							SpuId:         1001,
							ShopId:        10001,
							CategoryId:    100000,
							Price:         1200,
							Name:          "0.03mm",
							OriginalPrice: 1300,
							Stock:         5,
							Status:        productModel.Sku_Status_Del,
						},
					},
				}).Error

				if err != nil {
					return
				}

				return
			},
			input: &protos.GetSPUListRequest{
				ShopId:   10001,
				Page:     1,
				PageSize: 20,
			},
			want: []*protos.SPUInformation{
				{
					ShopId:                 10001,
					SpuId:                  1000,
					CategoryId:             100000,
					Name:                   "打印纸",
					Price:                  1000,
					Unit:                   "张",
					Images:                 []string{"/images/a/b/c.jpg", "/images/b/c/d.jpg"},
					DetailText:             "",
					Stock:                  5,
					Support7DayReturnGoods: true,
					Barcode:                "123-456-789",
					IsOnline:               true,
					SkuList: []*protos.SKUInformationStruct{
						{
							SkuId:         100001,
							OriginalPrice: 1000,
							Price:         1000,
							Name:          "0.01mm",
							Stock:         5,
						},
						{
							SkuId:         100002,
							OriginalPrice: 1100,
							Price:         1100,
							Name:          "0.02mm",
							Stock:         5,
						},
					},
				},
				{
					ShopId:                 10001,
					SpuId:                  1001,
					CategoryId:             100000,
					Name:                   "照片纸",
					Price:                  1100,
					Unit:                   "张",
					Images:                 []string{"/images/a/b/c.jpg", "/images/b/c/d.jpg"},
					DetailText:             "",
					Stock:                  5,
					Support7DayReturnGoods: true,
					Barcode:                "123-456-780",
					IsOnline:               true,
					SkuList: []*protos.SKUInformationStruct{
						{
							SkuId:         100101,
							OriginalPrice: 1100,
							Price:         1100,
							Name:          "0.01mm",
							Stock:         5,
						},
						{
							SkuId:         100102,
							OriginalPrice: 1200,
							Price:         1200,
							Name:          "0.02mm",
							Stock:         5,
						},
					},
				},
			},
			wantCount: 2,
		},
	}

	for _, test := range tests {

		if test.initData != nil {
			assert.Nil(t, test.initData(), test.input)
		}

		c := &Controller{}
		reply, err := c.GetSPUList(context.Background(), test.input)

		assert.Nil(t, err)

		if test.wantError != nil {
			assert.Equal(t, test.wantError, reply.Err, test.input)
			continue
		} else {
			assert.Nil(t, reply.Err, test.input)
		}

		assert.Equal(t, reply.Count, test.wantCount, test.input)

		assert.NotNil(t, reply.Spus, test.input)

		if reply.Spus == nil {
			continue
		}

		for _, spu := range reply.Spus {

			spu.CreateTime = 0
			spu.UpdateTime = 0
		}

		assert.Equal(t, test.want, reply.Spus, test.input)
	}
}
