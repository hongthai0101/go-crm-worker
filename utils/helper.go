package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"time"
)

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func ExtractClaims(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return token, nil

	})

	return "Error Parsing Token: ", err
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}
	return "unable to extract claims", nil
}

func Random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func Hash(s interface{}) string {
	var b bytes.Buffer
	_ = gob.NewEncoder(&b).Encode(s)
	return fmt.Sprintf("%x", md5.Sum(b.Bytes()))
}

func Shorthand[T string](first T, second T) T {
	if "" != first {
		return first
	}
	return second
}
