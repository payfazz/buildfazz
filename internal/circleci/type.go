package circleci

type Data struct {
	Version int `yaml:"version"`
	Jobs    Job `yaml:"jobs"`
}

type Job struct {
	Build struct {
		Dockers          []Docker `yaml:"docker"`
		WorkingDirectory string   `yaml:"working_directory"`
		//Steps            []Step   `yaml:"steps"`
	} `yaml:"build"`
}

type Docker struct {
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
}

//type Step struct {
//	Run struct {
//		Command     string            `yaml:"command"`
//		Environment map[string]string `yaml:"environment"`
//	} `yaml:"run"`
//}
