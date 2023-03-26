docker stop qnhd_product
docker rm qnhd_product
docker rmi qnhd-go:latest
docker build -t qnhd-go . && \
docker run -p 7013:7013 -v `pwd`/pages:/app/pages -v `pwd`/conf:/app/conf -v `pwd`/runtime:/app/runtime -v /etc/localtime:/etc/localtime:ro  --name=qnhd_product -d qnhd-go && \
docker logs -f qnhd_product
