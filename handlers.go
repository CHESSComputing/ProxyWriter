package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	services "github.com/CHESSComputing/golib/services"
	"github.com/gin-gonic/gin"
)

// helper function to exgract JSON dict from HTTP request
func parseRequest(c *gin.Context) (map[string]any, error) {
	var spec map[string]any
	defer c.Request.Body.Close()
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return spec, fmt.Errorf("[UserProxyWriter.main.parseRequest] io.ReadAll error: %w", err)
	}
	err = json.Unmarshal(body, &spec)
	if err != nil {
		return spec, fmt.Errorf("[UserProxyWriter.main.parseRequest] json.Unmarshal error: %w", err)
	}
	return spec, nil
}

// PostHandler handles POST upload of meta-data record
func PostHandler(c *gin.Context) {
	rec, err := parseRequest(c)
	if err != nil {
		log.Println("ERROR:", err)
		rec := services.Response("ProxyWriter", http.StatusInternalServerError, services.ParseError, err)
		c.JSON(http.StatusInternalServerError, rec)
		return
	}

	// insert record to meta-data database
	did, err := insertData(rec)
	if err != nil {
		log.Println("ERROR:", err)
		rec := services.Response("ProxyWriter", http.StatusInternalServerError, services.InsertError, err)
		c.JSON(http.StatusInternalServerError, rec)
		return
	}
	var records []map[string]any
	resp := services.Response("ProxyWriter", http.StatusOK, services.OK, nil)
	r := make(map[string]any)
	r["did"] = did
	records = append(records, r)
	resp.Results = services.ServiceResults{NRecords: 1, Records: records}
	c.JSON(http.StatusOK, resp)
}
