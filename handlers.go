package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)
// set a jwt secret key to a config file
//var jwtKey = []byte("my_secret_key")

// create two users and save to DB， as simple, save in the golang map
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
// Create a struct that models the structure of a claims(payload)
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Create a struct that models the structure of a json result
type JsonResult struct {
	Code int `json:"code"`
	Token string `json:"token"`
}


// *http.Request表示http请求对象， http.ResonseWriter是一个接口类型
func Signin(w http.ResponseWriter, r *http.Request)  {
	var jwtKey = []byte(conf.JwtKey)
	//creds := Credentials{Password:"123456", Username:"zyh"}
	// json.Marshal() 编码为json格式
	//str, _ := json.Marshal(creds)
	//fmt.Printf("%s\n", str)
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)	//&取址符号
	if err != nil {
		// if the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("%s\n", creds)
	fmt.Printf("creds's type: %T\n", creds)
	// Get the expected password from our in memory map
	exceptedPassword, ok := users[creds.Username]
	// If a password exists for the given user
	// And, if it is the same as the password we received, we can pass
	// if NOT, then we return an "Unauthorized" status
	if !ok || exceptedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Declare the expiration time of the token , we kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT payload(claims), which includes the username and expiry time
	claims :=&Claims{
		Username:       creds.Username,
		StandardClaims: jwt.StandardClaims{
			// in JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt:expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string use the jwtKey
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Finally, we return the tokenString as json to frontend
	msg, _ := json.Marshal(JsonResult{Code:200, Token:tokenString})
	_, _ = w.Write(msg)


	//_, _ = w.Write([]byte("Hello world"))
}

// 处理认证后的路由 /welcome
func Welcome (w http.ResponseWriter, r *http.Request) {
	var jwtKey = []byte(conf.JwtKey)
	// we can get the token from the request headers, which come with every request, string type
	token_in_header := r.Header.Get("token")
	//fmt.Printf("token_in_header:%v, type:%T", token_in_header, token_in_header)
	if token_in_header == "" {
		//fmt.Println("token is not in header! please check it!")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// expire? or valid?
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token_in_header, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	//fmt.Println(err)	// signature is invalid      or       token is expired by 2m30s
	//fmt.Printf("err type:%T\n", err)	// *jwt.ValidationError
	fmt.Println(tkn.Valid)
	//fmt.Println(err.Error()=="signature is invalid")	// err to string
	// expire? or valid?
	if err != nil{
		if err.Error()[:16] == "token is expired"{
			// tell fontend the expired flag
			_, _ = w.Write([]byte("token is expired"))
			return
		}else if err.Error() == "signature is invalid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	// if tkn is vaild, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s", claims.Username)))
}

// when fontend receive the token is expired flag, fontend refresh token and storage in local
func Refresh (w http.ResponseWriter, r *http.Request) {
	var jwtKey = []byte(conf.JwtKey)
	//fmt.Println(conf.JwtKey)
	// expire token is exist
	token_in_header := r.Header.Get("token")
	if token_in_header == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// expire? or valid?
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token_in_header, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	if err != nil{
		if err.Error() == "signature is invalid" {
			w.WriteHeader(http.StatusUnauthorized)
		}else if err.Error()[:16] == "token is expired" {
			expirationTime := time.Now().Add(5 * time.Minute)
			claims.ExpiresAt = expirationTime.Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// return the new tokenString to fountend
			msg, _ := json.Marshal(JsonResult{Code:200, Token:tokenString})
			_, _ = w.Write(msg)
		}
	}
}
