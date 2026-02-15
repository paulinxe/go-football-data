# Go Football Data
This is a small golang package that allows you to fetch matches and competitions from the Football Data API.
It uses the V4 version of the API.

## Installation

```bash
go get github.com/paulinxe/go-football-data
```

## Prerequisites

You'll need an API key from [Football-Data.org](https://www.football-data.org). You can register for a free tier account to get started.

## Usage

### Creating a Client

```go
import footballdata "github.com/paulinxe/go-football-data"

// Create a client with your API key and using the v4 endpoint
client := footballdata.New("YOUR_API_KEY")
```

If you want to use a custom endpoint, let's say for testing, you can:
```go
server := httptest.NewServer(
    //...
)
client := New("YOUR_API_KEY", WithBaseURL(server.URL))
```

### Fetching Competitions

Check `examples/competitions/main.go`

### Fetching Matches

Check `examples/matches/main.go`

For more information about the API, visit the [Football-Data.org API documentation](https://www.football-data.org/documentation/api).

## Error Handling

The package provides custom error types for better error handling:

Check `examples/error_handling/main.go`

### Rate Limits

Be aware of the API rate limits based on your subscription tier.

The client uses an `HTTPError` struct to provide detailed error information, including rate limiting. Here's how to handle it:

```go
func main() {
    matches, err := client.GetMatches(context.Background(), competitionId, filters, &mapTo)
    if err != nil {
        var httpErr *footballdata.HTTPError
        if errors.As(err, &httpErr) {
            switch httpErr.StatusCode {
            case 429:
                // Here goes your logic
        }
    }
}
```

## API Coverage

This SDK currently supports the Football-Data.org API V4 endpoints for:

- Competitions
- Matches

## License

This project is open source. Please check the repository for license information.

## Resources

- [Football-Data.org](https://www.football-data.org) - Official API website
- [API Documentation](https://www.football-data.org/documentation/api) - Detailed API documentation
- [Go Documentation](https://pkg.go.dev/github.com/paulinxe/go-football-data) - Package documentation