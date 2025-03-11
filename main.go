package main

import (
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
	"sync"
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
	var wg sync.WaitGroup
	for _, s := range server {
		wg.Add(1)
		go digCommand(&wg, &answers, req.Domain, s, req.TypeOfRecord)
	}
	wg.Wait()
	result := DigResult{Answer: answers}
	return result
}

func digCommand(wg *sync.WaitGroup, answers *[]Answer, domain string, server string, recordType string) {
	defer wg.Done()
	cmd := exec.Command("dig", domain, recordType, "@"+server, "+noall", "+answer")
	out, err := cmd.Output()
	if err != nil {
		*answers = append(*answers, Answer{
			RequestServer: server,
			Error:         err.Error(),
		})
		return
	}

	lines := strings.Split(string(out), "\n")
	haveAnswer := false
	for _, line := range lines {
		if strings.HasPrefix(line, ";") {
			continue
		}
		if line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				haveAnswer = true
				*answers = append(*answers, Answer{
					IP:            fields[4],
					RecordType:    fields[3],
					TTL:           fields[1],
					RequestServer: server,
				})
			}
		} else {
			if !haveAnswer {
				*answers = append(*answers, Answer{
					RequestServer: server,
					Error:         "No answer",
				})
			}
		}
	}
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
