package dao

import (
	"database/sql"
	"os"

	"github.com/baiyz0825/corp-webot/utils/xlog"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_NAME = "/db/data.db"
	DB_TYPE = "sqlite3"
)

var Db *sql.DB

func LoadDatabase() {
	xlog.Log.Info("初始化Db.....")
	WorkPath, _ := os.Getwd()
	// 打开数据库
	db, err := sql.Open(DB_TYPE, WorkPath+DB_NAME)
	if err != nil {
		xlog.Log.WithField("初始化DB:", "打开数据库失败").Fatal(err)
		panic("数据库初始化失败")
	}
	// 测试数据库是否连通
	if err = db.Ping(); err != nil {
		xlog.Log.Fatal(err)
	}
	// 创建表
	if err := createTable(db); err != nil {
		xlog.Log.WithField("初始化DB:", "创建表失败").Fatal(err)
	}
	Db = db
}

// CloseDb
// @Description: 关闭db
func CloseDb() {
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(Db)
}

// createTable
// @Description: 创建用户表
// @param db
// @return error
func createTable(db *sql.DB) error {
	sql := `create table if not exists 
    			"user" (
                    id INTEGER
                        primary key autoincrement,
                    name CHAR(50) not null ,
                    sys_prompt CHAR(512) DEFAULT '',
                    update_time BIGINT
				);
			create table if not exists 
			    "context" (
			    	id INTEGER
			    	    primary key autoincrement,
			    	name CHAR(50) not null ,
			    	context_msg CHAR(512) DEFAULT '',
			    	update_time BIGINT
				);`
	_, err := db.Exec(sql)
	return err
}
