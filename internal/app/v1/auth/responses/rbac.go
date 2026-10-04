package responses

type UserRolesResponse struct {
	Roles []RoleResponse `json:"roles"`
}

type RolePermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
}
