package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"net/http"
	"time"
)

// This package is a very basic fake Google OAuth server used to run our tests without any needs for a real Google Oauth client, neither a Google Account.
// It only has one endpoint : POST /. If users send a code value "my_user_is_granted", it will respond a 200 OK with a fake user with sub = 102950075881792615162.
// Otherwise, it will respond a 401 Unauthorized. All other fields will be ignored. This server is only meant for testing purposes, NEVER PUSH THIS TO PRODUCTION.

func main() {
	http.HandleFunc("/", handle)

	fmt.Println("=== Send request to POST / with code = my_user_is_granted to get a fake JWT token.")
	fmt.Println("=== User is John DOE <john.doe@mailhog.local> with sub 102950075881792615162")
	fmt.Println("=== Send request to POST / with code = my_user_is_admin to get a fake JWT token.")
	fmt.Println("=== User is Admin <my.admin@mailhog.local> with sub 102950075881792615000")
	fmt.Println("Mock server listening on port 8000")
	_ = http.ListenAndServe(":8000", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.FormValue("code") == "my_user_is_granted" {

		token := googleToken{
			AccessToken: generateRandomString(60),
			ExpiresIn:   3341,
			Scope:       "https://www.googleapis.com/auth/userinfo.email openid https://www.googleapis.com/auth/userinfo.profile",
			TokenType:   "Bearer",
			IdToken:     "",
		}

		generatedJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp":            time.Now().Add(time.Hour * time.Duration(24)).Unix(),
			"iat":            time.Now().Unix(),
			"family_name":    "DOE",
			"given_name":     "John",
			"picture":        "https://gravatar.com/avatar/f5c2434739e1cb18e31a7e754be5cdda?size=256",
			"name":           "John DOE",
			"at_hash":        generateRandomString(25),
			"email_verified": true,
			"email":          "john.doe@mailhog.local",
			"hd":             "google",
			"sub":            "102950075881792615162",
			"aud":            "874627747812-8s855g5olr2lks9eedka4vi6469bimte.apps.googleusercontent.com",
			"azp":            "874627747812-8s855g5olr2lks9eedka4vi6469bimte.apps.googleusercontent.com",
			"iss":            "https///accounts.google.com",
		})

		rsaKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
		token.IdToken, _ = generatedJwt.SignedString(rsaKey)

		data, _ := json.Marshal(token)
		_, _ = w.Write(data)

	} else if r.FormValue("code") == "my_user_is_admin" {
		token := googleToken{
			AccessToken: generateRandomString(60),
			ExpiresIn:   3341,
			Scope:       "https://www.googleapis.com/auth/userinfo.email openid https://www.googleapis.com/auth/userinfo.profile",
			TokenType:   "Bearer",
			IdToken:     "",
		}

		generatedJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"exp":            time.Now().Add(time.Hour * time.Duration(24)).Unix(),
			"iat":            time.Now().Unix(),
			"family_name":    "DOE",
			"given_name":     "John",
			"picture":        "https://gravatar.com/avatar/f5c2434739e1cb18e31a7e754be5cdda?size=256",
			"name":           "Admin",
			"at_hash":        generateRandomString(25),
			"email_verified": true,
			"email":          "john.doe@mailhog.local",
			"hd":             "google",
			"sub":            "102950075881792615000",
			"aud":            "874627747812-8s855g5olr2lks9eedka4vi6469bimte.apps.googleusercontent.com",
			"azp":            "874627747812-8s855g5olr2lks9eedka4vi6469bimte.apps.googleusercontent.com",
			"iss":            "https///accounts.google.com",
		})

		rsaKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
		token.IdToken, _ = generatedJwt.SignedString(rsaKey)

		data, _ := json.Marshal(token)
		_, _ = w.Write(data)
	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type googleToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

var key = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`
