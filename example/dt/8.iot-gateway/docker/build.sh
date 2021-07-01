registry="10.0.3.20:5005"
imagename="keti-iotgateway"
version="v1.0"
docker_id="openmcp"

# make image
#docker build -t $imagename:$version .

# add tag
#docker tag $imagename:$version $registry/$imagename:$version

# login
#docker login

# push image
#docker push $registry/$imagename:$version

#docker build -t $docker_registry_ip/$docker_id/$controller_name:v0.0.1 build && \
#docker push $docker_registry_ip/$docker_id/$controller_name:v0.0.1

docker build -t $registry/$docker_id/$imagename:$version . && \
docker push $registry/$docker_id/$imagename:$version




