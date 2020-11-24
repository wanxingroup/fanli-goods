package databases

import (
	"github.com/jinzhu/gorm"
)

const MaxPageSize = 100

func FindPage(db *gorm.DB, pageData map[string]uint64, out interface{}) (uint64, error) {

	// 先查询
	pageSize, has := pageData["pageSize"]

	if has {

		if pageSize < 0 {
			return 0, nil
		}

		if pageSize > MaxPageSize {
			db = db.Limit(MaxPageSize)
		} else {
			db = db.Limit(pageSize)
		}

	} else {
		db = db.Limit(MaxPageSize)
	}

	// 不传页码为通过最后一条 ID 查询
	if pageNum, has := pageData["pageNum"]; has {
		if pageNum < 0 {
			return 0, nil
		}

		db = db.Offset((pageNum - 1) * pageSize)

	}

	err := db.Find(out).Error

	if err != nil {
		return 0, err
	}

	var count uint64
	err = db.Count(&count).Error
	return count, err
}
