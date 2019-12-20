package src

// Project the area where all feachers are saved / linked
type Project struct {
	Model

	Name string

	Team []UserRoles `gorm:"foreignkey:TeamID"`
}
