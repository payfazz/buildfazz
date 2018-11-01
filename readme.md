# BuildFazz

_BuildFazz_ is the docker image builder and pusher. Easy to use, help you create a docker image for your project without any skills required.<br />
<br />

##### Installation

    go get github.com/payfazz/buildfazz/cmd/buildfazz

##### Usage

    buildfazz COMMAND [OPTIONS] [project-name]:[project-version]

##### Help
    
    buildfazz --help
    
   
    
### buildfazz.yml
You also need to specify the configuration file in **buildfazz.yml**. <br />
See an example for `buildfazz.yml` in [here](https://github.com/payfazz/buildfazz/blob/master/buildfazz.yml).

**fields**
- project (project name)
- base (project's base [golang/node/etc])
- main (source target, example: cmd/buildfazz _or_ cmd/{your-project})
- version (project version's, could be override in CLI)

### Command & Option

Command list :
- build -> Build docker image
- push -> Push docker image to registry server
    
### Build
    
    buildfazz build project:version
    
Build option list for build:<br />
- -p        Set buildfazz working directory, default: current directory

### Push

    buildfazz push name:tag --ssh user@127.0.0.1

Use `buildfazz COMMAND --help` to see detail use for each command.

Push option list for push:<br />
- -e	    Set environment for your workstation (put 'mac' for mac user, don't forget to see docker mac doc), default: none
- -t	    Target server (example: -t localhost), default: localhost
- -p	    Target port (example: -p 9080), default: 5000

