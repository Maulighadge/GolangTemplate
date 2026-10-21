package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	BaseModel
	Name     string `gorm:"size:255;not null" json:"name"`
	Email    string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     Role   `gorm:"size:50;default:'user';not null" json:"role"`
	IsActive bool   `gorm:"default:true;not null" json:"is_active"`

	// Relationships
	Products []Product `gorm:"foreignKey:UserID" json:"products,omitempty"`
}

func (User) TableName() string {
	return "users"
}
