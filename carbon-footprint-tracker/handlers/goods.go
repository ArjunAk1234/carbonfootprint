package handlers

import (
	"carbon-footprint-tracker/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GoodsPurchasedRequest struct {
	Date                time.Time `json:"date"`
	Location            string    `json:"location"`
	ItemName            string    `json:"item_name" binding:"required"`
	Category            string    `json:"category"`
	Quantity            int       `json:"quantity" binding:"required"`
	Unit                string    `json:"unit"`
	VendorName          string    `json:"vendor_name"`
	Origin              string    `json:"origin"`
	TransportMode       string    `json:"transport_mode"`
	TransportDistanceKM float64   `json:"transport_distance_km"`
	BillAmountINR       float64   `json:"bill_amount_inr" binding:"required"`
	BillAttachmentURL   string    `json:"bill_attachment_url"`
	PackagingType       string    `json:"packaging_type"`
	IsRecyclable        bool      `json:"is_recyclable"`
	Remarks             string    `json:"remarks"`
}

func AddGoodsPurchased(c *gin.Context) {
	var req GoodsPurchasedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Location == "" {
		req.Location = "Overall"
	}
	if req.Date.IsZero() {
		req.Date = time.Now()
	}

	g := models.GoodsPurchased{
		Date:                req.Date,
		Location:            req.Location,
		ItemName:            req.ItemName,
		Category:            toNullString(req.Category),
		Quantity:            req.Quantity,
		Unit:                toNullString(req.Unit),
		VendorName:          toNullString(req.VendorName),
		Origin:              toNullString(req.Origin),
		TransportMode:       toNullString(req.TransportMode),
		TransportDistanceKM: toNullFloat64(req.TransportDistanceKM),
		BillAmountINR:       req.BillAmountINR,
		BillAttachmentURL:   toNullString(req.BillAttachmentURL),
		PackagingType:       toNullString(req.PackagingType),
		IsRecyclable:        toNullBool(req.IsRecyclable),
		Remarks:             toNullString(req.Remarks),
	}

	if err := g.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Goods purchased added successfully", "id": g.ID})
}

func GetGoodsPurchased(c *gin.Context) {
	goods, err := models.GetAllGoodsPurchased()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve goods purchased", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, goods)
}

func UpdateGoodsPurchased(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	g, err := models.GetGoodsPurchasedByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Goods purchased entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve goods purchased", "details": err.Error()})
		return
	}

	var req GoodsPurchasedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g.Date = req.Date
	if req.Location != "" {
		g.Location = req.Location
	}
	if req.ItemName != "" {
		g.ItemName = req.ItemName
	}
	g.Category = toNullString(req.Category)
	if req.Quantity != 0 {
		g.Quantity = req.Quantity
	}
	g.Unit = toNullString(req.Unit)
	g.VendorName = toNullString(req.VendorName)
	g.Origin = toNullString(req.Origin)
	g.TransportMode = toNullString(req.TransportMode)
	g.TransportDistanceKM = toNullFloat64(req.TransportDistanceKM)
	if req.BillAmountINR != 0 {
		g.BillAmountINR = req.BillAmountINR
	}
	g.BillAttachmentURL = toNullString(req.BillAttachmentURL)
	g.PackagingType = toNullString(req.PackagingType)
	g.IsRecyclable = toNullBool(req.IsRecyclable)
	g.Remarks = toNullString(req.Remarks)

	if err := g.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goods purchased updated successfully"})
}

func DeleteGoodsPurchased(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteGoodsPurchased(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete goods purchased", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Goods purchased deleted successfully"})
}
