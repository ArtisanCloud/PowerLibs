package database

import "github.com/ArtisanCloud/PowerLibs/v2/object"

// TableName overrides the table name used by price_book to `profiles`
func (mdl *PowerOperationLog) TableName() string {
	return mdl.GetTableName(true)
}

// PowerOperationLog 数据表结构
type PowerOperationLog struct {
	*PowerCompactModel

	OperatorName  *string `gorm:"column:operatorName" json:"operatorName"`
	OperatorTable *string `gorm:"column:operatorTable" json:"operatorTable"`
	OperatorID    *string `gorm:"column:operatorID;index" json:"operatorID"`
	Module        *int16  `gorm:"column:module" json:"module"`
	Operate       *string `gorm:"column:operate" json:"operate"`
	Event         *int8   `gorm:"column:event" json:"event"`
	ObjectName    *string `gorm:"column:objectName" json:"objectName"`
	ObjectTable   *string `gorm:"column:objectTable" json:"objectTable"`
	ObjectID      *string `gorm:"column:objectID;index" json:"objectID"`
	Result        *int8   `gorm:"column:result" json:"result"`
}

const TABLE_NAME_OPERATION_LOG = "power_operation_log"
const OPERAION_LOG_UNIQUE_ID = COMPACT_UNIQUE_ID

func NewPowerOperationLog(mapObject *object.Collection) *PowerOperationLog {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	return &PowerOperationLog{
		PowerCompactModel: NewPowerCompactModel(),
		OperatorName:      mapObject.GetStringPointer("operatorName", ""),
		OperatorTable:     mapObject.GetStringPointer("operatorTable", ""),
		OperatorID:        mapObject.GetStringPointer("operatorID", ""),
		Module:            mapObject.GetInt16Pointer("module", 0),
		Operate:           mapObject.GetStringPointer("operate", ""),
		Event:             mapObject.GetInt8Pointer("event", 0),
		ObjectName:        mapObject.GetStringPointer("objectName", ""),
		ObjectTable:       mapObject.GetStringPointer("objectTable", ""),
		ObjectID:          mapObject.GetStringPointer("objectID", ""),
		Result:            mapObject.GetInt8Pointer("result", 0),
	}
}

// 获取当前 Model 的数据库表名称
func (mdl *PowerOperationLog) GetTableName(needFull bool) string {
	tableName := TABLE_NAME_OPERATION_LOG
	if needFull {
		tableName = "public." + tableName
	}
	return tableName
}

func (mdl *PowerOperationLog) SaveOps(operator ModelInterface, module int16, operate string, object ModelInterface) {

}
