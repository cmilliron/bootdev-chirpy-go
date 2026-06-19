package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Println(ss, err)
		return "", err
	}
	// fmt.Printf("New Token: %s\n", ss)
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, err	
	}	
	
	idStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, err	
	}	
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(err)
		return uuid.UUID{}, err	
	}	
	

	return id, nil 
}