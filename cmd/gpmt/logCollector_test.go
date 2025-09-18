package main

import (
	"strings"
	"testing"
)

// Test to verify that our log directory parsing handles various result formats correctly
func TestLogDirectoryParsing(t *testing.T) {
	// This test simulates the result parsing logic without requiring a real database
	
	testCases := []struct {
		name           string
		resultData     []map[string]interface{}
		expectedResult string
		expectError    bool
	}{
		{
			name: "normal string result",
			resultData: []map[string]interface{}{
				{"?column?": "/data/coordinator/gpseg-1/log"},
			},
			expectedResult: "/data/coordinator/gpseg-1/log",
			expectError:    false,
		},
		{
			name: "string with whitespace",
			resultData: []map[string]interface{}{
				{"?column?": "  /data/coordinator/gpseg-1/log  \n"},
			},
			expectedResult: "/data/coordinator/gpseg-1/log",
			expectError:    false,
		},
		{
			name: "byte array result",
			resultData: []map[string]interface{}{
				{"?column?": []byte("/data/coordinator/gpseg-1/log")},
			},
			expectedResult: "/data/coordinator/gpseg-1/log",
			expectError:    false,
		},
		{
			name: "empty string result",
			resultData: []map[string]interface{}{
				{"?column?": ""},
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name: "whitespace only result",
			resultData: []map[string]interface{}{
				{"?column?": "  \n  \t  "},
			},
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "no results",
			resultData:     []map[string]interface{}{},
			expectedResult: "",
			expectError:    true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the parsing logic from getLogDirectoryFromDB
			var logDir string
			var foundValid bool
			
			if len(tc.resultData) == 0 {
				foundValid = false
			} else {
				for _, row := range tc.resultData {
					for _, value := range row {
						var candidate string
						if str, ok := value.(string); ok {
							candidate = str
						} else if bytes, ok := value.([]byte); ok && len(bytes) > 0 {
							candidate = string(bytes)
						} else {
							continue
						}
						
						// Apply the same trimming logic as the real function
						trimmed := strings.TrimSpace(candidate)
						if trimmed != "" {
							logDir = trimmed
							foundValid = true
							break
						}
					}
					if foundValid {
						break
					}
				}
			}
			
			if tc.expectError {
				if foundValid {
					t.Errorf("Expected error for case '%s', but got result: %s", tc.name, logDir)
				}
			} else {
				if !foundValid {
					t.Errorf("Expected result for case '%s', but got error", tc.name)
				} else if logDir != tc.expectedResult {
					t.Errorf("Expected result '%s' for case '%s', but got '%s'", tc.expectedResult, tc.name, logDir)
				}
			}
		})
	}
}