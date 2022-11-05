package orm

import (
	"gorm.io/gorm"
)

type Backup struct {
	gorm.Model
	UserID   string `gorm:"index:backup_user_id"`
	Hash     string `gorm:"hash"`
	S3Path   string `gorm:"s3_path"`
	IsActive bool   `gorm:"-"`
}
