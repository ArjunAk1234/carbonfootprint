package main

import (
	"carbon-footprint-tracker/config"
	"carbon-footprint-tracker/handlers"
	"carbon-footprint-tracker/middleware"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration (DB connection, JWT secret)
	config.LoadConfig()
	defer config.DB.Close() // Ensure DB connection is closed when main exits

	// Initialize Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8000", "http://localhost:8000", "http://localhost:5500", "http://127.0.0.1:5500"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public Routes (No authentication required)
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handlers.RegisterUser)
		authRoutes.POST("/login", handlers.LoginUser)
	}

	// Authenticated Routes (Requires JWT token)
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthRequired())
	{
		// User Management (Admin only)
		userRoutes := authenticated.Group("/users")
		userRoutes.Use(middleware.AuthorizeRoles("admin"))
		{
			userRoutes.GET("", handlers.GetUsers)
			userRoutes.POST("", handlers.AddUser)
			userRoutes.PUT("/:id", handlers.UpdateUser)
			userRoutes.DELETE("/:id", handlers.DeleteUser)
		}

		// Electric Consumption
		electricRoutes := authenticated.Group("/electric")
		electricRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			electricRoutes.GET("", handlers.GetElectricConsumptions)
			electricRoutes.POST("", handlers.AddElectricConsumption)
			electricRoutes.PUT("/:id", handlers.UpdateElectricConsumption)
			electricRoutes.DELETE("/:id", handlers.DeleteElectricConsumption)
		}

		// Population
		populationRoutes := authenticated.Group("/population")
		populationRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			populationRoutes.GET("", handlers.GetPopulationStats)
			populationRoutes.POST("", handlers.AddPopulation)
			populationRoutes.PUT("/:id", handlers.UpdatePopulation)
			populationRoutes.DELETE("/:id", handlers.DeletePopulation)
		}

		// Transport
		transportRoutes := authenticated.Group("/transport")
		transportRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			transportRoutes.GET("", handlers.GetTransportData)
			transportRoutes.POST("", handlers.AddTransportData)
			transportRoutes.PUT("/:id", handlers.UpdateTransportData)
			transportRoutes.DELETE("/:id", handlers.DeleteTransportData)
		}

		// Water Consumption (Usage)
		waterConsumptionRoutes := authenticated.Group("/water_consumption")
		waterConsumptionRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			waterConsumptionRoutes.GET("", handlers.GetWaterConsumptions)
			waterConsumptionRoutes.POST("", handlers.AddWaterConsumption)
			waterConsumptionRoutes.PUT("/:id", handlers.UpdateWaterConsumption)
			waterConsumptionRoutes.DELETE("/:id", handlers.DeleteWaterConsumption)
		}

		// Water Treatment (NEW Module)
		waterTreatmentRoutes := authenticated.Group("/water_treatment")
		waterTreatmentRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			waterTreatmentRoutes.GET("", handlers.GetWaterTreatments)
			waterTreatmentRoutes.POST("", handlers.AddWaterTreatment)
			waterTreatmentRoutes.PUT("/:id", handlers.UpdateWaterTreatment)
			waterTreatmentRoutes.DELETE("/:id", handlers.DeleteWaterTreatment)
		}

		// Waste Generation
		wasteRoutes := authenticated.Group("/waste")
		wasteRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			wasteRoutes.GET("", handlers.GetWasteData)
			wasteRoutes.POST("", handlers.AddWasteEntry)
			wasteRoutes.PUT("/:id", handlers.UpdateWasteEntry)
			wasteRoutes.DELETE("/:id", handlers.DeleteWasteEntry)
		}

		// Accommodation
		accommodationRoutes := authenticated.Group("/accommodation")
		accommodationRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			accommodationRoutes.GET("", handlers.GetAccommodationData)
			accommodationRoutes.POST("", handlers.AddAccommodationData)
			accommodationRoutes.PUT("/:id", handlers.UpdateAccommodationData)
			accommodationRoutes.DELETE("/:id", handlers.DeleteAccommodationData)
		}

		// Goods Purchased
		goodsRoutes := authenticated.Group("/goods")
		goodsRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			goodsRoutes.GET("", handlers.GetGoodsPurchased)
			goodsRoutes.POST("", handlers.AddGoodsPurchased)
			goodsRoutes.PUT("/:id", handlers.UpdateGoodsPurchased)
			goodsRoutes.DELETE("/:id", handlers.DeleteGoodsPurchased)
		}

		// Food Consumption (NEW Module)
		foodConsumptionRoutes := authenticated.Group("/food_consumption")
		foodConsumptionRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			foodConsumptionRoutes.GET("", handlers.GetFoodConsumptions)
			foodConsumptionRoutes.POST("", handlers.AddFoodConsumption)
			foodConsumptionRoutes.PUT("/:id", handlers.UpdateFoodConsumption)
			foodConsumptionRoutes.DELETE("/:id", handlers.DeleteFoodConsumption)
		}

		// Dashboard (Admin and Viewer roles)
		dashboardRoutes := authenticated.Group("/dashboard")
		dashboardRoutes.Use(middleware.AuthorizeRoles("admin", "staff", "viewer"))
		{
			dashboardRoutes.GET("", handlers.GetDashboardSummary)
		}
	}

	// Run the server
	log.Fatal(router.Run(":8080")) // Listen and serve on 0.0.0.0:8080
}
