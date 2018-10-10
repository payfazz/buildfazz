package help

type Basic struct {}

func (b *Basic) GenerateHelp() string {
	return `
Usage: buildfazz COMMAND [OPTIONS] {docker-name}:[docker-tag]

Commands:
	build		Build docker image
`
}

func NewBasicHelp() HelperInterface{
	return &Basic{}
}
