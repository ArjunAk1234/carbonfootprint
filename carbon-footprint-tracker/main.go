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
	config.LoadConfig()
	defer config.DB.Close()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8000", "http://localhost:8000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", handlers.RegisterUser)
		authRoutes.POST("/login", handlers.LoginUser)
	}

	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthRequired())
	{

		userRoutes := authenticated.Group("/users")
		userRoutes.Use(middleware.AuthorizeRoles("admin"))
		{
			userRoutes.GET("", handlers.GetUsers)
			userRoutes.POST("", handlers.AddUser)
			userRoutes.PUT("/:id", handlers.UpdateUser)
			userRoutes.DELETE("/:id", handlers.DeleteUser)
		}

		electricRoutes := authenticated.Group("/electric")
		electricRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			electricRoutes.GET("", handlers.GetElectricConsumptions)
			electricRoutes.POST("", handlers.AddElectricConsumption)
			electricRoutes.PUT("/:id", handlers.UpdateElectricConsumption)
			electricRoutes.DELETE("/:id", handlers.DeleteElectricConsumption)
		}

		populationRoutes := authenticated.Group("/population")
		populationRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			populationRoutes.GET("", handlers.GetPopulationStats)
			populationRoutes.POST("", handlers.AddPopulation)
			populationRoutes.PUT("/:id", handlers.UpdatePopulation)
			populationRoutes.DELETE("/:id", handlers.DeletePopulation)
		}

		transportRoutes := authenticated.Group("/transport")
		transportRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			transportRoutes.GET("", handlers.GetTransportData)
			transportRoutes.POST("", handlers.AddTransportData)
			transportRoutes.PUT("/:id", handlers.UpdateTransportData)
			transportRoutes.DELETE("/:id", handlers.DeleteTransportData)
		}

		waterRoutes := authenticated.Group("/water")
		waterRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			waterRoutes.GET("", handlers.GetWaterReadings)
			waterRoutes.POST("", handlers.AddWaterReading)
			waterRoutes.PUT("/:id", handlers.UpdateWaterReading)
			waterRoutes.DELETE("/:id", handlers.DeleteWaterReading)
		}

		wasteRoutes := authenticated.Group("/waste")
		wasteRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			wasteRoutes.GET("", handlers.GetWasteData)
			wasteRoutes.POST("", handlers.AddWasteEntry)
			wasteRoutes.PUT("/:id", handlers.UpdateWasteEntry)
			wasteRoutes.DELETE("/:id", handlers.DeleteWasteEntry)
		}

		accommodationRoutes := authenticated.Group("/accommodation")
		accommodationRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			accommodationRoutes.GET("", handlers.GetAccommodationData)
			accommodationRoutes.POST("", handlers.AddAccommodationData)
			accommodationRoutes.PUT("/:id", handlers.UpdateAccommodationData)
			accommodationRoutes.DELETE("/:id", handlers.DeleteAccommodationData)
		}

		goodsRoutes := authenticated.Group("/goods")
		goodsRoutes.Use(middleware.AuthorizeRoles("admin", "staff"))
		{
			goodsRoutes.GET("", handlers.GetGoodsPurchased)
			goodsRoutes.POST("", handlers.AddGoodsPurchased)
			goodsRoutes.PUT("/:id", handlers.UpdateGoodsPurchased)
			goodsRoutes.DELETE("/:id", handlers.DeleteGoodsPurchased)
		}

		dashboardRoutes := authenticated.Group("/dashboard")
		dashboardRoutes.Use(middleware.AuthorizeRoles("admin", "staff", "viewer"))
		{
			dashboardRoutes.GET("", handlers.GetDashboardSummary)
		}
	}

	log.Fatal(router.Run(":8080"))
}
