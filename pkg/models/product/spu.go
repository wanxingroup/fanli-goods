package product

import (
	"dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"
)

type Spu struct {
	databases.BaseModel
	ShopId            uint64 `gorm:"column:shopId;type:bigint unsigned;not null;comment:'店铺ID'"`
	CategoryId        uint64 `gorm:"column:categoryId;type:bigint unsigned;not null;comment:'类目ID'"`
	Name              string `gorm:"column:name;type:varchar(50);not null;comment:'SPU 名称'"`
	BannerImages      string `gorm:"column:bannerImages;type:varchar(255);comment:'主图url列表，逗号分隔'"` // TODO Lucky 拆分这个字段到一个 1 对 n 的关联表
	BannerVideo       string `gorm:"column:bannerVideo;type:varchar(100);comment:'主视频'"`
	DetailText        string `gorm:"column:detailText;type:varchar(3000);not null;default:'';comment:'详情文本'"`
	Price             uint64 `gorm:"column:price;type:bigint unsigned;not null;comment:'一口价'"`
	SettlePrice       uint64 `gorm:"column:settlePrice;type:bigint unsigned;not null;comment:'结算价'"`
	CostPrice         uint64 `gorm:"column:costPrice;type:bigint unsigned;not null;comment:'成本价'"`
	VipPrice          uint64 `gorm:"column:vipPrice;type:bigint unsigned;not null;comment:'会员价'"`
	Point             uint64 `gorm:"column:point;type:bigint unsigned;not null;comment:'积分价格'"`
	Stock             uint64 `gorm:"column:stock;type:bigint unsigned;not null;default:'0';comment:'总库存'"`
	Unit              string `gorm:"column:unit;type:varchar(4);comment:'单位'"`
	Barcode           string `gorm:"column:barcode;type:varchar(100);comment:'第三方关联编码'"`
	Support7DayReturn bool   `gorm:"column:support7Return;type:tinyint unsigned;not null;default:'0';comment:'是否7天无理由退货'"`
	Sort              uint   `gorm:"column:sort;type:smallint unsigned;not null;default:'0';comment:'前端排序'"`
	Status            uint8  `gorm:"column:status;type:tinyint unsigned;not null;default:'1';comment:'1: 下架, 2: 上架, 9: 删除'"`
	IsHot             bool   `gorm:"column:isHot;type:tinyint unsigned;not null;default:'0';comment:'是否热门'"`
	Skus              []*Sku `gorm:"foreignkey:SpuId"`
}

func (c *Spu) TableName() string {
	return "spu"
}

const (
	SPUStatusOffline = uint8(1) // 下架
	SPUStatusOnline  = uint8(2) // 上架
	SPUStatusDelete  = uint8(9) // 删除
)
