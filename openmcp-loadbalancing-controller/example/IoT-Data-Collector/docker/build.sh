registry="10.0.6.230:5000"
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


docker build -t $docker_id/$imagename:$version .
docker push $docker_id/$imagename:$version




