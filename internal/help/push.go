package help

type Push struct {}

func (b *Build) GeneratePush() string {
	return `
Usage: buildfazz push [OPTIONS] {docker-name}:[docker-tag]

Options:	
	-e		Set environment for your workstation (put 'mac' for mac user)
`
}

func NewPushHelp() HelperInterface{
	return &Build{}
}
