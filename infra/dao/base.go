package dao

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite 驱动
)

const (
	ID = "id"
)

var (
	once       sync.Once
	dbInstance *gorm.DB
)

type BaseModel struct {
	ID         int        `json:"id" gorm:"primaryKey,column:id"`
	CreateTime *time.Time `json:"create_time" gorm:"autoCreateTime,column:create_time"`
	UpdateTime *time.Time `json:"update_time" gorm:"autoUpdateTime,column:update_time"`
}

func (b *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	now := time.Now()
	if b.CreateTime == nil {
		scope.SetColumn("CreateTime", &now)
	}
	if b.UpdateTime == nil {
		scope.SetColumn("UpdateTime", &now)
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(scope *gorm.Scope) error {
	now := time.Now()
	scope.SetColumn("UpdateTime", &now)
	return nil
}

// 获取数据库连接
func getDB() *gorm.DB {
	once.Do(func() {
		// db 初始化
		var err error
		dbInstance, err = gorm.Open("sqlite3", "data.db")
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
		// 设置连接池
		dbInstance.DB().SetMaxIdleConns(3)
		dbInstance.DB().SetMaxOpenConns(10)
		// 启用日志
		// dbInstance.LogMode(true)
		// 自动迁移模式
		dbInstance.AutoMigrate(&WfDef{})
		dbInstance.AutoMigrate(&WfIns{})
		dbInstance.AutoMigrate(&WfExec{})
		dbInstance.AutoMigrate(&WfTask{})
		dbInstance.AutoMigrate(&WfParticipant{})
	})
	return dbInstance
}

func WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	return getDB().Transaction(func(tx *gorm.DB) error {
		ctx = WithTx(ctx, tx)
		return fn(ctx, tx)
	})
}

// txKey is a custom type to use as context key to avoid collisions
type txKey struct{}

func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return getDB()
}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func Create[T any](db *gorm.DB, record *T) (int, error) {
	if record == nil {
		return 0, fmt.Errorf("record is nil")
	}
	err := db.Create(record).Error
	if err != nil {
		return 0, err
	}
	// 通过反射获取记录的 ID
	id := reflect.ValueOf(record).Elem().FieldByName("ID").Int()
	return int(id), nil
}

func Update[T any](db *gorm.DB, record *T) error {
	if record == nil {
		return fmt.Errorf("record is nil")
	}
	var t T
	return db.Model(&t).Updates(record).Error
}

func GetByCond[T any](db *gorm.DB, cond map[string]any) ([]*T, error) {
	var (
		t       T
		records []*T
	)
	err := db.Model(&t).Where(cond).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func Delete[T any](db *gorm.DB, id int) error {
	var t T
	return db.Model(&t).Where("id =?", id).Delete(&t).Error
}
