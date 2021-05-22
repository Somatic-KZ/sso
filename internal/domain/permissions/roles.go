package permissions

const (
	CreateRoles = "roles-create"
	ViewRoles   = "roles-view"
	UpdateRoles = "roles-update"
	DeleteRoles = "roles-delete"
)

var RolesPermissions = []string{CreateRoles, ViewRoles, UpdateRoles, DeleteRoles}
