package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type PetClaims struct {
	UserID string
	jwt.StandardClaims
}

var SecretKey = "pet-adoption"

func GetToken(uid string) string {
	expireToken := time.Now().Add(time.Hour * 24 * 30).Unix()
	claims := PetClaims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "pet-adoption",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(SecretKey))
	return signedToken
}

func GetUID(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &PetClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*PetClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("非法操作")

}

func JwtAuthFilter(ctx *context.Context) {
	token, err := request.ParseFromRequestWithClaims(ctx.Request,
		request.AuthorizationHeaderExtractor,
		&PetClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
	if err != nil || token.Claims.(*PetClaims).UserID == "" {
		ctx.Output.Status = 401
		ctx.Output.JSON(err, false, false)
		return
	}

	uid := token.Claims.(*PetClaims).UserID
	ctx.Input.SetData("uid", uid)
}
