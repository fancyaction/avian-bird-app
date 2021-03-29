package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type birdRecord struct {
	ID          int
	Name        string
	SpeciesInfo interface{}
}

// GetBird receives url to send to Flask API
func GetBird(c *gin.Context) {
	url := c.Query("image_url")

	prediction, err := getPrediction(url)
	if err != nil {
		handleErr(c, err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"success": false, "errors": err.Error()})
		return
	}

	speciesInfo, err := getBirdDetails(prediction.Name)

	if err != nil {
		handleErr(c, err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"success": false, "errors": err.Error()})
		return
	}

	data := gin.H{
		"success": true,
		"msg":     "prediction found",
		"data":    birdRecord{prediction.ID, prediction.Name, speciesInfo},
	}

	c.JSON(http.StatusOK, data)
}

// GetLocation receives lat,lng to send to eBird API
func GetLocation(c *gin.Context) {
	var data []Location
	lat := c.Query("lat")
	lng := c.Query("lng")
	token := os.Getenv("EBIRD_API_KEY")
	targetURL := fmt.Sprintf("https://api.ebird.org/v2/data/obs/geo/recent?lat=%s&lng=%s", lat, lng)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", targetURL, nil)
	req.Header.Set("x-ebirdapitoken", token)

	response, err := client.Do(req)
	if err != nil {
		handleErr(c, err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "errors": "Could not retrieve location details"})
	}

	_ = json.NewDecoder(response.Body).Decode(&data)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "location found",
		"data":    data,
	})
}
