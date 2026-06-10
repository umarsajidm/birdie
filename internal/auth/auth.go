package auth

import (
	"time"
	"fmt"
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


//HashPassword creates the hash using argon2id of the password which is given as str
func HashPassword(password string) (string, error) {
	
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

//CheckPasswordHash checks by comparing the password entered with the password that is saved and converted to hash
func CheckPasswordHash(password, hash string) (bool, error) {

	return argon2id.ComparePasswordAndHash(password, hash)
}

//MakeJWT creates new json webtoken for specific user
func MakeJWT(userID uuid.UUID, tokenSecret string, ExpiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		Issuer:	"birdie-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(ExpiresIn)),
		Subject: userID.String(),
	}
	//create the token using the H256 algo and our claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	//sign it using the secret key (converted to bytes)
	return token.SignedString([]byte(tokenSecret))
}
