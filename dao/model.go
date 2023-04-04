package dao

import (
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
func InsertUserContext(name, contextMsg string) error {
	query := `insert into context (name, context_msg,update_time) values (?,?,?)`
	stmt, err := Db.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, contextMsg, time.Now().UnixMilli())
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// UpdateContext
// @Description: 更新用户上下文
// @receiver c
// @param Db
// @return error
func UpdateContext(contextMsg, userName string) error {
	query := `update context set context_msg=?,update_time=? where id in{
				select id from context where name=? ORDER BY update_time limit 1}`
	stmt, err := Db.Prepare(query)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
	_, err = stmt.Exec(contextMsg, time.Now().UnixMilli(), userName)
	if err != nil {
		return errors.Wrap(err, "prepare数据库语句执行失败")
	}
	return nil
}

// GetLatestUserContext
// @Description:  获取最新用户上下文对象中上下文
// @receiver c
// @param Db
// @return *model.MessageContext
// @return error
func GetLatestUserContext(userName string) (*Context, error) {
	sql := `select * from context where name=? limit 1 Order By update_time DESC`
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
	rows, err := stmt.Query(userName)
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库语句执行失败")
	}
	if rows.Next() {
		temp := &Context{}
		err = rows.Scan(&temp.Id, &temp.Name, &temp.ContextMsg, &temp.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "prepare数据库语句执行结果赋值失败")
		}
		return temp, nil
	}
	return nil, nil
}

// DeleteHistoryContext
// @Description:  删除全部用户上下文
// @receiver c
// @param Db
// @return error
func DeleteHistoryContext(userName string) error {
	sql := `delete from context where name=?`
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
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
func InsertUser(name, sysPrompt string) error {
	sql := `insert into user (name,sys_prompt,update_time) values(?,?,?)`
	stmt, err := Db.Prepare(sql)
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
func GetUser(userName string) (*User, error) {
	sql := "select name,sys_prompt,update_time from user where name=?"
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
	row, err := stmt.Query(userName)
	defer row.Close()
	if err != nil {
		return nil, errors.Wrap(err, "prepare数据库语句执行失败")
	}
	if row.Next() {
		temp := &User{}
		err = row.Scan(&temp.Name, &temp.SysPrompt, &temp.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "读取用户db失败"+userName)
		}
		return temp, nil
	}
	return nil, nil
}

// UpdateUser
// @Description: 更新用户信息
// @receiver u
// @return error
func UpdateUser(sysPrompt, userName string) error {
	sql := `update user set sys_prompt=?,update_time=? where name=?`
	stmt, err := Db.Prepare(sql)
	if err != nil {
		return errors.Wrap(err, "prepare数据库查询语句")
	}
	defer stmt.Close()
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
	user, err := GetUser(userName)
	if err != nil {
		xlog.Log.WithField("用户名:", userName).Error(xconst.USER_DAO_SEARCH_ERR)
		return false
	}
	// 存在不创建
	if user != nil {
		return true
	}
	// 创建用户
	xlog.Log.WithField("用户名:", userName).Info(xconst.USER_DAO_FIRST_CREATE)
	err = InsertUser(userName, xconst.AI_DEFAULT_PROMPT)
	if err != nil {
		xlog.Log.WithField("用户名:", userName).Error(xconst.USER_DAO_INSERT_ERR)
		return false
	}
	return true
}
