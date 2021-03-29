package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//HandleErr //generic error handler, logs error and Os.Exit(1)
func HandleErr(c *gin.Context, err error) error {
	if err != nil {
		c.Error(err)
	}
	return err
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
			MatchAgainst: "allNames",
			Operator:     "equals",
		}},
	}

	postBody, _ := json.Marshal(requestQuery)

	return bytes.NewBuffer(postBody)
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
	if len(data.Results) == 0 {
		err := errors.New("species data not found")
		return SpeciesGlobal{}, err
	}

	speciesInfo := data.Results[0].SpeciesGlobal

	return speciesInfo, nil
}
