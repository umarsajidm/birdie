package auth

import (
	"github.com/alexedwards/argon2id"
)


//HashPassword creates the hash using argon2id of the password which is given as str
func HashPassword(password string) (string, error) {
	
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

//CheckPasswordHash checks by comparing the password entered with the password that is saved and converted to hash
func CheckPasswordHash(password, hash string) (bool, error) {

	return argon2id.ComparePasswordAndHash(password, hash)
}