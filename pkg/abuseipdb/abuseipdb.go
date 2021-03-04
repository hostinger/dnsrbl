package abuseipdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var baseURL = "https://api.abuseipdb.com/api/v2"

type Report struct {
	Data struct {
		IPAddress            string        `json:"ipAddress"`
		IsPublic             bool          `json:"isPublic"`
		IPVersion            int           `json:"ipVersion"`
		IsWhitelisted        bool          `json:"isWhitelisted"`
		AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
		CountryCode          string        `json:"countryCode"`
		UsageType            string        `json:"usageType"`
		Isp                  string        `json:"isp"`
		Domain               string        `json:"domain"`
		Hostnames            []interface{} `json:"hostnames"`
		TotalReports         int           `json:"totalReports"`
		NumDistinctUsers     int           `json:"numDistinctUsers"`
		LastReportedAt       *time.Time    `json:"lastReportedAt"`
	} `json:"data"`
}

type Client struct {
	Client *http.Client
	Key    string
}

func NewClient(key string) (*Client, error) {
	return &Client{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Key: key,
	}, nil
}

func (c *Client) Check(ip string) (Report, error) {
	if net.ParseIP(ip) == nil {
		return Report{}, fmt.Errorf("Argument must be a valid IP address.")
	}

	uri := fmt.Sprintf("%s/check", baseURL)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return Report{}, err
	}
	q := req.URL.Query()
	q.Set("ipAddress", ip)

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Key", c.Key)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return Report{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Report{}, fmt.Errorf("Failure from AbuseIPDB API: %s", body)
	}

	var report Report
	if err := json.Unmarshal(body, &report); err != nil {
		return Report{}, err
	}

	return report, nil
}
