package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"valve_database/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func GetSettings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Settings
		db.FirstOrCreate(&s, models.Settings{ID: 1})
		c.JSON(http.StatusOK, s)
	}
}

func UpdateSettings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Settings
		if err := db.FirstOrCreate(&s, models.Settings{ID: 1}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal load settings: " + err.Error()})
			return
		}

		var input struct {
			ModbusServerAddress string `json:"modbus_server_address"`
			CtrlxHost           string `json:"ctrlx_host"`
			CtrlxUsername       string `json:"ctrlx_username"`
			CtrlxPassword       string `json:"ctrlx_password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s.ModbusServerAddress = input.ModbusServerAddress
		s.CtrlxHost = input.CtrlxHost
		s.CtrlxUsername = input.CtrlxUsername
		if input.CtrlxPassword != "" {
			s.CtrlxPassword = input.CtrlxPassword
		}

		if err := db.Save(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan settings: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, s)
	}
}

func GenerateCtrlxToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Settings
		db.First(&s, 1)

		if s.CtrlxHost == "" || s.CtrlxUsername == "" || s.CtrlxPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ctrlX host/username/password belum diisi"})
			return
		}

		payload, _ := json.Marshal(map[string]string{
			"name":     s.CtrlxUsername,
			"password": s.CtrlxPassword,
		})

		client := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		url := fmt.Sprintf("https://%s/identity-manager/api/v1/auth/token", s.CtrlxHost)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "tidak bisa menghubungi ctrlX CORE: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		var result struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.AccessToken == "" {
			c.JSON(resp.StatusCode, gin.H{"error": "gagal generate token, cek username/password ctrlX"})
			return
		}

		expiresAt := time.Now().Add(1 * time.Hour) // fallback default kalau decode gagal
		parser := jwt.NewParser()
		claims := jwt.MapClaims{}
		if _, _, err := parser.ParseUnverified(result.AccessToken, claims); err == nil {
			if expFloat, ok := claims["exp"].(float64); ok {
				expiresAt = time.Unix(int64(expFloat), 0)
			}
		}

		s.CtrlxToken = result.AccessToken
		s.CtrlxTokenExpiresAt = &expiresAt
		db.Save(&s)

		c.JSON(http.StatusOK, gin.H{
			"token":      result.AccessToken,
			"expires_at": expiresAt,
		})
	}
}