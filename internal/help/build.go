package help

// Build ...
type Build struct {}

// GenerateHelp ...
func (b *Build) GenerateHelp() string {
	return `
Usage: buildfazz build [OPTIONS] {docker-name}:[docker-tag]

Options:	
	-p		Set buildfazz working directory, default: current directory
	-os		Set buildfazz default OS (options: debian/ubuntu/scratch/etc...), default: debian
`
}

// NewBuildHelp ...
func NewBuildHelp() HelperInterface{
	return &Build{}
}
