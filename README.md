# Docker Machine Router

This tool allows you to **assign an chosen ip address** that is **visible from the OS X host** to your containers.

## Installation

```
curl -L https://github.com/fntlnz/docker-machine-router/releases/download/v0.2.1/docker-machine-router > /usr/local/bin/docker-machine-router
chmod +x /usr/local/bin/docker-machine-router
```

## Usage

Docker machine router relies on the `DOCKER_HOST` environment variable as the Docker client does.

You can export that variable using the `docker-machine env <machine-name>` command.

```
MACHINE=dev
eval $(docker-machine env $MACHINE)
```

Now that all's ready you can create the routes by invoking `docker-machine-router`.
The `-E` option, tells `sudo` to get all the environment variables.

```
sudo -E docker-machine-router
```

#### Options

**cidr**:

This is the [Class Inter-Domain Routing](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) that `docker-machine-router` will use to create its routes.

The default cidr is `10.18.0.0/16` you can change it using the `-cidr` option:

```
sudo -E ./dist/docker-machine-router -cidr="10.20.0.0/16"
```

The default cidr: `10.18.0.0/16` allows allocation of ip addresses in the `10.18.0.0 - 10.18.255.255` range, for a total of 65536 addresses.

### Start a container using the `dmr` network and assigning a custom IP.

Here we are starting a container using the `nginx` image and the `dmr` network (The one provided by `docker-machine-router`).

Our container will have two ip addresses, one exposed to the OS X host `10.18.241.100` and one exposed only inside the VM `127.0.241.100`

Note that the ip exposed only in the VM  (`127.0.241.100`)is needed in order to allocate multiple `80` ports on the docker network.

```
docker run --net dmr --ip 10.18.241.100 -p 127.0.241.100:80:80 -it nginx
```

### The Same with docker-compose

**docker-compose.yml**

```
version: "2"
services:
  web:
    image: nginx
    ports:
      - "127.18.20.20:80:80"
    networks:
      dmr:
        ipv4_address: "10.18.20.20"
networks:
  dmr:
    external: true
```

This is the equivalent of what we did with the run command.
