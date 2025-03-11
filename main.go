package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
)

type DigResult struct {
	Answer []Answer `json:"answer"`
	Error  string   `json:"error,omitempty"`
}

type Answer struct {
	IP         string `json:"ip"`
	TTL        string `json:"ttl"`
	RecordType string `json:"record_type"`
}

type DigRequest struct {
	Domain       string `json:"domain"`
	TypeOfRecord string `json:"type"`
	Server       string `json:"server"`
}

func handler(c *gin.Context) {
	var req DigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}

	if req.Domain == "" {
		c.JSON(400, gin.H{"error": "Missing 'domain' parameter"})
		return
	}

	cmd := exec.Command("dig", req.Domain, req.TypeOfRecord, "@"+req.Server, "+noall", "+answer")
	out, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	lines := strings.Split(string(out), "\n")
	var answers []Answer

	for _, line := range lines {
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				answers = append(answers, Answer{
					IP:         fields[4],
					RecordType: fields[3],
					TTL:        fields[1],
				})
			}
		}
	}
	result := DigResult{Answer: answers, Error: ""}

	c.JSON(200, result)
}

func main() {
	r := gin.Default()
	r.POST("/dig", handler)
	_ = r.Run(":8080")
}
