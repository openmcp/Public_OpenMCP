package analyticEngine

import (
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"google.golang.org/grpc"
	"log"
	"net"
	"openmcp-analytic-engine/pkg/protobuf"
	"openmcp-analytic-engine/pkg/influx"
	"context"
)

type AnalyticEngineStruct struct {
	Influx         influx.Influx
}

func NewAnalyticEngine(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD string) *AnalyticEngineStruct{
	ae := &AnalyticEngineStruct{}
	ae.Influx = *influx.NewInflux(INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	return ae
}

func (ae *AnalyticEngineStruct) SendAnalysisResult(ctx context.Context, data *protobuf.RequestValue) (*protobuf.ResponseValue, error) {
	fmt.Println(data)
	//요청한 클러스터들의 매트릭 값 받아오기
	influxData := ae.Influx.SelectMetricsData()
	//매트릭 값을 기반으로 QoS 분석하기
	result := ae.QoSAnalysisResult(influxData)

	//결과값을 리턴값으로 넘겨주기
	return &protobuf.ResponseValue{
		TargetCluster:          result,
	}, nil
}

/*func (ae *AnalyticEngineStruct) SearchInfluxDB(){

}*/

func (ae *AnalyticEngineStruct) QoSAnalysisResult(result []client.Result) string{
	a := result

	fmt.Println("QoSAnalysisResult : ",a)

	return "cluster333"
}

func (ae *AnalyticEngineStruct) StartGRPC(GRPC_PORT string){
	log.Printf("Grpc Server Start at Port %s\n", GRPC_PORT)

	//manager = NewClusterManager()
	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	protobuf.RegisterSendAnalysisResultServer(grpcServer,ae)
	if err := grpcServer.Serve(l); err!=nil {
		log.Fatalf("fail to serve: %v", err)
	}

}