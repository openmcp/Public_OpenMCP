package mypdns

import (
	"context"
	"openmcp/openmcp/util/controller/logLevel"

	//"database/sql"
	//"github.com/dmportella/powerdns"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"fmt"
	"github.com/mittwald/go-powerdns"
	"github.com/mittwald/go-powerdns/apis/zones"
	ketiv1alpha1 "openmcp/openmcp/openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	//"github.com/joeig/go-powerdns/v2"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	PDNS_IP      = os.Getenv("PDNS_IP")      //"10.0.3.12"
	PDNS_PORT    = os.Getenv("PDNS_PORT")    //"8081"
	PDNS_API_KEY = os.Getenv("PDNS_API_KEY") // "1234"
)

func PdnsNewClient() (pdns.Client, error) {
	pdnsClient, err := pdns.New(
		pdns.WithBaseURL("http://"+PDNS_IP+":"+PDNS_PORT),
		pdns.WithAPIKeyAuthentication(PDNS_API_KEY),
	)
	return pdnsClient, err
}

func GetZone(pdnsClient pdns.Client, domain string) (*zones.Zone, error) {

	zone, err := pdnsClient.Zones().GetZone(context.TODO(), "localhost", domain+".")
	return zone, err

}
func GetZoneList(pdnsClient pdns.Client) ([]zones.Zone, error) {

	zoneList, err := pdnsClient.Zones().ListZones(context.TODO(), "localhost")
	return zoneList, err

}
func DeleteZone(pdnsClient pdns.Client, liveClient client.Client) error {
	instanceDNSEndpointList := &ketiv1alpha1.OpenMCPDNSEndpointList{}
	err := liveClient.List(context.TODO(), instanceDNSEndpointList, &client.ListOptions{})
	if err != nil {
		return err
	}
	zoneList, err := GetZoneList(pdnsClient)
	if err != nil {
		return err
	}

	var deleteZone zones.Zone
	for _, zone := range zoneList {
		find := false
		for _, instanceDNSEndpoint := range instanceDNSEndpointList.Items {
			for _, domain := range instanceDNSEndpoint.Spec.Domains {
				if zone.Name == domain+"." {
					find = true
					break
				}
			}
			if find {
				break
			}
		}
		if !find {
			deleteZone = zone
			err := pdnsClient.Zones().DeleteZone(context.TODO(), "localhost", zone.Name)
			if err != nil {
				for {
					logLevel.KetiLog(0, "[ERROR Retry Delete] ", err)
					err = pdnsClient.Zones().DeleteZone(context.TODO(), "localhost", zone.Name)
					if err == nil {
						break
					}
				}
			}

		}

	}
	logLevel.KetiLog(0, "[Deleted Pdns Zone] ", deleteZone.Name)
	return nil
}

func GetResourceRecordSets(domainName string, Endpoints []*ketiv1alpha1.Endpoint) []zones.ResourceRecordSet {
	ResourceRecordSets := []zones.ResourceRecordSet{}
	for _, endpoint := range Endpoints {

		startIndex := len(endpoint.DNSName) - len(domainName)
		if startIndex < 0 {
			continue
		}
		if domainName != endpoint.DNSName[startIndex:] {
			continue
		}
		records := []zones.Record{}

		for _, target := range endpoint.Targets {
			record := zones.Record{
				Content:  target,
				Disabled: false,
				SetPTR:   false,
			}
			records = append(records, record)
		}
		if len(records) == 0 {
			continue
		}

		ResourceRecordSet := zones.ResourceRecordSet{
			Name:       endpoint.DNSName + ".",
			Type:       endpoint.RecordType,
			TTL:        int(endpoint.RecordTTL),
			ChangeType: zones.ChangeTypeReplace,
			Records:    records,
			Comments:   nil,
		}

		ResourceRecordSets = append(ResourceRecordSets, ResourceRecordSet)

	}

	logLevel.KetiLog(1, "[Get RecordSets] ", ResourceRecordSets)
	return ResourceRecordSets
}
func UpdateZoneWithRecords(client pdns.Client, domainName string, resourceRecordSets []zones.ResourceRecordSet) error {
	for _, resourceRecordSet := range resourceRecordSets {
		err := client.Zones().AddRecordSetToZone(context.TODO(), "localhost", domainName+".", resourceRecordSet)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateZoneWithRecords(client pdns.Client, domainName string, resourceRecordSets []zones.ResourceRecordSet) error {
	_, err := client.Zones().CreateZone(context.Background(), "localhost", zones.Zone{
		Name: domainName + ".",
		Type: zones.ZoneTypeZone,
		Kind: zones.ZoneKindNative,
		Nameservers: []string{
			"ns1.example.com.",
			"ns2.example.com.",
		},
		ResourceRecordSets: resourceRecordSets,
	})
	if err != nil {
		return err
	}
	return nil
}

func SyncZone(pdnsClient pdns.Client, domainName string, Endpoints []*ketiv1alpha1.Endpoint) error {

	_, err := GetZone(pdnsClient, domainName)
	resourceRecordSets := GetResourceRecordSets(domainName, Endpoints)

	if err == nil {
		// Already Exist
		logLevel.KetiLog(0, "Update Zone ", domainName)
		err = UpdateZoneWithRecords(pdnsClient, domainName, resourceRecordSets)
		if err != nil {
			logLevel.KetiLog(0, "[OpenMCP External DNS Controller] : UpdateZone?  ", err)
		}
	} else {
		logLevel.KetiLog(0, "Create Zone ", domainName)
		err = CreateZoneWithRecords(pdnsClient, domainName, resourceRecordSets)
		if err != nil {
			logLevel.KetiLog(0, "[OpenMCP External DNS Controller] : CreateZone? ", err)
		}
	}
	return err
}
