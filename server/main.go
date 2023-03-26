package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sanchitdeora/calorie-tracker/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
    router.Use(CORS) 

	api := router.Group("/api")
	api.POST("/entry/create", routes.AddEntry)
	api.GET("/entries", routes.GetEntries)
	api.GET("/entry/:id", routes.GetEntryById)
	api.GET("/ingredient/:ingredient", routes.GetEntriesByIngredient)

	api.PUT("/entry/update/:id", routes.UpdateEntry)
	api.PUT("/ingredient/update/:id", routes.UpdateIngredient)
	api.DELETE("/entry/delete/:id", routes.DeleteEntry)

	router.Run(":" + port)

}

func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
        
		// Everytime we receive an OPTIONS request, 
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real 
		// request using any other method than OPTIONS
		c.AbortWithStatus(http.StatusOK)
	}
}