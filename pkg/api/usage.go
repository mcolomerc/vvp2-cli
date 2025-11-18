package api

import (
	"encoding/csv"
	"fmt"
	"strings"
)

type ResourceUsageReport struct {
	CSVData string
}

// ParseCSV parses the CSV data into a slice of maps
func (r *ResourceUsageReport) ParseCSV() ([]map[string]string, error) {
	lines := strings.Split(r.CSVData, "\n")
	var csvLines []string
	
	// Filter out comment lines (starting with #) and empty lines
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			csvLines = append(csvLines, line)
		}
	}
	
	if len(csvLines) == 0 {
		return []map[string]string{}, nil
	}
	
	csvReader := csv.NewReader(strings.NewReader(strings.Join(csvLines, "\n")))
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}
	
	if len(records) < 2 {
		return []map[string]string{}, nil
	}
	
	headers := records[0]
	var result []map[string]string
	
	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		result = append(result, row)
	}
	
	return result, nil
}

func (c *Client) GetResourceUsageReport(from, to string) (*ResourceUsageReport, error) {
       req := c.httpClient.R()
       if from != "" {
              req.SetQueryParam("from", from)
       }
       if to != "" {
              req.SetQueryParam("to", to)
       }
       
       endpoint := "/api/v1/status/resourceusage"
       resp, err := req.Get(endpoint)

       if resp != nil && resp.StatusCode() == 404 {
              return nil, fmt.Errorf("Resource usage report endpoint is not enabled on this platform (404)")
       }
       if err := handleResponse(resp, err); err != nil {
              return nil, err
       }
       
       return &ResourceUsageReport{CSVData: resp.String()}, nil
}
