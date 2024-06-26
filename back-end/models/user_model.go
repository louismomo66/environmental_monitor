package models

import "gorm.io/gorm"

type GormUserRepo struct{
	DB *gorm.DB
}

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

func NewUserRepo(db *gorm.DB) * GormUserRepo{
	return &GormUserRepo{DB : db}
}

type UserRepository interface {
	GetUserByEmail( email string) (User,error)
	CreateUser(user *User) error
	GetUserEmail( email string) (string,error)
	GetAllUsers() ([]User, error)
	UpdatePasswordByEmail(email, hashedPassword string) error 
	
}

func (repo *GormUserRepo) GetUserByEmail(email string) (User, error) {
	var user User
	err := repo.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// In your repository file
func (repo *GormUserRepo) GetAllUsers() ([]User, error) {
	var users []User
	err := repo.DB.Find(&users).Error
	return users, err
}

func (repo *GormUserRepo) GetUserEmail( email string) (string,error) {
	var user User
	err := repo.DB.Where("email=?",email).First(&user).Error
	if err != nil {
		return "",err
	}
	return user.Email, err
}
func (repo *GormUserRepo) CreateUser(user *User) error{
	return repo.DB.Create(user).Error
}

func (repo *GormUserRepo)UpdatePasswordByEmail(email, hashedPassword string) error {
	return repo.DB.Model(&User{}).Where("email = ?", email).Update("password",hashedPassword).Error
}