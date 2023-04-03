package dao

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/baiyz0825/corp-webot/model"
	"github.com/baiyz0825/corp-webot/utils/xlog"
	"github.com/baiyz0825/corp-webot/xconst"
	"github.com/pkg/errors"
)

type Context struct {
	Id         int64
	Name       string
	ContextMsg string
	UpdateTime int64
}

// SetContext
// @Description: 设置用户上下文对象中上下文
// @receiver c
// @param context
// @return bool
func (c *Context) SetContext(context model.MessageContext) bool {
	data, err := json.Marshal(context)
	if err != nil {
		xlog.Log.WithField("用户:", c.Name).WithError(err).Error("存储上下文序列化错误")
		return false
	}
	c.ContextMsg = string(data)
	return true
}

// InsertUserContext
// @Description: 插入用户上下文
// @receiver c
// @param db
// @return error
func InsertUserContext(name, contextMsg string, db *sql.DB) error {
	query := `insert into context (name, context_msg,update_time) values (?,?,?)`
	prepare, err := db.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	_, err = prepare.Exec(name, contextMsg, time.Now().UnixMilli())
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// UpdateContext
// @Description: 更新用户上下文
// @receiver c
// @param db
// @return error
func UpdateContext(contextMsg, userName string, db *sql.DB) error {
	query := `update context set context_msg=?,update_time=? where id in{
				select id from context where name=? ORDER BY update_time limit 1
				}`
	prepare, err := db.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	_, err = prepare.Exec(contextMsg, time.Now().UnixMilli(), userName)
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// GetLatestUserContext
// @Description:  获取最新用户上下文对象中上下文
// @receiver c
// @param db
// @return *model.MessageContext
// @return error
func GetLatestUserContext(userName string, db *sql.DB) (*Context, error) {
	sql := `select * from context where name=? limit 1 Order By update_time DESC`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库查询语句")
	}
	data, err := stmt.Query(userName)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库语句执行失败")
	}
	if data.Next() {
		temp := Context{}
		err = data.Scan(temp.Id, temp.Name, temp.ContextMsg, temp.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "prepare数据库语句执行结果赋值失败")
		}
	}
	return nil, errors.New("不存在数据")
}

// DeleteHistoryContext
// @Description:  删除全部用户上下文
// @receiver c
// @param db
// @return error
func DeleteHistoryContext(userName string, db *sql.DB) error {
	sql := `delete from context where name=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	_, err = stmt.Exec(userName)
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

type User struct {
	Id         int64
	Name       string
	SysPrompt  string
	UpdateTime int64
}

// InsertUser
// @Description: 插入用户
// @receiver u
// @return error
func InsertUser(name, sysPrompt string, db *sql.DB) error {
	sql := `insert user (name,sys_prompt,update_time) values(?,?,?)`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	_, err = stmt.Exec(name, sysPrompt, time.Now().UnixMilli())
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// GetUser
// @Description: 查询用户信息
// @receiver u
// @return error
func GetUser(userName string, db *sql.DB) (*User, error) {
	sql := `select (name,sys_prompt,update_time) from  user where name=? `
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库查询语句")
	}
	row, err := stmt.Query(userName)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库语句执行失败")
	}
	if row.Next() {
		temp := &User{}
		err = row.Scan(temp.Id, temp.SysPrompt, temp.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "反序列话用户"+userName+"上下文失败")
		}
		return temp, nil
	}
	return nil, errors.New("不存在数据")
}

// UpdateUser
// @Description: 更新用户信息
// @receiver u
// @return error
func UpdateUser(sysPrompt, userName string, db *sql.DB) error {
	sql := `update user set sys_prompt=?,update_time+? where name=?`
	stmt, err := db.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	_, err = stmt.Exec(sysPrompt, time.Now().UnixMilli(), userName)
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// CheckUserAndCreate
// @Description: 查询并创建用户数据
// @param userName
// @return bool
func CheckUserAndCreate(userName string) bool {
	// 查询用户是否存在
	user, err := GetUser(userName, DB)
	if err != nil || user == nil {
		xlog.Log.WithField("用户名:", userName).Error(xconst.USER_DAO_SEARCH_ERR)
		return false
	}
	// 创建用户
	xlog.Log.WithField("用户名:", userName).Info(xconst.USER_DAO_FIRST_CREATE)
	err = InsertUser(userName, xconst.AI_DEFAULT_PROMPT, DB)
	if err != nil {
		xlog.Log.WithField("用户名:", userName).Error(xconst.USER_DAO_INSERT_ERR)
		return false
	}
	return true
}
