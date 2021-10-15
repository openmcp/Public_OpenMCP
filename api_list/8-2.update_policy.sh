USERNAME="openmcp"
PASSWORD="keti"

IP="openmcp-apiserver.openmcp.default-domain.svc.openmcp.example.org"
PORT="8080"

POLICYNAME="post-scheduling-type"

URL="apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcppolicys/$POLICYNAME"
CONTEXT="openmcp"

echo -n | openssl s_client -connect $IP:$PORT | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > server.crt

TOKEN_JSON=`curl -XPOST \
        --cacert server.crt \
        --insecure \
        -H "Content-type: application/json" \
        --data "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
        https://$IP:$PORT/token`

TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

#policy name : log-level
#LOGVALUE="4" # value : 0~5
#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$LOGVALUE\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


#policy name : hpa-minmax-distribution-mode
#HPAMODE="Unequal" # value : Equal or Unequal
#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$HPAMODE\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


#policy name : metric-collector-period
#PERIOD="6" # value : seconds
#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$PERIOD\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


#policy name : analytic-metrics-weight
#CPU Index : 0 / Memory Index : 1 / FS Index : 2 / NET Index : 3 / LATENCY Index : 4
#CPUWEIGHT="0.2"
#MEMWEIGHT="0.3"
#FSWEIGHT="0.4"
#NETWEIGHT="0.8"
#LATENCYWEIGHT="0.5"
#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$CPUWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/1/value/0\", \"value\": \"$MEMWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/2/value/0\", \"value\": \"$FSWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/3/value/0\", \"value\": \"$NETWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/4/value/0\", \"value\": \"$LATENCYWEIGHT\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT


#policy name : lb-scoring-weight
#GeoRate Index : 0 / Period Index : 1 / RegionZoneMatchedScore Index : 2 / OnlyRegionMatchedScore Index : 3 / NoRegionZoneMatchedScore Index : 4
#GEOWEIGHT="0.5"
#PERIODWEIGHT="5.0"
#RZWEIGHT="100"
#ONLYRWEIGHT="60"
#NORZWEIGHT="20"
#curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
#--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$GEOWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/1/value/0\", \"value\": \"$PERIODWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/2/value/0\", \"value\": \"$RZWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/3/value/0\", \"value\": \"$ONLYRWEIGHT\"},{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/4/value/0\", \"value\": \"$NORZWEIGHT\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT

#policy name : post-scheduling-type
SCHEDULINGTYPE="Scoring"  # value : FIFO or Scoring
curl -X PATCH --cacert server.crt -H "Content-Type: application/json-patch+json" -H "Authorization: Bearer $TOKEN" \
--data "[{\"op\": \"replace\", \"path\": \"/spec/template/spec/policies/0/value/0\", \"value\": \"$SCHEDULINGTYPE\"}]" https://$IP:$PORT/$URL?clustername=$CONTEXT



rm server.crt