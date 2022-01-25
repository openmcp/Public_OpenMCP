kubectl get namespace $1 -o json --context $2 > temp.json

sed -i -e 's/"kubernetes"//' temp.json

kubectl replace --raw "/api/v1/namespaces/$1/finalize" -f ./temp.json --context $2

rm temp.json
