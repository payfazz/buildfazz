**BuildFazz**<br />
<br />
_BuildFazz_ is the docker image builder.<br />
<br />
How to install : <br />
`go get github.com/payfazz/buildfazz/cmd/buildfazz`<br />
<br />
How to use : <br />
`buildfazz COMMAND [OPTIONS] {project-name}:"[project-tag]`<br />

See the help in:<br />
`buildfazz --help`<br /><br />

Command exist :
- build<br />

Use `buildfazz COMMAND --help` to see detail use for each command.

Options exist for build:<br />
- -p         Specify path for project path
<br />

You also need to specify the configuration file in `buildfazz.yml` <br />
See an example for `buildfazz.yml` in [here](https://github.com/payfazz/buildfazz/blob/master/buildfazz.yml).
