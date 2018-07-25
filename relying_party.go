package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/mvndaai/webauthn"
	"github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver

const dbColletion = "users"

type (
	user struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	}

	dbItem struct {
		User         user   `json:"user"`
		Challenge    []byte `json:"challenge"`
		CredentialID string `json:"credentialId"`
		PublicKey    string `json:"publicKey"`
	}
)

func main() {
	initDatabase()

	e := echo.New()
	e.HideBanner = true
	e.GET("/", indexHandle)

	e.POST("/registration/start", startRegistration)
	e.POST("/registration/finish", finishRegistration)
	e.POST("/authentication/start", startAuthentication)
	// e.POST("/authentication/finish", finishAuthentication)

	e.GET("/users", listUsers)
	// e.DELETE("/user/:username", deleteUser)

	fmt.Println("Starting server on port :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func indexHandle(c echo.Context) error {
	return c.File("index.html")
}

type (
	startRegistrationResponse struct {
		Challenge string `json:"challenge"`
		User      user   `json:"user"`
	}
)

func startRegistration(c echo.Context) error {
	log.Println("Starting a registration")

	u := user{}
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		return err
	}
	if u.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username required")
	}
	// u.ID = uuid.New().String()
	u.ID = base64Encode([]byte(uuid.New().String()))

	chal, err := webauthn.NewChallenge()
	if err != nil {
		return err
	}

	log.Printf("Saving registration data %#v challenge:%s", u, base64Encode(chal))
	err = db.Write(dbColletion, u.Name, dbItem{User: u, Challenge: chal})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, startRegistrationResponse{
		User:      u,
		Challenge: base64Encode(chal),
	})
}

type (
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
)

func finishRegistration(c echo.Context) error {
	// b := finishRegistrationBody{}
	b := echo.Map{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	log.Printf("body\n%#v\n", b)

	entry := dbItem{}
	err := db.Read(dbColletion, "mvndaai", &entry)
	if err != nil {
		return err
	}
	entry.Challenge = []byte{}
	entry.CredentialID = "rawId"
	entry.PublicKey = "pubkey"
	log.Printf("entry %#v", entry)

	err = db.Write(dbColletion, "mvndaai", entry)
	if err != nil {
		return err
	}

	// log.Println("finish id:", b.ID)
	// log.Println("finish rawID:", b.RawID)
	// log.Println("finish type:", b.Type)

	// log.Printf("finish ClientDataJSON decode: %#v\n", decodeClientData(b.Response.ClientDataJSON))
	// log.Printf("finish AttestationObject: %#v\n", decodeAttestationObject(b.Response.AttestationObject))

	// decodeAttestationObject(b.Response.AttestationObject)

	// https://w3c.github.io/webauthn/#registering-a-new-credential

	return c.NoContent(http.StatusCreated)
}

type (
	startAuthBody struct {
		Username string `json:"username"`
	}

	startAuthResponse struct {
		Challenge    string `json:"challenge"`
		CredentialID string `json:"credentialId"`
	}
)

func startAuthentication(c echo.Context) error {
	b := startAuthBody{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	entry := dbItem{}
	err := db.Read(dbColletion, b.Username, &entry)
	if err != nil {
		return err
	}

	chal, err := webauthn.NewChallenge()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, startAuthResponse{
		Challenge:    base64Encode(chal),
		CredentialID: entry.CredentialID,
	})
}

func listUsers(c echo.Context) error {
	records, err := db.ReadAll(dbColletion)
	if err != nil {
		panic(err)
	}

	items := []dbItem{}
	for _, r := range records {
		item := dbItem{}
		if err := json.Unmarshal([]byte(r), &item); err != nil {
			panic(err)
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, items)
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
