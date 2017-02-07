package main

import (
	"encoding/json"
	//	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"time"
)

func GlobalServiceStatus(c *gin.Context) {
	Lock.RLock()
	defer Lock.RUnlock()

	for _, value := range Lock.State {
		if value != "OK" {
			c.Writer.WriteHeader(http.StatusServiceUnavailable)
		}
	}
	bs, err := json.Marshal(Lock.State)
	if err != nil {
		//		TODO..do not panic; use a recovery handler
		panic(err)
	}
	c.Writer.Write(bs)
}

func SingleServiceStatus(c *gin.Context) {
	Lock.RLock()
	defer Lock.RUnlock()

	name := c.Param("name")
	//ts := name["testService"]
	if val, ok := Lock.State[name]; !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.Writer.Write([]byte("☄ hey!, the requested test service could not be found."))

	} else {
		bs, err := json.Marshal(val)
		if err != nil {
			//		TODO..do not panic; use a recovery handler
			panic(err)
		}
		if val != "OK" {
			c.Writer.WriteHeader(http.StatusServiceUnavailable)
		}
		c.Writer.Write(bs)
	}
}

func Logrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		c.Next()
		end := time.Now().UTC()
		latency := end.Sub(start)
		Logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"duration":   latency,
			"user_agent": c.Request.UserAgent(),
		}).Info()
	}
}
