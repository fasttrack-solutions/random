package main

import (
	"fmt"
	"github.com/fasttrack-solutions/random"
	"github.com/fasttrack-solutions/random/internal/config"
	"github.com/gin-gonic/gin"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())

	ginEngine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong @ %s", time.Now().UTC().String())
	})

	ginEngine.GET("/getRandomFloat64", func(c *gin.Context) {
		number, errUniformFloat64 := random.UniformFloat64()
		if errUniformFloat64 != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error generating random float64: %s", errUniformFloat64))
			c.Abort()
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%v", number))
	})

	ginEngine.GET("/getRandomInt64", func(c *gin.Context) {
		minimum := int32(0)
		minimumAsStr := c.Query("min")
		if len(minimumAsStr) == 0 {
			c.String(http.StatusBadRequest, "min is missing")
			c.Abort()
			return
		} else if len(minimumAsStr) > 10 {
			// cap so min + max avoids overflow
			c.String(http.StatusBadRequest, "min must be less than 2,147,483,647")
			c.Abort()
			return
		} else {
			minimumAsNumber, errParseInt := strconv.ParseInt(minimumAsStr, 10, 64)
			if errParseInt != nil {
				c.String(http.StatusBadRequest, "unable to parse min as number")
				c.Abort()
				return
			} else if minimumAsNumber < 0 || minimumAsNumber >= math.MaxInt32 {
				c.String(http.StatusBadRequest, "min must be between 0 and 2,147,483,647")
				c.Abort()
				return
			} else {
				minimum = int32(minimumAsNumber)
			}
		}

		maximum := int32(0)
		maximumAsStr := c.Query("max")
		if len(maximumAsStr) == 0 {
			c.String(http.StatusBadRequest, "max is missing")
			c.Abort()
			return
		} else if len(maximumAsStr) > 10 {
			// cap so min + max avoids overflow
			c.String(http.StatusBadRequest, "max must be less than 2,147,483,647")
			c.Abort()
			return
		} else {
			maximumAsNumber, errParseInt := strconv.ParseInt(maximumAsStr, 10, 64)
			if errParseInt != nil {
				c.String(http.StatusBadRequest, "unable to parse max as number")
				c.Abort()
				return
			} else if maximumAsNumber < 0 || maximumAsNumber >= math.MaxInt32 {
				c.String(http.StatusBadRequest, "max must be between 0 and 2,147,483,647")
				c.Abort()
				return
			} else {
				maximum = int32(maximumAsNumber)
			}
		}

		number, errUniformInt64 := random.UniformInt64(minimum, maximum)
		if errUniformInt64 != nil {
			c.String(http.StatusBadRequest, errUniformInt64.Error())
			c.Abort()
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%v", number))
	})

	// start server
	slog.Info(fmt.Sprintf("http server listening on %v", *config.HTTPPort))
	errRun := ginEngine.Run(fmt.Sprintf(":%v", *config.HTTPPort))
	if errRun != nil {
		panic(errRun)
	}
}
