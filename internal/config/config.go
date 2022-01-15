package config

type Package struct {
	name    string
	configs *[]PackageConfig
}

type PackageConfig interface {
}

type FileConfig interface {
	PackageConfig
	getFile() string
	config(c *Prompt)
}

type Prompt struct{}

func (p *Prompt) Ask(message string) (string, error) {
	return "nope", nil
}
