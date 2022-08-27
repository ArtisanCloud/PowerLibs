package models

import (
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
	"gorm.io/gorm"
)

// TableName overrides the table name used by PermissionModule to `profiles`
func (mdl *PermissionModule) TableName() string {
	return mdl.GetTableName(true)
}

// PermissionModule 数据表结构
type PermissionModule struct {
	*database.PowerCompactModel

	Parent      *PermissionModule   `gorm:"ForeignKey:ParentID;references:UniqueID" json:"parent"`
	Children    []*PermissionModule `gorm:"ForeignKey:ParentID;references:UniqueID" json:"children"`
	Permissions []*Permission       `gorm:"ForeignKey:ModuleID;references:UniqueID" json:"permissions"`

	UniqueID    string  `gorm:"column:index_permission_module_id;index:,unique" json:"permissionModuleID"`
	Name        string  `gorm:"column:name" json:"name"`
	Description string  `gorm:"column:description" json:"description"`
	ParentID    *string `gorm:"column:parent_id;index" json:"parentID"`
}

const TABLE_NAME_PERMISSION_MODULE = "rbac_permission_modules"

const PERMISSION_MODULE_UNIQUE_ID = "index_permission_module_id"

var TABLE_FULL_NAME_PERMISSION_MODULE string = "public.ac_" + TABLE_NAME_PERMISSION_MODULE

func NewPermissionModule(mapObject *object.Collection) *PermissionModule {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newPermissionModule := &PermissionModule{
		PowerCompactModel: database.NewPowerCompactModel(),
		Name:              mapObject.GetString("name", ""),
		Description:       mapObject.GetString("description", ""),
		ParentID:          mapObject.GetStringPointer("parentID", ""),
	}
	newPermissionModule.UniqueID = newPermissionModule.GetComposedUniqueID()

	return newPermissionModule

}

// 获取当前 Model 的数据库表名称
func (mdl *PermissionModule) GetTableName(needFull bool) string {
	if needFull {
		return TABLE_FULL_NAME_PERMISSION_MODULE
	} else {
		return TABLE_NAME_PERMISSION_MODULE
	}
}
func (mdl *PermissionModule) SetTableFullName(tableName string) {
	TABLE_FULL_NAME_PERMISSION_MODULE = tableName
}

func (mdl *PermissionModule) GetForeignKey() string {
	return "index_permission_module_id"
}

func (mdl *PermissionModule) GetForeignValue() string {
	return mdl.UniqueID
}

func (mdl *PermissionModule) GetComposedUniqueID() string {

	strKey := *mdl.ParentID + "-" + mdl.Name
	//fmt2.Dump(strKey)
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *PermissionModule) GetRBACRuleName() string {

	//return mdl.Name + "-" + mdl.UniqueID[0:5]
	return mdl.UniqueID

}

func (mdl *PermissionModule) GetGroupList(db *gorm.DB, conditions *map[string]interface{}, preloads []string) (permissionModules []*PermissionModule, err error) {
	permissionModules = []*PermissionModule{}

	if preloads == nil {
		preloads = []string{"Permissions"}
	}

	if conditions == nil {
		conditions = &map[string]interface{}{}
	}
	if _, ok := (*conditions)["parent_id"]; !ok {
		(*conditions)["parent_id"] = ""
	}

	//db = db.Debug()
	err = database.GetAllList(db, conditions, &permissionModules, preloads)
	if err != nil {
		return nil, err
	}

	for _, module := range permissionModules {
		(*conditions)["parent_id"] = &module.UniqueID
		children, err := mdl.GetGroupList(db, conditions, preloads)
		if err != nil {
			return nil, err
		}

		module.Children = children
	}

	return permissionModules, err
}

func (mdl *PermissionModule) CheckPermissionModuleNameAvailable(db *gorm.DB) (err error) {

	result := db.
		//Debug().
		Where("name", mdl.Name).
		Where("index_permission_module_id != ?", mdl.UniqueID).
		Where("parent_id = ?", mdl.ParentID).
		First(&PermissionModule{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	err = errors.New("permission module name is not available")

	return err
}
