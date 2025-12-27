package store

import (
	"axe-backend/util"
	"database/sql"
	"strings"
)

// CREATE TABLE video (
//     id BIGINT PRIMARY KEY COMMENT '主键，章节id',
//     status INT COMMENT '1-未开始 2-进行中 3-成功 4-失败',
//     video_url VARCHAR(255) COMMENT '生成的视频链接',
//     create_ts BIGINT,
//     update_ts BIGINT,
//     extea TEXT COMMENT '扩展字段，json'
// );

type Task struct {
	ID        int64  `db:"id" json:"id"`
	Status    int    `db:"status" json:"status"`
	VideoURL  string `db:"video_url" json:"video_url"`
	CreateTs  int64  `db:"create_ts" json:"create_ts"`
	UpdateTs  int64  `db:"update_ts" json:"update_ts"`
	Extea     string `db:"extea" json:"extea"`
}

func (t *Task) TableName() string {
	return "video"
}

func GetTaskInfo(tid int64) (*Task, error) {
	task := &Task{}
	err := MainDB.Get(task, "SELECT * FROM "+task.TableName()+" WHERE id = ?", tid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return task, nil
}

func (ta *Task) Add() error {
	sql := "INSERT INTO `video` ("
	fields, values := util.GetStructFieldsAndValues(*ta)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ") on duplicate key update status = :status, video_url = :video_url, update_ts = :update_ts, extea = :extea"
	_, err := MainDB.Unsafe().NamedExec(query, ta)
	if err != nil {
		return err
	}
	return nil
}

func (ta *Task) Update() error {
	sql := "UPDATE `video` SET "
	fields, values := util.GetStructFieldsAndValues(*ta)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, ta)
	if err != nil {
		return err
	}
	return nil
}

func (ta *Task) Delete() error {
	query := "DELETE FROM `video` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, ta.ID)
	if err != nil {
		return err
	}
	return nil
}



func GetTaskByID(tid int64) (*Task, error) {
	task := &Task{}
	err := MainDB.Get(task, "SELECT * FROM "+task.TableName()+" WHERE id = ?", tid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func DeleteTaskByID(tid int64) error {
	query := "DELETE FROM `video` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, tid)
	if err != nil {
		return err
	}
	return nil
}