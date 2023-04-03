package dao

import (
	"database/sql"

	"github.com/baiyz0825/corp-webot/utils/xlog"
)

const (
	DB_NAME = "./data.db"
	DB_TYPE = "sqlite3"
)

var DB *sql.DB

func init() {
	xlog.Log.Info("初始化Db.....")
	// 打开数据库
	DB, err := sql.Open(DB_TYPE, DB_NAME)
	if err != nil {
		xlog.Log.WithField("初始化DB:", "打开数据库失败").Fatal(err)
	}
	// 创建表
	if err := createTable(DB); err != nil {
		xlog.Log.WithField("初始化DB:", "创建表失败").Fatal(err)
	}
}

// CloseDb
// @Description: 关闭db
func CloseDb() {
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(DB)
}

// createTable
// @Description: 创建用户表
// @param db
// @return error
func createTable(db *sql.DB) error {
	sql := `create table if not exists 
    			"users" (
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
