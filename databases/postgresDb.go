package databases

import (
	"fmt"
	"sync"

	"github.com/Witthaya22/golang-backend-itctc/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDb struct {
	*gorm.DB
}

var (
	PostgresDbInstance *postgresDb
	once               sync.Once
)

func NewPostgresDb(conf *config.Database) IDatabase {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			conf.Host,
			conf.Port,
			conf.User,
			conf.Password,
			conf.Dbname,
			conf.Sslmode,
		)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Connected to postgres database: %s", conf.Dbname)

		PostgresDbInstance = &postgresDb{conn}
	})

	return PostgresDbInstance
}

func (db *postgresDb) ConnectionGetting() *gorm.DB {
	return db.DB
}
