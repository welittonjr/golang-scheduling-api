#!/bin/bash

failed_containers=$(docker ps -a -f status=exited --format "{{.ID}}")

image_ids=$(docker images --filter "dangling=true" -q)

if [ -z "$failed_containers" ]; then
    echo "Não foram encontrados containers que falharam no build."
fi

for container_id in $failed_containers; do
    echo "Removendo container: $container_id"
    docker rm $container_id
done

if [[ -z "$image_ids" ]]; then
  echo "Não foram encontradas imagens sem tags."
else
  docker rmi $image_ids
  echo "Imagens sem tags removidas com sucesso!"
fi