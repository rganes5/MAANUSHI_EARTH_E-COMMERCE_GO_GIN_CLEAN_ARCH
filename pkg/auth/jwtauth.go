package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rganes5/maanushi_earth_e-commerce/pkg/config"
)

type Claims struct {
	Email string
	ID    uint
	jwt.RegisteredClaims
}

func GenerateJWT(email string, ID uint) (string, error) {

	// expireTime := time.Now().Add(60 * time.Minute)
	expiryTime := time.Now().Add(10 * time.Minute)

	// create token with expire time and claims id as user id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email: email,
		ID:    ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	})

	// convert the token into signed string
	tokenString, err := token.SignedString([]byte(config.GetJWTCofig()))

	if err != nil {
		return "", err
	}
	// refresh token add next time
	return tokenString, nil
}

// func ValidateToken(tokenString string) (Claims, error) {
// 	claims := Claims{}
// 	token, err := jwt.ParseWithClaims(tokenString, &claims,
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(config.GetJWTCofig()), nil
// 		},
// 	)
// 	//checking the expiry of the token
// 	if time.Now().Unix() > claims.ExpiresAt.Unix() {
// 		return claims, errors.New("token expired re-login")
// 	}
// 	if err != nil || !token.Valid {
// 		return claims, errors.New("not valid token")
// 	}
// 	return claims, nil
// }

func ValidateToken(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTCofig()), nil
		},
	)
	if err != nil || !token.Valid {
		return claims, errors.New("not valid token")
	}
	return claims, nil
}

//There are some changes to be made,
//Tokenn has to be set so that if user reloads or refreshes within the expiry , then it should be exteneded and if not then it
//should ask the user to relogin.
