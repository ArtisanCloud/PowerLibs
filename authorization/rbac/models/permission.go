package models

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/database"
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

	Parent   *Permission   `gorm:"ForeignKey:ParentID;references:ID" json:"parent"`
	Children []*Permission `gorm:"ForeignKey:ParentID;references:ID" json:"children"`

	UniqueID     string `gorm:"column:index_permission_id;index:,unique" json:"permissionID"`
	SubjectAlias string `gorm:"column:subject_alias" json:"subjectAlias"`
	SubjectValue string `gorm:"column:subject_value; not null;" json:"subjectValue"`
	Action       string `gorm:"column:action; not null;" json:"action"`
	Description  string `gorm:"column:description" json:"description"`
	ParentID     int32  `gorm:"column:parent_id" json:"parentID"`
	Type         int8   `gorm:"column:type" json:"type"`
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
		ParentID:          mapObject.GetInt32("parentID", 0),
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
	return "permission_id"
}

func (mdl *Permission) GetForeignValue() string {
	return mdl.UniqueID
}

func (mdl *Permission) GetComposedUniqueID() string {

	strKey := fmt.Sprintf("%d", mdl.ParentID) + "-" + mdl.Action + "-" + mdl.SubjectValue
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Permission) GetTreeList(db *gorm.DB, conditions *map[string]interface{}, preloads []string,
	parentID *string, needQueryChildren bool,
) (permissions []*Permission, err error) {
	permissions = []*Permission{}
	if parentID != nil {
		if conditions == nil {
			conditions = &map[string]interface{}{}
		}
		(*conditions)["parent_id"] = parentID
	}

	err = database.GetAllList(db, conditions, &permissions, preloads)
	if err != nil {
		return nil, err
	}

	if needQueryChildren {
		for _, permission := range permissions {
			children, err := mdl.GetTreeList(db, conditions, preloads, &permission.UniqueID, needQueryChildren)
			if err != nil {
				return nil, err
			}

			permission.Children = children
		}
	}

	return permissions, err
}
