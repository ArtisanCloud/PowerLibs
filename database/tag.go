package database

import (
	"database/sql"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
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
	*PowerModel

	ParentTag *Tag `gorm:"foreignKey:ParentTagUUID;references:UUID" json:"parentTag"`

	Name          string            `gorm:"column:name;unique" json:"name"`
	ParentTagUUID object.NullString `gorm:"column:parent_tag_uuid" json:"parentTagUUID"`
	Type          int8              `gorm:"column:type" json:"type"`
}

const TABLE_NAME_TAG = "tags"

const TAG_UNIQUE_ID = "name"

const TAG_TYPE_NORMAL int8 = 1

func NewTag(mapObject *object.Collection) *Tag {

	if mapObject == nil {
		mapObject = &object.Collection{}
	}

	return &Tag{
		PowerModel:    NewPowerModel(),
		Name:          mapObject.GetString("name", ""),
		ParentTagUUID: mapObject.GetNullString("parentTagUUID", ""),
		Type:          mapObject.GetInt8("type", TAG_TYPE_NORMAL),
	}

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

// 通过 UUID 或者 Name 查看tag数据
func (mdl *Tag) WhereTagName(uuidOrName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid=@value OR name=@value", sql.Named("value", uuidOrName))
	}
}

func (mdl *Tag) ClearAssociations(db *gorm.DB) (err error) {

	err = ClearAssociations(db, mdl, mdl.GetForeignKey(), &RTagToObject{})
	if err != nil {
		return err
	}

	return nil
}

func (mdl *Tag) GetTagUUIDsFromTags(tags []*Tag) []string {
	tagUUIDs := []string{}
	for _, tag := range tags {
		tagUUIDs = append(tagUUIDs, tag.UUID)
	}
	return tagUUIDs
}
