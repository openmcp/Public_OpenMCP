
kubectl label nodes kube1-master topology.kubernetes.io/region=KR --context=cluster1 --overwrite
kubectl label nodes kube1-master topology.kubernetes.io/zone=Seoul --context=cluster1 --overwrite

kubectl label nodes kube1-worker topology.kubernetes.io/region=KR --context=cluster1 --overwrite
kubectl label nodes kube1-worker topology.kubernetes.io/zone=Seoul --context=cluster1 --overwrite

kubectl label nodes kube2-master topology.kubernetes.io/region=KR --context=cluster2 --overwrite
kubectl label nodes kube2-master topology.kubernetes.io/zone=Seoul --context=cluster2 --overwrite

kubectl label nodes kube2-worker1 topology.kubernetes.io/region=KR --context=cluster2 --overwrite
kubectl label nodes kube2-worker1 topology.kubernetes.io/zone=Gyeonggi-do --context=cluster2 --overwrite

