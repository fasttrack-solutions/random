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
	"strings"
	"time"
)

func main() {
	seed := *config.SEEDHEX
	if len(seed) != 64 {
		panic("seed must be 64 hex characters")
	} else if seed == "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("a unique seed value is required")
	}

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

	ginEngine.GET("/getDeterministicRandom", func(c *gin.Context) {
		sequence := int64(0)
		sequenceAsStr := c.Query("s")
		if len(sequenceAsStr) == 0 {
			c.String(http.StatusBadRequest, "sequence is missing")
			c.Abort()
			return
		} else {
			sequenceAsNumber, errParseInt := strconv.ParseInt(sequenceAsStr, 10, 64)
			if errParseInt != nil {
				c.String(http.StatusBadRequest, "unable to parse sequence as number")
				c.Abort()
				return
			} else if sequenceAsNumber < 0 || sequenceAsNumber >= math.MaxInt64 {
				c.String(http.StatusBadRequest, "sequence must be between 0 and 9,223,372,036,854,775,806")
				c.Abort()
				return
			} else {
				sequence = sequenceAsNumber
			}
		}

		var probabilities []float64
		probabilitiesAsStr := c.Query("p")
		if len(probabilitiesAsStr) == 0 {
			c.String(http.StatusBadRequest, "probabilities are missing")
			c.Abort()
			return
		} else if len(probabilitiesAsStr) > 300 {
			c.String(http.StatusBadRequest, "string of probabilities must be less than 300 characters")
			c.Abort()
			return
		} else {
			probabilitiesAsStrList := strings.Split(probabilitiesAsStr, ",")
			if len(probabilitiesAsStrList) == 0 {
				c.String(http.StatusBadRequest, "invalid probabilities, use comma separated list (i.e. 0.01,0.09,0.9")
				c.Abort()
				return
			}

			for _, v := range probabilitiesAsStrList {
				probability, errParse := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if errParse != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("invalid probability: %s", strings.TrimSpace(v)))
					c.Abort()
					return
				}
				probabilities = append(probabilities, probability)
			}
		}

		number, errDeterministicRandom := random.DeterministicRandom(seed, sequence, probabilities)
		if errDeterministicRandom != nil {
			c.String(http.StatusBadRequest, errDeterministicRandom.Error())
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
