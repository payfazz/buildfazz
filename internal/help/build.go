package help

type Build struct {}

func (b *Build) GenerateHelp() string {
	return `
Usage: buildfazz build [OPTIONS] {docker-name}:[docker-tag]

Options:	
	-p		Set buildfazz working directory
`
}

func NewBuildHelp() HelperInterface{
	return &Build{}
}
