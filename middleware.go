package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// ? Generate a 6-digit random number
func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}

// ? JWT Secret
var JWT_SECRET = os.Getenv("JWT_SECRET")

var jwtSecret = []byte(JWT_SECRET)

// ? Generate a JWT token
func generateJWTToken(user User) (string, error) {
	id := user.ID
	//? Create the claims containing the email
	claims := jwt.MapClaims{
		"id": id,
		// StandardClaims: jwt.StandardClaims{
		// 	ExpiresAt: time.Now().Add(time.Hour).Unix(), // Token expiry time
		// },
	}

	//? Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//? Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ? Verify a JWT token
func verifyJWTToken(tokenString string) (*MyCustomClaims, error) {
	//? Create an instance of MyCustomClaims with embedded StandardClaims as a pointer
	claims := &MyCustomClaims{
		StandardClaims: &jwt.StandardClaims{},
	}

	//? Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	//? Check if the token is valid and the claims were successfully extracted
	if !token.Valid || claims == nil || claims.StandardClaims == nil {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ? Cache middleware
func (c *CacheStruct) Set(key string, value interface{}) error {
	c.cache.Store(key, value)
	return nil
}

func (c *CacheStruct) Get(key string) (interface{}, error) {
	value, ok := c.cache.Load(key)
	if !ok {
		return nil, fmt.Errorf("Key not found: %s", key)
	}
	return value, nil
}

func (c *CacheStruct) Delete(key string) {
	c.cache.Delete(key)
}

func setOtpVerification(cache Cache, email string, verification *OTPVerification) error {
	err := cache.Set(email, verification)
	if err != nil {
		return fmt.Errorf("Failed to set OTP verification: %v", err)
	}
	return nil
}

func getOtpVerification(cache Cache, email string) (*OTPVerification, error) {
	value, err := cache.Get(email)
	if err != nil {
		return nil, fmt.Errorf("Failed to get OTP verification: %v", err)
	}

	otpVerification, ok := value.(*OTPVerification)
	if !ok {
		return nil, fmt.Errorf("Invalid data type for OTP verification")
	}

	return otpVerification, nil
}
