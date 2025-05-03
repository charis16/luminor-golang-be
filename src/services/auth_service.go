package services

import (
	"errors"

	"github.com/charis16/luminor-golang-be/src/config"
	"github.com/charis16/luminor-golang-be/src/models"
	"github.com/charis16/luminor-golang-be/src/utils"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Cocokkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func Login(userID, role string) (string, string, error) {
	accessToken, err := utils.GenerateAccessToken(userID, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(userID, role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshToken(refreshToken string) (string, error) {
	_, claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Ambil ulang access token berdasarkan claim
	newAccessToken, err := utils.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return "", errors.New("failed to generate new access token")
	}

	return newAccessToken, nil
}

func VerifyAccessToken(token string) (*utils.CustomClaims, error) {
	_, claims, err := utils.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
