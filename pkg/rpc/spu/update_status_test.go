package spu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/rpc/protos"
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/spu"
)

func TestController_UpdateStatus(t *testing.T) {

	tests := []struct {
		initData  func() error
		input     *protos.UpdateStatusRequest
		wantError *protos.Error
	}{
		{
			initData: func() error {

				return spu.Create(
					&product_model.Spu{
						BaseModel: databases.BaseModel{
							ID: 100003,
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
					},
				)
			},
			input: &protos.UpdateStatusRequest{
				ShopId: 10000,
				SpuId:  100003,
				Status: uint64(product_model.SPUStatusOnline),
			},
		},
	}

	for _, test := range tests {

		if test.initData != nil {
			assert.Nil(t, test.initData(), test.input)
		}

		c := &Controller{}
		reply, err := c.UpdateStatus(context.Background(), test.input)
		assert.Nil(t, err)
		if test.wantError != nil || reply.Err != nil {
			assert.Equal(t, test.wantError, reply.Err)
			continue
		}
	}
}
