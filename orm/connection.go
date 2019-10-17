package orm

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/url"
	"runtime"
)

func mysqlConn(c *config) *sql.DB {
	sqlMeta := c.sqlMeta()
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		sqlMeta.User, sqlMeta.Password, sqlMeta.Host, sqlMeta.Port, sqlMeta.Database)

	db, err := sql.Open("mysql", mysqlURI)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return db
}

func mssqlConn(c *config) *sql.DB {
	u := &url.URL{
		Scheme:     "sqlserver",
		User:       url.UserPassword(c.Mssql.User, c.Mssql.Password),
		Host:       fmt.Sprintf("%s:%d", c.Mssql.Host, c.Mssql.Port),
		ForceQuery: true,
		RawQuery:   fmt.Sprintf("database=%s", c.Mssql.Database),
	}
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	return db
}

func redisConn(passwd string) *redis.Client {
	var redisAddress string

	if runtime.GOOS != "linux" {
		redisAddress = "localhost:6379"
		passwd = ""
	} else {
		redisAddress = "redis:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: passwd,
		DB:       0,
	})

	return client
}

var (
	MssqlDB   *sql.DB
	MysqlDB   *sql.DB
	RedisConn *redis.Client
	_config    *config
)

func init() {
	_config := GetConfig()
	MssqlDB = mssqlConn(_config)
	MysqlDB = mysqlConn(_config)
	RedisConn = redisConn(_config.Secret)
}
