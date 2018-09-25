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

	//TOOD temp
	// b := webauthn.ParsedRegistrationResponse{
	// 	Type:              "public-key",
	// 	CredentialID:      "AJNp99XP4fp1UXFPgGsVNP9xjHf2XkGS73PZ+gdDms9leSLXnROcx7YXEh/BqY7bW6/k4bFvCJEveJ21tdKk16wnOKnUn+khWLZt8xuJS2U=",
	// 	ClientDataJSON:    "eyJjaGFsbGVuZ2UiOiJOalF6WldZM01EZ3RPR0ZqT1MwMFl6WTVMV0UzTnpRdE9XTTNPR1F5TTJGbE1qYzUiLCJvcmlnaW4iOiJodHRwOi8vbG9jYWxob3N0OjgwODAiLCJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIn0=",
	// 	AttestationObject: "o2NmbXRoZmlkby11MmZnYXR0U3RtdKJjc2lnWEcwRQIhAJ9IEGmFulCkcRaamRp9OyGPM84Pye2royyqs62XgZfhAiBswQUYuuu9yUIVoeulbpplBdwa0oii/k4RyuQ1UFWog2N4NWOAaGF1dGhEYXRhWNRJlg3liA6MaHQ0Fw9kdmBbj+SuuaKGMseZXPO6gx2XY0VbXRJ5rc4AAjW8xgpkiwsl8fBVAwBQAJNp99XP4fp1UXFPgGsVNP9xjHf2XkGS73PZ+gdDms9leSLXnROcx7YXEh/BqY7bW6/k4bFvCJEveJ21tdKk16wnOKnUn+khWLZt8xuJS2WlAQIDJiABIVggRJ5Wbf462wADcoZm7N4GnsRZGUfgkqNy3afGujC/mHQiWCB/HbeQep7fe++SJZ/NcRH9k2mu4fGuvx2snhNqoryG5Q==",
	// }
	// chal := []byte{0xa, 0xa, 0xb, 0xd3, 0x61, 0x7a, 0x8a, 0xfa, 0x52, 0x25, 0x64, 0xa9, 0x65, 0x96, 0x18, 0x6, 0x31, 0xcd, 0xca, 0x78, 0xab, 0x41, 0x16, 0x16, 0x77, 0xd4, 0x93, 0xe9, 0x88, 0x54, 0xf9, 0xeb}
	// if _, err := webauthn.IsValidRegistration(b, chal, "http://localhost:8080", false); err != nil {
	// 	log.Println("Error", err)
	// }

	// if true {
	// 	return
	// }
	//TODO temp..

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

	// log.Printf("Saving registration data %#v challenge:%s", u, base64Encode(chal))
	err = db.Write(dbColletion, u.Name, dbItem{User: u, Challenge: chal})
	if err != nil {
		return err
	}

	u.ID = uuid.New().String()
	r := webauthn.RegistrationParts{
		ToArrayBuffter: webauthn.BuildToArrayBuffer(chal, u.ID),
		PublicKey: webauthn.PublicKeyCredentialOptions{
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
	// finishRegistrationBody struct {
	// 	ID       string                        `json:"id"`
	// 	RawID    string                        `json:"rawID"`
	// 	Type     string                        `json:"type"`
	// 	Response finishRegistrationBodyResonse `json:"response"`
	// }
	// finishRegistrationBody struct {
	// 	//Type ?
	// 	CredentialID      string `json:"credentialId"`
	// 	ClientDataJSON    string `json:"clientDataJSON"`
	// 	AttestationObject string `json:"attestationObject"`
	// }

	finishRegistrationBodyResonse struct {
		AttestationObject string `json:"attestationObject"`
		ClientDataJSON    string `json:"clientDataJSON"`
	}
)

type registionResponse struct {
	webauthn.ParsedRegistrationResponse
	User webauthn.UserEntity `json:"user"`
}

func finishRegistration(c echo.Context) error {
	b := registionResponse{}
	// b := echo.Map{}
	if err := c.Bind(&b); err != nil {
		return err
	}
	log.Printf("body\n%#v\n", b)

	entry := dbItem{}
	err := db.Read(dbColletion, b.User.Name, &entry)
	if err != nil {
		return err
	}

	err = webauthn.ValidateRegistration(b.ParsedRegistrationResponse, entry.Challenge, "http://localhost:8080", false)
	if err != nil {
		return err
	}

	entry.Challenge = []byte{}
	entry.CredentialID = string(b.CredentialID)
	entry.PublicKey = "pubkey"
	log.Printf("entry %#v", entry)

	err = db.Write(dbColletion, b.User.Name, entry)
	if err != nil {
		return err
	}
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

type finishAuthBody struct {
	Username string `json:"username"`
}

func finishAuthentication(c echo.Context) error {
	// b := echo.Map{}
	b := finishAuthBody{}
	if err := c.Bind(&b); err != nil {
		return err
	}

	log.Println("finish body", b)

	entry := dbItem{}
	err := db.Read(dbColletion, b.Username, &entry)
	if err != nil {
		return err
	}

	err = webauthn.ValidateAuthentication()
	if err != nil {
		return err
	}

	entry.Challenge = []byte{}
	err = db.Write(dbColletion, b.Username, entry)
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
