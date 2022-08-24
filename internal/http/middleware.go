package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"tests2/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
)

var debugUser = models.UserFromRequest{
	ID:       1,
	Username: "admin",
}

func (h *Handler) VerifyUser(ctx *fiber.Ctx) error {
	var identity models.UserFromRequest

	if !h.cfg.DebugAuth {
		// get user from request
		user, err := h.GetUserFromRequest(ctx)
		if err != nil {
			return h.Response(ctx, http.StatusUnauthorized, nil, err)
		}

		// Отправить запрос в гейтвей на проверку токена
		resp, err := h.middleware.GatewayValidateToken(ctx.UserContext(), user.AccessToken)
		if err != nil {
			return h.Response(ctx, http.StatusUnauthorized, nil, err)
		}

		identity = resp.User
	} else {
		identity = debugUser
	}

	// Положить информацию о пользователе в контекст
	ctxBack := ctx.UserContext()
	ctxBack, err := h.utilsUser.PutUserToContext(ctxBack, &identity)
	if err != nil {
		return err
	}
	ctx.SetUserContext(ctxBack)

	return ctx.Next()
}

// GetUserFromRequest возвращает данные из токена Bearer
func (h *Handler) GetUserFromRequest(ctx *fiber.Ctx) (*models.UserFromRequest, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("empty auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return nil, errors.New("wrong auth header")
	}

	if headerParts[0] != "Bearer" {
		return nil, errors.New("auth method is not Bearer")
	}

	claims, err := ParseToken(headerParts[1], []byte(h.cfg.TokenSignedKey))
	if err != nil {
		return nil, err
	}

	var user = models.UserFromRequest{
		ID:          claims.ID,
		Username:    claims.UserName,
		AccessToken: headerParts[1],
	}
	return &user, nil
}

func ParseToken(accessToken string, signedKey []byte) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signedKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidAccessToken
}
