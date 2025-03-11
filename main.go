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
	IP            string `json:"ip"`
	TTL           string `json:"ttl"`
	RecordType    string `json:"record_type"`
	RequestServer string `json:"request_server"`
}

type DigRequest struct {
	Domain       string `json:"query"`
	TypeOfRecord string `json:"type"`
}

func dig(req DigRequest) (DigResult, error) {
	var server = []string{"1.1.1.1", "8.8.8.8", "114.114.114.114", "2400:3200:baba::1", "2402:4e00::"}
	var answers []Answer
	for _, s := range server {
		cmd := exec.Command("dig", req.Domain, req.TypeOfRecord, "@"+s, "+noall", "+answer")
		out, err := cmd.Output()
		if err != nil {
			return DigResult{}, err
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if line != "" {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					answers = append(answers, Answer{
						IP:            fields[4],
						RecordType:    fields[3],
						TTL:           fields[1],
						RequestServer: s,
					})
				}
			}
		}
	}
	result := DigResult{Answer: answers, Error: ""}
	return result, nil
}

func handler(c *gin.Context) {
	var req DigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}
	digResult, err := dig(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(200, digResult)
}

func main() {
	r := gin.Default()
	r.POST("/dig", handler)
	_ = r.Run(":8080")
}
