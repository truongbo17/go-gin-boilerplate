package types

type ListRolesInput struct {
	Search  string `json:"search"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

type CreateRoleInput struct {
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type UpdateRoleInput struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type DeleteRoleInput struct {
	ID uint `json:"id"`
}

type AssignPermissionToRoleInput struct {
	RoleID        uint   `json:"role_id"`
	PermissionIDs []uint `json:"permission_ids"`
}

type AssignRoleToUserInput struct {
	UserID  uint   `json:"user_id"`
	RoleIDs []uint `json:"role_ids"`
}
