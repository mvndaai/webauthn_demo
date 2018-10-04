package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/mvndaai/webauthn"
	"github.com/nanobox-io/golang-scribble"
)

var port = flag.String("port", ":8080", "Port the server starts on")
var origin = flag.String("origin", "http://localhost:8080", "Origin used in verification")
var timeout = flag.Int("timeout", 6000, "Time till auth timeout in ms")

var db *scribble.Driver

const dbColletion = "users"

type (
	dbItem struct {
		User         webauthn.UserEntity `json:"user"`
		Challenge    []byte              `json:"challenge"`
		CredentialID string              `json:"credentialId"`
		PublicKey    string              `json:"publicKey"`
	}
)

func main() {
	flag.Parse()
	initDatabase()

	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.String(he.Code, fmt.Sprint(he.Message))
		} else {
			c.NoContent(http.StatusInternalServerError)
		}
		c.Logger().Error(err)
	}

	e.GET("/", indexHandle)
	e.GET("/localstorage", func(c echo.Context) error { return c.File("localstorage.html") })

	e.POST("/registration/start", startRegistration)
	e.POST("/registration/finish", finishRegistration)
	e.POST("/authentication/start", startAuthentication)
	e.POST("/authentication/finish", finishAuthentication)

	e.GET("/users", listUsers)
	e.DELETE("/users/:username", deleteUser)

	e.Logger.Fatal(e.Start(*port))
}

func indexHandle(c echo.Context) error {
	return c.File("index.html")
}

type (
	startRegistrationResponse struct {
		Challenge string              `json:"challenge"`
		User      webauthn.UserEntity `json:"user"`
	}
)

func startRegistration(c echo.Context) error {
	u := webauthn.UserEntity{}
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		return err
	}
	if u.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username required")
	}

	log.Println("Starting registation for:", u.Name)
	chal, err := webauthn.NewChallenge()
	if err != nil {
		return err
	}

	u.ID = []byte(uuid.New().String())
	log.Printf("Saving registration data %#v challenge:%s", u, base64Encode(chal))
	err = db.Write(dbColletion, u.Name, dbItem{User: u, Challenge: chal})
	if err != nil {
		return err
	}

	r := webauthn.RegistrationParts{
		PublicKey: webauthn.PublicKeyCredentialOptions{
			Challenge: chal,
			RP: webauthn.RpEntity{
				Name: "mvndaai-webauth-demo",
			},
			PubKeyCredParams: []webauthn.Parameters{
				webauthn.Parameters{
					Type: webauthn.PublicKeyCredentialTypePublicKey,
					Alg:  -7,
				},
			},
			Timeout:     50000,
			User:        u,
			Attestation: "direct",
		},
	}

	return c.JSON(http.StatusCreated, r)
}

type (
	finishResponse struct {
		webauthn.PublicKeyCredential
		User webauthn.UserEntity `json:"user"`
	}
)

func finishRegistration(c echo.Context) error {
	b := finishResponse{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	entry := dbItem{}
	err := db.Read(dbColletion, b.User.Name, &entry)
	if err != nil {
		return err
	}

	err = webauthn.ValidateRegistration(b.PublicKeyCredential, entry.Challenge, *origin, false)
	if err != nil {
		db.Delete(dbColletion, b.User.Name)
		log.Println("Registation Validation failed", err)
		// return err //TODO enanable once registartion is complete
	}

	entry.Challenge = []byte{}
	entry.CredentialID = string(b.RawID)
	// entry.PublicKey = "TODO-PUBKEY"

	err = db.Write(dbColletion, b.User.Name, entry)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

type (
	startAuthResponse struct {
		Challenge    string              `json:"challenge"`
		CredentialID string              `json:"credentialId"`
		User         webauthn.UserEntity `json:"user"`
	}
)

func startAuthentication(c echo.Context) error {
	b := webauthn.UserEntity{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	entry := dbItem{}
	err := db.Read(dbColletion, b.Name, &entry)
	if err != nil {
		return err
	}

	chal, err := webauthn.NewChallenge()
	if err != nil {
		return err
	}
	entry.Challenge = chal
	err = db.Write(dbColletion, b.Name, entry)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, startAuthResponse{
		Challenge:    base64Encode(chal),
		CredentialID: entry.CredentialID,
	})
}

func finishAuthentication(c echo.Context) error {
	// b := echo.Map{}
	b := finishResponse{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	entry := dbItem{}
	err := db.Read(dbColletion, b.User.Name, &entry)
	if err != nil {
		return err
	}
	chal := entry.Challenge

	// Cleanup challenge
	entry.Challenge = []byte{}
	err = db.Write(dbColletion, b.User.Name, entry)
	if err != nil {
		return err
	}

	err = webauthn.ValidateAuthentication(b.PublicKeyCredential, chal, *origin, string(entry.User.ID))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func listUsers(c echo.Context) error {
	records, err := db.ReadAll(dbColletion)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(records) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "No users found")
	}
	items := []dbItem{}
	for _, r := range records {
		item := dbItem{}
		if err := json.Unmarshal([]byte(r), &item); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		items = append(items, item)
	}
	return c.JSON(http.StatusOK, items)
}

func deleteUser(c echo.Context) error {
	username := c.Param("username")
	log.Println("Deleting user:", username)
	if err := db.Delete(dbColletion, username); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
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
