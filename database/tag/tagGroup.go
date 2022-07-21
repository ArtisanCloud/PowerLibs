package tag

import (
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
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

func (mdl *TagGroup) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_TAG_GROUP
	if needFull {
		tableName = "public" + "." + tableName
	}
	return tableName
}

func (mdl *TagGroup) GetComposedUniqueID() string {
	strKey := mdl.GroupName + "-" + mdl.OwnerType

	hashKey := security.HashStringData(strKey)

	return hashKey
}

/**
 *  Relationships
 */
