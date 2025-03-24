package config

var Headers []string

type CustomHeaders []string

func (h *CustomHeaders) String() string {
	return "Custom headers"
}

func (h *CustomHeaders) Set(value string) error {
	*h = append(*h, value)
	return nil
}