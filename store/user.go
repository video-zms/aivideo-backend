package store

import (
	"database/sql"
	"strings"

	"axe-backend/util"

	"github.com/sirupsen/logrus"
)

type User struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	CreateTs  int64  `db:"create_ts" json:"create_ts"`
	UpdateTs  int64  `db:"update_ts" json:"update_ts"`
	Privilege int    `db:"privilege" json:"privilege"`
	Coin      int    `db:"coin" json:"coin"`
	Extra     string `db:"extra" json:"extra"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserInfo(uid int64) (*User, error) {
	user := &User{}
	err := MainDB.Get(user, "SELECT * FROM "+user.TableName()+" WHERE id = ?", uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (cr *User) Add() error {
	// 从结构体生成字段与占位符
	fields, values := util.GetStructFieldsAndValues(*cr)

	// 过滤掉 id 字段，让 MySQL 使用自增值
	newFields := make([]string, 0, len(fields))
	newValues := make([]string, 0, len(values))
	for i, f := range fields {
		// 去掉可能的反引号并忽略大小写比较
		name := strings.Trim(f, "`")
		if strings.EqualFold(strings.TrimSpace(name), "id") {
			continue
		}
		newFields = append(newFields, f)
		newValues = append(newValues, values[i])
	}

	query := "INSERT INTO `user` (" + strings.Join(newFields, ",") + ") values (" + strings.Join(newValues, ",") + ")"
	res, err := MainDB.Unsafe().NamedExec(query, cr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"collect_record": cr, "uid": cr.ID, "err": err}).Error("save collect record error")
		return err
	}

	// 读取并设置自动生成的 id
	if res != nil {
		if lastID, err2 := res.LastInsertId(); err2 == nil {
			cr.ID = lastID
		}
	}

	return nil
}

func (cr *User) Update() error {
	sql := "UPDATE `user` SET "
	fields, values := util.GetStructFieldsAndValues(*cr)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, cr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"collect_record": cr, "uid": cr.ID, "err": err}).Error("update collect record error")
	}
	return err
}

func (cr *User) Delete() error {
	query := "DELETE FROM `user` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, cr.ID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"collect_record": cr, "uid": cr.ID, "err": err}).Error("delete collect record error")
	}
	return err
}

func GetUserByUsernameOrEmail(username, email string) (*User, error) {
	user := &User{}
	var err error
	if username != "" {
		err = MainDB.Get(user, "SELECT * FROM "+user.TableName()+" WHERE username = ?", username)
	} else {
		err = MainDB.Get(user, "SELECT * FROM "+user.TableName()+" WHERE email = ?", email)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
