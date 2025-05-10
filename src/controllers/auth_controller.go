package controllers

import (
	"fmt"
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
		utils.RespondError(c, http.StatusBadRequest, "Invalid login payload")
		return
	}

	// Autentikasi lewat service
	user, err := services.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, refreshToken, err := services.Login(user.UUID, user.Role)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	setTokenCookies(c, accessToken, refreshToken)

	utils.RespondSuccess(c, gin.H{
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
	utils.RespondError(c, http.StatusNotImplemented, "Register not implemented")
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	fmt.Println("Refresh token:", refreshToken)
	if err != nil || refreshToken == "" {
		utils.RespondError(c, http.StatusUnauthorized, "Missing refresh token in cookie")
		return
	}

	newAccessToken, err := services.RefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	setTokenCookies(c, newAccessToken, refreshToken)

	utils.RespondSuccess(c, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": refreshToken,
	})

}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	utils.RespondSuccess(c, gin.H{
		"message": "Logged out successfully",
	})

}

func VerifyToken(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Access token not found in cookie")
		return
	}

	claims, err := services.VerifyAccessToken(token)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	utils.RespondSuccess(c, gin.H{
		"valid":   true,
		"user_id": claims.UserID,
		"role":    claims.Role,
		"expires": claims.ExpiresAt.Time,
	})

}

func ForgotPassword(c *gin.Context) {
	// TODO: implementasi kirim email reset password
	utils.RespondError(c, http.StatusNotImplemented, "Forgot password not implemented")
}

func ResetPassword(c *gin.Context) {
	// TODO: implementasi set password baru (dengan token reset)
	utils.RespondError(c, http.StatusNotImplemented, "Reset password not implemented")
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
