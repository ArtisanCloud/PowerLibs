package models

import (
	"github.com/ArtisanCloud/PowerLibs/v3/database"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
)

// TableName overrides the table name used by User to `profiles`
func (mdl *Recipient) TableName() string {
	//return config.DatabaseConn.Schemas["option"] + "." + TABLE_NAME_CAMPAIGN
	return mdl.GetTableName(true)
}

type Recipient struct {
	*database.PowerPivot

	Email string `gorm:"column:email" json:"email"`
	Phone string `gorm:"column:phone" json:"phone"`

	//common fields
	OwnerID   object.NullString `gorm:"column:owner_id;not null;index:owner_id" json:"ownerID"`
	OwnerType string            `gorm:"column:owner_type" json:"ownerType"`
}

func NewRecipient(mapObject *object.Collection) *Recipient {
	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	email := mapObject.GetString("email", "")
	phone := mapObject.GetString("phone", "")
	if email == "" || phone == "" {
		return nil
	}

	ownerID := mapObject.GetString("ownerID", "")
	ownerType := mapObject.GetString("ownerType", "")
	if ownerID == "" && ownerType == "" {
		return nil
	}

	return &Recipient{
		PowerPivot: database.NewPowerPivot(),
		Email:      email,
		Phone:      phone,
		OwnerID:    object.NewNullString(ownerID, true),
		OwnerType:  ownerType,
	}
}

const TABLE_NAME_RECIPEINT string = "recipients"
const OBJECT_NAME_RECIPEINT string = "Recipient"

func (mdl *Recipient) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_RECIPEINT
	if needFull {
		tableName = "public." + tableName
	}
	return tableName
}

func (mdl *Recipient) GetID() int32 {
	return mdl.ID
}
func (mdl *Recipient) GetUUID() string {
	return ""
}
func (mdl *Recipient) GetPrimaryKey() string {
	return "id"
}

func (mdl *Recipient) GetForeignKey() string {
	return "recipient_id"
}
