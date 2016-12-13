package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/martinekvili/Wolverhampton/datacontract"
)

func findUserByName(users *mgo.Collection, userName string) (user datacontract.User, found bool) {
	err := users.Find(bson.M{"name": userName}).One(&user)
	if err != nil {
		return user, false
	}

	return user, true
}

func createUser(userName string, password string, userType datacontract.UserType) (user datacontract.User, err error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return user, err
	}

	passwordHash, err := scrypt.Key([]byte(password), saltBytes, 16384, 8, 1, 32)
	if err != nil {
		return user, err
	}

	user = datacontract.User{
		Name:         userName,
		UserType:     userType,
		PasswordSalt: base64.StdEncoding.EncodeToString(saltBytes),
		PasswordHash: base64.StdEncoding.EncodeToString(passwordHash),
	}

	return user, nil
}

func CreateUser(userName string, password string, userType datacontract.UserType) bool {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Printf("Couldn't access database: %v\n", err)
		return false
	}
	defer session.Close()

	users := session.DB("WolverhamptonDB").C("User")

	user, found := findUserByName(users, userName)
	if found {
		return false
	}

	user, err = createUser(userName, password, userType)
	if err != nil {
		return false
	}

	err = users.Insert(user)
	if err != nil {
		log.Printf("Couldn't insert user: %v\n", err)
		return false
	}

	return true
}

func ChangePassword(userName string, password string) bool {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Printf("Couldn't access database: %v\n", err)
		return false
	}
	defer session.Close()

	users := session.DB("WolverhamptonDB").C("User")

	user, found := findUserByName(users, userName)
	if !found {
		return false
	}

	userID := user.ID
	user, err = createUser(userName, password, user.UserType)
	if err != nil {
		return false
	}

	err = users.UpdateId(userID, user)
	if err != nil {
		log.Printf("Couldn't insert user: %v\n", err)
		return false
	}

	return true
}

func CheckPassword(userName string, password string) (datacontract.User, bool) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Printf("Couldn't access database: %v\n", err)
		return datacontract.User{}, false
	}
	defer session.Close()

	users := session.DB("WolverhamptonDB").C("User")

	user, found := findUserByName(users, userName)
	if !found {
		return datacontract.User{}, false
	}

	saltBytes, err := base64.StdEncoding.DecodeString(user.PasswordSalt)
	if err != nil {
		return datacontract.User{}, false
	}

	passwordHash, err := scrypt.Key([]byte(password), saltBytes, 16384, 8, 1, 32)
	if err != nil {
		return datacontract.User{}, false
	}

	passwordString := base64.StdEncoding.EncodeToString(passwordHash)
	if passwordString != user.PasswordHash {
		return datacontract.User{}, false
	}

	return user, true
}

func CreateAdminUser() bool {
	return CreateUser("admin", "password", datacontract.Admin)
}
