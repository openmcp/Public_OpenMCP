const fs = require("fs"); //database.json파일 접근
const express = require("express");
const bodyParser = require("body-parser");
const app = express();
var os = require("os");
var path = require("path");
const data = fs.readFileSync("./config.json");
const conf = JSON.parse(data);

const port = process.env.PORT || 5000;
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

app.get("/api/hello", (req, res) => {
  res.send({ messge: "Hello Express!" });
});

// const apiServer = "http://192.168.0.51:4885"; //로컬 API 서버
const apiServer = conf.api.url; //로컬 API 서버
// const apiServer = "http://192.168.0.4:4885"; //kvm API 서버
// const apiServer = "http://10.0.3.40:4885";

//데이터베이스 접속 설정
const { Client } = require("pg");
const { toNamespacedPath } = require("path");

const connection = new Client({
  user: conf.db.user,
  host: conf.db.host,
  database: conf.db.database,
  password: conf.db.password,
  port: conf.db.port,
});

//데이터베이스 접속
connection.connect();

//데이터베이스에서 데이터 가져오기
// app.get("/api/customers", (req, res) => {
//   // res.send()
//   connection.query("SELECT * FROM CUSTOMER", (err, result) => {
//     res.send(result.rows);
//   });
// });




function getDateTime() {
  var d = new Date();
  d = new Date(d.getTime());
  var date_format_str =
    d.getFullYear().toString() +
    "-" +
    ((d.getMonth() + 1).toString().length == 2
      ? (d.getMonth() + 1).toString()
      : "0" + (d.getMonth() + 1).toString()) +
    "-" +
    (d.getDate().toString().length == 2
      ? d.getDate().toString()
      : "0" + d.getDate().toString()) +
    " " +
    (d.getHours().toString().length == 2
      ? d.getHours().toString()
      : "0" + d.getHours().toString()) +
    ":" +
    // ((parseInt(d.getMinutes() / 5) * 5).toString().length == 2
    //   ? (parseInt(d.getMinutes() / 5) * 5).toString()
    //   : "0" + (parseInt(d.getMinutes() / 5) * 5).toString()) +
    // ":00";
    (d.getMinutes().toString().length == 2
      ? d.getMinutes().toString()
      : "0" +d.getMinutes().toString()) +
    ":" + 
    (d.getSeconds().toString().length == 2
      ? d.getSeconds().toString()
      : "0" +d.getSeconds().toString());    
  // console.log(date_format_str);
  return date_format_str;
}

///////////////////////
// Write Log
///////////////////////
app.post("/apimcp/portal-log", (req, res) => {
  const bcrypt = require("bcrypt");
  var created_time = getDateTime();

  connection.query(
    `insert into tb_portal_logs values ('${req.body.userid}','${req.body.code}','${created_time}');`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Update success",
      };

      if (err !== null) {
        console.log(err)
        result_set = {
          data: [],
          message: "Update log failed : " + err,
        };
      } 

      res.send(result_set);
    }
  );
});

///////////////////////
// Login
///////////////////////

app.post("/user_login", (req, res) => {
  const bcrypt = require("bcrypt");

  connection.query(
    `select * from tb_accounts where user_id = '${req.body.userid}';`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Please check your Password",
      };
      console.log(result)

      if (result.rows.length < 1) {
        result_set = {
          data: [],
          message: "There is no user, please check your account",
        };
        res.send(result_set);
      } else {
        const hashPassword = result.rows[0].user_password;
        bcrypt.compare(req.body.password, hashPassword).then(function (r) {
          if (r) {
            // console.log("compare", r, result_set)
            result_set = {
              data: result,
              message: "Login Successful !!",
            };
            // console.log("compare", r, result_set);
          }
          res.send(result_set);
        });
      }
    }
  );
});


///////////////////////
// Dashboard APIs 
///////////////////////
app.get("/dashboard", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/dashboard.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);

  // var request = require("request");
  // var options = {
  //   uri: `${apiServer}/apis/dashboard`,
  //   method: "GET",
  //   // headers: {
  //   //   Authorization:
  //   //     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMxMDQ4NzcsImlhdCI6MTYwMzEwMTI3NywidXNlciI6Im9wZW5tY3AifQ.mgO5hRruyBioZLTJ5a3zwZCkNBD6Bg2T05iZF-eF2RI",
  //   // },
  // };

  // request(options, function (error, response, body) {
  //   if (!error && response.statusCode == 200) {
  //     // console.log("result", body);
  //     res.send(body);
  //   } else {
  //     console.log("error", error);
  //     return error;
  //   }
  // });
});


app.get("/dashboard-master-cluster", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/dashboard_master_cluster.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});


let token = "";
// Projects 리스트 가져오기
app.get("/api/projects", (req, res) => {
  var request = require("request");
  // var url = "http://192.168.0.152:31635/token?username=openmcp&password=keti";
  // var uri ="http://192.168.0.152:31635/api/v1/namespaces/kube-system/pods?clustername=cluster1";

  var options = {
    uri:
      "http://192.168.0.152:31635/api/v1/namespaces/kube-system/pods?clustername=cluster1",
    method: "GET",
    headers: {
      Authorization:
        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMxMDQ4NzcsImlhdCI6MTYwMzEwMTI3NywidXNlciI6Im9wZW5tY3AifQ.mgO5hRruyBioZLTJ5a3zwZCkNBD6Bg2T05iZF-eF2RI",
    },
  };

  var options = {
    uri:
      "http://192.168.0.152:31635/api/v1/namespaces/kube-system/pods?clustername=cluster1",
    method: "GET",
    headers: {
      Authorization:
        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMxMDQ4NzcsImlhdCI6MTYwMzEwMTI3NywidXNlciI6Im9wZW5tY3AifQ.mgO5hRruyBioZLTJ5a3zwZCkNBD6Bg2T05iZF-eF2RI",
    },
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      console.log("result", body);
    } else {
      console.log("error", error);
      return error;
    }
  });

  //   request(url, function (error, response, body) {
  //     if (!error && response.statusCode == 200) {
  //         console.log(body);
  //         token = body.token;
  //     } else {
  //         return error
  //     }
  //   });

  connection.query("SELECT * FROM PROJECT_LIST", (err, result) => {
    res.send(result.rows);
  });
});

///////////////////////
// Projects APIs 
///////////////////////

// Prjects
app.get("/projects", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  // console.log("projects")

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/projects`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });

});

// Prjects > overview
app.get("/projects/:project/overview", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_overview.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}`,
    method: "GET",
  };

  // console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});


// Prjects > create
app.post("/projects/create", (req, res) => {
  requestData = {
    project : req.body.project,
    clusters : req.body.clusters,
  }

  
  var data = JSON.stringify(requestData);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/projects/create`,
    method: "POST",
    body: data
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
      console.log(body)
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// // Prjects > get Clusters Names
// app.get("/clusters/name", (req, res) => {
//   let rawdata = fs.readFileSync(
//     "./json_data/clusters_name.json"
//   );
//   let overview = JSON.parse(rawdata);
//   //console.log(overview);
//   res.send(overview);
// });

// Prjects > Resources > Workloads > Deployments
app.get("/projects/:project/resources/workloads/deployments", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_deployments.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.params.project}/deployments`,
    method: "GET",
  };


  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });


  
});

// Prjects > Resources > Workloads > Deployments > detail
app.get("/projects/:project/resources/workloads/deployments/:deployment",
  (req, res) => {
    // let rawdata = fs.readFileSync(
    //   "./json_data/projects_deployment_detail.json"
    // );
    // let overview = JSON.parse(rawdata);
    // // //console.log(overview);
    // res.send(overview);


    var request = require("request");
    var options = {
      uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.params.project}/deployments/${req.params.deployment}`,
    };

    console.log(options.uri)


    request(options, function (error, response, body) {
      if (!error && response.statusCode == 200) {
        res.send(body);
      } else {
        console.log("error", error);
        return error;
      }
    });

  }
);

// Prjects > Resources > Workloads > Deployments > detail > replica status
app.get("/projects/:project/resources/workloads/deployments/:deployment/replica_status", (req, res) => {

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.params.project}/deployments/${req.params.deployment}/replica_status`,
  };

  console.log(options.uri)


  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });


    // connection.query(
    //   "select * from tb_replica_status order by cluster asc, created_time desc, status desc",
    //   (err, result) => {
    //     var result2 = result.rows.reduce(
    //       (obj, { cluster, status, pod, created_time }, index) => {
    //         if (!obj[cluster]) {
    //           obj[cluster] = { cluster: cluster, pods: [] };
    //         }

    //         obj[cluster].pods.push({
    //           status: status,
    //           name: pod,
    //           created_time: created_time,
    //         });
    //         return obj;
    //       },
    //       {}
    //     );

    //     var arr = [];
    //     for (i = 0; i < Object.keys(result2).length; i++) {
    //       arr.push(result2[Object.keys(result2)[i]]);
    //       // console.log(result2[Object.keys(result2)[i]]);
    //     }
    //     // console.log(arr)

    //     res.send(arr);
    //   }
    // );
    // let rawdata = fs.readFileSync("./json_data/projects_deployment_detail.json");
    // let overview = JSON.parse(rawdata);
    // //console.log(overview);
    // res.send(overview);



  }
);

app.post(
  "/projects/:project/resources/workloads/deployments/:deployment/replica_status/add_pod",
  (req, res) => {
    var create_time = getDateTime();
    var podName = Math.random().toString(36).substring(10);
    connection.query(
      `insert into tb_replica_status values ('${req.body.cluster}', '${podName}','config','${create_time}');`,
      (err, result) => {
        res.send(result);
      }
    );
  }
);

app.delete(
  "/projects/:project/resources/workloads/deployments/:deployment/replica_status/del_pod",
  (req, res) => {
    // console.log("delete", req.body);
    connection.query(
      `delete from tb_replica_status where ctid IN (select ctid from tb_replica_status where cluster = '${req.body.cluster}' order by created_time desc limit 1)`,
      (err, result) => {
        res.send(result);
      }
    );
  }
);

// Deployments 상세부터 구현해나가야 함
// Prjects > Resources > Workloads > Statefulsets
app.get("/projects/:project/resources/workloads/statefulsets", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_statefulsets.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.params.project}/statefulsets`,
    method: "GET",
  };




  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/projects/:project/resources/workloads/statefulsets/:statefulset", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_statefulsets.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.params.project}/statefulsets/${req.params.statefulset}`,
    method: "GET",
  };

  console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > pods
app.get("/projects/:project/resources/pods", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_pods.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/pods`,
    method: "GET",
  };

  // console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });

});

// Prjects > Resources > Pods Detail
app.get("/projects/:project/resources/pods/:pod", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_pod_detail.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/pods/${req.params.pod}?cluster=${req.query.cluster}&project=${req.query.project}`,
    method: "GET",
  };

  // console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > Services
app.get("/projects/:project/resources/services", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_services.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/services`,
    method: "GET",
  };

  // console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > Services Detail
app.get("/projects/:project/resources/services/:service", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_service_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  // console.log(`${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/services/${req.params.service}`)
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/services/${req.params.service}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > Ingress
app.get("/projects/:project/resources/ingress", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_ingress.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/ingress`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > Ingress Detail
app.get("/projects/:project/resources/ingress/:ingress", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_ingress_detail.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/ingress/${req.params.ingress}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > volumes
app.get("/projects/:project/volumes", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_volumes.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/volumes`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > volumes Detail
app.get("/projects/:project/volumes/:volume", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_volume_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/volumes/${req.params.volume}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Config > Secrets
app.get("/projects/:project/config/secrets", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_secrets.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/secrets`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Config > Secrets Detail
app.get("/projects/:project/config/secrets/:secret", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_secret_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/secrets/${req.params.secret}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });

});

// Prjects > Config > ConfigMaps
app.get("/projects/:project/config/config_maps", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_config_maps.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/configmaps`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Config > ConfigMaps Detail
app.get("/projects/:project/config/config_maps/:config_map", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_config_map_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${req.query.cluster}/projects/${req.params.project}/configmaps/${req.params.config_map}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Settings > Members
app.get("/projects/:project/settings/members", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/projects_members.json");
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});


//////////////////////////
// Deployments
/////////////////////////

// Deployments
app.get("/deployments", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/projects_deployments.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/deployments`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/deployments/:deployment", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/deployment_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.query.project}/deployments/${req.params.deployment}`,
    method: "GET",
  };

  // console.log(`${apiServer}/apis/clsuters/${req.query.cluster}/projects/${req.query.project}/deployments/${req.params.deployment}`)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.post("/deployments/migration", (req, res) => {
  const YAML = req.body.yaml
  // console.log(YAML)
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/migration`,
    method: "POST",
    body: YAML
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
      console.log(body)
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.post("/deployments/create", (req, res) => {
  const YAML = req.body.yaml
  
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/yamlapply`,
    method: "POST",
    body: YAML
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});


app.get("/snapshots", (req, res) => {
  let deployment = req.query.deployment
  let rawdata = fs.readFileSync("./json_data/snapshots.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);

  // var request = require("request");
  // var options = {
  //   uri: `${apiServer}/apis/deployments`,
  //   method: "GET",
  // };

  // request(options, function (error, response, body) {
  //   if (!error && response.statusCode == 200) {
  //     res.send(body);
  //   } else {
  //     console.log("error", error);
  //     return error;
  //   }
  // });
});

///////////////////////
// Clusters APIs
///////////////////////
app.get("/clusters", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/clusters.json");
  // let rawdata = fs.readFileSync("./json_data/clusters2_warning.json");
  // let rawdata = fs.readFileSync("./json_data/clusters3-1_normal.json");
  // let rawdata = fs.readFileSync("./json_data/clusters3-1_50.json");
  // let rawdata = fs.readFileSync("./json_data/clusters3-1_70.json");
  // let rawdata = fs.readFileSync("./json_data/clusters3-1_80.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  // console.log("cluster")

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/clusters-joinable", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/clusters_joinable.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/joinableclusters`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Clusters > overview
app.get("/clusters/:cluster/overview", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/clusters_overview.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/overview?clustername=${req.params.cluster}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Clusters joinable > overview
app.get("/clusters-joinable/:cluster/overview", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/clusters_joinable_overview.json");
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});

// Clusters > nodes
app.get("/clusters/:cluster/nodes", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/clusters_nodes.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  const clusterName = req.params.cluster
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${clusterName}/nodes`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Clusters > nodes > detail ///// not use
// app.get("/clusters/:cluster/nodes/:node", (req, res) => {
//   let rawdata = fs.readFileSync("./json_data/clusters_node_detail.json");
//   let overview = JSON.parse(rawdata);
//   //console.log(overview);
//   res.send(overview);
// });

// Clusters > pods
app.get("/clusters/:cluster/pods", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/clusters_pods.json");
  // let overview = JSON.parse(rawdata);
  // //console.log(overview);
  // res.send(overview);

  const clusterName = req.params.cluster
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/clusters/${clusterName}/pods`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Clusters > pods > detail
app.get("/clusters/:cluster/pods/:pod", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/clusters_pod_detail.json");
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});

// Clusters > Storage Class
app.get("/clusters/:cluster/storage_class", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/clusters_storage_class.json");
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});

// Clusters > Storage Class > detail
app.get("/clusters/:cluster/storage_class/:storage_class", (req, res) => {
  let rawdata = fs.readFileSync(
    "./json_data/clusters_storage_class_detail.json"
  );
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});


// cluster > join
app.post("/cluster/join", (req, res) => {
  requestData = {
		clusterName : req.body.clusterName,
    clusterAddress : req.body.clusterAddress
  }

  var data = JSON.stringify(requestData);
  var request = require("request");
  var options = {
    // https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster2?clustername=openmcp
    uri: `${apiServer}/apis/clusters/join`,
    method: "POST",
    body: data
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
      console.log(body)
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// cluster > unjoin
app.post("/cluster/unjoin", (req, res) => {
  console.log("cluster/unjoin");
  requestData = {
		clusterName : req.body.clusterName
  }

  var data = JSON.stringify(requestData);
  var request = require("request");
  var options = {
    // https://192.168.0.152:30000/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/openmcpclusters/cluster2?clustername=openmcp
    uri: `${apiServer}/apis/clusters/unjoin`,
    method: "POST",
    body: data
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
      console.log(body)
    } else {
      console.log("error", error);
      return error;
    }
  });
});



/////////////
// Nodes
/////////////
app.get("/nodes", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/nodes.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/nodes`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

//asdasd
app.post("/nodes/add/eks", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_eks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the EKS Auth Informations.\n
          Settings > Config > Public cloud Auth > EKS`,
        };
        console.log(result_set);
        return res.send(result_set);
      }

      requestData = {
        region : req.body.region,
        cluster:req.body.cluster,
        nodePool: req.body.nodePool,
        desiredCnt:req.body.desiredCnt,
        accessKey : result.rows[0].accessKey,
        secretKey : result.rows[0].secretKey,
      }

      console.log(requestData)
      var data = JSON.stringify(requestData);

      var request = require("request");
      var options = {
        // uri: `${apiServer}/apis/addeksnode`,
        uri: `${apiServer}/apis/changeeksnode`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
          console.log(body)
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
  
});

app.post("/nodes/add/aks", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_aks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        cluster : req.body.cluster,
        desiredCnt : req.body.desiredCnt,
        nodePool : req.body.nodePool,
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId
      }

      console.log("addNodeAKS : ", requestData);
      var data = JSON.stringify(requestData);

      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/addaksnode`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
          console.log(body)
        } else {
          console.log("error", error);
          return error;
        }
      });

    }
  );
});

app.post("/nodes/add/gke", (req, res) => {
  connection.query(
    `select * 
     from tb_config_gke
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the GKE Auth Informations.\n
          Settings > Config > Public cloud Auth > GKE`,
        };
        return res.send(result_set);
      }

      requestData = {
        projectId: result.rows[0].projectID,
        clientEmail: result.rows[0].clientEmail,
        privateKey: result.rows[0].privateKey,
        cluster : req.body.cluster,
        nodePool: req.body.nodePool,
        desiredCnt : req.body.desiredCnt,
      }
     
      // console.log("gke/addnode",requestData)
      var data = JSON.stringify(requestData);

      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/gkechangenodecount`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
          console.log("body",body)
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/add/kvm", (req, res) => {
  connection.query(
    `select * 
     from tb_config_kvm
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the KVM Auth Informations.\n
          Settings > Config > Public cloud Auth > KVM`,
        };
        return res.send(result_set);
      }

      requestData = {
        agentURL : result.rows[0].agentURL,
        master : result.rows[0].mClusterName,
        mpass : result.rows[0].mClusterPwd,
        newvm: req.body.newvm,
        template : req.body.template,
        wpass : req.body.newVmPassword,
        cluster : req.body.cluster,
      }

      // console.log(requestData);

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/createkvmnode`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log(response);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/delete/kvm", (req, res) => {
  connection.query(
    `select * 
     from tb_config_kvm
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the KVM Auth Informations.\n
          Settings > Config > Public cloud Auth > KVM`,
        };
        return res.send(result_set);
      }

      requestData = {
        agentURL : result.rows[0].agentURL,
        targetvm : req.body.node,
        mastervm : result.rows[0].mClusterName,
        mastervmpwd : result.rows[0].mClusterPwd
      }

      // console.log(requestData, result.rows[0])

      // res.send(result.rows);

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/deletekvmnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});


//////////////////
// Nodes > datail
//////////////////
app.get("/nodes/:node", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/nodes_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/nodes/${req.params.node}?clustername=${req.query.clustername}&provider=${req.query.provider}`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.post("/nodes/eks/start", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_eks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the EKS Auth Informations.\n
          Settings > Config > Public cloud Auth > EKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        akid : result.rows[0].accessKey,
        secretKey : result.rows[0].secretKey,
        region : req.body.region,
        node : req.body.node,
      }

      // console.log(requestData, result.rows[0])

      // res.send(result.rows);

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/starteksnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/eks/stop", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_eks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the EKS Auth Informations.\n
          Settings > Config > Public cloud Auth > EKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        akid : result.rows[0].accessKey,
        secretKey : result.rows[0].secretKey,
        region : req.body.region,
        node : req.body.node,
      }

      console.log(requestData)

      // res.send(result.rows);

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/stopeksnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          console.log(body);
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/eks/change", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_eks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the EKS Auth Informations.\n
          Settings > Config > Public cloud Auth > EKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        akid : result.rows[0].accessKey,
        secretKey : result.rows[0].secretKey,
        region : req.body.region,
        type : req.body.type,
        node : req.body.node,
      }

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/changeekstype`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          console.log(body);
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/aks/start", (req, res) => {
  connection.query(
    `select * 
     from tb_config_aks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        cluster : req.body.cluster,
        node : req.body.node,
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId
      }

      
      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/startaksnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/aks/stop", (req, res) => {
  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_aks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        cluster : req.body.cluster,
        node : req.body.node,
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId
      }

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/stopaksnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/clusters/aks/change", (req, res) => {
  connection.query(
    `select * 
     from tb_config_aks
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId,
        cluster : req.body.cluster,
        poolName : req.body.poolName,
        skuTierStr : req.body.tier,
        skuNameStr : req.body.type,
      }

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/akschangevmss`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/kvm/stop", (req, res) => {
  connection.query(
    `select * 
     from tb_config_kvm
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the KVM Auth Informations.\n
          Settings > Config > Public cloud Auth > KVM`,
        };
        return res.send(result_set);
      }

      requestData = {
        node : req.body.node,
        agentURL : result.rows[0].agentURL
      }

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/stopkvmnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/kvm/start", (req, res) => {
  connection.query(
    `select * 
     from tb_config_kvm
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the KVM Auth Informations.\n
          Settings > Config > Public cloud Auth > KVM`,
        };
        return res.send(result_set);
      }

      requestData = {
        node : req.body.node,
        agentURL : result.rows[0].agentURL
      }

      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/startkvmnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.post("/nodes/kvm/change", (req, res) => {
  connection.query(
    `select * 
     from tb_config_kvm
     where cluster='${req.body.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the KVM Auth Informations.\n
          Settings > Config > Public cloud Auth > KVM`,
        };
        return res.send(result_set);
      }

      requestData = {
        agentURL : result.rows[0].agentURL,
        node : req.body.node,
        cpu : req.body.cpu,
        memory : req.body.memory,
      }

      console.log(requestData)
      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/changekvmnode`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});


/////////////////////////
// Public Cloud Cluster
//////////////////////////
app.get("/aws/eks-type", (req, res) => {
  connection.query(
    `select * from tb_codes where kinds='EKS-TYPE' order by etc;`,
    (err, result) => {
      res.send(result.rows);
    }
  );
});

app.get("/azure/aks-type", (req, res) => {
  connection.query(
    `select * from tb_codes where kinds='AKS-TYPE' order by etc;`,
    (err, result) => {

      res.send(result.rows);
    }
  );
});

app.get("/azure/pool/:cluster", (req, res) => {
  var cluster = req.params.cluster
  console.log(req.params.cluster);
  connection.query(
    `select * 
     from tb_config_aks
     where cluster='${req.params.cluster}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      console.log(result.rows[0].clientId)

      requestData = {
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId
      }
      
      var data = JSON.stringify(requestData);
      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/aksgetallres`,
        method: "POST",
        body : data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          var clusterInfo = {};
          console.log(body)
          for (let value of JSON.parse(body)){
            // if(value.name == cluster){ //임시로 막음(일치하는 클러스터가 없음)
            if(value.name === "aks-cluster-01"){ //임시로 하드코딩함
              clusterInfo = value;
            }
          }
          res.send(clusterInfo);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.get("/eks/clusters", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/eks_clusters.json");
  let overview = JSON.parse(rawdata);
  //console.log(overview);
  res.send(overview);
});

app.get("/eks/clusters/workers", (req, res) => {
  var clusterName = req.query.clustername;
  console.log(req.query)
  // let rawdata = fs.readFileSync("./json_data/eks_workers.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_eks
     where cluster='${clusterName}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the EKS Auth Informations.\n
          Settings > Config > Public cloud Auth > EKS`,
        };
        console.log(result_set);
        return res.send(result_set);
      }

      requestData = {
        region : result.rows[0].region,
        accessKey : result.rows[0].accessKey,
        secretKey : result.rows[0].secretKey,
      }

      var data = JSON.stringify(requestData);

      var request = require("request");
      var options = {
        // uri: `${apiServer}/apis/addeksnode`,
        uri: `${apiServer}/apis/geteksclusterinfo`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          var result = JSON.parse(body)
          result.map((item)=> {
            if(item.name == clusterName){
              res.send(item.nodegroups);
            }
          })

          // res.send(body);
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );

});

app.get("/gke/clusters", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/gke_clusters.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});

app.get("/gke/clusters/pools", (req, res) => {
  var clusterName = req.query.clustername;
  console.log(clusterName)
  // let rawdata = fs.readFileSync("./json_data/gke_node_pools.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    `select * 
     from tb_config_gke
     where cluster='${clusterName}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the GKE Auth Informations.\n
          Settings > Config > Public cloud Auth > GKE`,
        };
        return res.send(result_set);
      }

      requestData = {
        projectId: result.rows[0].projectID,
        clientEmail: result.rows[0].clientEmail,
        privateKey: result.rows[0].privateKey,
      }
     
      var data = JSON.stringify(requestData);
      // console.log("gke/addnode",data)

      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/getgkeclusters`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          var result = JSON.parse(body)
          result.map((item)=> {
            if(item.clusterName == clusterName){
              res.send(item.nodePools);
            }
          })
        } else {
          console.log("error", error);
          return error;
        }
      });
      
    }
  );
});

app.get("/aks/clusters", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/aks_clusters.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});

app.get("/aks/clusters/pools", (req, res) => {
  var clusterName = req.query.clustername;
  // let rawdata = fs.readFileSync("./json_data/aks_node_pools.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    // tb_auth_eks > seq,cluster,accessKey,secretKey
    `select * 
     from tb_config_aks
     where cluster='${clusterName}';`,
    (err, result) => {
      if (result.rows.length  === 0){
        const result_set = {
          error : true,
          message: `Auth Information does not Exist.\nPlease Enter the AKS Auth Informations.\n
          Settings > Config > Public cloud Auth > AKS`,
        };
        return res.send(result_set);
      }

      requestData = {
        clientId : result.rows[0].clientId,
        clientSec : result.rows[0].clientSec,
        tenantId : result.rows[0].tenantId,
        subId : result.rows[0].subId
      }

      console.log("addNodeAKS : ", requestData);

      var data = JSON.stringify(requestData);
      // console.log("aks/addnode",data)

      var request = require("request");
      var options = {
        uri: `${apiServer}/apis/aksgetallres`,
        method: "POST",
        body: data
      };

      request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
          var result = JSON.parse(body)
          result.map((item)=> {
            if(item.name == clusterName){
              res.send(item.agentpools);
            }
          })
        } else {
          console.log("error", error);
          return error;
        }
      });

    }
  );


});

app.get("/kvm/clusters", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/kvm_clusters.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});

//////////////////////////
// Pods
/////////////////////////

// Pods
app.get("/pods", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/pods.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/pods`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Pods > detail
app.get("/pods/:pod", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/pods_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/pods/${req.params.pod}?cluster=${req.query.cluster}&project=${req.query.project}`,
    method: "GET",
  };

  console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/pods/:pod/physicalResPerMin", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/pods_detail.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/pods/${req.params.pod}/physicalResPerMin?cluster=${req.query.cluster}&project=${req.query.project}`,
    method: "GET",
  };

  console.log(options.uri)

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});


app.get("/hpa", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/hpa.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/hpa`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/vpa", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/vpa.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/vpa`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});



////////////////////////
// Ingress
////////////////////////

// Prjects > Resources > Ingress
app.get("/ingress", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/ingress.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/ingress`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Prjects > Resources > Ingress Detail
app.get("/ingress/:ingress", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/ingress_detail.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});


////////////////////////
// Services
////////////////////////

// Services
app.get("/services", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/services.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/services`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// Services Detail
app.get("/services/:service", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/service_detail.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);

});



////////////////////////
// Network > DNS
////////////////////////

// DNS > Services
app.get("/dns", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/dns.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/dns`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

// DNS > Detail
app.get("/dns/:dns", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/dns_detail.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);
});


///////////////////////
//Settings > Accounts
///////////////////////
app.get("/settings/accounts", (req, res) => {
  connection.query(`select
  user_id, 
  role_id,
  array(
      select role_name 
      from tb_account_role t 
      where t.role_id = ANY(u.role_id)
      ) as role,
  last_login_time,
  created_time
  from tb_accounts u`, (err, result) => {
    res.send(result.rows);
  });
});

app.post("/create_account", (req, res) => {
  const bcrypt = require("bcrypt");
  const saltRounds = 10;
  var password = "";

  bcrypt.genSalt(saltRounds, function (err, salt) {
    bcrypt.hash(req.body.password, salt, function (err, hash_password) {
      var create_time = getDateTime();
      connection.query(
        `insert into tb_accounts values ('${req.body.userid}', '${hash_password}','${req.body.role}','${create_time}','${create_time}');`,
        (err, result) => {
          if (err !== "null") {
            const result_set = {
              data: [],
              message: "Account creation was successful !!",
            };
            res.send(result_set);
          } else {
            const result_set = {
              data: [],
              message: "Account creation was faild, please check account",
            };
            res.send(result_set);
          }
        }
      );
    });
  });
});

app.get("/account-roles", (req, res) => {
  connection.query(
    `select * from tb_account_role;`,
    (err, result) => {
      // var result_set = {
      //   data: [],
      //   message: "Update success",
      // };

      // if (err !== "null") {
      //   console.log(err)
      //   const result_set = {
      //     data: [],
      //     message: "Update log failed : " + err,
      //   };
      // } 
      res.send(result.rows);
    }
  );
});

app.put("/update/account-roles", (req, res) => {
  // console.log(req.body);
  connection.query(
    `update tb_accounts set roles = '{"${req.body.role}"}' where user_id = '${req.body.userid}';`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check account : " + err,
        };
        res.send(result_set);
      }
    }
  );


  // bcrypt.genSalt(saltRounds, function (err, salt) {
  //   bcrypt.hash(req.body.password, salt, function (err, hash_password) {
  //     var create_time = getDateTime();
      
  //   });
  // });
});

//////////////////////////
// Settings > Groups Role
//////////////////////////
app.get("/settings/group-role", (req, res) => {
  connection.query(`select
  ga.group_id,
  ga.group_name,
  array(select role_name 
    from tb_account_role t 
    where t.role_id = ANY(ga.role_id)) as role,
    projects,
  ga.description,
  ga.member
  from tb_group_role ga
  order by group_id`, (err, result) => {
    res.send(result.rows);
  });
});

app.post("/settings/group-role", (req, res) => {
  let query = `
  INSERT INTO tb_group_role (group_name, description, role_id, member, projects)
  VALUES ('${req.body.groupName}', '${req.body.description}', '{${req.body.role_id}}', '{${req.body.user_id}}', '{${req.body.projects}}')
  `
  console.log(query);
  connection.query(query, 
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Group role is saved !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Save was faild, please check policy : " + err,
        };
        res.send(result_set);
      }
    });
});

app.put("/settings/group-role", (req, res) => {
  let query =  `update tb_group_role 
      set group_name='${req.body.groupName}', 
          description='${req.body.description}',
          role_id='{${req.body.role_id}}',
          member='{${req.body.user_id}}',
          projects = '{${req.body.projects}}'
      where group_id = ${req.body.group_id}`

  console.log(query);
  connection.query(query, 
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Group role is updated !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check policy : " + err,
        };
        res.send(result_set);
      }
    });
});

app.delete("/settings/group-role", (req, res) => {
  let query = `delete from tb_group_role 
  where "group_id" = ${req.body.group_id};`;
  console.log(query);
  connection.query(
    query,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});


//////////////////////////
// Settings > Policy
//////////////////////////
app.get("/settings/policy/openmcp-policy", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/settings_policy.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  // let sql =`select  policy_id, policy_name,
  //                   rate, period
  //           from tb_policy`

  // connection.query(sql, (err, result) => {
  //   res.send(result.rows);
  // });
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/policy/openmcp`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.post("/settings/policy/openmcp-policy", (req, res) => {
  requestData = {
    policyName : req.body.policyName,
    values : req.body.values,
  }

  var data = JSON.stringify(requestData);
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/policy/openmcp/edit`,
    method: "POST",
    body: data
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/settings/policy/project-policy", (req, res) => {
  let sql =`select  project, cluster, cls_cpu_trh_r, cls_mem_trh_r,
            pod_cpu_trh_r, pod_mem_trh_r, updated_time
            from tb_policy_projects`

  connection.query(sql, (err, result) => {
    res.send(result.rows);
  });
});

app.put("/settings/policy/project-policy", (req, res) => {
  var updated_time = getDateTime();

  let sql = `
  INSERT INTO tb_policy_projects (project, cluster, cls_cpu_trh_r, cls_mem_trh_r, pod_cpu_trh_r, pod_mem_trh_r, updated_time)
  VALUES ('${req.body.project}', '${req.body.cluster}', ${req.body.cls_cpu_trh_r}, ${req.body.cls_mem_trh_r}, ${req.body.pod_cpu_trh_r}, ${req.body.pod_mem_trh_r}, '${updated_time}')
  ON CONFLICT (project, cluster) DO
  UPDATE 
    SET project='${req.body.project}',
    cluster='${req.body.cluster}',
    cls_cpu_trh_r=${req.body.cls_cpu_trh_r},
    cls_mem_trh_r=${req.body.cls_mem_trh_r},
    pod_cpu_trh_r=${req.body.pod_cpu_trh_r},
    pod_mem_trh_r=${req.body.pod_mem_trh_r},
    updated_time='${updated_time}'
  `
  connection.query(sql, (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check data : " + err,
        };
        res.send(result_set);
      }
    }
  );
});


// Settings > Config > Public Cloud Auth
app.get("/settings/config/pca/eks", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    // tb_auth_eks > seq,cluster,type,accessKey,secretKey
    `select * from tb_config_eks;`,
    (err, result) => {
      res.send(result.rows);
    }
  );
});

app.post("/settings/config/pca/eks", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  connection.query(
    `insert into tb_config_eks (cluster,"accessKey","secretKey","region") values ('${req.body.cluster}','${req.body.accessKey}','${req.body.secretKey}','${req.body.region}');`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Insert success",
      };

      if (err !== null) {
        console.log(err)
        result_set = {
          data: [],
          message: "Insert log failed : " + err,
        };
      } 

      res.send(result_set);
    }
  );
});

app.put("/settings/config/pca/eks", (req, res) => {
  connection.query(
    `update tb_config_eks set 
      "cluster" = '${req.body.cluster}',
      "accessKey" = '${req.body.accessKey}',
      "secretKey" = '${req.body.secretKey}',
      "region" = '${req.body.region}'
    where seq = ${req.body.seq};`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.delete("/settings/config/pca/eks", (req, res) => {
  connection.query(
    `delete from tb_config_eks 
      where "seq" = '${req.body.seq}' and
            "cluster" = '${req.body.cluster}'`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.get("/settings/config/pca/gke", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    // tb_auth_eks > seq,cluster,type,accessKey,secretKey
    `select * from tb_config_gke;`,
    (err, result) => {
      res.send(result.rows);
    }
  );
});

app.post("/settings/config/pca/gke", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  connection.query(
    `insert into tb_config_gke (cluster,"type","clientEmail","projectID","privateKey") values ('${req.body.cluster}','${req.body.type}','${req.body.clientEmail}','${req.body.projectID}','${req.body.privateKey}');`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Insert success",
      };

      if (err !== null) {
        console.log(err)
        result_set = {
          data: [],
          message: "Insert log failed : " + err,
        };
      } 

      res.send(result_set);
    }
  );
});

app.put("/settings/config/pca/gke", (req, res) => {
  connection.query(
    `update tb_config_gke set 
      "cluster" = '${req.body.cluster}',
      "type" = '${req.body.type}',
      "clientEmail" = '${req.body.clientEmail}',
      "projectID" = '${req.body.projectID}',
      "privateKey" = '${req.body.privateKey}'
    where seq = ${req.body.seq};`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.delete("/settings/config/pca/gke", (req, res) => {
  connection.query(
    `delete from tb_config_gke 
      where "seq" = '${req.body.seq}' and
            "cluster" = '${req.body.cluster}'`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

// tb_auth_aks > seq,cluster,clientId,clientSec,tenantId,subId
app.get("/settings/config/pca/aks", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);

  connection.query(
    // tb_auth_eks > seq,cluster,type,accessKey,secretKey
    `select * from tb_config_aks;`,
    (err, result) => {
      res.send(result.rows);
    }
  );
});

app.post("/settings/config/pca/aks", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  connection.query(
    `insert into tb_config_aks (cluster,"clientId","clientSec","tenantId","subId") values ('${req.body.cluster}','${req.body.clientId}','${req.body.clientSec}','${req.body.tenantId}','${req.body.subId}');`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Insert success",
      };

      if (err !== null) {
        console.log(err)
        result_set = {
          data: [],
          message: "Insert log failed : " + err,
        };
      } 

      res.send(result_set);
    }
  );
});


app.put("/settings/config/pca/aks", (req, res) => {
  connection.query(
    `update tb_config_aks set 
      "cluster" = '${req.body.cluster}',
      "clientId" = '${req.body.clientId}',
      "clientSec" = '${req.body.clientSec}',
      "tenantId" = '${req.body.tenantId}',
      "subId" = '${req.body.subId}'
    where seq = ${req.body.seq};`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.delete("/settings/config/pca/aks", (req, res) => {
  console.log("ddd",req.body);
  connection.query(
    `delete from tb_config_aks
      where "seq" = '${req.body.seq}' and
            "cluster" = '${req.body.cluster}'`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.get("/settings/config/pca/kvm", (req, res) => {
  connection.query(
    // tb_auth_kvm > seq,cluster,
    `select * from tb_config_kvm;`,
    (err, result) => {
      res.send(result.rows);
    }
  );
});

app.post("/settings/config/pca/kvm", (req, res) => {
  // let rawdata = fs.readFileSync("./json_data/eks_auth.json");
  // let overview = JSON.parse(rawdata);
  // res.send(overview);
  connection.query(
    `insert into tb_config_kvm (cluster,"agentURL","mClusterName","mClusterPwd") values ('${req.body.cluster}','${req.body.agentURL}','${req.body.mClusterName}','${req.body.mClusterPwd}');`,
    (err, result) => {
      var result_set = {
        data: [],
        message: "Insert success",
      };

      if (err !== null) {
        console.log(err)
        result_set = {
          data: [],
          message: "Insert log failed : " + err,
        };
      } 

      res.send(result_set);
    }
  );
});

app.put("/settings/config/pca/kvm", (req, res) => {
  connection.query(
    `update tb_config_kvm set 
      "cluster" = '${req.body.cluster}',
      "agentURL" = '${req.body.agentURL}',
      "mClusterName" = '${req.body.mClusterName}',
      "mClusterPwd" = '${req.body.mClusterPwd}'
    where seq = ${req.body.seq};`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Update was successful !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.delete("/settings/config/pca/kvm", (req, res) => {
  connection.query(
    `delete from tb_config_kvm 
      where "seq" = '${req.body.seq}' and
            "cluster" = '${req.body.cluster}'`,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

//////////////////////////
// Settings > Alert
//////////////////////////
app.get("/settings/threshold", (req, res) => {
  var create_time = getDateTime();
  connection.query(`select
    ht.node_name,
    ht.cluster_name,
    ht.cpu_warn,
    ht.cpu_danger,
    ht.ram_warn,
    ht.ram_danger,
    ht.storage_warn,
    ht.storage_danger,
    ht.created_time,
    ht.updated_time
    from tb_host_threshold ht
    order by cluster_name, node_name;`, (err, result) => {
    res.send(result.rows);
  });
});

app.post("/settings/threshold", (req, res) => {
  var now = getDateTime();
  let query = `
  INSERT INTO public.tb_host_threshold(
    node_name, cluster_name, cpu_warn, cpu_danger, ram_warn, ram_danger, storage_warn, storage_danger, created_time, updated_time)
    VALUES ('${req.body.nodeName}', '${req.body.clusterName}', ${req.body.cpuWarn}, ${req.body.cpuDanger}, ${req.body.ramWarn}, ${req.body.ramDanger}, ${req.body.storageWarn}, ${req.body.stroageDanger}, '${now}', '${now}');
  `
  console.log(query);
  connection.query(query, 
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Host Threshold is saved !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Save was faild, please check Host Threshold : " + err,
        };
        res.send(result_set);
      }
    });
});

app.put("/settings/threshold", (req, res) => {
  var now = getDateTime();
  let query =  `
  UPDATE public.tb_host_threshold
	SET cpu_warn=${req.body.cpuWarn}, cpu_danger=${req.body.cpuDanger}, ram_warn=${req.body.ramWarn}, ram_danger=${req.body.ramDanger}, storage_warn=${req.body.storageWarn}, storage_danger=${req.body.storageDanger}, updated_time='${now}'	WHERE cluster_name='${req.body.clusterName}' AND node_name='${req.body.nodeName}';
  `

  console.log(query);
  connection.query(query, 
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Host Threshold is updated !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Update was faild, please check threshold : " + err,
        };
        res.send(result_set);
      }
    });
});

app.delete("/settings/threshold", (req, res) => {
  let query = `delete from tb_host_threshold 
  where node_name = '${req.body.node}' and cluster_name = '${req.body.cluster}';`;
  // console.log(query);
  connection.query(
    query,
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Delete was successful!!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Delete was faild, please check error : " + err,
        };
        res.send(result_set);
      }
    }
  );
});

app.get("/settings/threshold/log", (req, res) => {
  var create_time = getDateTime();
  connection.query(`select
    tl.node_name,
    tl.cluster_name,
    tl.created_time,
    tl.status,
    tl.message,
    tl.resource
    from tb_threshold_log tl
    order by created_time desc, node_name;`, (err, result) => {
    res.send(result.rows);
  });
});

app.post("/settings/threshold/log", (req, res) => {
  var now = getDateTime();
  let query = `

  INSERT INTO public.tb_threshold_log(
    cluster_name, node_name, created_time, status, message, resource)
    VALUES ('${req.body.clusterName}', '${req.body.nodeName}', '${now}', '${req.body.status}', '${req.body.message}', '${req.body.resource}');
  `
  console.log(query);
  connection.query(query, 
    (err, result) => {
      if (err !== "null") {
        const result_set = {
          data: [],
          message: "Host Threshold is saved !!",
        };
        res.send(result_set);
      } else {
        const result_set = {
          data: [],
          message: "Save was faild, please check Host Threshold : " + err,
        };
        res.send(result_set);
      }
    });
});

//all nodes metrics
app.get("/apis/nodes_metric", (req, res) => {
  var request = require("request");
  var options = {
    uri: `${apiServer}/apis/nodes_metric`,
    method: "GET",
  };

  request(options, function (error, response, body) {
    if (!error && response.statusCode == 200) {
      res.send(body);
    } else {
      console.log("error", error);
      return error;
    }
  });
});

app.get("/apis/metering", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/metering.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);

  // var request = require("request");
  // var options = {
  //   uri: `${apiServer}/apis/dashboard`,
  //   method: "GET",
  //   // headers: {
  //   //   Authorization:
  //   //     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMxMDQ4NzcsImlhdCI6MTYwMzEwMTI3NywidXNlciI6Im9wZW5tY3AifQ.mgO5hRruyBioZLTJ5a3zwZCkNBD6Bg2T05iZF-eF2RI",
  //   // },
  // };

  // request(options, function (error, response, body) {
  //   if (!error && response.statusCode == 200) {
  //     // console.log("result", body);
  //     res.send(body);
  //   } else {
  //     console.log("error", error);
  //     return error;
  //   }
  // });
});

app.get("/apis/metering/bill", (req, res) => {
  let rawdata = fs.readFileSync("./json_data/metering_bill.json");
  let overview = JSON.parse(rawdata);
  res.send(overview);

  // var request = require("request");
  // var options = {
  //   uri: `${apiServer}/apis/dashboard`,
  //   method: "GET",
  //   // headers: {
  //   //   Authorization:
  //   //     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDMxMDQ4NzcsImlhdCI6MTYwMzEwMTI3NywidXNlciI6Im9wZW5tY3AifQ.mgO5hRruyBioZLTJ5a3zwZCkNBD6Bg2T05iZF-eF2RI",
  //   // },
  // };

  // request(options, function (error, response, body) {
  //   if (!error && response.statusCode == 200) {
  //     // console.log("result", body);
  //     res.send(body);
  //   } else {
  //     console.log("error", error);
  //     return error;
  //   }
  // });
});


app.listen(port, () => console.log(`Listening on port ${port}`));


