package store

import (
	"axe-backend/db"
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
	sql := "INSERT INTO `user` ("
	fields, values := util.GetStructFieldsAndValues(*cr)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ") on duplicate key update is_collect = :is_collect, update_time = :update_time"
	_, err := MainDB.Unsafe().NamedExec(query, cr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"collect_record": cr, "uid": cr.ID, "err": err}).Error("save collect record error")
	}
	return err
}
func (cr *User) Update() error {
	sql := "UPDATE `user` SET "
	fields, values := util.GetStructFieldsAndValues(*cr)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := db.MainDB.Unsafe().NamedExec(query, cr)
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
