package db

import (
	"gorm.io/gorm"
)

func NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func WithDeleted(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}

func ByID(id uint64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func ByField(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" = ?", value)
	}
}

func KeywordInFields(keyword string, fields ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(fields) == 0 {
			return db
		}
		query := db
		for i, field := range fields {
			if i == 0 {
				query = query.Where(field+" LIKE ?", "%"+keyword+"%")
			} else {
				query = query.Or(field+" LIKE ?", "%"+keyword+"%")
			}
		}
		return query
	}
}

func OrderBy(field string, desc bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if desc {
			return db.Order(field + " DESC")
		}
		return db.Order(field + " ASC")
	}
}

func In(field string, values ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}
		return db.Where(field+" IN (?)", values)
	}
}

func NotIn(field string, values ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(values) == 0 {
			return db
		}
		return db.Where(field+" NOT IN (?)", values)
	}
}

func GTE(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" >= ?", value)
	}
}

func GT(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" > ?", value)
	}
}

func LTE(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" <= ?", value)
	}
}

func LT(field string, value interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" < ?", value)
	}
}

func Between(field string, min, max interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(field+" BETWEEN ? AND ?", min, max)
	}
}

func DateRange(field string, start, end string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start != "" && end != "" {
			return db.Where(field+" BETWEEN ? AND ?", start, end)
		}
		if start != "" {
			return db.Where(field+" >= ?", start)
		}
		if end != "" {
			return db.Where(field+" <= ?", end)
		}
		return db
	}
}
