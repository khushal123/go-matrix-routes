# Matrix Operations Service

A Go web service that performs various operations on square matrices provided
via CSV files.

## Features

The service provides five endpoints, each accepting a CSV file containing a
square matrix:

1. `/echo` - Returns the matrix in its original format
   ```
   Input: 1,2,3\n4,5,6\n7,8,9
   Output: 1,2,3
           4,5,6
           7,8,9
   ```

2. `/invert` - Returns the matrix with rows and columns inverted
   ```
   Input: 1,2,3\n4,5,6\n7,8,9
   Output: 1,4,7
           2,5,8
           3,6,9
   ```

3. `/flatten` - Returns the matrix as a single line of comma-separated values
   ```
   Input: 1,2,3\n4,5,6\n7,8,9
   Output: 1,2,3,4,5,6,7,8,9
   ```

4. `/sum` - Returns the sum of all integers in the matrix
   ```
   Input: 1,2,3\n4,5,6\n7,8,9
   Output: 45
   ```

5. `/multiply` - Returns the product of all integers in the matrix
   ```
   Input: 1,2,3\n4,5,6\n7,8,9
   Output: 362880
   ```

## Requirements

- Go 1.x or higher

## Installation

1. Clone the repository
2. Navigate to the project directory

## Running the Service

Start the server:

```bash
go run .
```

The server will start on port 8000.

## Making Requests

Use curl to send requests with a CSV file:

```bash
curl -F 'file=@matrix.csv' "localhost:8000/echo"
curl -F 'file=@matrix.csv' "localhost:8000/invert"
curl -F 'file=@matrix.csv' "localhost:8000/flatten"
curl -F 'file=@matrix.csv' "localhost:8000/sum"
curl -F 'file=@matrix.csv' "localhost:8000/multiply"
```

### Input Format Requirements

- File must be in CSV format
- Matrix must be square (equal number of rows and columns)
- All values must be integers
- No header row
- Non-empty cells

Example valid input (matrix.csv):

```
1,2,3
4,5,6
7,8,9
```

## Error Handling

The service includes robust error handling for:

- Invalid HTTP methods (only POST allowed)
- Missing or invalid files
- Non-square matrices
- Non-integer values
- Empty matrices
- Wrong file types
- Server panics (with recovery)

## Running Tests

Run all tests:

```bash
go test
```

Run tests with verbose output:

```bash
go test -v
```

## Implementation Details

- Written in Go
- Uses standard library http package
- Includes panic recovery middleware
- Comprehensive test suite with both success and failure cases
- Concurrent request handling
