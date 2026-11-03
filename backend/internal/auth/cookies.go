package auth

import (
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
)

type humaContext = huma.Context

type setCookieContext struct {
	humaContext
}

func (c setCookieContext) AppendHeader(name string, value string) {
	if strings.EqualFold(name, "Set-Cookie") {
		humafiber.Unwrap(c.humaContext).Response().Header.Add(name, value)
		return
	}
	c.humaContext.AppendHeader(name, value)
}

func (c setCookieContext) Unwrap() huma.Context {
	return c.humaContext
}

func SetCookieMiddleware(ctx huma.Context, next func(huma.Context)) {
	next(setCookieContext{ctx})
}

const (
	AccessTokenCookieName  = "jwt"
	RefreshTokenCookieName = "refresh_token"

	refreshTokenMaxAge = 30 * 24 * 60 * 60
)

func AccessTokenCookie(accessToken string, expiresIn int) http.Cookie {
	return http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		MaxAge:   expiresIn,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func RefreshTokenCookie(refreshToken string) http.Cookie {
	return http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   refreshTokenMaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func ExpiredRefreshTokenCookie() http.Cookie {
	return http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
