package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/majorfi/ydaemon/internal/ethereum"
	"github.com/majorfi/ydaemon/internal/logs"
	"github.com/majorfi/ydaemon/internal/models"
	"github.com/majorfi/ydaemon/internal/utils"
)

func graphQLRequestForReports(strategyAddress string, c *gin.Context) *graphql.Request {
	return graphql.NewRequest(`{
		strategy(id: "` + strings.ToLower(strategyAddress) + `") {
			` + utils.GetStrategyReports() + `
		}
	}`)
}

//GetReports will, for a given strategy on a given chainID, return a list of reports
func (y controller) GetReports(c *gin.Context) {
	chainID, err := strconv.ParseUint(c.Param("chainID"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid chainID")
		return
	}
	strategyAddress := c.Param("address")
	if strategyAddress == `` {
		c.String(http.StatusBadRequest, "invalid address")
		return
	}

	client := graphql.NewClient(ethereum.GetGraphURI(chainID))
	request := graphQLRequestForReports(strategyAddress, c)
	var response models.TReportsFromGraph
	if err := client.Run(context.Background(), request, &response); err != nil {
		logs.Error(err)
		c.String(http.StatusInternalServerError, "Impossible to fetch subgraph")
		return
	}

	c.JSON(http.StatusOK, response.Strategy.Reports)
}
