package store

import (
	"axe-backend/util"
	"database/sql"
	"strings"
)

// CREATE TABLE shot (
//     id BIGINT PRIMARY KEY COMMENT '主键，章节id',
//     shot_id VARCHAR(64) COMMENT '分镜id',
//     shot_number INT COMMENT '分镜号',
//     duration INT COMMENT '时长',
//     scene_type VARCHAR(64) COMMENT '场景类型',
//     camera_movement VARCHAR(64) COMMENT '摄像头移动类型',
//     `desc` VARCHAR(255) COMMENT '描述',
//     dialogue VARCHAR(255) COMMENT '对白',
//     notes VARCHAR(255) COMMENT '备注',
//     timestamp BIGINT,
//     scene_id BIGINT COMMENT '场景id',
//     audio_id BIGINT COMMENT '音频id',
//     charater_ids VARCHAR(1024) COMMENT '角色id列表',
//     create_ts BIGINT,
//     update_ts BIGINT,
//     extea TEXT COMMENT '扩展字段，json'
// );

type Shot struct {
	ID             int64  `db:"id" json:"id"`
	ShotID         string `db:"shot_id" json:"shot_id"`
	ShotNumber     int    `db:"shot_number" json:"shot_number"`
	Duration       int    `db:"duration" json:"duration"`
	SceneType      string `db:"scene_type" json:"scene_type"`
	CameraMovement string `db:"camera_movement" json:"camera_movement"`
	Desc           string `db:"desc" json:"desc"`
	Dialogue       string `db:"dialogue" json:"dialogue"`
	Notes          string `db:"notes" json:"notes"`
	Timestamp      int64  `db:"timestamp" json:"timestamp"`
	SceneID        int64  `db:"scene_id" json:"scene_id"`
	AudioID        int64  `db:"audio_id" json:"audio_id"`
	CharaterIDs    string `db:"charater_ids" json:"charater_ids"`
	CreateTs       int64  `db:"create_ts" json:"create_ts"`
	UpdateTs       int64  `db:"update_ts" json:"update_ts"`
	Extra          string `db:"extra" json:"extra"`
}

func (s *Shot) TableName() string {
	return "shot"
}

func GetShotInfo(sid int64) (*Shot, error) {
	shot := &Shot{}
	err := MainDB.Get(shot, "SELECT * FROM "+shot.TableName()+" WHERE id = ?", sid)
	if err != nil {
		return nil, err
	}
	return shot, nil
}

func (sh *Shot) Add() error {
	sql := "INSERT INTO `shot` ("
	fields, values := util.GetStructFieldsAndValues(*sh)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ")"
	rsp, err := MainDB.Unsafe().NamedExec(query, sh)
	if err != nil {
		return err
	}
	id, err := rsp.LastInsertId()
	if err != nil {
		return err
	}
	sh.ID = id
	return nil
}

func (sh *Shot) Update() error {
	sql := "UPDATE `shot` SET "
	fields, values := util.GetStructFieldsAndValues(*sh)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, sh)
	if err != nil {
		return err
	}
	return nil
}

func (sh *Shot) Delete() error {
	query := "DELETE FROM `shot` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, sh.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetShotsBySceneID(sceneID int64) ([]*Shot, error) {
	shots := []*Shot{}
	err := MainDB.Select(&shots, "SELECT * FROM "+(&Shot{}).TableName()+" WHERE scene_id = ?", sceneID)
	if err != nil {
		return nil, err
	}
	return shots, nil
}

func GetShotByID(id int64) (*Shot, error) {
	shot := &Shot{}
	err := MainDB.Get(shot, "SELECT * FROM "+shot.TableName()+" WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return shot, nil
}

func GetShotByShotID(shotID string) (*Shot, error) {
	shot := &Shot{}
	err := MainDB.Get(shot, "SELECT * FROM "+shot.TableName()+" WHERE shot_id = ?", shotID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return shot, nil
}
