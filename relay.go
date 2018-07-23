package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/nanobox-io/golang-scribble"
	"github.com/ugorji/go/codec"
)

var db *scribble.Driver

const dbColletion = "users"

type (
	user struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	startRegistrationResponse struct {
		Challenge string `json:"challenge"`
		User      user   `json:"user"`
	}
	finishRegistrationBody struct {
		ID       string                        `json:"id"`
		RawID    string                        `json:"rawID"`
		Type     string                        `json:"type"`
		Response finishRegistrationBodyResonse `json:"response"`
	}
	finishRegistrationBodyResonse struct {
		AttestationObject string `json:"attestationObject"`
		ClientDataJSON    string `json:"clientDataJSON"`
	}

	dbItem struct {
		user      user
		challenge []byte
	}
)

func main() {
	initDatabase()

	e := echo.New()
	e.HideBanner = true
	e.GET("/", indexHandle)
	e.POST("/registration/start", startRegistration)
	e.POST("/registration/finish", finishRegistration)

	//Handle finish registration
	//Handle start authentication
	// Handle finish authentication

	e.GET("/users", listUsers)

	fmt.Println("Starting server on port :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func logError(err error) error {
	log.Println(err)
	return err
}

func indexHandle(c echo.Context) error {
	return c.File("index.html")
}

func startRegistration(c echo.Context) error {
	log.Println("Starting a registration")

	u := user{}
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		return logError(err)
	}
	if u.Name == "" {
		return logError(errors.New("username required"))
	}
	// u.ID = uuid.New().String()
	u.ID = base64Encode([]byte(uuid.New().String()))

	chal := createChallenge()
	log.Printf("Saving registration data %#v challenge:%s", u, base64Encode(chal))
	db.Write(dbColletion, u.Name, dbItem{user: u, challenge: chal})
	r := startRegistrationResponse{
		User:      u,
		Challenge: base64Encode(chal),
	}

	return c.JSON(http.StatusCreated, r)
}

func finishRegistration(c echo.Context) error {
	b := finishRegistrationBody{}
	if err := c.Bind(&b); err != nil {
		return logError(err)
	}
	log.Println("finish id:", b.ID)
	log.Println("finish rawID:", b.RawID)
	log.Println("finish type:", b.Type)

	log.Printf("finish ClientDataJSON decode: %#v\n", decodeClientData(b.Response.ClientDataJSON))
	log.Printf("finish AttestationObject: %#v\n", decodeAttestationObject(b.Response.AttestationObject))

	decodeAttestationObject(b.Response.AttestationObject)

	// https://w3c.github.io/webauthn/#registering-a-new-credential

	return c.NoContent(http.StatusCreated)
}

type (
	// https://developer.mozilla.org/en-US/docs/Web/API/AuthenticatorResponse/clientDataJSON
	clientData struct {
		Type      string `json:"type"`      // "webauthn.create" or "webauthn.get"
		Challenge string `json:"challenge"` // base64 encoded String containing the original challenge
		Origin    string `json:"origin"`    // the window.origin
	}
)

func decodeClientData(s string) clientData {
	c := clientData{}
	b := base64Decode(s)
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	return c
}

type (
	attStmt struct {
		Sig []uint8       `json:"sig"`
		X5c []interface{} `json:"x5c"`
	}

	// https://developer.mozilla.org/en-US/docs/Web/API/AuthenticatorAttestationResponse/attestationObject
	attestation struct {
		Fmt      string  `json:"fmt"`
		AuthData []byte  `json:"authData"`
		AttStmt  attStmt `json:"attStmt"`
	}
)

func decodeAttestationObject(s string) attestation {
	cbor := codec.CborHandle{}
	a := attestation{}
	dec := codec.NewDecoder(bytes.NewReader(base64Decode(s)), &cbor)
	err := dec.Decode(&a)
	if err != nil {
		panic(err)
	}
	return a
}

func listUsers(c echo.Context) error {
	users := []dbItem{}

	all, err := db.ReadAll(dbColletion)
	if err != nil {
		return logError(err)
	}
	for _, r := range all {
		u := dbItem{}
		if err := db.Read(dbColletion, r, &u); err != nil {
			return logError(err)
		}
		users = append(users, u)
	}
	return c.JSON(http.StatusOK, users)
}

func createChallenge() []byte {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func base64Decode(str string) []byte {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return data
}

func initDatabase() {
	var err error
	db, err = scribble.New("data", &scribble.Options{})
	if err != nil {
		panic(err)
	}
}
