package datacontract

type UserType int

const (
	Admin UserType = iota
	Teacher
	Student
)

type User struct {
	ID           string `bson:"_id,omitempty"`
	Name         string
	UserType     UserType
	PasswordSalt string
	PasswordHash string
}

type LoginCredentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool
	User    User
}
