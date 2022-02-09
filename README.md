# zt53
zt53 is a simple app to push ZeroTier hosts to AWS route53.

Setup a .env file with the following environmental variables.

```
ZT_NETWORK=
ZT_TOKEN=
R53_ZONE=
DOMAIN=
```

`$ go run zt53.go`