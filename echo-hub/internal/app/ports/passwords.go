package ports

type PasswordManager interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password, hashedPassword string) error
}
