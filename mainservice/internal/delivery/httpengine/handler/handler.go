package handler

import (
	"csu-lessons/internal/core/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	orderClient OrderClient
}

func NewHandler(client OrderClient) *Handler {
	return &Handler{orderClient: client}
}

type orderRequest struct {
	UserId int64  `json:"user_id" validate:"required"`
	Items  []Item `json:"items" validate:"required"`
}

type Item struct {
	ProductId int   `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required"`
}

func (o *orderRequest) toDomain() entity.Order {
	items := make([]entity.Item, len(o.Items))
	for i, item := range o.Items {
		items[i] = entity.Item(item)
	}

	return entity.Order{
		UserId: o.UserId,
		Items:  items,
	}
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req orderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	//v := &validator.Validate{}
	//err := v.Struct(req)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//}

	err := h.orderClient.CreateOrder(c.Request.Context(), req.toDomain())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusCreated)

}
