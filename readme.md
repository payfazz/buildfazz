# BuildFazz

_BuildFazz_ is the docker image builder and pusher. Easy to use, help you create a docker image for your project without any skills required.<br />
<br />
How to **install** : <br />

    go get github.com/payfazz/buildfazz/cmd/buildfazz

How to **use** : <br />

    buildfazz COMMAND [OPTIONS] {project-name}:"[project-tag]

See the **help** in:<br />
    
    buildfazz --help
    
    

### Command & Option

**Command** list :
- build<br />

Use `buildfazz COMMAND --help` to see detail use for each command.

#### Option

**Option** list for build:<br />
- -p         Specify path for project path.

# buildfazz.yml
You also need to specify the configuration file in **buildfazz.yml**. <br />
See an example for `buildfazz.yml` in [here](https://github.com/payfazz/buildfazz/blob/master/buildfazz.yml).

**fields**
- base (project's base [golang/node/etc])
- main (source target, example: cmd/buildfazz _or_ cmd/{your-project})
