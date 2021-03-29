package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"success": false, "errors": err.Error()})
}

//send image and get prediction from Pytorch API server
func getPrediction(fileURL string) (Prediction, error) {
	var data Prediction

	baseURL := os.Getenv("FLASK_API_BASE_URL")
	targetURL := baseURL + "/predict"

	postBody, _ := json.Marshal(map[string]string{
		"value": fileURL,
	})

	response, err := http.Post(targetURL, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return data, err
	} else if response.StatusCode != 200 {
		err := errors.New("Prediction not found")
		return data, err
	}

	defer response.Body.Close()

	_ = json.NewDecoder(response.Body).Decode(&data)

	return data, nil
}

func getPostBody(name string) *bytes.Buffer {
	requestQuery := NatureServeParams{
		CriteriaType: "species",
		TextCriteria: []TextCriteria{{
			ParamType:    "textSearch",
			SearchToken:  name,
			MatchAgainst: "primaryCommonName",
			Operator:     "similarTo",
		}},
	}

	postBody, _ := json.Marshal(requestQuery)

	return bytes.NewBuffer(postBody)
}

func getSpeciesInfo(results []Result, matchAgainst string) (SpeciesGlobal, error) {
	if len(results) == 0 {
		return SpeciesGlobal{}, errors.New("species data not found")
	}

	for _, record := range results {
		if strings.Contains(record.PrimaryCommonName, matchAgainst) {
			return record.SpeciesGlobal, nil
		}
	}

	return SpeciesGlobal{}, errors.New("species data not found")
}

// send bird name and get bird details from NatureServe api
func getBirdDetails(name string) (SpeciesGlobal, error) {
	var data NatureServeAPIResponse
	targetURL := "https://explorer.natureserve.org/api/data/speciesSearch"

	postBody := getPostBody(name)
	response, err := http.Post(targetURL, "application/json", postBody)

	if err != nil {
		return SpeciesGlobal{}, err
	} else if response.StatusCode != 200 {
		err := errors.New("species data not found")
		return SpeciesGlobal{}, err
	}

	defer response.Body.Close()

	_ = json.NewDecoder(response.Body).Decode(&data)

	speciesInfo, err := getSpeciesInfo(data.Results, name)
	if err != nil {
		return SpeciesGlobal{}, err
	}

	return speciesInfo, nil
}
