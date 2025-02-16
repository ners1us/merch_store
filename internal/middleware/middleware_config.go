package middleware

var jwtSecret []byte

func Init(secret []byte) {
	jwtSecret = secret
}
