package newsfeed

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var hmacSampleSecret []byte

func loadSecret() {
	// Load sample key data
	if keyData, e := ioutil.ReadFile("secret"); e == nil {
		hmacSampleSecret = keyData
	} else {
		panic(e)
	}
}

// loadSecret()

func getTokenString() (string, error) {
	ttl := 60 * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		"exp": ttl,
	})

	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

var token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"foo": "bar",
	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
})

func validateToken() {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

}
