package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/services"
	"github.com/charis16/luminor-golang-be/src/utils"
	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid login payload")
		return
	}

	// Autentikasi lewat service
	user, err := services.AuthenticateAdminUser(req.Email, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	accessToken, refreshToken, err := services.Login(user.UUID, user.Role)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	AdminSetTokenCookies(c, accessToken, refreshToken, user)

	utils.RespondSuccess(c, gin.H{
		"admin_access_token":  accessToken,
		"admin_refresh_token": refreshToken,
	})
}

func AdminRefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("admin_refresh_token")

	if err != nil || refreshToken == "" {
		utils.RespondError(c, http.StatusUnauthorized, "Missing refresh token in cookie")
		return
	}

	_, claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Missing refresh token in cookie")
		return
	}

	newAccessToken, err := services.RefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	user, err := services.GetUserByUUID(claims.UserID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to get user")
		return
	}

	AdminSetTokenCookies(c, newAccessToken, refreshToken, &user)

	utils.RespondSuccess(c, gin.H{
		"admin_access_token":  newAccessToken,
		"admin_refresh_token": refreshToken,
	})
}

func AdminLogout(c *gin.Context) {
	c.SetCookie("admin_access_token", "", -1, "/", "", false, true)
	c.SetCookie("admin_refresh_token", "", -1, "/", "", false, true)
	c.SetCookie("admin_user", "", -1, "/", "", false, true)

	utils.RespondSuccess(c, gin.H{
		"message": "Logged out successfully",
	})

}

func AdminVerifyToken(c *gin.Context) {
	token, err := c.Cookie("admin_access_token")
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

func AdminResetPassword(c *gin.Context) {
	// TODO: implementasi set password baru (dengan token reset)
	utils.RespondError(c, http.StatusNotImplemented, "Reset password not implemented")
}

func AdminSetTokenCookies(c *gin.Context, accessToken string, refreshToken string, user *models.User) {
	accessTokenAge := utils.GetEnvAsDurationInSeconds("JWT_EXPIRATION", "15m")
	refreshTokenAge := utils.GetEnvAsDurationInSeconds("JWT_REFRESH_EXPIRATION", "7d")

	secure := utils.IsProduction() // true only in production
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "admin_access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   accessTokenAge,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "admin_refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   refreshTokenAge,
	})

	userJSON, err := json.Marshal(user)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to serialize user data")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "admin_user",
		Value:    url.QueryEscape(string(userJSON)),
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
		SameSite: sameSite,
		MaxAge:   refreshTokenAge,
	})
}
