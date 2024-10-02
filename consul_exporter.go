package main

import (
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
)

const consulURL = "http://localhost:8500/v1/agent/metrics?format=prometheus"

func main() {
    // Get the Consul token from environment variables
    consulToken := os.Getenv("CONSUL_TOKEN")
    if consulToken == "" {
        log.Fatal("CONSUL_TOKEN environment variable is required")
    }

    // Create a new Gin router
    r := gin.Default()

    r.GET("/metrics", func(c *gin.Context) {
        // Make the request to Consul
        req, err := http.NewRequest("GET", consulURL, nil)
        if err != nil {
            c.String(http.StatusInternalServerError, "Error creating request: %v", err)
            return
        }

        // Add the Consul token to the request header
        req.Header.Add("X-Consul-Token", consulToken)

        // Perform the request
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            c.String(http.StatusInternalServerError, "Error making request to Consul: %v", err)
            return
        }
        defer resp.Body.Close()

        // Check the response status code
        if resp.StatusCode != http.StatusOK {
            c.String(resp.StatusCode, "Error fetching metrics from Consul: %s", resp.Status)
            return
        }

        // Read the response body
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            c.String(http.StatusInternalServerError, "Error reading response body: %v", err)
            return
        }

        // Return the response body as the response
        c.Data(http.StatusOK, "text/plain", body)
    })

    // Start the HTTP server
    r.Run(":9107")
}

