package handlers

import (
	pbidentity "black-shop/api/proto/identity/v1"
)

type IdentityHandler struct {
	client pbidentity.IdentityServiceClient
}

func NewIdentityHandler(client pbidentity.IdentityServiceClient) *IdentityHandler {
	return &IdentityHandler{client: client}
}

//func (h *IdentityHandler) Register(c *gin.Context) {
//	var req pbidentity.RegisterRequest
//	if err := c.ShouldBindJSON(&req); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	res, err := h.client.Register(c.Request.Context(), &req)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, res)
//}
