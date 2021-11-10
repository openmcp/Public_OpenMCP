

```
#delete openmcp snapshot resources.
kubectl delete pv,pvc,job  --selector=openmcp=snapshot -n openmcp  --context=cluster1



```