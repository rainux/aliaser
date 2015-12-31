# aliaser
Create git style alias for any command line utility.

## Usage

Create a configuration file `~/.aliaser` with the following contents:

```
[core]
    d   = docker
    dc  = docker-compose
    dm  = docker-machine

[docker]
    i       = images
    rma     = rm `docker ps -a -q`
    rmdi    = rmi `docker images -f "dangling=true" -q`
    rme     = rm `docker ps -a -f "status=exited" -q`

[docker-compose]
    rma     = kill
    rma     = rm
```

Run `aliaser install` will create symbolic links `d`, `dc` and `dm` to
`aliser`, now you can call them instead of the corresponding original commands,
with the bonus ability to use aliases configured in the above sections.
