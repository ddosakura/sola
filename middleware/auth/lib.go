package auth

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/dgrijalva/jwt-go"
)

func defaultSuccessFunc(c sola.Context) error {
	return c.String(http.StatusOK, "auth success")
}

func defaultSuccess(next sola.Handler) sola.Handler {
	return defaultSuccessFunc
}

func loadAuthCache(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		r := c.Request()
		cache, err := r.Cookie(authCookieCacheKey)
		tmp := r.Header.Get("Authorization")
		if tmp == "" && err == nil {
			r.Header.Set("Authorization", cache.Value)
		}
		return next(c)
	}
}

// New Auth Middleware
func New(signAuth, pre, success sola.Middleware) sola.Middleware {
	if pre == nil {
		pre = nextFn
	}
	if success == nil {
		success = defaultSuccess
	}
	return sola.Merge(loadAuthCache, pre, signAuth, success)
}

// NewFunc Auth Handler
func NewFunc(signAuth, pre sola.Middleware, success sola.Handler) sola.Handler {
	if pre == nil {
		pre = nextFn
	}
	if success == nil {
		success = defaultSuccessFunc
	}
	return sola.MergeFunc(success, loadAuthCache, pre, signAuth)
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
		username, ok1 := c.Get(CtxUsername).(string)
		password, ok2 := c.Get(CtxPassword).(string)
		if !ok1 || !ok2 {
			return c.Handle(http.StatusBadRequest)(c)
		}
		c.SetCookie(&http.Cookie{
			Path:  "/",
			Name:  authCookieCacheKey,
			Value: "Basic " + basicAuth(username, password),
		})
		return next(c)
	}
}

func signJWT(key interface{}) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			tmp := c.Get(CtxClaims)
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
				Path:  "/",
				Name:  authCookieCacheKey,
				Value: jwtAuthPrefix + t,
			})
			c.Set(CtxToken, t)
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
				w.Header().Add("WWW-Authenticate", "Basic realm=\"sola\"")
				return c.Handle(http.StatusUnauthorized)(c)
			}
			if check(username, password) {
				return next(c)
			}
			return c.Handle(http.StatusForbidden)(c)
		}
	}
}

const unexpectedSigningMethod = "Unexpected signing method: %v"

func jwtParse(key interface{}, tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(unexpectedSigningMethod, token.Header["alg"])
		}
		return key, nil
	})
}

func authJWT(key interface{}) sola.Middleware {
	return func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			r := c.Request()
			w := c.Response()

			auth := r.Header.Get("Authorization")
			tokenString, ok := parseBearerAuth(auth)
			if !ok {
				w.Header().Add("WWW-Authenticate", jwtAuthPrefix)
				return c.Handle(http.StatusUnauthorized)(c)
			}

			token, _ := jwtParse(key, tokenString)

			if token == nil {
				w.Header().Add("WWW-Authenticate", jwtAuthPrefix)
				return c.Handle(http.StatusUnauthorized)(c)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var tmp map[string]interface{} = claims
				c.Set(CtxClaims, tmp)
				c.Set(CtxToken, tokenString)
				return next(c)
			}
			return c.Handle(http.StatusForbidden)(c)
		}
	}
}
