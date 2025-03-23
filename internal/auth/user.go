package auth

type User struct {
	ID       string
	Username string
	Password string // Store hashed passwords in real apps
}
