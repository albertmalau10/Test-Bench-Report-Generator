package controllers

import(
	"net/http"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"valve_database/models"
	"gorm.io/gorm"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ping",
	})
}

func CreateValve(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newValve models.Valve

		if err := c.ShouldBindJSON(&newValve); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result := db.Create(&newValve); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, newValve)
	}
}

func GetValves(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var valves []models.Valve

		if result := db.Find(&valves); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, valves)
	}
}

func GetValveByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var v models.Valve

		result := db.First(&v, id)
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "valve tidak ditemukan"})
			return
		} else if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, v)
	}
}

func UpdateValve(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var valve models.Valve
		if result := db.First(&valve, id); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "valve tidak ditemukan"})
			return
		}

		var updatedValve models.Valve
		if err := c.ShouldBindJSON(&updatedValve); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		valve.Manufacturer = updatedValve.Manufacturer
		valve.PartNumber = updatedValve.PartNumber
		valve.ComponentSeries = updatedValve.ComponentSeries
		valve.ValveType = updatedValve.ValveType
		valve.Size_NG = updatedValve.Size_NG
		valve.Weight = updatedValve.Weight
		valve.RatedFlow = updatedValve.RatedFlow
		valve.MaxFlow = updatedValve.MaxFlow
		valve.CommandValue = updatedValve.CommandValue
		valve.CommandType = updatedValve.CommandType
		valve.MaxPressure = updatedValve.MaxPressure
		valve.ImagePath = updatedValve.ImagePath
		valve.DatasheetPath = updatedValve.DatasheetPath

		db.Save(&valve)

		c.JSON(http.StatusOK, valve)
	}
}

func DeleteValve(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var valve models.Valve
		if result := db.First(&valve, id); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "valve tidak ditemukan"})
			return
		}

		db.Delete(&valve)

		c.JSON(http.StatusOK, gin.H{"message": "valve berhasil dihapus"})
	}
}

func UploadValveImage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var valve models.Valve
		if result := db.First(&valve, id); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "valve tidak ditemukan"})
			return
		}

		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file gambar tidak ditemukan"})
			return
		}

		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%s_%d%s", id, time.Now().Unix(), ext)
		savePath := filepath.Join("images", newFileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
			return
		}

		valve.ImagePath = "/images/" + newFileName
		db.Save(&valve)

		c.JSON(http.StatusOK, valve)
	}
}

func UploadValveDatasheet(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var valve models.Valve
		if result := db.First(&valve, id); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "valve tidak ditemukan"})
			return
		}

		file, err := c.FormFile("datasheet")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file datasheet tidak ditemukan"})
			return
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file harus berformat PDF"})
			return
		}

		newFileName := fmt.Sprintf("%s_%d%s", id, time.Now().Unix(), ext)
		savePath := filepath.Join("datasheets", newFileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
			return
		}

		valve.DatasheetPath = "/datasheets/" + newFileName
		db.Save(&valve)

		c.JSON(http.StatusOK, valve)
	}
}