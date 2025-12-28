package store

import (
	"axe-backend/util"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CREATE TABLE story (
//     id BIGINT PRIMARY KEY COMMENT '主键，章节id',
//     project_id BIGINT COMMENT '属于那个项目的章节',
//     story_title VARCHAR(255) COMMENT '内容标题',
//     story_scene VARCHAR(255) COMMENT '场景',
//     story TEXT COMMENT '剧本',
//     create_ts BIGINT,
//     update_ts BIGINT,
//     story_shots VARCHAR(1024) COMMENT '分镜id集合',
//     extea TEXT COMMENT '扩展字段，json'
// );

type Chapter struct {
	ID         int64  `db:"id" json:"id"`
	ProjectID  int64  `db:"project_id" json:"project_id"`
	StoryTitle string `db:"story_title" json:"story_title"`
	StoryScene string `db:"story_scene" json:"story_scene"`
	Story      string `db:"story" json:"story"`
	CreateTs   int64  `db:"create_ts" json:"create_ts"`
	UpdateTs   int64  `db:"update_ts" json:"update_ts"`
	StoryShots string `db:"story_shots" json:"story_shots"`
	Extra      string `db:"extra" json:"extra"`
}

func (c *Chapter) TableName() string {
	return "chapter"
}

func GetChapterInfo(cid int64) (*Chapter, error) {
	chapter := &Chapter{}
	err := MainDB.Get(chapter, "SELECT * FROM "+chapter.TableName()+" WHERE id = ?", cid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return chapter, nil
}

func (ch *Chapter) Add() error {
	sql := "INSERT INTO `chapter` ("
	fields, values := util.GetStructFieldsAndValues(*ch)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ")"
	res, err := MainDB.Unsafe().NamedExec(query, ch)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	ch.ID = id
	return nil
}

func (ch *Chapter) Update() error {
	sql := "UPDATE `chapter` SET "
	fields, values := util.GetStructFieldsAndValues(*ch)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, ch)
	if err != nil {
		return err
	}
	return nil
}

func (ch *Chapter) Delete() error {
	query := "DELETE FROM `chapter` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, ch.ID)
	if err != nil {
		return err
	}
	return nil
}

func ListChaptersByProject(projectID int64) ([]*Chapter, error) {
	chapters := []*Chapter{}
	err := MainDB.Select(&chapters, "SELECT * FROM "+(&Chapter{}).TableName()+" WHERE project_id = ?", projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return chapters, nil
		}
		return nil, err
	}
	return chapters, nil
}

func ListChaptersByProjectIDs(projectIDs []int64) ([]*Chapter, error) {
	chapters := []*Chapter{}
	query, args, err := sqlx.In("SELECT * FROM "+(&Chapter{}).TableName()+" WHERE project_id IN (?)", projectIDs)
	if err != nil {
		return nil, err
	}
	err = MainDB.Select(&chapters, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return chapters, nil
		}
		return nil, err
	}
	return chapters, nil
}

func GetChapterByID(cid int64) (*Chapter, error) {
	chapter := &Chapter{}
	err := MainDB.Get(chapter, "SELECT * FROM "+chapter.TableName()+" WHERE id = ?", cid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return chapter, nil
}

func ListChaptersByProjectID(projectID int64) ([]*Chapter, error) {
	chapters := []*Chapter{}
	err := MainDB.Select(&chapters, "SELECT * FROM "+(&Chapter{}).TableName()+" WHERE project_id = ?", projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return chapters, nil
		}
		return nil, err
	}
	return chapters, nil
}
