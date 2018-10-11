# Webauthn Demo

Example of WebAuthN for a presentation.

## Server setup

### Downloading from Git

```bash
git clone https://github.com/mvndaai/webauthn_demo.git
pushd webauthn_demo
```

### Run from Dockerfile

```bash
docker build -t webauthndemo . && docker run -p 127.0.0.1:8080:8080 -it webauthndemo
```

### Run using local Golang

```bash
go get -v ./...
go run . -p 8080
```

## Share & Test

And you can use [Telebit](https://telebit.cloud/) to make it avaliable outside [localhost](http://localhost:8080/):

```bash
telebit http 8080 webauthn

> Forwarding https://webauthn.YOURDOMAIN.telebit.io -> localhost:8080
```

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

### Authentication

https://developer.mozilla.org/en-US/docs/Web/API/Web_Authentication_API#Authentication

Google IO demo https://youtu.be/kGGMgEfSzMw?t=30m4s

### Other Demos

* Duo Labs has a working Golang demo if this doesn't suffice ([code](https://github.com/duo-labs/webauthn))
* Google [Java demo](https://webauthndemo.appspot.com/) ([code](https://github.com/google/webauthndemo))
* https://webauthn.org/
* https://demo.yubico.com/webauthn/
