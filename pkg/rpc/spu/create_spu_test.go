package spu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
)

func TestController_CreateSPU(t *testing.T) {

	tests := []struct {
		input     *protos.CreateSPURequest
		want      *protos.SPUInformation
		wantError *protos.Error
	}{
		{
			input: &protos.CreateSPURequest{
				ShopId:                 10000,
				CategoryId:             10001,
				Name:                   "测试商品",
				Price:                  100,
				Unit:                   "件",
				Images:                 []string{"/test/a1.jpg"},
				DetailText:             "test detail text",
				Stock:                  10,
				SettlePrice:            50,
				CostPrice:              10,
				IsHot:                  false,
				Support7DayReturnGoods: false,
				Barcode:                "123-23-43",
			},
			want: &protos.SPUInformation{
				ShopId:                 10000,
				CategoryId:             10001,
				Name:                   "测试商品",
				Price:                  100,
				Unit:                   "件",
				Images:                 []string{"/test/a1.jpg"},
				DetailText:             "test detail text",
				Stock:                  10,
				SettlePrice:            50,
				CostPrice:              10,
				IsHot:                  false,
				Support7DayReturnGoods: false,
				Barcode:                "123-23-43",
				IsOnline:               false,
				CreateTime:             0,
				UpdateTime:             0,
				SkuList: []*protos.SKUInformationStruct{
					{
						OriginalPrice: 100,
						Price:         100,
						Stock:         10,
					},
				},
			},
		},
	}

	for _, test := range tests {

		c := &Controller{}
		reply, err := c.CreateSPU(context.Background(), test.input)
		assert.Nil(t, err)

		if test.wantError != nil {
			assert.Equal(t, test.wantError, reply.Err, test.input)
			continue
		} else {
			assert.Nil(t, reply.Err, test.input)
		}

		assert.NotNil(t, reply.Spu, test.input)

		if reply.Spu == nil {
			continue
		}

		test.want.SpuId = reply.Spu.SpuId
		test.want.CreateTime = reply.Spu.CreateTime
		test.want.UpdateTime = reply.Spu.UpdateTime

		assert.Greater(t, reply.Spu.SpuId, uint64(0), test.input)
		assert.Greater(t, reply.Spu.CreateTime, int64(0), test.input)
		assert.Greater(t, reply.Spu.UpdateTime, int64(0), test.input)

		assert.Len(t, reply.Spu.SkuList, len(test.want.SkuList), test.input)

		if len(reply.Spu.SkuList) != len(test.want.SkuList) {

			continue
		}

		for index, sku := range test.want.SkuList {

			sku.SkuId = reply.Spu.SkuList[index].SkuId
		}
	}
}
