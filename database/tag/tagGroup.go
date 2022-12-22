package tag

import (
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/database"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerLibs/v3/security"
	"gorm.io/gorm"
)

// TableName overrides the table name used by TagGroup to `profiles`
func (mdl *TagGroup) TableName() string {
	return mdl.GetTableName(true)
}

type TagGroup struct {
	*database.PowerCompactModel

	Tags []*Tag `gorm:"foreignKey:GroupID;references:UniqueID" json:"tags"`

	UniqueID  string `gorm:"column:index_tag_group_id;index:,unique" json:"tagGroupID"`
	GroupName string `gorm:"column:group_name;index:index_group_name" json:"groupName"`
	OwnerType string `gorm:"column:owner_type;index:index_owner_type" json:"ownerType"`
}

const TABLE_NAME_TAG_GROUP = "tag_groups"
const TAG_GROUP_UNIQUE_ID = "index_tag_group_id"

var TABLE_FULL_NAME_TAG_GROUP = "public.ac_" + TABLE_NAME_TAG_GROUP

const DEFAULT_OWNER_TYPE = "default"
const DEFAULT_GROUP_NAME = "默认组"

func NewTagGroup(mapObject *object.Collection) *TagGroup {
	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	tagsInterface := mapObject.Get("tags", nil)
	tags := []*Tag{}
	if tagsInterface != nil {
		tags = tagsInterface.([]*Tag)
	}

	tagGroup := &TagGroup{
		PowerCompactModel: database.NewPowerCompactModel(),
		GroupName:         mapObject.GetString("groupName", ""),
		OwnerType:         mapObject.GetString("ownerType", ""),
		Tags:              tags,
	}
	tagGroup.UniqueID = tagGroup.GetComposedUniqueID()

	return tagGroup
}

func GetDefaultTagGroup(db *gorm.DB) (defaultTagGroup *TagGroup, err error) {

	defaultTagGroup = &TagGroup{}

	conditions := &map[string]interface{}{
		"group_name": DEFAULT_GROUP_NAME,
		"owner_type": DEFAULT_OWNER_TYPE,
	}

	err = database.GetFirst(db, conditions, defaultTagGroup, nil)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		defaultTagGroup = &TagGroup{
			GroupName: DEFAULT_GROUP_NAME,
			OwnerType: DEFAULT_OWNER_TYPE,
		}
		defaultTagGroup.UniqueID = defaultTagGroup.GetComposedUniqueID()

		result := db.Create(defaultTagGroup)
		err = result.Error
		if err != nil {
			return nil, err
		}
	}

	return defaultTagGroup, err
}

// 获取当前 Model 的数据库表名称
func (mdl *TagGroup) GetTableName(needFull bool) string {
	if needFull {
		return TABLE_FULL_NAME_TAG_GROUP
	} else {
		return TABLE_NAME_TAG_GROUP
	}
}

func (mdl *TagGroup) SetTableFullName(tableName string) {
	TABLE_FULL_NAME_TAG_GROUP = tableName
}

func (mdl *TagGroup) GetComposedUniqueID() string {
	strKey := mdl.GroupName + "-" + mdl.OwnerType

	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *TagGroup) CheckTagGroupNameAvailable(db *gorm.DB) (err error) {

	result := db.
		//Debug().
		Where("group_name", mdl.GroupName).
		Where("owner_type", mdl.OwnerType).
		Where("index_tag_group_id != ?", mdl.UniqueID).
		First(&TagGroup{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	err = errors.New("tag group name is not available")

	return err
}

/**
 *  Relationships
 */
