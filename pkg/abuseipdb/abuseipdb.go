package abuseipdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Report struct {
	Data struct {
		Hostnames            []interface{} `json:"hostnames"`
		IPAddress            string        `json:"ipAddress"`
		CountryCode          string        `json:"countryCode"`
		UsageType            string        `json:"usageType"`
		Isp                  string        `json:"isp"`
		Domain               string        `json:"domain"`
		IPVersion            int           `json:"ipVersion"`
		AbuseConfidenceScore int           `json:"abuseConfidenceScore"`
		TotalReports         int           `json:"totalReports"`
		NumDistinctUsers     int           `json:"numDistinctUsers"`
		LastReportedAt       *time.Time    `json:"lastReportedAt"`
		IsWhitelisted        bool          `json:"isWhitelisted"`
		IsPublic             bool          `json:"isPublic"`
	} `json:"data"`
}

type Client struct {
	Client  *http.Client
	BaseURL string
	Key     string
}

func NewClient(key string) (*Client, error) {
	return &Client{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: "https://api.abuseipdb.com/api/v2",
		Key:     key,
	}, nil
}

func (c *Client) Check(ip string) (Report, error) {
	if net.ParseIP(ip) == nil {
		return Report{}, fmt.Errorf("argument must be a valid IP address")
	}

	uri := fmt.Sprintf("%s/check", c.BaseURL)
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
		return Report{}, fmt.Errorf("failure from AbuseIPDB API: %s", string(body))
	}

	var report Report
	if err := json.Unmarshal(body, &report); err != nil {
		return Report{}, err
	}

	return report, nil
}
