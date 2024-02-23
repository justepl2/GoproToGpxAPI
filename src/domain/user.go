package domain

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/justepl2/gopro_to_gpx_api/interfaces/request"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "ROLE_ADMIN"
	RoleUser  Role = "ROLE_USER"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"column:email"`
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	FirstName string    `gorm:"column:firstname"`
	LastName  string    `gorm:"column:lastname"`
	Role      Role      `gorm:"column:role"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "user"
}

func (du *User) FromRequest(s request.Signup) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	du.Email = s.Email
	du.Username = s.Username
	du.Password = string(hashedPassword)
	du.FirstName = s.FirstName
	du.LastName = s.LastName
	du.Role = RoleUser

	return nil
}

func (du *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(du.Password), []byte(password))
}

func CreateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
