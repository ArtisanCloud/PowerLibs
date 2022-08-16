package models

import (
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TableName overrides the table name used by Role to `profiles`
func (mdl *Role) TableName() string {
	return mdl.GetTableName(true)
}

// Role 数据表结构
type Role struct {
	*database.PowerCompactModel

	Parent   *Role   `gorm:"ForeignKey:ParentID;references:UniqueID" json:"parent"`
	Children []*Role `gorm:"ForeignKey:ParentID;references:UniqueID" json:"children"`

	UniqueID string  `gorm:"column:index_role_id;index:,unique" json:"roleID"`
	Name     string  `gorm:"column:name" json:"name"`
	ParentID *string `gorm:"column:parent_id;index" json:"parentID"`
	Type     int8    `gorm:"column:type" json:"type"`
}

const TABLE_NAME_ROLE = "roles"

const ROLE_UNIQUE_ID = "index_role_id"

const ROLE_TYPE_ALL int8 = 0
const ROLE_TYPE_SYSTEM int8 = 1
const ROLE_TYPE_NORMAL int8 = 2

const ROLE_ROOT_NAME string = "超级管理员"
const ROLE_ADMIN_NAME string = "管理员"
const ROLE_EMPLOYEE_NAME string = "普通员工"

func NewRole(mapObject *object.Collection) *Role {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newRole := &Role{
		PowerCompactModel: database.NewPowerCompactModel(),
		Name:              mapObject.GetString("name", ""),
		ParentID:          mapObject.GetStringPointer("parentID", ""),
		Type:              mapObject.GetInt8("type", ROLE_TYPE_NORMAL),
	}
	newRole.UniqueID = newRole.GetComposedUniqueID()

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
	return "role_id"
}

func (mdl *Role) GetForeignValue() string {
	return mdl.UniqueID
}

func (mdl *Role) GetComposedUniqueID() string {

	strKey := *mdl.ParentID + "-" + mdl.Name
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Role) GetRootComposedUniqueID() string {
	strKey := "" + "-" + ROLE_ROOT_NAME
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Role) GetAdminComposedUniqueID() string {
	strKey := "" + "-" + ROLE_ADMIN_NAME
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Role) GetEmployeeComposedUniqueID() string {
	strKey := "" + "-" + ROLE_EMPLOYEE_NAME
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Role) GetTreeList(db *gorm.DB, conditions *map[string]interface{}, preloads []string,
	roleType int8, parentID *string, needQueryChildren bool,
) (roles []*Role, err error) {
	roles = []*Role{}

	if conditions == nil {
		conditions = &map[string]interface{}{}
	}

	if parentID != nil {
		(*conditions)["parent_id"] = parentID
	}
	if roleType != ROLE_TYPE_ALL {
		(*conditions)["type"] = roleType
	}

	err = database.GetAllList(db, conditions, &roles, preloads)
	if err != nil {
		return nil, err
	}

	if needQueryChildren {
		for _, role := range roles {
			children, err := mdl.GetTreeList(db, conditions, preloads, roleType, &role.UniqueID, needQueryChildren)
			if err != nil {
				return nil, err
			}

			role.Children = children
		}
	}

	return roles, err
}

func (mdl *Role) CheckRoleNameAvailable(db *gorm.DB) (err error) {

	result := db.
		//Debug().
		Where("name", mdl.Name).
		Where("index_role_id != ?", mdl.UniqueID).
		First(&Role{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	err = errors.New("role name is not available")

	return err
}
