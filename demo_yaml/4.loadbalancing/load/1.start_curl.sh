function back_ground_process() {
	for ((;;))
	do
		curl http://keti.productpage.openmcp.in:8181/productpage
	done
}

cnt=10
# @ means all elements in array
for ((i=0;i<=cnt;i++)); do
    # run back_ground_process function in background
    # pass element of array as argument
    # make log file
    back_ground_process $i > ~/log_${i}.txt &
done

