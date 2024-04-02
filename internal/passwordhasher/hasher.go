package passwordhasher

type PasswordHasher interface {
	HashPassword(password string) string
	ComparePasswordAndHash(password, encodedHash string) (bool, error)
}
