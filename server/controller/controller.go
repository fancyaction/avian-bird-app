package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

// GetBird receives url to send to Flask API
func GetBird(c *gin.Context) {
	errors := make(chan error, 2)
	predictChan := make(chan Prediction)
	natureServeChan := make(chan NatureServeAPIResponse)

	url := c.Query("image_url")
	var prediction Prediction
	var natureServeData NatureServeAPIResponse
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		res, err := getPrediction(url)
		predictChan <- res

		if err != nil {
			HandleErr(c, err)
			errors <- err
		}
	}()

	prediction = <-predictChan

	go func() {
		defer wg.Done()

		res, err := getBirdDetails(prediction.Name)
		natureServeChan <- res
		if err != nil {
			HandleErr(c, err)
			errors <- err
		}
	}()

	natureServeData = <-natureServeChan

	go func() {
		wg.Wait()
		close(predictChan)
		close(natureServeChan)
		close(errors)
	}()

	for err := range errors {
		fmt.Println("requestError: ", err)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"success": false, "errors": err.Error()})
		return
	}

	var payload SendBirdPayload
	if len(natureServeData.Results) != 0 {
		payload = SendBirdPayload{
			ID:          prediction.ID,
			Name:        prediction.Name,
			SpeciesInfo: natureServeData.Results[0].SpeciesGlobal,
		}

	} else {
		payload = SendBirdPayload{
			ID:   prediction.ID,
			Name: prediction.Name,
		}
	}

	data := gin.H{
		"success": true,
		"msg":     "prediction found",
		"data":    payload,
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
		HandleErr(c, err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "errors": "Could not retrieve location details"})
	}

	_ = json.NewDecoder(response.Body).Decode(&data)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "location found",
		"data":    data,
	})
}
