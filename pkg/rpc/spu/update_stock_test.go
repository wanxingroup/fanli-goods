package spu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	stock_service "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/sku/stock"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
)

func TestController_UpdateStock(t *testing.T) {

	tests := []struct {
		initData    func() error
		input       *protos.UpdateStockRequest
		want        *protos.UpdateStockReply
		wantError   *protos.Error
		returnStock func() (int, error)
	}{
		{
			initData: func() error {

				return spu.Create(
					&product_model.Spu{
						BaseModel: databases.BaseModel{
							ID: 100002,
						},
						ShopId:            10000,
						CategoryId:        10001,
						DetailText:        "non test",
						Price:             200,
						Stock:             0,
						Unit:              "份",
						Barcode:           "123-456-789",
						Support7DayReturn: false,
						Status:            product_model.SPUStatusOffline,
						Skus: []*product_model.Sku{
							{
								BaseModel: databases.BaseModel{
									ID: 10000201,
								},
								SpuId:         100002,
								ShopId:        10000,
								CategoryId:    10001,
								Price:         200,
								OriginalPrice: 200,
								Stock:         0,
							},
						},
					},
				)
			},
			input: &protos.UpdateStockRequest{
				ShopId: 10000,
				SpuId:  100002,
				Stock:  30,
			},
			want: &protos.UpdateStockReply{
				Stock: 30,
			},
			returnStock: func() (int, error) {
				return stock_service.GetRealStock(100002, 10000201)
			},
		},
	}

	for _, test := range tests {

		if test.initData != nil {
			assert.Nil(t, test.initData(), test.input)
		}

		c := &Controller{}
		reply, err := c.UpdateStock(context.Background(), test.input)
		assert.Nil(t, err, test.input)
		if reply.Err != nil {
			assert.Equal(t, test.wantError, reply.Err, test.input)
			continue

		} else if test.wantError != nil {

			assert.Equal(t, test.wantError, reply.Err, test.input)
			continue
		}

		assert.Equal(t, test.want, reply)
		if test.returnStock != nil {
			stock, getStockError := test.returnStock()
			assert.Nil(t, getStockError, test.input)
			assert.Equal(t, test.want.GetStock(), uint64(stock))
		}
	}
}
