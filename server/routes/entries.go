package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/sanchitdeora/calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var entryCollection *mongo.Collection = OpenCollection(Client, "calories")

func AddEntry (c *gin.Context) {
	// c.Header("Access-Control-Allow-Origin", "http://localhost.com")
    // c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")

	fmt.Println("here:: gin context-", c)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry

	if err := c.BindJSON(&entry); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
        return
    }

	fmt.Println("entry file", entry)
	validationErr := validate.Struct(entry)
    if validationErr!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
        fmt.Println(validationErr)
        return
    }

	entry.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, entry)
	if insertErr != nil {
		msg := fmt.Sprintf("New entry was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
        fmt.Println(insertErr)
        return
	}
	
	defer cancel()
	c.JSON(http.StatusCreated, result)
}

func GetEntries (c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{})

	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
		return
    }

	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }

	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntriesByIngredient (c *gin.Context) {
	ingredient := c.Params.ByName("ingredient")
	fmt.Println(ingredient)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var entries []bson.M
	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient})

	if err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
		return
    }

	if err = cursor.All(ctx, &entries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }

	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntryById (c *gin.Context) {
	entryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)
	
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	
	var entry bson.M
	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
		return
    }
	
	defer cancel()
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

func UpdateEntry (c *gin.Context) {
	entryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)
	
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry models.Entry

	if err := c.BindJSON(&entry); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
        return
    }

	validationErr := validate.Struct(entry)
    if validationErr!= nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
        fmt.Println("validation error", validationErr)
        return
    }

	entry.ID = primitive.NewObjectID()
	result, err := entryCollection.ReplaceOne(ctx, bson.M{"_id": docID}, 
		bson.M{
			"dish": entry.Dish,
			"fat": entry.Fat,
			"ingredients": entry.Ingredients,
			"calories": entry.Calories,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}
	
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func UpdateIngredient (c *gin.Context) {
	entryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}
	var ingredient Ingredient

	if err := c.BindJSON(&ingredient); err!= nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
        return
    }

	fmt.Println("Ingredients: ", ingredient.Ingredients)

	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID}, 
		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func DeleteEntry (c *gin.Context) {
	entryId := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	    return
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}