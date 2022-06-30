package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/semirm-dev/myt/product"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetProducts(client product.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var priceLessThan int
		priceLessThanParam, ok := c.GetQuery("priceLessThan")
		if ok {
			plt, err := strconv.Atoi(priceLessThanParam)
			if err != nil {
				plt = 0
			}
			priceLessThan = plt
		}

		byCategory, _ := c.GetQuery("category")

		products, err := client.GetProductsByFilter(c.Request.Context(), &product.Filter{
			PriceLessThan: priceLessThan,
			ByCategory:    byCategory,
		})
		if err != nil {
			logrus.Error(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(
			http.StatusOK,
			products,
		)
	}
}
