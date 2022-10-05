package models

import (
	"github.com/laterius/service_architecture_hw3/app/modules/hash"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

// NewUserService func for creating connection to database
func NewUserService(db *gorm.DB) (UserService, error) {
	ug, err := newUserGorm(db)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(HmacSecret)
	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}
	return &userService{
		UserDB: uv,
	}, nil
}

// Implementing and returning the userGorm type
func newUserGorm(db *gorm.DB) (*userGorm, error) {
	return &userGorm{
		db: db,
	}, nil
}

// ByID method to get a user by ID
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id).First(&user)
	err := first(db, &user)
	return &user, err
}

// ByEmail method to get a user by Email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email).First(&user)
	err := first(db, &user)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This mdethos expects the remember
// token to be already hashed.
// Errors handeled as the same done by the ByEmail.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// will query the gorm.DB and get the first item from db and place it into
// dst , if nothing is found , it will return error.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Create method to create a user in database
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update method to update a user in database
func (ug *userGorm) Update(user *User) error {
	//return ug.db.Model(&user).Where("name = ?", &user.Name).Update(&user).Error TODO
	return nil
}

// Delete method to delete a user in database with the provided ID only
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Where("id = ?", id).Delete(&user).Error
}

// Authenticate Method is used for Authenticate and Validate login
func (us *userService) Authenticate(email, password string) (*User, error) {
	// Vlidate if the user is existed in the database or no
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	// Compare the login based in the Hash value
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+UserPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		case nil:
			return nil, err
		default:
			return nil, err
		}
	}
	return foundUser, nil
}
