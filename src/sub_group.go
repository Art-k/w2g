package src

// SubGroup the second layer of groups
type SubGroup struct {
	Model

	Name string

	Team []UserRoles `gorm:"foreignkey:TeamID"`
}
