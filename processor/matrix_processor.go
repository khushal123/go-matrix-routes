package processor

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
)

// Custom error types for better error handling
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// validateFileUpload performs comprehensive file validation
func ValidateFileUpload(file multipart.File, header *multipart.FileHeader) error {
	if header == nil {
		return &ValidationError{"No file was uploaded"}
	}

	// Check file size (e.g., max 5MB)
	if header.Size > 5*1024*1024 {
		return &ValidationError{"File size exceeds maximum limit of 5MB"}
	}

	// Check file extension
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		return &ValidationError{"Invalid file format. Only CSV files are accepted"}
	}

	return nil
}

type MatrixProcessor struct {
	matrix [][]int
}

// NewMatrixProcessor creates a new processor from CSV data
func NewMatrixProcessor(records [][]string) (*MatrixProcessor, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("empty matrix provided")
	}

	rows := len(records)
	cols := len(records[0])
	if rows != cols {
		return nil, fmt.Errorf("matrix must be square, got %dx%d", rows, cols)
	}

	matrix := make([][]int, rows)
	for i := range matrix {
		if len(records[i]) != cols {
			return nil, fmt.Errorf("inconsistent row length at row %d", i)
		}
		matrix[i] = make([]int, cols)
		for j, val := range records[i] {
			num, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return nil, fmt.Errorf("invalid integer at position [%d,%d]: %s", i, j, val)
			}
			matrix[i][j] = num
		}
	}

	return &MatrixProcessor{matrix: matrix}, nil
}

// Echo returns the matrix as a string in matrix format
func (mp *MatrixProcessor) Echo() string {
	var result strings.Builder
	for i, row := range mp.matrix {
		if i > 0 {
			result.WriteString("\n")
		}
		// Convert each number to string and join with commas
		nums := make([]string, len(row))
		for j, num := range row {
			nums[j] = strconv.Itoa(num)
		}
		result.WriteString(strings.Join(nums, ","))
	}
	return result.String()
}

// Invert returns the matrix with columns and rows inverted
func (mp *MatrixProcessor) Invert() string {
	size := len(mp.matrix)
	inverted := make([][]int, size)
	for i := range inverted {
		inverted[i] = make([]int, size)
		for j := range inverted[i] {
			inverted[i][j] = mp.matrix[j][i]
		}
	}

	var result strings.Builder
	for i, row := range inverted {
		if i > 0 {
			result.WriteString("\n")
		}
		nums := make([]string, len(row))
		for j, num := range row {
			nums[j] = strconv.Itoa(num)
		}
		result.WriteString(strings.Join(nums, ","))
	}
	return result.String()
}

// Flatten returns the matrix as a 1 line string
func (mp *MatrixProcessor) Flatten() string {
	var nums []string
	for _, row := range mp.matrix {
		for _, num := range row {
			nums = append(nums, strconv.Itoa(num))
		}
	}
	return strings.Join(nums, ",")
}

// Sum returns the sum of all integers in the matrix
func (mp *MatrixProcessor) Sum() int {
	sum := 0
	for _, row := range mp.matrix {
		for _, num := range row {
			sum += num
		}
	}
	return sum
}

// Multiply returns the product of all integers in the matrix
func (mp *MatrixProcessor) Multiply() int {
	product := 1
	for _, row := range mp.matrix {
		for _, num := range row {
			product *= num
		}
	}
	return product
}
