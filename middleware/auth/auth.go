package auth

// Type of Auth
type Type uint8

// AuthType
const (
	AuthBase Type = iota
	AuthJWT
)

// BaseCheck for Base Auth
type BaseCheck func(username, password string) bool

func basePass(string, string) bool {
	return true
}

const jwtAuthPrefix = "Bearer "
const authCookieCacheKey = "Sola-Authorization"
