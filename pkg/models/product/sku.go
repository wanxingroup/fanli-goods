package product

import (
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
)

type Sku struct {
	databases.BaseModel
	SpuId         uint64 `gorm:"column:spuId;type:bigint unsigned;not null;index:spuId;comment:'SPU ID'"`
	ShopId        uint64 `gorm:"column:shopId;type:bigint unsigned;not null;comment:'商户ID'"`
	CategoryId    uint64 `gorm:"column:categoryId;type:bigint unsigned;not null;comment:'类目ID'"`
	Price         uint64 `gorm:"column:price;type:bigint unsigned;not null;comment:'现价'"`
	VipPrice      uint64 `gorm:"column:vipPrice;type:bigint unsigned;not null;comment:'会员价'"`
	Point         uint64 `gorm:"column:point;type:bigint unsigned;not null;comment:'积分价格'"`
	Name          string `gorm:"column:name;type:varchar(50);not null;default:'';comment:'SKU 名称'"`
	OriginalPrice uint64 `gorm:"column:marketPrice;type:bigint unsigned;comment:'原价'"`
	Stock         uint64 `gorm:"column:stock;type:bigint unsigned;not null;default:'0';comment:'库存'"`
	Attrs         string `gorm:"column:attrs;type:varchar(255);comment:'销售属性串，方便定位sku'"`
	Status        uint8  `gorm:"column:status;type:tinyint unsigned;not null;default:'1';comment:'1: 正常, 9: 删除'"`
}

func (c *Sku) TableName() string {
	return "prd_sku"
}

const (
	Sku_Status_Normal = uint8(1)
	Sku_Status_Del    = uint8(9)
)
