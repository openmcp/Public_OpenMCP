docker_id="openmcp"
imagename="keti-http-generator"
version="v1.1"

# make image
docker build -t $docker_id/$imagename:$version . && \

# push image
docker push $docker_id/$imagename:$version
