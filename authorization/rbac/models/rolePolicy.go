package models

const (
	RBAC_CONTROL_ALL    = "all"
	RBAC_CONTROL_WRITE  = "write"
	RBAC_CONTROL_READ   = "read"
	RBAC_CONTROL_DELETE = "delete"
	RBAC_CONTROL_NONE   = "none"
)

type RolePolicy struct {
	RoleID   string `gorm:"column:policy_id;not null;index" json:"roleID"`
	ObjectID string `gorm:"column:object_id;not null;index" json:"objectID"`
	Control  string `gorm:"column:control" json:"control"`
}
