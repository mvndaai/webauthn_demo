# Webauthn Demo

Example of WebAuthn for a presentation.

## Server setup

### Downloading from Git

```bash
git clone https://github.com/mvndaai/webauthn_demo.git
pushd webauthn_demo
```

### Run from Dockerfile

```bash
docker build -t webauthndemo . && docker run -p 127.0.0.1:8080:8080 -it webauthndemo -port :8080 -origin https://example.com
```

Note: port and origin should match your configuration

### Run using local Golang

```bash
go get -v ./...
go run . -port :8080 -origin https://example.com
```
Note: port and origin should match your configuration

## Share & Test

And you can use [Telebit](https://telebit.cloud/) to make it avaliable outside [localhost](http://localhost:8080/):

```bash
telebit http 8080 webauthn

> Forwarding https://webauthn.YOURDOMAIN.telebit.io -> localhost:8080
```

## Enabling WebAuthn in Chrome

Chrome has flags that my need to be enabled. Paste this into the omnibar:

chrome://flags/#enable-web-authentication-api

MacOS Touch ID:
chrome://flags/#enable-web-authentication-touch-id

## Spec Variables

https://w3c.github.io/webauthn/#idl-index

### Other Demos

* Duo Labs  [golang code](https://github.com/duo-labs/webauthn)
* Google https://webauthndemo.appspot.com/ ([java code](https://github.com/google/webauthndemo))
* WebAuthn org https://webauthn.org
* Yubico https://demo.yubico.com/webauthn/
