package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gonce/config"
	"github.com/gonce/models"
	"github.com/gonce/utils"
	"github.com/olekukonko/tablewriter"
)

// PerformService struct to aggregate the metrics data and PerformServiceMethods
type PerformService struct {
	Metrics               models.Metrics
	PerformServiceMethods PerformServiceMethods
}

// PerformServiceMethods containing the method definations
type PerformServiceMethods interface {
	Perform(config.ConditionConfig, utils.HTTPClient)
	processCurrentMetric(models.CurrentMetrics)
	CreateTable()
}

// Perform method contains the core performance test logic
// Accepts the config and the initialised http client
// Creates three chans:
// 1. Create a thread rampup chan, 2. Timeout chan, 3. Current metric chan
// Select used to individually process the chan data
func (i *PerformService) Perform(c config.ConditionConfig, httpClient utils.HTTPClient) {
	validateInput(c)
	ch := make(chan int64)
	go rampUpThreads(c, ch)
	timeoutCh := time.After(time.Duration(c.ExecutionTimeInSeconds) * time.Second)
	currentMetricsCh := make(chan models.CurrentMetrics)
	for {
		select {
		case <-timeoutCh:
			log.Print("Execution completed")
			return
		case threadStart := <-ch:
			log.Printf("Thread number %d spawned", threadStart)
			go func() {
				for {
					httpTest(c.HTTPRequest, httpClient, currentMetricsCh)
				}
			}()
		case currentMetric := <-currentMetricsCh:
			i.processCurrentMetric(currentMetric)
		}
	}
}

func validateInput(c config.ConditionConfig) {
	if c.ExecutionTimeInSeconds < c.RampUpTimeInSeconds {
		log.Fatalf("Execution time less than the Ramp Up time : %d < %d", c.ExecutionTimeInSeconds, c.RampUpTimeInSeconds)
	}
}

func rampUpThreads(c config.ConditionConfig, ch chan<- int64) {
	waitTime := c.RampUpTimeInSeconds / c.Threads
	var i int64
	for i = 0; i < c.Threads; i++ {
		// log.Printf("Thread count : %d", i)
		ch <- i
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func httpTest(h config.HTTPRequest, httpClient utils.HTTPClient, currentMetricsCh chan<- models.CurrentMetrics) {
	var currentMetric models.CurrentMetrics
	payloadBytes, err := json.Marshal(h.PayLoad)
	if err != nil {
		log.Print("Unable to convert payload to bytes : ", err.Error())
	}
	payload := bytes.NewReader(payloadBytes)
	response, responseTimeInMillieconds, err := httpClient.SendRequest(h.Method, h.URLWithEndpoint, h.QueryParams, h.Headers, payload)
	currentMetric.ResponseTime = responseTimeInMillieconds
	if err != nil {
		log.Print("Error making a HTTP request : ", err.Error())
	}
	val, ok := h.SuccessStatusCodes[response.StatusCode]
	if val == true && ok {
		// log.Print("Success : ", response.Body)
		currentMetric.Error = false
		currentMetricsCh <- currentMetric
		return
	}
	currentMetric.Error = true
	currentMetricsCh <- currentMetric
	return
}

func (i *PerformService) processCurrentMetric(c models.CurrentMetrics) {
	i.Metrics.TotalRequests++
	i.Metrics.AverageResponseTime = (i.Metrics.AverageResponseTime*(i.Metrics.TotalRequests-1) + c.ResponseTime) / i.Metrics.TotalRequests
	if c.ResponseTime > i.Metrics.PeakResponseTime {
		i.Metrics.PeakResponseTime = c.ResponseTime
	}
	if c.Error == true {
		i.Metrics.ErrorCount++
	}
}

// CreateTable to render a result table
func (i *PerformService) CreateTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Avg Response Time", "Peak Response Time", "Total Requests", "Error Count"})
	table.SetBorder(true)
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgBlackColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.BgRedColor})
	table.Append([]string{fmt.Sprintf("%d ms", i.Metrics.AverageResponseTime), fmt.Sprintf("%d ms", i.Metrics.PeakResponseTime), fmt.Sprintf("%d", i.Metrics.TotalRequests), fmt.Sprintf("%d", i.Metrics.ErrorCount)})
	table.Render()
	return
}
