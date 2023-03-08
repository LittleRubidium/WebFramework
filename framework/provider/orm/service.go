package orm

import (
	"context"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"sync"
	"time"
)

type HadeGorm struct {
	container framework.Container
	dbs       map[string]*gorm.DB

	lock *sync.RWMutex
}

func NewHadeGorm(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	lock := &sync.RWMutex{}
	return &HadeGorm{
		container: container,
		dbs:       dbs,
		lock:      lock,
	}, nil
}

func (g *HadeGorm) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	logger := g.container.MustMake(contract.LogKey).(contract.Log)

	//读取默认配置
	config := GetBaseConfig(g.container)

	logService := g.container.MustMake(contract.LogKey).(contract.Log)

	//设置logger
	ormLogger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: ormLogger,
	}

	//option对opt进行修改
	for _, opt := range option {
		if err := opt(g.container, config); err != nil {
			return nil, err
		}
	}

	//如果最终的config没有设置dsn，就生成dns
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	//判断是否已经实例化gorm.DB
	g.lock.RLock()
	if db, ok := g.dbs[config.Dsn]; ok {
		g.lock.RUnlock()
		return db, nil
	}
	g.lock.RUnlock()

	g.lock.RLock()
	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}
	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}
	if config.ConnMaxLifetime != "" {
		lifeTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max life time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(lifeTime)
		}
	}
	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}
	if err != nil {
		g.dbs[config.Dsn] = db
	}
	return db, err
}
