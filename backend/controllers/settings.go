package controllers

import (
	"net/http"
	"valve_database/models"

	"github.com/gin-gonic/gin"
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
			OpcUaAddress        string `json:"opc_ua_address"`
			OpcNodePressure     string `json:"opc_node_pressure"`
			OpcNodeCommand      string `json:"opc_node_command"`
			OpcNodeFeedback     string `json:"opc_node_feedback"`
			OpcNodeFlow         string `json:"opc_node_flow"`
			OpcNodeOutput       string `json:"opc_node_output"`
			OpcUaUsername       string `json:"opc_ua_username"`
			OpcUaPassword       string `json:"opc_ua_password"`
			OpcUaSecurityPolicy string `json:"opc_ua_security_policy"`
			OpcUaSecurityMode   string `json:"opc_ua_security_mode"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s.OpcUaAddress = input.OpcUaAddress
		s.OpcNodePressure = input.OpcNodePressure
		s.OpcNodeCommand = input.OpcNodeCommand
		s.OpcNodeFeedback = input.OpcNodeFeedback
		s.OpcNodeFlow = input.OpcNodeFlow
		s.OpcNodeOutput = input.OpcNodeOutput
		s.OpcUaUsername = input.OpcUaUsername
		s.OpcUaSecurityPolicy = input.OpcUaSecurityPolicy
		s.OpcUaSecurityMode = input.OpcUaSecurityMode

		if input.OpcUaPassword != "" {
			s.OpcUaPassword = input.OpcUaPassword
		}

		if err := db.Save(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal simpan settings: " + err.Error()})
			return
		}
		
		ResetOpcConnection()
		c.JSON(http.StatusOK, s)
	}
}