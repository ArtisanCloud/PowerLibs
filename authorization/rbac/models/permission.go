package models

import (
	"errors"
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

	PermissionModule *PermissionModule `gorm:"ForeignKey:ModuleID;references:UniqueID" json:"permissionModule"`

	UniqueID     string  `gorm:"column:index_permission_id;index:,unique" json:"permissionID"`
	SubjectAlias string  `gorm:"column:subject_alias" json:"subjectAlias"`
	SubjectValue string  `gorm:"column:subject_value; not null;" json:"subjectValue"`
	Action       string  `gorm:"column:action; not null;" json:"action"`
	Description  string  `gorm:"column:description" json:"description"`
	ModuleID     *string `gorm:"column:module_id" json:"moduleID"`
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
		ModuleID:          mapObject.GetStringPointer("moduleID", ""),
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

	strKey := mdl.Action + "-" + mdl.SubjectValue
	//fmt2.Dump(strKey)
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Permission) CheckPermissionNameAvailable(db *gorm.DB) (err error) {

	result := db.
		//Debug().
		Where("subject_alias", mdl.SubjectAlias).
		Where("index_permission_id != ?", mdl.UniqueID).
		First(&Permission{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	err = errors.New("permission name is not available")

	return err
}
