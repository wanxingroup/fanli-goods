package product

import "dev-gitlab.wanxingrowth.com/fanli/goods/v2/pkg/databases"

type Category struct {
	databases.BaseModel
	ParentId    uint64 `gorm:"column:parentId;type:bigint unsigned;comment:'父类目ID'"`
	ShopId      uint64 `gorm:"column:shopId;type:bigint unsigned;not null;comment:'店铺ID'"`
	Name        string `gorm:"column:name;type:varchar(30);not null;comment:'名称'"`
	Key         string `gorm:"column:key;type:varchar(30);default:'';comment:'标记值'"`
	Level       uint8  `gorm:"column:level;type:tinyint;default:'1';comment:'层级'"`
	Description string `gorm:"column:desc;type:varchar(255);default:'';comment:'描述'"`
	Path        string `gorm:"column:path;type:varchar(255);default:'';comment:'结点路径，逗号分隔'"`
	Unit        string `gorm:"column:unit;type:varchar(4);default:'';comment:'类目下商品单位'"`
	Sort        uint16 `gorm:"column:sort;type:smallint;not null;default:'0';comment:'排序'"`
	Status      uint8  `gorm:"column:status;type:tinyint;not null;default:'1';comment:'1: 正常, 9: 删除'"`
}

func (c *Category) TableName() string {
	return "prd_category"
}

const (
	Category_Status_Normal = uint8(1)
	Category_Status_Del    = uint8(9)
)
