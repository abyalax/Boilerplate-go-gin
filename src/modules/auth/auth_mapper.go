package auth

type UserRolePermission struct {
	Name        string
	Email       string
	Roles       []Role
	Permissions []Permission
}

func MapUser(rows []GetUserWithPermissionsRow) *UserRolePermission {
	if len(rows) == 0 {
		return nil
	}

	user := &UserRolePermission{
		Name:        rows[0].UserName,
		Email:       rows[0].UserEmail,
		Roles:       []Role{},
		Permissions: []Permission{},
	}

	roleMap := map[int32]bool{}
	permMap := map[int32]bool{}

	for _, r := range rows {

		if r.RoleID.Valid && !roleMap[r.RoleID.Int32] {
			user.Roles = append(user.Roles, Role{
				ID:   r.RoleID.Int32,
				Name: r.RoleName.String,
			})
			roleMap[r.RoleID.Int32] = true
		}

		if r.PermissionID.Valid && !permMap[r.PermissionID.Int32] {
			user.Permissions = append(user.Permissions, Permission{
				ID:   r.PermissionID.Int32,
				Key:  r.PermissionKey.String,
				Name: r.PermissionName.String,
			})
			permMap[r.PermissionID.Int32] = true
		}
	}

	return user
}
