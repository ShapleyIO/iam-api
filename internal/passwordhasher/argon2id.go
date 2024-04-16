package passwordhasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/ShapleyIO/iam-api/internal/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/argon2"
)

var _ PasswordHasher = (*ArgonHasher)(nil)

type ArgonHasher struct {
	saltLength  int
	keyLength   uint32
	time        uint32
	memory      uint32
	threads     uint8
	hashVariant string
}

func NewArgonHasher(cfg *config.Config) *ArgonHasher {
	return &ArgonHasher{
		saltLength:  cfg.PasswordHasher.SaltLength,
		keyLength:   cfg.PasswordHasher.KeyLength,
		time:        cfg.PasswordHasher.Time,
		memory:      cfg.PasswordHasher.Memory,
		threads:     cfg.PasswordHasher.Threads,
		hashVariant: cfg.PasswordHasher.HashVariant,
	}
}

func (ph *ArgonHasher) HashPassword(password string) string {
	salt, err := generateSalt(ph.saltLength)
	if err != nil {
		log.Panic().Err(err).Msg("failed to generate salt")
	}
	hash := argon2.IDKey([]byte(password), salt, ph.time, ph.memory, ph.threads, ph.keyLength)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	format := "$%s$v=%d$m=%d,t=%d,p=%d$%s$%s"
	fullHash := fmt.Sprintf(format, ph.hashVariant, argon2.Version, ph.memory, ph.time, ph.threads, b64Salt, b64Hash)
	return fullHash
}

func (ph *ArgonHasher) Compare(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}
	c := &argon2Config{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.memory, &c.time, &c.threads)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	otherHash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateSalt(saltLength int) ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

type argon2Config struct {
	time    uint32
	memory  uint32
	threads uint8
}
