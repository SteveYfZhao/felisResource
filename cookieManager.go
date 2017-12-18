package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const keyLength int = 16

var keyInitialized = false
var currentEncryKey []byte
var previousEncryKey []byte
var lastKeyUpdateTime time.Time = time.Time{}

const refreshCycle = 2 * time.Hour
const separator = " "

func Login3rdParty(loginType string) bool {

	return true
}

func GenerateNewCookie(w http.ResponseWriter, cookiekey string, cookievalue string) {
	ct := time.Now()
	fmt.Println(ct)
	ck := &http.Cookie{Name: cookiekey, Value: cookievalue}
	http.SetCookie(w, ck)
	//expiration := time.Now().Add(365 * 24 * time.Hour)
	//cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
	//http.SetCookie(w, &cookie)
	fmt.Println("Set cookie", cookiekey, cookievalue)
}

func EncodeUserDataToCipherCookie(w http.ResponseWriter, uData map[string]string) {
	plainJsonString, err := json.Marshal(uData)
	if err != nil {
		log.Fatal(err)
	}
	cipherString := encryptString(string(plainJsonString))
	ct := time.Now()
	fmt.Println(ct)
	ck := &http.Cookie{Name: "data", Value: cipherString}
	http.SetCookie(w, ck)
	//expiration := time.Now().Add(365 * 24 * time.Hour)
	//cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
	//http.SetCookie(w, &cookie)
	fmt.Println("Set encrypted cookie", "data", cipherString)
}

func DecodeCipherCookieToUserData(r *http.Request) (map[string]string, error) {
	cipherString, err := GetUserCookie(r, "data")
	if err != nil {
		fmt.Println("failed to get cookie", err)
		return nil, err
	}
	uData := make(map[string]string)
	plainJsonString, err := decryptString(cipherString)
	err = json.Unmarshal([]byte(plainJsonString), &uData)
	if err != nil {
		fmt.Println("failed to get cookie", err)
		return nil, err
	}
	return uData, nil
}

func GetUserCookie(r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie(key)

	if err != nil {
		fmt.Println("failed to get cookie", err)
		return "", err
	} else {
		fmt.Println("receive cookie")
		if cookie != nil {
			fmt.Println(cookie.Name)
			fmt.Println(cookie.Value)
			return cookie.Value, nil
		}
		return "", nil
	}
	/*
			fmt.Println("List of all cookies")
			for _, ce := range r.Cookies() {
		        fmt.Println(ce.Name)
			}
	*/
}

func GetUserNamefromCookie(r *http.Request) (string, error) {
	//return GetUserCookie(r, "uid")
	uData, err := DecodeCipherCookieToUserData(r)
	return uData["uid"], err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Try loading the master key for DB from "masterKey.txt"
func loadMasterDBKey() (string, error) {
	filename := "masterKey.txt"
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(key), nil
}

func encryptString(value string) string {
	nBig, err := rand.Int(rand.Reader, big.NewInt(37))
	if err != nil {
		panic(err)
	}
	nBig2, err := rand.Int(rand.Reader, big.NewInt(37))
	if err != nil {
		panic(err)
	}
	n := int(nBig.Int64()) + 16
	n2 := int(nBig2.Int64()) + 16
	prefix, err := GenerateRandomString(n)
	if err != nil {
		panic(err)
	}
	appendix, err := GenerateRandomString(n2)
	if err != nil {
		panic(err)
	}
	key, _ := GetEncryptionKeys()
	plainText := prefix + separator + value + separator + appendix
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainText))

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decryptString(cryptoText string) (string, error) {
	currentKey, prevKey := GetEncryptionKeys()
	rawResult := decryptStringCore(cryptoText, currentKey)
	trialResult, err := unpackDecryptedString(rawResult)
	if err != nil || !isJSON(trialResult) {
		// decryption failed, try prevkey
		rawResult = decryptStringCore(cryptoText, prevKey)
		trialResult, err = unpackDecryptedString(rawResult)
	}
	return trialResult, err
}

func unpackDecryptedString(target string) (string, error) {
	startIndex := strings.Index(target, separator)
	if startIndex == -1 {
		return "", errors.New("Decrypted String is Invalid")
	}
	endIndex := strings.LastIndex(target, separator)
	if endIndex == -1 {
		return "", errors.New("Decrypted String is Invalid")
	}
	content := target[startIndex:endIndex]
	return content, nil
}

func decryptStringCore(cryptoText string, key []byte) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

// cookie = aes (json (rand1 string, username string, rand2 string)) + hmac
// key for aes = rand3 string // update every 8 hours, last one key is saved.
// cookies encrypted with the old key gets re-generated when the aes key changes.
// if cookie is generated two key cycles ago, cookie gets cleared, user needs to re-login
// research if iv needs random refresh as well.

// to avoid race condition, consider using a slice to store aes keys. cookie functions should
// fetch last two values. or use updating flags and sleep to avoid racing condition.

func GetEncryptionKeys() (currnt []byte, previous []byte) {
	if !keyInitialized {
		previousEncryKey, _ = GenerateRandomBytes(keyLength)
		currentEncryKey, _ = GenerateRandomBytes(keyLength)
		lastKeyUpdateTime = time.Now()
		keyInitialized = true
	}

	if lastKeyUpdateTime.Add(time.Duration(refreshCycle)).Before(time.Now()) {
		// regen keys

		newKey, err := GenerateRandomBytes(keyLength)
		if err == nil {
			previousEncryKey = currentEncryKey
			currentEncryKey = newKey
			lastKeyUpdateTime = time.Now()
		} else {
			panic(err)
		}

	}

	return currentEncryKey, previousEncryKey
}
