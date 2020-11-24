package databases

import "time"

type BaseModel struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment:false;comment:'主键ID'"`
	CreatedBy uint64    `gorm:"column:createdBy;type:bigint unsigned;comment:'创建人ID'"`
	CreatedAt time.Time `gorm:"column:createdAt;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:'创建时间'"`
	UpdatedBy uint64    `gorm:"column:updatedBy;type:bigint unsigned;comment:'更新人ID'"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:'更新时间'"`
}
