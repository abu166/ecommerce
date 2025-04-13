package apigateway

import (
	"ecommerce/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Server) SetupRoutes(r *gin.Engine) {
	r.Use(s.Logger(), s.Auth())

	r.POST("/products", s.createProduct)
	r.GET("/products/:id", s.getProduct)
	r.PATCH("/products/:id", s.updateProduct)
	r.DELETE("/products/:id", s.deleteProduct)
	r.GET("/products", s.listProducts)

	r.POST("/orders", s.createOrder)
	r.GET("/orders/:id", s.getOrder)
	r.PATCH("/orders/:id", s.updateOrder)
	r.GET("/orders", s.listOrders)

	r.POST("/users/register", s.registerUser)
	r.POST("/users/login", s.login)
	r.GET("/users/:id", s.getUser)
}

func (s *Server) createProduct(c *gin.Context) {
	var req proto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.invClient.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getProduct(c *gin.Context) {
	id := c.Param("id")
	resp, err := s.invClient.GetProduct(c.Request.Context(), &proto.GetProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) updateProduct(c *gin.Context) {
	var req proto.UpdateProductRequest
	req.Id = c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.invClient.UpdateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) deleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := s.invClient.DeleteProduct(c.Request.Context(), &proto.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func (s *Server) listProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")
	resp, err := s.invClient.ListProducts(c.Request.Context(), &proto.ListProductsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Category: category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) createOrder(c *gin.Context) {
	var req proto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	req.UserId = userID.(string)
	resp, err := s.ordClient.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getOrder(c *gin.Context) {
	id := c.Param("id")
	resp, err := s.ordClient.GetOrder(c.Request.Context(), &proto.GetOrderRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) updateOrder(c *gin.Context) {
	var req proto.UpdateOrderRequest
	req.Id = c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.ordClient.UpdateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) listOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	resp, err := s.ordClient.ListOrders(c.Request.Context(), &proto.ListOrdersRequest{
		UserId:   userID.(string),
		Page:     int32(page),
		PageSize: int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) registerUser(c *gin.Context) {
	var req proto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.usrClient.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) login(c *gin.Context) {
	var req proto.AuthenticateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := s.usrClient.AuthenticateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getUser(c *gin.Context) {
	id := c.Param("id")
	resp, err := s.usrClient.GetUserProfile(c.Request.Context(), &proto.GetUserProfileRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
