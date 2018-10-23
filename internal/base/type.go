package base

// Data ...
type Data struct {
	ProjectName string `yaml:"project"`
	Base        string `yaml:"base"`
	Main        string `yaml:"main"`
	Version     string `yaml:"version"`
	Pwd         string
}
