

omcpctl delete -n openmcp job sns-1605230624-job -c cluster1
omcpctl delete -n openmcp pvc sns-1605230624-epvc -c cluster1
omcpctl delete -n openmcp pvc sns-1605230624-pvc -c cluster1
omcpctl delete pv sns-1605230624-pv -c cluster1
omcpctl delete pv sns-1605230624-epv -c cluster1
omcpctl delete deploy snapshot-test-dp -c cluster1
omcpctl delete pvc snapshot-test-pvc -c cluster1
omcpctl delete pv snapshot-test-pv -c cluster1



omcpctl delete -n openmcp job snr-1605230624-job -c cluster1
omcpctl delete -n openmcp pvc snr-1605230624-epvc -c cluster1
omcpctl delete -n openmcp pvc snr-1605230624-pvc -c cluster1
omcpctl delete pv snr-1605230624-pv -c cluster1
omcpctl delete pv snr-1605230624-epv -c cluster1


 kubectl delete -n openmcp job sns-1605230624-job
 kubectl delete -n openmcp pvc sns-1605230624-epvc
 kubectl delete -n openmcp pvc sns-1605230624-pvc
 kubectl delete pv sns-1605230624-pv
 kubectl delete pv sns-1605230624-epv
 kubectl delete deploy snapshot-test-dp
 kubectl delete pvc snapshot-test-pvc
 kubectl delete pv snapshot-test-pv



 kubectl delete -n openmcp job snr-1605230624-job
 kubectl delete -n openmcp pvc snr-1605230624-epvc
 kubectl delete -n openmcp pvc snr-1605230624-pvc
 kubectl delete pv snr-1605230624-pv
 kubectl delete pv snr-1605230624-epv

