package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

func genRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

func hashPassword(pass string, salt []byte) string {
	var passBytes = []byte(pass)

	var sha512Hasher = sha512.New()
	passBytes = append(passBytes, salt...)

	sha512Hasher.Write(passBytes)
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func isPassMatch(hashedPass string, curPass string, salt []byte) bool {
	var hashedCurPass = hashPassword(curPass, salt)
	return hashedCurPass == hashedPass
}