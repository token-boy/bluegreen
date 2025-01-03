docker run --detach \
  --name=$project_name-$tag \
  --label=traefik.enable=true \
  --label=traefik.http.routers.$project_name.rule=Host\(\`$project_name.mxsyx.site\`\) \
  --label=traefik.http.routers.$project_name.tls.certResolver=letsencrypt \
  --label=traefik.http.services.$project_name.loadbalancer.server.port=8000 \
  --label=traefik.http.services.$project_name.loadbalancer.sticky.cookie.name=sticky \
  --label=traefik.http.services.$project_name.loadbalancer.healthcheck.path=/healthcheck \
  --label=traefik.http.services.$project_name.loadbalancer.healthcheck.interval=30s \
  --network=main \
  ghcr.io/$repository:$tag
