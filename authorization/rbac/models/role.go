package models

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
)

// TableName overrides the table name used by Role to `profiles`
func (mdl *Role) TableName() string {
	return mdl.GetTableName(true)
}

// Role 数据表结构
type Role struct {
	*database.PowerCompactModel

	Name     string `gorm:"column:name" json:"name"`
	ParentID int32  `gorm:"column:parent_id" json:"parentID"`
	Type     int8   `gorm:"column:type" json:"type"`
}

const TABLE_NAME_ROLE = "roles"

const ROLE_UNIQUE_ID = "id"

const ROLE_TYPE_SYSTEM int8 = 1
const ROLE_TYPE_NORMAL int8 = 2

func NewRole(mapObject *object.Collection) *Role {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newRole := &Role{
		PowerCompactModel: database.NewPowerCompactModel(),
		Name:              mapObject.GetString("name", ""),
		ParentID:          mapObject.GetInt32("parentID", 0),
		Type:              mapObject.GetInt8("type", ROLE_TYPE_NORMAL),
	}

	return newRole

}

// 获取当前 Model 的数据库表名称
func (mdl *Role) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_ROLE
	if needFull {
		tableName = "public." + tableName
	}
	return tableName
}

func (mdl *Role) GetForeignKey() string {
	return "id"
}

func (mdl *Role) GetForeignValue() string {
	return fmt.Sprintf("%d", mdl.ID)
}
