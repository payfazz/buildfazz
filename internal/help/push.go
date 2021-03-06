package help

// Push ...
type Push struct {}

// GenerateHelp ...
func (b *Push) GenerateHelp() string {
	return `
Usage: 		buildfazz push --ssh {user@server} [OPTIONS] {docker-name}:[docker-tag]
Example: 	buildfazz push --ssh root@127.0.0.1 myImage:latest

Options:	
	-e		Set environment for your workstation (put 'mac' for mac user, don't forget to see docker mac doc), default: none
	-t		Target server (example: -t localhost), default: localhost
	-p		Target port (example: -p 9080), default: 5000
`
}

// NewPushHelp ...
func NewPushHelp() HelperInterface{
	return &Push{}
}
