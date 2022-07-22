package tag

import (
	"database/sql"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/database"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/ArtisanCloud/PowerLibs/v2/security"
	"gorm.io/gorm"
)

type TagObjectInterface interface {
	SetTagUUID(tagUUID string) (err error)
}

// TableName overrides the table name used by Tag to `profiles`
func (mdl *Tag) TableName() string {
	return mdl.GetTableName(true)
}

// Tag 数据表结构
type Tag struct {
	*database.PowerCompactModel

	TagGroup *TagGroup `gorm:"foreignKey:GroupID;references:UniqueID" json:"tagGroup"`

	UniqueID string `gorm:"column:index_tag_id;index:,unique" json:"tagID"`
	Name     string `gorm:"column:name;" json:"name"`
	GroupID  string `gorm:"column:group_id" json:"groupID"`
	Type     int8   `gorm:"column:type" json:"type"`
}

const TABLE_NAME_TAG = "tags"

const TAG_UNIQUE_ID = "index_tag_id"

const TAG_TYPE_NORMAL int8 = 1
const TAG_TYPE_STAGE int8 = 2

func NewTag(mapObject *object.Collection) *Tag {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newTag := &Tag{
		PowerCompactModel: database.NewPowerCompactModel(),
		Name:              mapObject.GetString("name", ""),
		GroupID:           mapObject.GetString("groupID", ""),
		Type:              mapObject.GetInt8("type", TAG_TYPE_NORMAL),
	}
	newTag.UniqueID = newTag.GetComposedUniqueID()

	return newTag

}

// 获取当前 Model 的数据库表名称
func (mdl *Tag) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_TAG
	if needFull {
		tableName = "public." + tableName
	}
	return tableName
}

func (mdl *Tag) GetForeignKey() string {
	return "tag_uuid"
}

func (mdl *Tag) GetComposedUniqueID() string {

	strKey := fmt.Sprintf("%d", mdl.Type) + "-" + mdl.GroupID + "-" + mdl.Name
	hashKey := security.HashStringData(strKey)

	return hashKey
}

// 通过 UUID 或者 Name 查看tag数据
func (mdl *Tag) WhereTagName(uuidOrName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid=@value OR name=@value", sql.Named("value", uuidOrName))
	}
}

func (mdl *Tag) ClearAssociations(db *gorm.DB) (err error) {

	err = database.ClearAssociations(db, mdl, mdl.GetForeignKey(), &RTagToObject{})
	if err != nil {
		return err
	}

	return nil
}

func (mdl *Tag) GetTagUniqueIDsFromTags(tags []*Tag) []string {
	tagUniqueIDs := []string{}
	for _, tag := range tags {
		tagUniqueIDs = append(tagUniqueIDs, tag.UniqueID)
	}
	return tagUniqueIDs
}
