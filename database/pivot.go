package database

import (
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type PowerRelationship struct {
	ID        int32     `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at; ->;<-:create " json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func NewPowerRelationship() *PowerRelationship {
	now := time.Now()
	return &PowerRelationship{
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// --------------------------------------------------------------------
func (mdl *PowerRelationship) GetTableName(needFull bool) string {
	return ""
}

func (mdl *PowerRelationship) GetPowerModel() ModelInterface {
	return mdl
}
func (mdl *PowerRelationship) GetID() int32 {
	return mdl.ID
}

func (mdl *PowerRelationship) GetUUID() string {
	return ""
}

func (mdl *PowerRelationship) GetPrimaryKey() string {
	return "id"
}
func (mdl *PowerRelationship) GetForeignKey() string {
	return "model_id"
}

func GetPivotComposedUniqueID(foreignValue string, joinValue string) object.NullString {
	if foreignValue != "" && joinValue != "" {
		strUniqueID := foreignValue + "-" + joinValue
		return object.NewNullString(strUniqueID, true)
	} else {
		return object.NewNullString("", false)
	}
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

func ClearAssociations(db *gorm.DB, object ModelInterface, foreignKey string, pivot ModelInterface) error {
	result := db.Exec("DELETE FROM "+pivot.GetTableName(true)+" WHERE "+foreignKey+"=?", object.GetID())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// --------------------------------------------------------------------
func AppendAssociates(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValues []string) (err error) {

	return AppendMorphAssociates(db, pivot, foreignKey, foreignValue, joinKey, joinValues, "", "")
}
func SyncAssociates(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValues []string) (err error) {

	return SyncMorphAssociates(db, pivot, foreignKey, foreignValue, joinKey, joinValues, "", "")
}

func SelectPivots(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string) (result *gorm.DB) {

	return SelectMorphPivots(db, pivot, foreignKey, foreignValue, "", "")
}

func SelectPivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string) (result *gorm.DB) {

	return SelectMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValue, "", "")
}

func SavePivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string) (err error) {
	return SaveMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValue, "", "")
}

func UpdatePivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string) (err error) {
	return UpdateMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValue, "", "")
}

// --------------------------------------------------------------------

func AppendMorphAssociates(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValues []string,
	ownerKey string, ownerValue string,
) (err error) {

	var result = &gorm.DB{}

	err = db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(joinValues); i++ {

			result = SelectMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValues[i], ownerKey, ownerValue)
			if result.Error != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if result.RowsAffected == 0 || result.Error == gorm.ErrRecordNotFound {
				err = SaveMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValues[i], ownerKey, ownerValue)
				if err != nil {
					return err
				}
			} else {
				err = UpdateMorphPivot(db, pivot, foreignKey, foreignValue, joinKey, joinValues[i], ownerKey, ownerValue)
				if err != nil {
					return err
				}
			}
		}
		return result.Error
	})

	return err
}

func SyncMorphAssociates(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValues []string,
	ownerKey string, ownerValue string,
) (err error) {

	err = db.Transaction(func(tx *gorm.DB) error {

		err = ClearPivots(db, pivot, foreignKey, foreignValue)
		if err != nil {
			return err
		}
		err = AppendMorphAssociates(db, pivot, foreignKey, foreignValue, joinKey, joinValues, ownerKey, ownerValue)

		return err
	})

	return err
}

func SelectMorphPivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string,
	ownerKey string, ownerValue string,
) (result *gorm.DB) {

	result = &gorm.DB{}

	// join foreign type if exists
	db = SelectMorphPivots(db, pivot, foreignKey, foreignValue, ownerKey, ownerValue)

	result = db.Where(joinKey, joinValue)

	return result
}

func SelectMorphPivots(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	ownerKey string, ownerValue string,
) (result *gorm.DB) {

	result = &gorm.DB{}

	// join foreign type if exists
	strWhereOwner := ""
	strWhere := " WHERE " + foreignKey + "=?"
	if ownerKey != "" && ownerValue != "" {
		strWhereOwner = " AND " + ownerKey + "=" + ownerValue
		strWhere += " ?"
		result = db.
			Debug().
			Exec("select * from "+pivot.GetTableName(true)+strWhere, foreignValue, strWhereOwner)
	} else {
		result = db.
			Debug().
			Exec("select * from "+pivot.GetTableName(true)+strWhere, foreignValue)
	}

	return result
}

func SaveMorphPivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string,
	ownerKey string, ownerValue string,
) (err error) {
	now := time.Now()

	if ownerKey != "" && ownerValue != "" {
		strValue := " (" + foreignKey + ", " + joinKey + ", " + ownerKey + ", created_at, updated_at ) VALUES (?, ?, ?, ?, ?)"
		db = db.
			Debug().
			Exec("INSERT INTO "+pivot.GetTableName(true)+strValue, foreignValue, joinValue, ownerValue, now, now)
	} else {
		strValue := " (" + foreignKey + ", " + joinKey + ", created_at, updated_at ) VALUES (?, ?, ?, ?)"
		db = db.
			Debug().
			Exec("INSERT INTO "+pivot.GetTableName(true)+strValue, foreignValue, joinValue, now, now)
	}

	return db.Error
}

func UpdateMorphPivot(db *gorm.DB, pivot ModelInterface,
	foreignKey string, foreignValue string,
	joinKey string, joinValue string,
	ownerKey string, ownerValue string,
) (err error) {
	now := time.Now()

	strSet := " SET updated_at = ?"

	// join foreign type if exists
	strWhere := " WHERE " + foreignKey + " = ? AND " + joinKey + "=?"
	strWhereOwner := ""
	if ownerKey != "" && ownerValue != "" {
		strWhereOwner = " AND " + ownerKey + "=" + ownerValue
		strWhere += " ?"
		db = db.
			Debug().
			Exec("UPDATE "+pivot.GetTableName(true)+strSet+strWhere, now, foreignValue, joinValue, strWhereOwner)
	} else {
		db = db.
			Debug().
			Exec("UPDATE "+pivot.GetTableName(true)+strSet+strWhere, now, foreignValue, joinValue)
	}

	return db.Error
}

func ClearPivots(db *gorm.DB, pivot ModelInterface, foreignKey string, foreignValue string) (err error) {
	result := db.
		Debug().
		Exec("DELETE FROM "+pivot.GetTableName(true)+" WHERE "+foreignKey+"=?", foreignValue)
	if result.Error != nil {
		return result.Error
	}

	return nil

}
