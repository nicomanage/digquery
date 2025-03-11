package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
)

type DigResult struct {
	Answer []Answer `json:"answer"`
}

type Answer struct {
	IP            string `json:"ip,omitempty"`
	TTL           string `json:"ttl,omitempty"`
	RecordType    string `json:"record_type,omitempty"`
	RequestServer string `json:"request_server,omitempty"`
	Error         string `json:"error,omitempty"`
}

type DigRequest struct {
	Domain       string `json:"query"`
	TypeOfRecord string `json:"type"`
}

func dig(req DigRequest) DigResult {
	var server = []string{"1.1.1.1", "8.8.8.8", "114.114.114.114", "2400:3200:baba::1", "2402:4e00::"}
	var answers []Answer
	for _, s := range server {
		cmd := exec.Command("dig", req.Domain, req.TypeOfRecord, "@"+s, "+noall", "+answer")
		out, err := cmd.Output()
		if err != nil {
			answers = append(answers, Answer{
				RequestServer: s,
				Error:         err.Error(),
			})
			continue
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, ";") {
				continue
			}
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
			} else {
				answers = append(answers, Answer{
					RequestServer: s,
					Error:         "No answer",
				})
			}
		}
	}
	result := DigResult{Answer: answers}
	return result
}

func handler(c *gin.Context) {
	var req DigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON body"})
		return
	}
	digResult := dig(req)
	c.JSON(200, digResult)
}

func main() {
	r := gin.Default()
	r.POST("/dig", handler)
	_ = r.Run(":8080")
}
