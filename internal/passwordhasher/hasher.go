package passwordhasher

type PasswordHasher interface {
	HashPassword(password string) string
	Compare(password, encodedHash string) (bool, error)
}
