package help

// Basic ...
type Basic struct {}

// GenerateHelp ...
func (b *Basic) GenerateHelp() string {
	return `
Usage: buildfazz COMMAND [OPTIONS] {docker-name}:[docker-tag]

Commands:
	build		Build docker image
	push		Push docker image to registry server
`
}

// NewBasicHelp ...
func NewBasicHelp() HelperInterface{
	return &Basic{}
}
