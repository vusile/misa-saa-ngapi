package repository

import (
	"gorm.io/gorm"
)

type FindAllPage struct {
	Size    int
	PageNum int
}

func Paginate(pg FindAllPage) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if pg.PageNum <= 0 {
			pg.PageNum = 1
		}

		switch {
		case pg.Size > 100:
			pg.Size = 100
		case pg.Size <= 0:
			pg.Size = 10
		}

		offset := (pg.PageNum - 1) * pg.Size
		return db.Offset(offset).Limit(pg.Size)
	}
}
