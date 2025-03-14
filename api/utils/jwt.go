package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rmsubekti/indico/config"
)

type Claim struct {
	ID         uint    `json:"id"`
	Role       string  `json:"role"`
	ExpireDays float32 `json:"expire_days,omitempty"`
	Token      string  `json:"token,omitempty"`
}

var secret = config.APP.JWT_SECRET

func (p *Claim) CreateToken() (err error) {
	var jwtToken string
	mySigningKey := []byte(secret)

	if p.ExpireDays == 0 {
		p.ExpireDays++
	}

	// Create the Claim
	claim := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add((time.Hour * 24) * time.Duration(p.ExpireDays))),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ID:        strconv.Itoa(int(p.ID)),
		Issuer:    p.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	if jwtToken, err = token.SignedString(mySigningKey); err != nil {
		return
	}
	p.Token = jwtToken
	return
}

func (p *Claim) Parse() error {

	if len(p.Token) < 1 {
		return errors.New("no token provided")
	}

	token := strings.SplitN(p.Token, " ", 2)

	if (len(token) < 2) || (token[0] != "Bearer") {
		return errors.New("incorrect format authorization header")
	}

	key, err := jwt.ParseWithClaims(token[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if !key.Valid && err != nil {
		return err
	}

	if claim, ok := key.Claims.(*jwt.RegisteredClaims); ok {
		cid, _ := strconv.Atoi(claim.ID)
		p.ID = uint(cid)
		p.Role = claim.Issuer
		return nil
	}

	return fmt.Errorf("invalid token %s ", err.Error())
}
