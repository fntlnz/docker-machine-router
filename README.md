# Docker Machine Router

This tool allows you to reach the container's internal ip address from the host by just using the `-p` option available in the docker `run` command.

## How does it work

The [Docker documentation](https://docs.docker.com/engine/reference/run/#expose-incoming-ports) states:

```
-p=[]      : Publish a container᾿s port or a range of ports to the host
               format: ip:hostPort:containerPort | ip::containerPort | hostPort:containerPort | containerPort
               Both hostPort and containerPort can be specified as a
               range of ports. When specifying ranges for both, the
               number of container ports in the range must match the
               number of host ports in the range, for example:
                   -p 1234-1236:1234-1236/tcp
```

Unfortunately when using docker machine we can't directly use the format `ip:hostPort:containerPort` because the port is not being opened
on the host machine but on the Boot2Docker VM.

If for example we start an NGINX container like this:

```
docker run -p 10.0.0.40:80:80 nginx:latest
```

The IP will be opened inside the VM therefore we will not able to reach the http server from the host.

`docker-machine-router` has the task of waiting for the creation of new containers and then routes the container IP to the Docker Machine.

## Usage
```
sudo -E docker-machine-router
```
