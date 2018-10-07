# Webauthn Demo

Example of WebAuthN for a presentation.

## Server setup

To run the server just run this docker command

```bash
docker build -t webauthndemo . && docker run -p 127.0.0.1:8080:8080 -it webauthndemo
```

[Telebit](https://telebit.cloud/) can be used to make this avaliable outside [localhost](http://localhost:8080/)

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