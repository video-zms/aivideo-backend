package store

import (
	"axe-backend/util"
	"database/sql"
	"strings"
)

// CREATE TABLE asset (
//     id BIGINT PRIMARY KEY COMMENT '主键，章节id',
//     asset_type INT COMMENT '1-场景 2-人物 3-道具',
//     detail TEXT COMMENT '不同资产对应的json数据',
//     create_ts BIGINT,
//     update_ts BIGINT,
//     extea TEXT COMMENT '扩展字段，json'
// );

type Asset struct {
	ID        int64  `db:"id" json:"id"`
	AssetType int    `db:"asset_type" json:"asset_type"`
	Detail    string `db:"detail" json:"detail"`
	CreateTs  int64  `db:"create_ts" json:"create_ts"`
	UpdateTs  int64  `db:"update_ts" json:"update_ts"`
	Extra     string `db:"extra" json:"extra"`
}

func (a *Asset) TableName() string {
	return "asset"
}

func GetAssetInfo(aid int64) (*Asset, error) {
	asset := &Asset{}
	err := MainDB.Get(asset, "SELECT * FROM "+asset.TableName()+" WHERE id = ?", aid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return asset, nil
}

func (as *Asset) Add() error {
	sql := "INSERT INTO `asset` ("
	fields, values := util.GetStructFieldsAndValues(*as)
	query := sql + strings.Join(fields, ",") + ") values (" + strings.Join(values, ",") + ")"
	rsp, err := MainDB.Unsafe().NamedExec(query, as)
	if err != nil {
		return err
	}
	id, err := rsp.LastInsertId()
	if err != nil {
		return err
	}
	as.ID = id
	return nil
}

func (as *Asset) Update() error {
	sql := "UPDATE `asset` SET "
	fields, values := util.GetStructFieldsAndValues(*as)
	setParts := make([]string, len(fields))
	for i, field := range fields {
		setParts[i] = field + "=" + values[i]
	}
	query := sql + strings.Join(setParts, ",") + " WHERE id = :id"
	_, err := MainDB.Unsafe().NamedExec(query, as)
	if err != nil {
		return err
	}
	return nil
}

func (as *Asset) Delete() error {
	query := "DELETE FROM `asset` WHERE id = ?"
	_, err := MainDB.Unsafe().Exec(query, as.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetAssetsByType(assetType int64) ([]*Asset, error) {
	assets := []*Asset{}
	err := MainDB.Select(&assets, "SELECT * FROM "+(&Asset{}).TableName()+" WHERE asset_type = ?", assetType)
	if err != nil {
		if err == sql.ErrNoRows {
			return assets, nil
		}
		return nil, err
	}
	return assets, nil
}
