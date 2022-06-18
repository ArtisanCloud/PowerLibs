package database

import (
	"database/sql"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"math"
	"reflect"
	"sync"
	"time"
)

type ModelInterface interface {
	GetTableName(needFull bool) string
	GetID() interface{}
	GetPrimaryKey() string
	GetForeignKey() string
}

type PowerModel struct {
	ID   int32  `gorm:"autoIncrement:true;unique; column:id; ->;<-:create" json:"-"`
	UUID string `gorm:"primaryKey;autoIncrement:false;unique; column:uuid; ->;<-:create " json:"uuid" sql:"index"`

	CreatedAt time.Time `gorm:"column:created_at; ->;<-:create " json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type PowerRelationship struct {
	ID        int32     `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at; ->;<-:create " json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

const UNIQUE_ID = "uuid"

const MODEL_STATUS_DRAFT int8 = 0
const MODEL_STATUS_ACTIVE int8 = 1
const MODEL_STATUS_CANCELED int8 = 2
const MODEL_STATUS_PENDING int8 = 3
const MODEL_STATUS_INACTIVE int8 = 4

const APPROVAL_STATUS_DRAFT int8 = 0
const APPROVAL_STATUS_PENDING int8 = 1
const APPROVAL_STATUS_APPROVED int8 = 3
const APPROVAL_STATUS_REJECTED int8 = 4

var ArrayModelFields *object.HashMap = &object.HashMap{}

func NewPowerModel() *PowerModel {
	now := time.Now()
	return &PowerModel{
		UUID:      uuid.New().String(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewPowerRelationship() *PowerRelationship {
	now := time.Now()
	return &PowerRelationship{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (mdl *PowerModel) GetID() interface{} {
	return mdl.UUID
}

func (mdl *PowerModel) GetPrimaryKey() string {
	return "uuid"
}
func (mdl *PowerModel) GetForeignKey() string {
	return "model_uuid"
}

func (mdl *PowerRelationship) GetID() interface{} {
	return ""
}

func (mdl *PowerRelationship) GetPrimaryKey() string {
	return "id"
}
func (mdl *PowerRelationship) GetForeignKey() string {
	return "model_id"
}

/**
 * Scope Where Conditions
 */
func WhereUUID(uuid string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid=?", uuid)
	}
}

func WhereAccountUUID(uuid string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("account_uuid=@value", sql.Named("value", uuid))
	}
}

func WhereCampaignUUID(uuid string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("campaign_uuid=@value", sql.Named("value", uuid))
	}
}

func GetFirst(db *gorm.DB, conditions *map[string]interface{}, model interface{}, preloads []string) (err error) {

	if conditions != nil {
		db = db.Where(*conditions)
	}

	// add preloads
	if len(preloads) > 0 {
		for _, preload := range preloads {
			if preload != "" {
				db.Preload(preload)
			}
		}
	}

	result := db.First(model)

	return result.Error
}

func GetList(db *gorm.DB, conditions *map[string]interface{},
	models interface{}, preloads []string,
	page int, pageSize int) (paginator *Pagination, err error) {

	// add pagination
	paginator = NewPagination(page, pageSize, "")
	var totalRows int64
	db.Model(models).Count(&totalRows)
	paginator.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(paginator.Limit)))
	paginator.TotalPages = totalPages

	db = db.Scopes(
		Paginate(page, pageSize),
	)

	if conditions != nil {
		db = db.Where(*conditions)
	}

	// add preloads
	if len(preloads) > 0 {
		for _, preload := range preloads {
			if preload != "" {
				db.Preload(preload)
			}
		}
	}

	// chunk datas
	result := db.Find(models)
	err = result.Error
	if err != nil {
		return paginator, err
	}

	paginator.Data = models

	return paginator, nil
}

func GetDownloadList(db *gorm.DB, conditions *map[string]interface{},
	models interface{}, preloads []string) (err error) {

	if conditions != nil {
		db = db.Where(*conditions)
	}

	// add preloads
	if len(preloads) > 0 {
		for _, preload := range preloads {
			if preload != "" {
				db.Preload(preload)
			}
		}
	}

	// chunk datas
	result := db.Find(models)
	err = result.Error
	if err != nil {
		return err
	}

	return nil
}

/**
 * Association Relationship
 */
func AssociationRelationship(db *gorm.DB, conditions *map[string]interface{}, mdl interface{}, relationship string, withClauseAssociations bool) *gorm.Association {

	tx := db.Model(mdl)

	if withClauseAssociations {
		tx.Preload(clause.Associations)
	}

	if conditions != nil {
		tx = tx.Where(*conditions)
	}

	return tx.Association(relationship)
}

func ClearAssociation(db *gorm.DB, object ModelInterface, foreignKey string, pivot ModelInterface) error {
	result := db.Exec("DELETE FROM "+pivot.GetTableName(true)+" WHERE "+foreignKey+"=?", object.GetID())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func GetModelFields(model interface{}) (fields []string) {

	// check if it has been loaded
	modelType := reflect.TypeOf(model)
	modelName := modelType.String()
	if (*ArrayModelFields)[modelName] != nil {
		return (*ArrayModelFields)[modelName].([]string)
	}

	fmt.Printf("parse object ~%s~ model fields \n", modelName)
	gormSchema, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		println(err)
		return fields
	}

	fields = []string{}
	for _, field := range gormSchema.Fields {
		if field.DBName != "" && !field.PrimaryKey && !field.Unique && field.Updatable {
			fields = append(fields, field.DBName)
		}
	}
	(*ArrayModelFields)[modelName] = fields
	fmt.Printf("parsed object ~%s~ model fields and fields count is %d \n\n", modelName, len(fields))

	return fields
}

func (mdl *PowerModel) IsModelLoaded() bool {
	return mdl != nil && mdl.UUID != ""
}
