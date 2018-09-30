# Webauthn Demo

Example of WebAuthN for a presentation.

## Server setup

To run the server just run this docker command

```bash
docker build -t webauthndemo . && docker run -p 127.0.0.1:8080:8080 -it webauthndemo
```

Use [telebit](https://telebit.cloud/) or howereve else you want to make the site avaliable on a mobile device

## Enabling WebAuthN in Chrome

Chrome has flags that my need to be enabled. Paste this into the omnibar:

chrome://flags/#enable-web-authentication-api

MacOS Touch ID:
chrome://flags/#enable-web-authentication-touch-id

## Spec Variables

https://w3c.github.io/webauthn/#idl-index


## Notes

### Registration

https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API#Registration

Google IO demo https://youtu.be/kGGMgEfSzMw?t=27m22s

1. Server creates a challenge associated with the user stores it and sends it to the webapp.
2. Browser sends User info, callenge and authoritative domain name to an authenticatore
    * The domain name prevents phishing sites
3. The authenticator stores the credential id, a private key, the domanin name, and the user info
4. The authenticate sends the credential id, the public key and a signature to the web app which forwards them to the server
5. The server validates the signature, stores the credential id and the public key, then invalidates the challenge

### Authentication

https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API#Authentication

Google IO demo https://youtu.be/kGGMgEfSzMw?t=30m4s

1. Server creates a challenge and sends it with the credential id to the web app
2. The browser extracts the domain name and sends it with the challend and the crendential id to the authenticator.
3. The authenticator checks the domain name then creates a signature and sends it back through the web app to the server
4. The server uses the signature to check if the public key and challenge match, and if so invalidates the challenge and considers the user logged in.