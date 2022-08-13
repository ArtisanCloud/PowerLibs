package models

import (
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	fmt2 "github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
	"gorm.io/gorm"
)

// TableName overrides the table name used by Permission to `profiles`
func (mdl *Permission) TableName() string {
	return mdl.GetTableName(true)
}

// Permission 数据表结构
type Permission struct {
	*database.PowerCompactModel

	Parent   *Permission   `gorm:"ForeignKey:ParentID;references:UniqueID" json:"parent"`
	Children []*Permission `gorm:"ForeignKey:ParentID;references:UniqueID" json:"children"`

	UniqueID     string  `gorm:"column:index_permission_id;index:,unique" json:"permissionID"`
	SubjectAlias string  `gorm:"column:subject_alias" json:"subjectAlias"`
	SubjectValue string  `gorm:"column:subject_value; not null;" json:"subjectValue"`
	Action       string  `gorm:"column:action; not null;" json:"action"`
	Description  string  `gorm:"column:description" json:"description"`
	ParentID     *string `gorm:"column:parent_id;index" json:"parentID"`
	Type         int8    `gorm:"column:type" json:"type"`
}

const TABLE_NAME_PERMISSION = "rbac_permissions"

const PERMISSION_UNIQUE_ID = "index_permission_id"

const PERMISSION_TYPE_NORMAL int8 = 1
const PERMISSION_TYPE_MODULE int8 = 2

func NewPermission(mapObject *object.Collection) *Permission {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newPermission := &Permission{
		PowerCompactModel: database.NewPowerCompactModel(),
		SubjectAlias:      mapObject.GetString("subjectAlias", ""),
		SubjectValue:      mapObject.GetString("subjectValue", ""),
		Action:            mapObject.GetString("action", ""),
		Description:       mapObject.GetString("description", ""),
		ParentID:          mapObject.GetStringPointer("parentID", ""),
		Type:              mapObject.GetInt8("type", PERMISSION_TYPE_NORMAL),
	}
	newPermission.UniqueID = newPermission.GetComposedUniqueID()

	return newPermission

}

// 获取当前 Model 的数据库表名称
func (mdl *Permission) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_PERMISSION
	if needFull {
		tableName = "public." + tableName
	}
	return tableName
}

func (mdl *Permission) GetForeignKey() string {
	return "index_permission_id"
}

func (mdl *Permission) GetForeignValue() string {
	return mdl.UniqueID
}

func (mdl *Permission) GetComposedUniqueID() string {

	strKey := ""
	if mdl.Type == PERMISSION_TYPE_MODULE {
		strKey = *mdl.ParentID + "-" + mdl.Action + "-" +
			mdl.SubjectAlias + mdl.SubjectValue
	} else {
		strKey = *mdl.ParentID + "-" + mdl.Action + "-" +
			mdl.SubjectAlias + mdl.SubjectValue
	}
	fmt2.Dump(strKey)
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Permission) GetGroupList(db *gorm.DB, conditions *map[string]interface{}, preloads []string) (groupedPermissions map[string]*Permission, err error) {
	permissions := []*Permission{}

	err = database.GetAllList(db, conditions, &permissions, preloads)
	if err != nil {
		return nil, err
	}

	for _, permission := range permissions {
		if permission.ParentID != nil {
			groupedPermissions[*permission.ParentID] = permission
		} else {
			groupedPermissions["unGrouped"] = permission
		}
	}

	return groupedPermissions, err
}
