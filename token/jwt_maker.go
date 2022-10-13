package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct{
	secretKey string
}

func NewJWTMaker(secretKey string)(Maker, error){
	if len(secretKey) < minSecretKeySize{
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration)(string,*Payload,error){
	payload ,err := NewPayload(username, duration)
	if err!=nil{
		return "", payload,err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256,payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string)(*Payload,error){
	//校验token是不是HMAC方式加密的
	keyFunc := func(token *jwt.Token)(interface{},error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token,&Payload{},keyFunc)
	
	if err != nil{
		verr, ok:= err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner,ErrExpiredToken){
			return nil,ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload,ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload,nil
}