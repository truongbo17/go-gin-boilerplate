package auth

// Stable auth error codes shared by the auth use cases and HTTP adapter.
const (
	ErrAuthLoginFailed       = 2000
	ErrAuthUserNotFound      = 2001
	ErrAuthWrongPassword     = 2002
	ErrAuthGenerateToken     = 2003
	ErrAuthUserExists        = 2004
	ErrAuthLogoutFailed      = 2005
	ErrAuthRegisterFailed    = 2006
	ErrUserListInternalError = 2007
	ErrAuthResetUnavailable  = 2008
	ErrAuthResetInvalid      = 2009
	ErrAuthResetFailed       = 2010

	ErrChangePass   = 13008
	ErrRefreshToken = 13009

	ErrFindUserFailed   = 13010
	ErrUserNotFound     = 13011
	ErrUpdateUserFailed = 13012

	ErrRoleInternalError    = 14000
	ErrRoleNotFound         = 14001
	ErrRoleCreateFailed     = 14002
	ErrRoleUpdateFailed     = 14003
	ErrRoleDeleteFailed     = 14004
	ErrRoleAssignPermission = 14005
	ErrRoleAssignUser       = 14006

	ErrPermissionInternalError = 15000
	ErrPermissionNotFound      = 15001
	ErrPermissionCreateFailed  = 15002
	ErrPermissionUpdateFailed  = 15003
	ErrPermissionDeleteFailed  = 15004
)
