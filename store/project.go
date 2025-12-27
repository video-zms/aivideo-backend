package store

import (
	"axe-backend/util"
	"database/sql"
	"strings"
)

// CREATE TABLE project (
//     id BIGINT PRIMARY KEY COMMENT '主键，项目id',
//     `desc` VARCHAR(255) COMMENT '项目描述',
//     creator VARCHAR(255) COMMENT '项目创建人，用户表中email',
//     create_ts BIGINT,
//     update_ts BIGINT,
//     extea TEXT COMMENT '扩展字段，json'
// );

type Project struct {
	ID       int64  `db:"id" json:"id"`
	Desc     string `db:"desc" json:"desc"`
	Creator  string `db:"creator" json:"creator"`
	CreateTs int64  `db:"create_ts" json:"create_ts"`
	UpdateTs int64  `db:"update_ts" json:"update_ts"`
	Extea    string `db:"extea" json:"extea"`
}

func (p *Project) TableName() string {
	return "project"
}

func GetProjectInfo(pid int64) (*Project, error) {
	project := &Project{}
	err := MainDB.Get(project, "SELECT * FROM "+project.TableName()+" WHERE id = ?", pid)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (pr *Project) Add() error {
	sql := "INSERT INTO `project` ("
	fields, values := util.GetStructFieldsAndValues(*pr)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ") on duplicate key update desc = :desc, update_ts = :update_ts, extea = :extea"
	_, err := MainDB.Unsafe().NamedExec(query, pr)
	if err != nil {
		return err
	}
	return nil
}

func (pr *Project) Update() error {
	sql := "UPDATE `project` SET "
	fields, values := util.GetStructFieldsAndValues(*pr)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, pr)
	if err != nil {
		return err
	}
	return nil
}

func (pr *Project) Delete() error {
	query := "DELETE FROM `project` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, pr.ID)
	if err != nil {
		return err
	}
	return nil
}

func ListProjectsByCreator(creator string) ([]*Project, error) {
	projects := []*Project{}
	err := MainDB.Select(&projects, "SELECT * FROM "+(&Project{}).TableName()+" WHERE creator = ?", creator)
	if err != nil {
		if err == sql.ErrNoRows {
			return projects, nil
		}
		return nil, err
	}
	return projects, nil
}


func ListAllProjects() ([]*Project, error) {
	projects := []*Project{}
	err := MainDB.Select(&projects, "SELECT * FROM "+(&Project{}).TableName())
	if err != nil {
		if err == sql.ErrNoRows {
			return projects, nil
		}
		return nil, err
	}
	return projects, nil
}

func GetProjectById(pid int64) (*Project, error) {
	return GetProjectInfo(pid)
}

func GetProjectsByCreator(creator string) ([]*Project, error){
	return ListProjectsByCreator(creator)
}