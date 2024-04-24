package entity

// The User domain it's only for example purpose...
type User struct {
	Id       int    `gorm:"column:id; primary_key; not null" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Email    string `gorm:"column:email;index:idx_email,unique" json:"email"`
	Password []byte `gorm:"column:password;" json:"password"`
	IsActive bool   `gorm:"default:true; column:is_active" json:"isActive"`
	Role     string `gorm:"default:ROLE_TEST; column:role" json:"role"`
	BaseEntity
}
