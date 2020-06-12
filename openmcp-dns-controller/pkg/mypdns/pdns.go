package mypdns

import (
	"context"
	//"database/sql"
	"fmt"
	//"github.com/dmportella/powerdns"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"fmt"
	"github.com/mittwald/go-powerdns"
	"github.com/mittwald/go-powerdns/apis/zones"
	ketiv1alpha1 "openmcp-dns-controller/pkg/apis/keti/v1alpha1"
	//"github.com/joeig/go-powerdns/v2"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	PDNS_IP      = os.Getenv("PDNS_IP")      //"10.0.3.12"
	PDNS_PORT    = os.Getenv("PDNS_PORT")    //"8081"
	PDNS_API_KEY = os.Getenv("PDNS_API_KEY") // "1234"
)

func PdnsNewClient() (pdns.Client, error) {
	client, err := pdns.New(
		pdns.WithBaseURL("http://"+PDNS_IP+":"+PDNS_PORT),
		pdns.WithAPIKeyAuthentication(PDNS_API_KEY),
	)
	return client, err
}

func GetZone(client pdns.Client, domain string) (*zones.Zone, error) {

	zone, err := client.Zones().GetZone(context.TODO(), "localhost", domain+".")
	return zone, err

}
func GetZoneList(client pdns.Client) ([]zones.Zone, error) {

	zoneList, err := client.Zones().ListZones(context.TODO(), "localhost")
	return zoneList, err

}
func DeleteZone(client pdns.Client, liveClient client.Client) error {
	instanceDNSEndpointList := &ketiv1alpha1.OpenMCPDNSEndpointList{}
	err := liveClient.List(context.TODO(), nil, instanceDNSEndpointList)
	if err != nil {
		return err
	}
	zoneList, err := GetZoneList(client)
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
			err := client.Zones().DeleteZone(context.TODO(), "localhost", zone.Name)
			if err != nil {
				for {
					fmt.Println("[ERROR Retry Delete] ", err)
					err = client.Zones().DeleteZone(context.TODO(), "localhost", zone.Name)
					if err == nil {
						break
					}
				}
			}

		}

	}
	fmt.Println("[Deleted Pdns Zone] ", deleteZone.Name)
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

	fmt.Println("[Get RecordSets] ", ResourceRecordSets)
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

//func UpdateZoneWithRecords(domainName string, endpoints []*ketiv1alpha1.Endpoint) error{
//	pdnsClient, err := powerdns.NewClient("http://"+PDNS_IP+":"+PDNS_PORT+"/api/v1", PDNS_API_KEY)
//	if err != nil {
//		fmt.Println(err)
//	}
//	zoneInfos, err := pdnsClient.ListZones()
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println("check")
//	for _, zoneInfo := range zoneInfos {
//		fmt.Println(zoneInfo.Name)
//		if zoneInfo.Name == domainName + "." {
//			for _, record := range zoneInfo.Records {
//				if record.Type != "A" {
//					continue
//				}
//				err = pdnsClient.DeleteRecordSet(zoneInfo.Name, record.Name, record.Type)
//				if err != nil {
//					fmt.Println(err)
//				}
//			}
//
//			for _, endpoint := range endpoints{
//				for _, target := range endpoint.Targets {
//					record := powerdns.Record{
//						Name:     endpoint.DNSName + ".",
//						Type:     endpoint.RecordType,
//						Content:  target,
//						TTL:      int(endpoint.RecordTTL),
//						Disabled: false,
//					}
//					p, err := pdnsClient.CreateRecord(zoneInfo.Name, record)
//					fmt.Println(p, err)
//				}
//
//
//			}
//
//
//
//			break
//		}
//	}
//	fmt.Println("check2")
//	return nil
//
//}
//func UpdateZoneWithRecords(client pdns.Client, domainName string,  resourceRecordSets []zones.ResourceRecordSet) error{
//
//
//	db, err := sql.Open("mysql", "root:ketilinux@tcp(10.0.3.12:3306)/powerdns")
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer db.Close()
//	tx, err := db.Begin() // 트랜잭션 Begin
//	defer tx.Rollback()
//
//	var domain_id int
//	err = db.QueryRow("SELECT id FROM domains WHERE name = ?",domainName).Scan(&domain_id)
//
//	tx.Exec("DELETE FROM records WHERE type='A'")
//
//	for _, resourceRecordSet := range resourceRecordSets{
//		for _, record := range resourceRecordSet.Records{
//			_, err = tx.Exec("INSERT INTO records (domain_id, name, type, content, ttl, prio) VALUES (?, ?, ?, ?, ?, ?)", domain_id, resourceRecordSet.Name, resourceRecordSet.Type, record.Content, resourceRecordSet.TTL, 0)
//			fmt.Println(err)
//		}
//
//
//		//err := client.Zones().AddRecordSetToZone(context.TODO(), "localhost", domainName+".", resourceRecordSet)
//		//if err != nil{
//		//	fmt.Println("[Update Err] " ,err, resourceRecordSet)
//		//	continue
//		//	// return err
//		//}
//	}
//	err = tx.Commit()
//
//
//	return nil
//}
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
		fmt.Println("Update Zone ", domainName)
		err = UpdateZoneWithRecords(pdnsClient, domainName, resourceRecordSets)
		if err != nil {
			fmt.Println("[OpenMCP External DNS Controller] : UpdateZone?  ", err)
		}
	} else {
		fmt.Println("Create Zone ", domainName)
		err = CreateZoneWithRecords(pdnsClient, domainName, resourceRecordSets)
		if err != nil {
			fmt.Println("[OpenMCP External DNS Controller] : CreateZone? ", err)
		}
	}
	return err
}
