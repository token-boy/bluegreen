# bluegreen

> 零停机部署

[Docker Swarm](https://docs.docker.com/engine/swarm/)

[traefikdockerlbswarm](https://doc.traefik.io/traefik/routing/providers/docker/#traefikdockerlbswarm)


```shell
docker swarm init
```

```shell
docker stack deploy -c /etc/traefik/deploy.yml traefik
```
