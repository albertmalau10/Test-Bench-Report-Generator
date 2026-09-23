package controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"io"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"valve_database/models"
)

func writeBoolToCtrlX(
	host string,
	token string,
	node string,
	value bool,
) error {

	payload := map[string]interface{}{
    "type":  "bool8",
    "value": value,
	}

	body, err := json.Marshal(payload)
	if err != nil {
    	return err
	}


	url := fmt.Sprintf(
		"https://%s/automation/api/v2/nodes/%s",
		host,
		node,
	)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewBuffer(body),
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
    return fmt.Errorf(
        "ctrlX returned status %d: %s",
        resp.StatusCode,
        string(responseBody),
    )
}

	return nil
}

func StartOutput(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var s models.Settings

		db.First(&s, 1)

		err := writeBoolToCtrlX(
			s.CtrlxHost,
			s.CtrlxToken,
			"plc/app/Application/sym/PLC_PRG/output",
			true,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "active",
		})
	}
}
func StopOutput(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var s models.Settings

		db.First(&s, 1)

		err := writeBoolToCtrlX(
			s.CtrlxHost,
			s.CtrlxToken,
			"plc/app/Application/sym/PLC_PRG/output",
			false,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "stopped",
		})
	}
}

