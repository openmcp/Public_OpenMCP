package db

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/configor"
	_ "github.com/lib/pq"
)

var Config = struct {
	DB struct {
		Host     string
		User     string
		Password string
		Port     string
	}
}{}

func initDBConfig() {
	configor.Load(&Config, "dbconfig.yml")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func InsertReadyNode(cluster string, nodenm string, publicIPAddress string, status string, provider string) {
	initDBConfig()
	var connectionString string = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", Config.DB.Host, Config.DB.User, Config.DB.Password, "portal-controller", Config.DB.Port)
	db, err := sql.Open("postgres", connectionString)
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	sqlStatement := `INSERT INTO readynode(cluster_nm,node_nm,status,ip_addr,provider)
											VALUES ($1,$2,$3,$4,$5) ON
											CONFLICT (cluster_nm,node_nm) where coalesce(cluster_nm, node_nm) is not null 
											DO UPDATE
											SET
												status = $3,
												ip_addr = $4;`
	_, err = db.Exec(sqlStatement, cluster, nodenm, status, publicIPAddress, provider)
	checkError(err)
	fmt.Println("Updated data")

	defer db.Close()
}
