package core

import (
	"os"
	"time"

	"github.com/d1vshar/splitgo/dao"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App interface {
	DB() *dao.AppDao
	Logger() *zap.SugaredLogger
}

type BaseApp struct {
	dao           *dao.AppDao
	sugaredLogger *zap.SugaredLogger
}

func NewApp() *BaseApp {
	db := ConnectDB()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	return &BaseApp{
		dao:           dao.NewAppDao(db),
		sugaredLogger: sugar,
	}
}

func ConnectDB() *gorm.DB {
	DbUrl := os.Getenv("DB_URL")
	pgDb, pgErr := gorm.Open(postgres.Open(DbUrl), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if pgErr != nil {
		panic(pgErr)
	}

	db, dbErr := pgDb.DB()

	if dbErr != nil {
		panic(dbErr)
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(time.Duration(1) * time.Hour)

	return pgDb
}

func (a *BaseApp) Logger() *zap.SugaredLogger {
	return a.sugaredLogger
}

func (a *BaseApp) DB() *dao.AppDao {
	return a.dao
}

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
	CreatedBy *string    `json:"createdBy,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

type Tabler interface {
	TableName() string
}
