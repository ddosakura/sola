package auth

import (
	"fmt"
	"net/http"

	"github.com/ddosakura/sola"
	"github.com/ddosakura/sola/middleware"
	"github.com/dgrijalva/jwt-go"
)

// New Auth Middleware
func New(signOrAuth middleware.Middleware, pre middleware.Middleware, success middleware.Middleware) middleware.Middleware {
	if pre == nil {
		pre = nextFn
	}
	if success == nil {
		success = func(c middleware.Context, next middleware.Next) {
			c[sola.Response].(http.ResponseWriter).Write([]byte("auth success"))
		}
	}
	return middleware.Merge(func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		cache, err := r.Cookie(authCookieCacheKey)
		tmp := r.Header.Get("Authorization")
		if tmp == "" && err == nil {
			r.Header.Set("Authorization", cache.Value)
		}
		next()
	}, pre, signOrAuth, success)
}

// Sign Token
func Sign(t Type, key interface{}) middleware.Middleware {
	switch t {
	case AuthJWT:
		return signJWT(key)
	}
	return signBase
}

func signBase(c middleware.Context, next middleware.Next) {
	// 可取代浏览器默认弹窗的方式
	w := c[sola.Response].(http.ResponseWriter)
	username, ok1 := c[CtxUsername].(string)
	password, ok2 := c[CtxPassword].(string)
	if !ok1 || !ok2 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  authCookieCacheKey,
		Value: "Basic " + basicAuth(username, password),
	})
	next()
}

func signJWT(key interface{}) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		w := c[sola.Response].(http.ResponseWriter)
		tmp := c[CtxClaims]
		if tmp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		claims := tmp.(map[string]interface{})
		var tmp2 jwt.MapClaims = claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, tmp2)
		t, err := token.SignedString(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  authCookieCacheKey,
			Value: jwtAuthPrefix + t,
		})
		next()
	}
}

// Auth Token
func Auth(t Type, key interface{}) middleware.Middleware {
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

func authBase(check BaseCheck) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)
		username, password, ok := r.BasicAuth()
		if !ok {
			// TODO: custom
			w.Header().Add("WWW-Authenticate", "Basic realm=\"sola\"")
			sola.Text(c, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if check(username, password) {
			next()
		} else {
			// TODO: custom
			sola.Text(c, "Forbidden", http.StatusForbidden)
		}
	}
}

func authJWT(key interface{}) middleware.Middleware {
	return func(c middleware.Context, next middleware.Next) {
		r := c[sola.Request].(*http.Request)
		w := c[sola.Response].(http.ResponseWriter)

		auth := r.Header.Get("Authorization")
		tokenString, ok := parseBearerAuth(auth)
		if !ok {
			// TODO: custom
			w.Header().Add("WWW-Authenticate", jwtAuthPrefix)
			sola.Text(c, "Unauthorized", http.StatusUnauthorized)
			return
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
			sola.Text(c, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var tmp map[string]interface{} = claims
			c[CtxClaims] = tmp
			next()
		} else {
			// TODO: custom
			sola.Text(c, "Forbidden", http.StatusForbidden)
		}
	}
}
