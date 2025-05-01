package handlers

import (
	"fmt"
	"net/http"

	"google-trends-api/src/services"

	"github.com/gin-gonic/gin"
)

// var aiclient = &http.Client{}

func Handle_empty(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		c.Writer.Write([]byte("Get Response"))
	case http.MethodPost:
		fmt.Fprint(c.Writer, "POST request")
	default:
		fmt.Fprintf(c.Writer, "Unsupported method: %s", c.Request.Method)
	}
}

func Handle_kigo(c *gin.Context) {
	fmt.Println(c.Request.Method)

	// if GET method, send this
	if c.Request.Method == http.MethodGet {
		c.Writer.Write([]byte("Ei cholche Go :)"))
	}

	// if POST method, do this
	// if c.Request.Method == http.MethodPost {

	// 	// read req body as json
	// 	jsondata, err := io.ReadAll(c.Request.Body)
	// 	if err != nil {
	// 		fmt.Println("Error reading request body: ", err)
	// 	}

	// 	// json to payload data
	// 	data := &services.Payload{}
	// 	err = json.Unmarshal(jsondata, data)
	// 	if err != nil {
	// 		fmt.Println("Error parsing json data: ", err)
	// 	}

	// 	// request to AI server and get response
	// 	resp, err := services.Ask_AI(config.APP_CONFIG.AI_CLIENT_API_ENDPOINT, data, aiclient)
	// 	if err != nil {
	// 		fmt.Println("Error from the AI server: ", err)
	// 	}

	// 	// send response to client
	// 	// fmt.Fprint(c.Writer, resp)
	// 	// println(resp)
	// 	c.JSON(http.StatusOK, resp)

	// }
}

func GetGoogleTrends(c *gin.Context) {
	// Read the 'geo' query parameter
	geo := c.Query("geo") // Returns "" if not present

	// TODO: Use the 'geo' parameter to filter/fetch data from services
	// For now, just printing it
	fmt.Printf("Received geo parameter: %s\n", geo)

	// services.ExportGoogleTrends() // Assuming this might need the geo param later
	// data := services.ExtractGoogleTrends() // Assuming this might need the geo param later
	// c.Writer.Write([]byte("Google Trends exported"))
	data := services.Data // This likely needs to be filtered based on 'geo'
	// data := services.RawData
	// data := services.RawHTML
	// fmt.Println("Google Trends exported", data)

	c.JSON(http.StatusOK, gin.H{
		"data": data, // Return potentially filtered data
	})
}

func GetGoogleTrendsFiltered(c *gin.Context) {
	// Read the request body
	var requestBody map[string]string
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract the 'geo' parameter from the request body
	geo, exists := requestBody["geo"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'geo' parameter"})
		return
	}

	fmt.Printf("Received geo parameter: %s\n", geo)

	data := services.Data[geo]
	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for the specified geo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data, // Return potentially filtered data
	})
}
