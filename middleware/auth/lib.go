package auth

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/dgrijalva/jwt-go"
)

// New Auth Middleware
func New(signOrAuth sola.Middleware, pre sola.Middleware, success sola.Middleware) sola.Middleware {
	if pre == nil {
		pre = nextFn
	}
	if success == nil {
		success = func(next sola.Handler) sola.Handler {
			return func(c sola.Context) error {
				return c.String(http.StatusOK, "auth success")
			}
		}
	}
	return sola.Merge(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			cache, err := r.Cookie(authCookieCacheKey)
			tmp := r.Header.Get("Authorization")
			if tmp == "" && err == nil {
				r.Header.Set("Authorization", cache.Value)
			}
			return next(c)
		}
	}, pre, signOrAuth, success)
}

// NewFunc Auth Handler
func NewFunc(signOrAuth sola.Middleware, pre sola.Middleware, success sola.Handler) sola.Handler {
	if pre == nil {
		pre = nextFn
	}
	if success == nil {
		success = func(c sola.Context) error {
			return c.String(http.StatusOK, "auth success")
		}
	}
	return sola.MergeFunc(success, func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			cache, err := r.Cookie(authCookieCacheKey)
			tmp := r.Header.Get("Authorization")
			if tmp == "" && err == nil {
				r.Header.Set("Authorization", cache.Value)
			}
			return next(c)
		}
	}, pre, signOrAuth)
}

// Sign Token
func Sign(t Type, key interface{}) sola.Middleware {
	switch t {
	case AuthJWT:
		return signJWT(key)
	}
	return signBase
}

func signBase(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		// 可取代浏览器默认弹窗的方式
		username, ok1 := c[CtxUsername].(string)
		password, ok2 := c[CtxPassword].(string)
		if !ok1 || !ok2 {
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		c.SetCookie(&http.Cookie{
			Name:  authCookieCacheKey,
			Value: "Basic " + basicAuth(username, password),
		})
		return next(c)
	}
}

func signJWT(key interface{}) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			tmp := c[CtxClaims]
			if tmp == nil {
				return ErrNoClaims
			}
			claims := tmp.(map[string]interface{})
			var tmp2 jwt.MapClaims = claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tmp2)
			t, err := token.SignedString(key)
			if err != nil {
				return err
			}
			c.SetCookie(&http.Cookie{
				Name:  authCookieCacheKey,
				Value: jwtAuthPrefix + t,
			})
			return next(c)
		}
	}
}

// Auth Token
func Auth(t Type, key interface{}) sola.Middleware {
	switch t {
	case AuthJWT:
		return authJWT(key)
	}
	fn, ok := key.(BaseCheck)
	if !ok {
		fn = basePass
	}
	return authBase(fn)
}

func authBase(check BaseCheck) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()
			username, password, ok := r.BasicAuth()
			if !ok {
				// TODO: custom
				w.Header().Add("WWW-Authenticate", "Basic realm=\"sola\"")
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}
			if check(username, password) {
				return next(c)
			}
			// TODO: custom
			return c.String(http.StatusForbidden, "Forbidden")
		}
	}
}

func authJWT(key interface{}) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()

			auth := r.Header.Get("Authorization")
			tokenString, ok := parseBearerAuth(auth)
			if !ok {
				// TODO: custom
				w.Header().Add("WWW-Authenticate", jwtAuthPrefix)
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}

			token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return key, nil
			})

			if token == nil {
				// TODO: custom
				w.Header().Add("WWW-Authenticate", jwtAuthPrefix)
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var tmp map[string]interface{} = claims
				c[CtxClaims] = tmp
				return next(c)
			}
			// TODO: custom
			return c.String(http.StatusForbidden, "Forbidden")
		}
	}
}
