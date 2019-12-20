package src

import (
// "github.com/google/uuid"
)

// SetAllPermissionsToByRoleIDifNotExists set all permissions to one value
func SetAllPermissionsToByRoleIDifNotExists(roleID string, route string, permission bool) {

	var RP RolePermission
	Db.Where("role_id = ?", roleID).Where("end_point = ?", route).Last(&RP)
	if RP.EndPoint == "" {
		RP.Post = permission
		RP.Get = permission
		RP.Patch = permission
		RP.Delete = permission
		RP.EndPoint = route
		RP.RoleID = roleID
		Db.Create(&RP)
	}

}

// IfUserHasPermission check if user has permission to do ...
func IfUserHasPermission(currentUser User, route string, method string) bool {

	var up UserPermission
	Db.Where("user_id = ?", currentUser.ID).Where("end_point = ?", route).Last(&up)

	var Answer bool = false
	if up.EndPoint != "" {

		switch method {
		case "OPTIONS":
			Answer = true
		case "GET":
			Answer = up.Get
		case "POST":
			Answer = up.Post
		case "PATCH":
			Answer = up.Patch
		case "DELETE":
			Answer = up.Delete
		default:
			Answer = false
		}
	} else {
		var rp RolePermission
		Db.Where("role_id = ?", currentUser.RoleID).Where("end_point = ?", route).Last(&rp)

		if rp.EndPoint != "" {

			switch method {
			case "OPTIONS":
				Answer = true
			case "GET":
				Answer = rp.Get
			case "POST":
				Answer = rp.Post
			case "PATCH":
				Answer = rp.Patch
			case "DELETE":
				Answer = rp.Delete
			default:
				Answer = false
			}
		}
	}

	return Answer
}
