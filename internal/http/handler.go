package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func (r *Server) handleBalance(c *gin.Context) {
	var userID int64
	if v, err := strconv.ParseInt(c.Param("user_id"), 0, 64); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	} else {
		userID = v
	}

	v := r.balanceService.GetByUserID(userID)

	output := &struct {
		Balance string `json:"balance"`
	}{
		Balance: v.String(),
	}

	c.JSON(http.StatusOK, output)
}

func (r *Server) handleDeposit(c *gin.Context) {
	input := &struct {
		UserID int64  `json:"user_id"`
		Value  string `json:"value"`
	}{}
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	v, err := decimal.NewFromString(input.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if v.LessThan(decimal.NewFromInt(0)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid value")
		return
	}

	if err := r.balanceService.Deposit(input.UserID, v); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid value")
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (r *Server) handleWithdraw(c *gin.Context) {
	input := &struct {
		UserID int64  `json:"user_id"`
		Value  string `json:"value"`
	}{}
	if err := c.BindJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	v, err := decimal.NewFromString(input.Value)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	if v.LessThan(decimal.NewFromInt(0)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid value")
		return
	}

	if err := r.balanceService.Withdraw(input.UserID, v); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid value")
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
