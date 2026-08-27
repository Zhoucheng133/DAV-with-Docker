package routes

import (
	"database/sql"

	"dav_docker/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleInit(c *fiber.Ctx) error {
	var count int
	err := utils.DB.QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}
	needsInit := count == 0
	return JSONResponse(c, true, needsInit)
}

func HandleRegister(c *fiber.Ctx) error {
	var req UserRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONResponse(c, false, "invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return JSONResponse(c, false, "username and password are required")
	}

	var count int
	if err := utils.DB.QueryRow("SELECT COUNT(*) FROM user").Scan(&count); err != nil {
		return JSONResponse(c, false, err.Error())
	}
	if count > 0 {
		return JSONResponse(c, false, "user already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return JSONResponse(c, false, "failed to hash password")
	}

	userID := utils.GenerateID()
	_, err = utils.DB.Exec("INSERT INTO user (id, username, password) VALUES (?, ?, ?)", userID, req.Username, string(hashedPassword))
	if err != nil {
		return JSONResponse(c, false, err.Error())
	}

	return JSONResponse(c, true, "")
}

func HandleLogin(c *fiber.Ctx) error {
	var req UserRequest
	if err := c.BodyParser(&req); err != nil {
		return JSONResponse(c, false, "invalid request body")
	}

	var id, username, hashedPassword string
	err := utils.DB.QueryRow("SELECT id, username, password FROM user LIMIT 1").Scan(&id, &username, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return JSONResponse(c, false, "user not found")
		}
		return JSONResponse(c, false, err.Error())
	}

	if req.Username != username {
		return JSONResponse(c, false, "invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return JSONResponse(c, false, "invalid username or password")
	}

	accessToken, err := utils.GenerateAccessToken(username)
	if err != nil {
		return JSONResponse(c, false, "failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefreshToken(username)
	if err != nil {
		return JSONResponse(c, false, "failed to generate refresh token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "dav_docker_refresh_token",
		Value:    refreshToken,
		Path:     "/api/refresh",
		MaxAge:   30 * 24 * 60 * 60,
		HTTPOnly: true,
		Secure:   false,
	})

	return JSONResponse(c, true, accessToken)
}

func HandleRefresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("dav_docker_refresh_token")
	if refreshToken == "" {
		return JSONResponse(c, false, "missing refresh token")
	}

	username, err, isExpired := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		if isExpired {
			return JSONResponse(c, false, "expired")
		}
		return JSONResponse(c, false, err.Error())
	}

	accessToken, err := utils.GenerateAccessToken(username)
	if err != nil {
		return JSONResponse(c, false, "failed to generate access token")
	}

	return JSONResponse(c, true, accessToken)
}

func HandleAuth(c *fiber.Ctx) error {
	return JSONResponse(c, true, "")
}

func HandleLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "dav_docker_refresh_token",
		Value:    "",
		Path:     "/api/refresh",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   false,
	})
	return JSONResponse(c, true, "")
}
