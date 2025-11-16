package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ForecastRequest represents a forecasting request
type ForecastRequest struct {
	ItemName  string `json:"itemName"`
	Days      int    `json:"days"`      // Number of days to forecast
	Historical bool  `json:"historical"` // Whether to include historical data
}

// ForecastResponse represents the forecasting result
type ForecastResponse struct {
	ItemName     string           `json:"itemName"`
	Predictions  []DayPrediction  `json:"predictions"`
	Confidence   string           `json:"confidence"`
	GeneratedAt  time.Time        `json:"generatedAt"`
}

// DayPrediction represents a single day's prediction
type DayPrediction struct {
	Date     string  `json:"date"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

// GeminiRequest represents request to Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents response from Gemini API
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Handler for AI-powered demand forecasting
func forecastHandler(w http.ResponseWriter, r *http.Request) {
	var req ForecastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	
	if req.ItemName == "" {
		http.Error(w, "itemName is required", 400)
		return
	}
	
	if req.Days == 0 {
		req.Days = 7 // Default to 7 days
	}
	
	// Get historical sales data for the item
	historicalData, err := getHistoricalSales(r.Context(), req.ItemName, 30) // Last 30 days
	if err != nil {
		http.Error(w, "failed to fetch historical data: "+err.Error(), 500)
		return
	}
	
	// Call Google Gemini API for forecasting
	predictions, confidence, err := callGeminiForForecast(req.ItemName, req.Days, historicalData)
	if err != nil {
		http.Error(w, "forecasting failed: "+err.Error(), 500)
		return
	}
	
	response := ForecastResponse{
		ItemName:    req.ItemName,
		Predictions: predictions,
		Confidence:  confidence,
		GeneratedAt: time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getHistoricalSales fetches historical sales data for an item
func getHistoricalSales(ctx interface{}, itemName string, days int) (string, error) {
	type historicalRecord struct {
		Date     string
		Quantity float64
	}
	
	// Query sales from sale_items joined with transactions
	rows, err := db.Query(`
		SELECT DATE(t.date) as sale_date, SUM(si.quantity) as total_qty
		FROM sale_items si
		JOIN transactions t ON si.transaction_id = t.id
		WHERE si.item_name = $1
		AND t.date >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY DATE(t.date)
		ORDER BY sale_date ASC
	`, itemName)
	
	if err != nil {
		return "", err
	}
	defer rows.Close()
	
	records := []historicalRecord{}
	for rows.Next() {
		var r historicalRecord
		if err := rows.Scan(&r.Date, &r.Quantity); err != nil {
			continue
		}
		records = append(records, r)
	}
	
	// Format historical data as text
	if len(records) == 0 {
		return "No historical data available for this item.", nil
	}
	
	var result string
	result = fmt.Sprintf("Historical sales data for %s (last %d days):\n", itemName, days)
	for _, r := range records {
		result += fmt.Sprintf("- %s: %.2f units\n", r.Date, r.Quantity)
	}
	
	return result, nil
}

// callGeminiForForecast calls Google Gemini API to generate forecast
func callGeminiForForecast(itemName string, days int, historicalData string) ([]DayPrediction, string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// Return a stub forecast if API key is not configured
		return generateStubForecast(itemName, days), "low", nil
	}
	
	// Construct prompt for Gemini
	prompt := fmt.Sprintf(`You are a demand forecasting expert for a perishable goods vendor management system.

Item: %s
Forecast Period: Next %d days

%s

Based on the historical data above, please provide a demand forecast for the next %d days.
Respond ONLY with a JSON array of predictions in this exact format:
[
  {"date": "YYYY-MM-DD", "quantity": X.XX, "unit": "kg"},
  ...
]

Consider:
- Historical trends and patterns
- Seasonal variations
- Day of week effects
- Typical demand patterns for perishable goods

Respond with ONLY the JSON array, no other text.`, itemName, days, historicalData, days)
	
	// Prepare Gemini API request
	geminiReq := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}
	
	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, "", err
	}
	
	// Call Gemini API
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("gemini API error: %d - %s", resp.StatusCode, string(body))
	}
	
	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, "", err
	}
	
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return generateStubForecast(itemName, days), "low", nil
	}
	
	// Parse the JSON response from Gemini
	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	
	var predictions []DayPrediction
	if err := json.Unmarshal([]byte(responseText), &predictions); err != nil {
		// If parsing fails, return stub forecast
		return generateStubForecast(itemName, days), "low", nil
	}
	
	return predictions, "medium", nil
}

// generateStubForecast generates a simple stub forecast when Gemini is not available
func generateStubForecast(itemName string, days int) []DayPrediction {
	predictions := []DayPrediction{}
	baseQuantity := 10.0 // Base prediction
	
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, i+1).Format("2006-01-02")
		// Simple pattern: slightly varying quantities
		quantity := baseQuantity + float64(i%3)*2.5
		
		predictions = append(predictions, DayPrediction{
			Date:     date,
			Quantity: quantity,
			Unit:     "kg",
		})
	}
	
	return predictions
}
