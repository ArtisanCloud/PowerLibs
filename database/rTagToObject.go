package database

import (
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"gorm.io/gorm"
)

// TableName overrides the table name used by RTagToObject to `profiles`
func (mdl *RTagToObject) TableName() string {
	return mdl.GetTableName(true)
}

// r_tag_to_object 数据表结构
type RTagToObject struct {
	*PowerRelationship

	//common fields
	UniqueID          object.NullString `gorm:"index:index_taggable_object_id;index:index_taggable_id;index;column:index_tag_to_object_id;unique"`
	TaggableOwnerType object.NullString `gorm:"column:taggable_owner_type;not null" json:"taggableOwnerType"`
	TaggableObjectID  object.NullString `gorm:"column:taggable_object_id;not null;index:index_taggable_object_id" json:"taggableObjectID"`
	TaggableID        object.NullString `gorm:"column:tag_id;not null;index:index_taggable_id" json:"taggableID"`
}

const TABLE_NAME_R_TAG_TO_OBJECT = "r_tag_to_object"

const R_TAG_TO_OJECT_UNIQUE_ID = "index_tag_to_object_id"

const R_TAG_TO_OJECT_FOREIGN_KEY = "taggable_object_id"
const R_TAG_TO_OJECT_OWNER_KEY = "taggable_owner_type"
const R_TAG_TO_OJECT_JOIN_KEY = "tag_id"

func (mdl *RTagToObject) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_R_TAG_TO_OBJECT
	if needFull {
		tableName = "public" + "." + tableName
	}
	return tableName
}

func (mdl *RTagToObject) GetPivots(db *gorm.DB, foreignValue string, joinValue string, ownerValue string) ([]*RTagToObject, error) {
	pivots := []*RTagToObject{}

	db = SelectMorphPivot(db, mdl,
		R_TAG_TO_OJECT_FOREIGN_KEY, foreignValue,
		R_TAG_TO_OJECT_JOIN_KEY, joinValue,
		R_TAG_TO_OJECT_OWNER_KEY, ownerValue,
	)

	result := db.Find(&pivots)

	return pivots, result.Error

}
