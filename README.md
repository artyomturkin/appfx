# appfx
Opinionated application configuration with Uber.Fx

# Configuration

Use `github.com/uber-go/config` for config provisioning. Example:

```go
type config struct {
	Name string
}

func setupHandler(config config.Provider) (http.Handler, error) {
	var cfg config
	if err := config.Get("section").Populate(&cfg); err != nil {
		return nil, err
	}
..................
}
