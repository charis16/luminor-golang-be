package controllers

import (
	"net/http"

	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid login payload"})
		return
	}

	// Autentikasi lewat service
	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := services.Login(user.UUID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.UUID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func Register(c *gin.Context) {
	// TODO: implementasi pendaftaran user baru
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Register not implemented"})
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing refresh token in cookie"})
		return
	}

	newAccessToken, err := services.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	setTokenCookies(c, newAccessToken, refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": refreshToken,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func VerifyToken(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Access token not found in cookie"})
		return
	}

	claims, err := services.VerifyAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"user_id": claims.UserID,
		"role":    claims.Role,
		"expires": claims.ExpiresAt.Time,
	})
}

func ForgotPassword(c *gin.Context) {
	// TODO: implementasi kirim email reset password
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Forgot password not implemented"})
}

func ResetPassword(c *gin.Context) {
	// TODO: implementasi set password baru (dengan token reset)
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Reset password not implemented"})
}

func setTokenCookies(c *gin.Context, accessToken string, refreshToken string) {
	accessTokenAge := utils.GetEnvAsDurationInSeconds("JWT_EXPIRATION", "15m")
	refreshTokenAge := utils.GetEnvAsDurationInSeconds("JWT_REFRESH_EXPIRATION", "7d")

	secure := utils.IsProduction() // true only in production
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   accessTokenAge,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   refreshTokenAge,
	})
}
