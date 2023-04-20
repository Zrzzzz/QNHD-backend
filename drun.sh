docker stop qnhd_run
docker rm qnhd_run
docker rmi qnhd
docker build . -t qnhd && 	\
	cd DockerRuntime && 	\
	docker run -d -p 7013:7013 --name=qnhd_run -v `pwd`/DockerRuntime/runtime:/qnhd/runtime qnhd &&\
	docker logs -f qnhd_run
