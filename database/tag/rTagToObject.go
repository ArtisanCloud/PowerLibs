package tag

import (
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
	"gorm.io/gorm"
)

// TableName overrides the table name used by RTagToObject to `profiles`
func (mdl *RTagToObject) TableName() string {
	return mdl.GetTableName(true)
}

// r_tag_to_object 数据表结构
type RTagToObject struct {
	*database.PowerPivot

	//common fields
	UniqueID          object.NullString `gorm:"index:index_taggable_object_id;index:index_taggable_id;index;column:index_tag_to_object_id;unique"`
	TaggableOwnerType object.NullString `gorm:"column:taggable_owner_type;not null" json:"taggableOwnerType"`
	TaggableObjectID  object.NullString `gorm:"column:taggable_object_id;not null;index:index_taggable_object_id" json:"taggableObjectID"`
	TaggableID        object.NullString `gorm:"column:tag_id;not null;index:index_taggable_id" json:"taggableID"`
}

const TABLE_NAME_R_TAG_TO_OBJECT = "r_tag_to_object"

const R_TAG_TO_OJECT_UNIQUE_ID = "index_tag_to_object_id"

var TABLE_FULL_NAME_R_TAG_TO_OBJECT string = "public.ac_" + TABLE_NAME_R_TAG_TO_OBJECT

const R_TAG_TO_OJECT_FOREIGN_KEY = "taggable_object_id"
const R_TAG_TO_OJECT_OWNER_KEY = "taggable_owner_type"
const R_TAG_TO_OJECT_JOIN_KEY = "tag_id"

func (mdl *RTagToObject) GetTableName(needFull bool) string {
	if needFull {
		return TABLE_FULL_NAME_R_TAG_TO_OBJECT
	} else {
		return TABLE_NAME_R_TAG_TO_OBJECT
	}
}

func (mdl *RTagToObject) SetTableFullName(tableName string) {
	TABLE_FULL_NAME_R_TAG_TO_OBJECT = tableName
}

func (mdl *RTagToObject) GetForeignKey() string {
	return R_TAG_TO_OJECT_FOREIGN_KEY
}

func (mdl *RTagToObject) GetForeignValue() string {
	return mdl.TaggableObjectID.String
}

func (mdl *RTagToObject) GetJoinKey() string {
	return R_TAG_TO_OJECT_JOIN_KEY
}

func (mdl *RTagToObject) GetJoinValue() string {
	return mdl.TaggableID.String
}

func (mdl *RTagToObject) GetOwnerKey() string {
	return R_TAG_TO_OJECT_OWNER_KEY
}

func (mdl *RTagToObject) GetOwnerValue() string {
	return mdl.TaggableOwnerType.String
}

func (mdl *RTagToObject) GetPivotComposedUniqueID() string {
	strKey := mdl.GetOwnerValue() + "-" + mdl.GetForeignValue() + "-" + mdl.GetJoinValue()
	hashKey := security.HashStringData(strKey)

	return hashKey
}

// ---------------------------------------------------------------------------------------------------------------------

func (mdl *RTagToObject) MakePivotsFromObjectAndTags(obj database.ModelInterface, tags []*Tag) ([]database.PivotInterface, error) {
	pivots := []database.PivotInterface{}
	for _, tag := range tags {
		pivot := &RTagToObject{
			TaggableOwnerType: object.NewNullString(obj.GetTableName(true), true),
			TaggableObjectID:  object.NewNullString(obj.GetForeignReferValue(), true),
			TaggableID:        object.NewNullString(tag.UniqueID, true),
		}
		pivot.UniqueID = object.NewNullString(pivot.GetPivotComposedUniqueID(), true)
		pivots = append(pivots, pivot)
	}
	return pivots, nil
}

func (mdl *RTagToObject) GetPivots(db *gorm.DB) ([]*RTagToObject, error) {
	pivots := []*RTagToObject{}

	db = database.SelectMorphPivot(db, mdl)

	result := db.Find(&pivots)

	return pivots, result.Error

}
