package permissions

const (
	CreateActions = "actions-create"
	ViewActions   = "actions-view"
	UpdateActions = "actions-update"
	DeleteActions = "actions-delete"
)

var ActionsPermissions = []string{CreateActions, ViewActions, UpdateActions, DeleteActions}
