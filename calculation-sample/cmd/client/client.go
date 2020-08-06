package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	calc "github.com/rinosukmandityo/grpc-sample/calculation-sample/calculation"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	addr := ":8989"
	cc, e := grpc.Dial(addr, grpc.WithInsecure())
	if e != nil {
		log.Println("could not connect to %s: %v", addr, e)
	}

	client := calc.NewCalculationServiceClient(cc)

	r := gin.Default()

	r.GET("/add/:a/:b", func(ctx *gin.Context) {
		a, e := strconv.ParseUint(ctx.Param("a"), 10, 64)
		if e != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
			return
		}
		b, e := strconv.ParseUint(ctx.Param("b"), 10, 64)
		if e != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
			return
		}

		if resp, e := client.Add(ctx, &calc.Request{A: int64(a), B: int64(b)}); e == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(resp.GetResult())})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		}

	})

	r.GET("/mult/:a/:b", func(ctx *gin.Context) {
		a, e := strconv.ParseUint(ctx.Param("a"), 10, 64)
		if e != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
			return
		}
		b, e := strconv.ParseUint(ctx.Param("b"), 10, 64)
		if e != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
			return
		}

		if resp, e := client.Multiply(ctx, &calc.Request{A: int64(a), B: int64(b)}); e == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(resp.GetResult())})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
		}
	})

	addrAPI := ":8080"
	log.Fatal(r.Run(addrAPI))
}
