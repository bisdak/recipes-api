package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bisdak/recipes-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// @Summary get all items in recipe list
// @id recipes
// @Produce json
// @Success 200 {object} Recipe
// @Router /recipes [get]
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(handler.ctx)

	recipes := make([]models.Recipe, 0)
	for cur.Next(handler.ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

// @Summary add new recipe to the list
// @id create-recipe
// @Produce json
// @accept json
// @Param data body Recipe true "recipe data"
// @Success 200 {object} Recipe
// @Failure 400
// @Router /recipes [post]
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe

	// Check if json body is valid
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new reciper"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// @Summary delete a recipe item by ID
// @ID delete-recipe-by-id
// @Produce json
// @Param id path string true "recipe ID"
// @Success 200 {object} Recipe
// @Failure 404
// @Router /recipes/{id} [delete]
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
}

// @Summary get recipe by ID
// @ID get-recipe-by-id
// @Produce json
// @Param id path string true "recipe ID"
// @Success 200 {object} Recipe
// @Failure 404
// @Router /recipes/{id} [get]
func (handler *RecipesHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": objectId,
	})
	var recipe models.Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// @Summary update recipe by ID
// @id update-recipe
// @Produce json
// @Accept json
// @Param id path string true "update ID"
// @Param id body Recipe true "data values"
// @Success 200 {object} Recipe
// @Failure 404
// @Router /recipes/{id} [put]
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// // @Summary get recipe by tag name
// // @Produce json
// // @Param tag query string false "recipe search by tag"
// // @Success 200 {array} Recipe
// // @Failure 404
// // @Router /recipes/search [get]
// func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]models.Recipe, 0)

// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 				break
// 			}
// 		}
// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}

// 	if len(listOfRecipes) == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "No recipe matched that tag."})
// 		return
// 	}
// 	c.JSON(http.StatusOK, listOfRecipes)
// }
