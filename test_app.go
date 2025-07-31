package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	// Wait for server to start
	time.Sleep(2 * time.Second)
	
	fmt.Println("Testing HTML Sharer Application...")
	fmt.Println("==================================")
	
	// Test 1: Home page loads
	fmt.Println("Test 1: Testing home page...")
	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		fmt.Printf("❌ Home page test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if strings.Contains(bodyStr, "HTML Sharer") && 
		   strings.Contains(bodyStr, "textarea") && 
		   strings.Contains(bodyStr, "input type=\"file\"") &&
		   strings.Contains(bodyStr, "Create Link") {
			fmt.Println("✅ Home page loads correctly with all required elements")
		} else {
			fmt.Println("❌ Home page missing required elements")
		}
	} else {
		fmt.Printf("❌ Home page returned status %d\n", resp.StatusCode)
	}
	
	// Test 2: API endpoint for sharing HTML
	fmt.Println("\nTest 2: Testing HTML sharing API...")
	testHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Hello World!</h1>
    <p>This is a test HTML page.</p>
</body>
</html>`
	
	shareReq := map[string]string{
		"html_content": testHTML,
	}
	
	jsonData, _ := json.Marshal(shareReq)
	resp, err = http.Post("http://localhost:8080/api/share", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Share API test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	var shareResp map[string]string
	json.NewDecoder(resp.Body).Decode(&shareResp)
	
	if resp.StatusCode == 200 && shareResp["url"] != "" {
		fmt.Printf("✅ HTML sharing successful, got URL: %s\n", shareResp["url"])
		
		// Test 3: Access shared content
		fmt.Println("\nTest 3: Testing shared content access...")
		sharedURL := "http://localhost:8080" + shareResp["url"]
		resp, err = http.Get(sharedURL)
		if err != nil {
			fmt.Printf("❌ Shared content access failed: %v\n", err)
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			bodyStr := string(body)
			if strings.Contains(bodyStr, "Hello World!") && strings.Contains(bodyStr, "Test Page") {
				fmt.Println("✅ Shared content serves correctly")
			} else {
				fmt.Println("❌ Shared content doesn't match original HTML")
			}
		} else {
			fmt.Printf("❌ Shared content returned status %d\n", resp.StatusCode)
		}
	} else {
		fmt.Printf("❌ Share API failed with status %d\n", resp.StatusCode)
	}
	
	// Test 4: 404 for non-existent slug
	fmt.Println("\nTest 4: Testing 404 for non-existent content...")
	resp, err = http.Get("http://localhost:8080/shared/nonexistent")
	if err != nil {
		fmt.Printf("❌ 404 test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 404 {
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		if strings.Contains(bodyStr, "404") && strings.Contains(bodyStr, "Not Found") {
			fmt.Println("✅ 404 page works correctly")
		} else {
			fmt.Println("❌ 404 page content incorrect")
		}
	} else {
		fmt.Printf("❌ Expected 404, got status %d\n", resp.StatusCode)
	}
	
	// Test 5: Empty content validation
	fmt.Println("\nTest 5: Testing empty content validation...")
	emptyReq := map[string]string{
		"html_content": "",
	}
	
	jsonData, _ = json.Marshal(emptyReq)
	resp, err = http.Post("http://localhost:8080/api/share", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Empty content test failed: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 400 {
		fmt.Println("✅ Empty content validation works correctly")
	} else {
		fmt.Printf("❌ Expected 400 for empty content, got status %d\n", resp.StatusCode)
	}
	
	fmt.Println("\n==================================")
	fmt.Println("Testing completed!")
}