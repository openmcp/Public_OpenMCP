kubectl create secret generic "regcred" \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=keti --context cluster01

kubectl create secret generic "regcred" \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=keti --context cluster02

kubectl create secret generic "regcred" \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=keti --context cluster03

kubectl create secret generic "regcred" \
    --from-file=.dockerconfigjson=/root/.docker/config.json \
    --type=kubernetes.io/dockerconfigjson \
    --namespace=keti --context cluster04
