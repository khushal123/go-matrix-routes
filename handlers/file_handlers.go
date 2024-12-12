package handlers

import (
	"assignment/processor"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

func HandleMatrix(w http.ResponseWriter, r *http.Request, operation string) {
	// Validate HTTP method
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate content type
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := processor.ValidateFileUpload(file, header); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read CSV with custom reader settings
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields for better error handling
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing CSV: %s", err.Error()), http.StatusBadRequest)
		return
	}

	processor, err := processor.NewMatrixProcessor(records)
	if err != nil {
		http.Error(w, "Internal server error processing matrix", http.StatusInternalServerError)
	}

	var result string
	switch operation {
	case "echo":
		result = processor.Echo()
	case "invert":
		result = processor.Invert()
	case "flatten":
		result = processor.Flatten()
	case "sum":
		result = fmt.Sprintf("%d", processor.Sum())
	case "multiply":
		// Check for potential integer overflow
		product := processor.Multiply()
		if product == 0 {
			w.Header().Set("Warning", "Possible integer overflow detected")
		}
		result = fmt.Sprintf("%d", product)
	default:
		http.Error(w, "Invalid operation", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
}
