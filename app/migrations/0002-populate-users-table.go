package mixtures

import (
	"github.com/ezn-go/mixture"
	"github.com/go-gormigrate/gormigrate/v2"
)

type User0002 struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

func (u User0002) TableName() string {
	return "users"
}

func init() {

	users := []User0002{
		{Username: "j.smith", FirstName: "James", LastName: "Smith", Email: "j.smith@mail.com", Phone: "+11111", Password: "qwerty"},
		{Username: "j.johnson", FirstName: "John", LastName: "Johnson", Email: "j.johnson@mail.com", Phone: "+22222", Password: "12345"},
		{Username: "r.williams", FirstName: "Robert", LastName: "Williams", Email: "r.williams@mail.com", Phone: "+33333", Password: "lostpass"},
	}

	mx := &gormigrate.Migration{
		ID:       "0002",
		Migrate:  mixture.CreateBatchM(users),
		Rollback: mixture.DeleteBatchR(users),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
