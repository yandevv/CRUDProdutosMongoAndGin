package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	_ "crudProdutosMongoAndGin/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Structs
type Product struct {
	ID           bson.ObjectID `json:"id" example:"1" bson:"_id"`                                                            // Unique identifier for the product
	Name         string        `json:"name" example:"Paperclip" bson:"name"`                                                 // Name of the product
	Description  string        `json:"description" example:"It's a object used to hold papers together." bson:"description"` // Description of the product
	Color        string        `json:"color" example:"gray" bson:"color"`                                                    // Color of the product
	Weight       float32       `json:"weight" example:"0.05" bson:"weight"`                                                  // Weight of the product in kilograms
	Type         string        `json:"type" example:"office_supplies" bson:"type"`                                           // Type of the product (e.g., office supplies, electronics, etc.)
	Price        float32       `json:"price" example:"0.1" bson:"price"`                                                     // Price of the product in the local currency
	RegisterDate time.Time     `json:"registerDate" example:"2025-06-22T19:57:53.788Z" bson:"register_date"`                 // Register date of the product
}

type CreateProductDTO struct {
	Name        string  `json:"name"  example:"Paperclip" binding:"required"`                      // Name of the product
	Description string  `json:"description" example:"It's a object used to hold papers together."` // Description of the product
	Color       string  `json:"color" example:"gray"`                                              // Color of the product
	Weight      float32 `json:"weight" example:"0.05" binding:"required"`                          // Weight of the product in kilograms
	Type        string  `json:"type" example:"office_supplies" binding:"required"`                 // Type of the product (e.g., office supplies, electronics, etc.)
	Price       float32 `json:"price" example:"0.1" binding:"required"`                            // Price of the product in the local currency
}

type UpdateProductDTO struct {
	Name         string    `json:"name" example:"Paperclip" bson:"name"`                                                 // Updated name of the product
	Description  string    `json:"description" example:"It's a object used to hold papers together." bson:"description"` // Updated description of the product
	Color        string    `json:"color" example:"gray" bson:"color"`                                                    // Updated color of the product
	Weight       float32   `json:"weight" example:"0.05" bson:"weight"`                                                  // Updated weight of the product
	Type         string    `json:"type" example:"office_supplies" bson:"type"`                                           // Updated type of the product
	Price        float32   `json:"price" example:"0.1" bson:"price"`                                                     // Updated price of the product
	RegisterDate time.Time `json:"registerDate" example:"2025-06-22T19:57:53.788Z" bson:"register_date"`                 // Updated register date of the product
}

type GetAllProductsResponse struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Data       []Product `json:"data"`
}

type GetProductByIdResponse struct {
	StatusCode int     `json:"statusCode"`
	Message    string  `json:"message"`
	Data       Product `json:"data"`
}

type DefaultResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// GetProducts godoc
//
//	@Summary		Get all products
//	@Description	Get a list of all products
//	@ID				get-all-products
//	@Tags			products
//	@Produce		json
//	@Success		200	{object}	GetAllProductsResponse	"Products fetched successfully"
//	@Failure		500	{object}	DefaultResponse			"Failed to fetch products"
//	@Router			/products [get]
func getAllProductsHandler(c *gin.Context, mongoClient *mongo.Client) {
	productsColl := mongoClient.Database("defaultDatabase").Collection("products")

	foundProductsCursor, err := productsColl.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"message":    "Failed to fetch products",
		})
		log.Printf("Error fetching products: %v", err)
		return
	}

	products := make([]Product, 0)

	err = foundProductsCursor.All(context.TODO(), &products)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"message":    "Failed to fetch products",
		})
		log.Printf("Error decoding products: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "Products fetched successfully",
		"data":       products,
	})
}

// GetProductById godoc
//
//	@Summary		Get a product by ID
//	@Description	Get a product by its ID
//	@ID				get-product-by-id
//	@Tags			products
//	@Param			id	path	string	true	"Product ID"
//	@Produce		json
//	@Success		200	{object}	GetProductByIdResponse	"Product fetched successfully"
//	@Failure		400	{object}	DefaultResponse			"Invalid product ID"
//	@Failure		404	{object}	DefaultResponse			"Product not found"
//	@Failure		500	{object}	DefaultResponse			"Failed to fetch product"
//	@Router			/products/{id} [get]
func getProductByIdHandler(c *gin.Context, mongoClient *mongo.Client) {
	paramId := c.Params.ByName("id")

	productId, err := bson.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid product ID",
		})
		return
	}

	productsColl := mongoClient.Database("defaultDatabase").Collection("products")

	var product Product
	err = productsColl.FindOne(context.TODO(), bson.D{{Key: "_id", Value: productId}}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": 404,
				"message":    "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"message":    "Failed to fetch product",
		})
		log.Printf("Error fetching product: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "Product fetched successfully",
		"data":       product,
	})
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Create a new product with the provided details
//	@ID				create-product
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		CreateProductDTO	true	"Product details"
//	@Success		201		{object}	DefaultResponse		"Product created successfully"
//	@Failure		400		{object}	DefaultResponse		"Invalid request body data"
//	@Failure		500		{object}	DefaultResponse		"Failed to create product"
//	@Router			/products [post]
func createProductHandler(c *gin.Context, mongoClient *mongo.Client) {
	var createProductDTO CreateProductDTO

	if c.Bind(&createProductDTO) == nil {
		productsColl := mongoClient.Database("defaultDatabase").Collection("products")

		insertData := bson.M{
			"name":          createProductDTO.Name,
			"description":   createProductDTO.Description,
			"color":         createProductDTO.Color,
			"weight":        createProductDTO.Weight,
			"type":          createProductDTO.Type,
			"price":         createProductDTO.Price,
			"register_date": time.Now(),
		}

		_, err := productsColl.InsertOne(context.TODO(), insertData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": 500,
				"message":    "Failed to create product",
			})
			log.Printf("Error creating product: %v", err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"statusCode": 201,
			"message":    "Product created successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid request body data",
		})
		return
	}
}

// UpdateProduct godoc
//
//	@Summary		Update an existing product
//	@Description	Update an existing product with the provided details
//	@ID				update-product
//	@Tags			products
//	@Param			id	path	string	true	"Product ID"
//	@Accept			json
//	@Produce		json
//	@Param			product	body		UpdateProductDTO	true	"Product details"
//	@Success		200		{object}	DefaultResponse		"Product updated successfully"
//	@Failure		400		{object}	DefaultResponse		"Invalid product ID or request body data"
//	@Failure		404		{object}	DefaultResponse		"Product not found"
//	@Failure		500		{object}	DefaultResponse		"Failed to update product"
//	@Router			/products/{id} [put]
func updateProductHandler(c *gin.Context, mongoClient *mongo.Client) {
	paramId := c.Params.ByName("id")

	productId, err := bson.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid product ID",
		})
		return
	}

	var updateProductDTO UpdateProductDTO

	if c.Bind(&updateProductDTO) == nil {
		productsColl := mongoClient.Database("defaultDatabase").Collection("products")

		_, err := productsColl.UpdateOne(
			context.TODO(),
			bson.D{{Key: "_id", Value: productId}},
			bson.D{{Key: "$set", Value: updateProductDTO}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": 500,
				"message":    "Failed to update product",
			})
			log.Printf("Error updating product: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"statusCode": 200,
			"message":    "Product updated successfully",
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid request body data",
		})
		return
	}
}

// DeleteProductById godoc
//
//	@Summary		Delete a product by ID
//	@Description	Delete a product by its ID
//	@ID				delete-product-by-id
//	@Tags			products
//	@Param			id	path	string	true	"Product ID"
//	@Produce		json
//	@Success		200	{object}	DefaultResponse	"Product deleted successfully"
//	@Failure		400	{object}	DefaultResponse	"Invalid product ID"
//	@Failure		404	{object}	DefaultResponse	"Product not found"
//	@Failure		500	{object}	DefaultResponse	"Failed to delete product"
//	@Router			/products/{id} [delete]
func deleteProductByIdHandler(c *gin.Context, mongoClient *mongo.Client) {
	paramId := c.Params.ByName("id")

	productId, err := bson.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid product ID",
		})
		return
	}

	productsColl := mongoClient.Database("defaultDatabase").Collection("products")

	deleteProductCursor, err := productsColl.DeleteOne(
		context.TODO(),
		bson.D{{Key: "_id", Value: productId}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"message":    "Failed to delete product",
		})
		log.Printf("Error deleting product: %v", err)
		return
	}

	if deleteProductCursor.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": 404,
			"message":    "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "Product deleted successfully",
	})
}

func setupRouter(mongoClient *mongo.Client) *gin.Engine {
	r := gin.Default()
	//	@BasePath	/

	r.GET("/products", func(c *gin.Context) { getAllProductsHandler(c, mongoClient) })

	r.GET("/products/:id", func(c *gin.Context) { getProductByIdHandler(c, mongoClient) })

	r.POST("/products", func(c *gin.Context) { createProductHandler(c, mongoClient) })

	r.PUT("/products/:id", func(c *gin.Context) { updateProductHandler(c, mongoClient) })

	r.DELETE("/products/:id", func(c *gin.Context) { deleteProductByIdHandler(c, mongoClient) })

	return r
}

// @title			Product's CRUD API
// @version		1.0
// @description	This is a sample CRUD API for managing products.
// @host			localhost:8080
// @schemes		http
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	mongoConnectionUri := os.Getenv("MONGODB_CONNECTION_URI")

	if mongoConnectionUri == "" {
		log.Fatal("Environment variable MONGODB_CONNECTION_URI is not set")
		os.Exit(1)
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(mongoConnectionUri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := setupRouter(client)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := os.Getenv("PORT")

	if port != "" {
		log.Printf("Starting server on port %s", port)
	} else {
		port = "8080"
		log.Println("PORT environment variable not set, using default port 8080")
	}
	router.Run(":" + port)
}
