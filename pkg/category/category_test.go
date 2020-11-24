package category

import (
	"fmt"
	"testing"

	"dev-gitlab.wanxingrowth.com/wanxin-go-micro/base/data/database"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/constants"
	product_model "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/models/product"
)

func TestCreateLevel1(t *testing.T) {
	type args struct {
		category *product_model.Category
	}
	tests := []struct {
		name         string
		args         args
		wantCategory *product_model.Category
		wantErr      error
	}{
		{
			name: "create_category_level1",
			args: args{
				&product_model.Category{
					ShopId:      100001,
					Name:        "类目1",
					Key:         "key1",
					Description: "desc1",
					Unit:        "个",
					Sort:        10,
					Status:      product_model.Category_Status_Normal,
				},
			},
			wantCategory: &product_model.Category{
				ParentId:    0,
				ShopId:      100001,
				Name:        "类目1",
				Key:         "key1",
				Level:       1,
				Description: "desc1",
				Path:        ",",
				Unit:        "个",
				Sort:        10,
				Status:      product_model.Category_Status_Normal,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Create(tt.args.category)
			if !assert.Equal(t, tt.wantErr, err) {
				return
			}

			actualCategory := new(product_model.Category)
			err = database.GetDB(constants.DatabaseConfigKey).Find(&actualCategory, tt.args.category.ID).Error
			if !assert.Nil(t, err) {
				return
			}

			tt.wantCategory.ID = actualCategory.ID
			tt.wantCategory.CreatedAt = actualCategory.CreatedAt
			tt.wantCategory.UpdatedAt = actualCategory.UpdatedAt
			if !assert.Equal(t, tt.wantCategory, actualCategory) {
				return
			}
		})
	}
}

func TestCreateLevel2(t *testing.T) {
	type args struct {
		category *product_model.Category
	}
	tests := []struct {
		name         string
		args         args
		wantCategory *product_model.Category
		wantErr      error
	}{
		{
			name: "create_category_level1",
			args: args{
				&product_model.Category{
					ShopId:      100001,
					Name:        "类目2",
					Key:         "key2",
					Description: "desc2",
					Unit:        "个",
					Sort:        20,
					Status:      product_model.Category_Status_Normal,
				},
			},
			wantCategory: &product_model.Category{
				ParentId:    0,
				ShopId:      100001,
				Name:        "类目2",
				Key:         "key2",
				Level:       2,
				Description: "desc2",
				Unit:        "个",
				Sort:        20,
				Status:      product_model.Category_Status_Normal,
			},
			wantErr: nil,
		},
	}

	existCategory := product_model.Category{
		ShopId:      100001,
		Name:        "类目x",
		Key:         "keyx",
		Description: "descx",
		Unit:        "个",
		Sort:        20,
		Status:      product_model.Category_Status_Normal,
	}
	err := Create(&existCategory)
	if !assert.Nil(t, err) {
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.category.ParentId = existCategory.ID
			err := Create(tt.args.category)
			if !assert.Equal(t, tt.wantErr, err) {
				return
			}

			actualCategory := new(product_model.Category)
			err = database.GetDB(constants.DatabaseConfigKey).Find(&actualCategory, tt.args.category.ID).Error
			if !assert.Nil(t, err) {
				return
			}

			tt.wantCategory.ID = actualCategory.ID
			tt.wantCategory.ParentId = existCategory.ID
			tt.wantCategory.Path = fmt.Sprintf("%s%v,", existCategory.Path, existCategory.ID)
			tt.wantCategory.CreatedAt = actualCategory.CreatedAt
			tt.wantCategory.UpdatedAt = actualCategory.UpdatedAt

			if !assert.Equal(t, tt.wantCategory, actualCategory) {
				return
			}
		})
	}
}
