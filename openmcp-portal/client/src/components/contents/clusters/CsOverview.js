import React, { Component } from "react";
import { NavLink} from 'react-router-dom';
import CircularProgress from "@material-ui/core/CircularProgress";
import { NavigateNext} from '@material-ui/icons';
import Paper from "@material-ui/core/Paper";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import SelectBox from '../../modules/SelectBox';
import PieReChart2 from '../../modules/PieReChart2';
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
// import LineChart from './../../modules/LineChart';
// import PieHalfReChart from './../../modules/PieHalfReChart';
// import PieReChart from './../../modules/PieReChart';
// import line_chart_sample from './../../../json/line_chart_sample.json'
import ChangeAKSReource from '../modal/chageResource/ChangeAKSReource';

let apiParams = "";
class CsOverview extends Component {
  state = {
    rows:"",
    completed: 0,
    reRender : ""
  }

  componentWillMount() {
    const result = {
      menu : "clusters",
      title : this.props.match.params.cluster,
      pathParams : {
        cluster : this.props.match.params.cluster,
        // state : this.props.location.state.data
      },
    }
    this.props.menuData(result);
    apiParams = this.props.match.params.cluster;
  }

  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, 'log-CL-VW02');
  }  

  callApi = async () => {
    var param = this.props.match.params.cluster;
    const response = await fetch(`/clusters/${param}/overview`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };


  onRefresh = () => {
    console.log("onClick")
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  };

  render() {
    return (
      <div>
        <div className="content-wrapper cluster-overview">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
            Overview
              <small>{this.props.match.params.cluster}</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                <NavLink to="/clusters">Clusters</NavLink>
              </li>
              <li className="active">
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                Overview
              </li>
            </ol>
          </section>

          {/* 내용부분 */}
          <section className="content">
          {this.state.rows ? (
            [
            <BasicInfo rowData={this.state.rows.basic_info}
            //  provider={this.props.location.state.data == undefined ? "-" : this.props.location.state.data.provider }
            //  region={this.props.location.state.data == undefined ? "-" : this.props.location.state.data.region}
             />,
            <div style={{display:"flex"}}>
              <ProjectUsageTop5 rowData={this.state.rows.project_usage_top5}/>
              <NodeUsageTop5 rowData={this.state.rows.node_usage_top5}/>
            </div>,
            <ClusterResourceUsage rowData={this.state.rows.cluster_resource_usage} onRefresh={this.onRefresh} 
            clusterInfo = {this.state.rows.basic_info}/>,
            <KubernetesStatus rowData={this.state.rows.kubernetes_status}/>,
            <Events rowData={this.state.rows.events}/>
            ]
          ) : (
            <CircularProgress
              variant="determinate"
              value={this.state.completed}
              style={{ position: "absolute", left: "50%", marginTop: "20px" }}
            ></CircularProgress>
          )}
          </section>
        </div>
      </div>
    );
  }
}

class BasicInfo extends Component {
  render(){
    return (
      <div className="content-box">
        <div className="cb-header">Basic Info</div>
        <div className="cb-body">
          <div>
            <span>Name : </span>
            <strong>{this.props.rowData.name}</strong>
          </div>
          <div>
            <span>Provider : </span>
            {this.props.rowData.provider}
          </div>
          <div>
            <span>Kubernetes Version : </span>
            {this.props.rowData.kubernetes_version}
          </div>
          <div>
            <span>Status : </span>
            {this.props.rowData.status}
          </div>
          <div>
            <span>Region : </span>
            {this.props.rowData.region}
          </div>
          <div>
            <span>Zone : </span>
            {this.props.rowData.zone}
          </div>
        </div>
      </div>
    );
  }
}

class ProjectUsageTop5 extends Component {
  state = {
    columns: [
      { name: "name", title: "Name" },
      { name: "usage", title: "Usage" },
    ],
    rows : this.props.rowData.cpu,
  }

  callApi = async () => {
    const response = await fetch(`/clusters/${apiParams}/overview`);
    const body = await response.json();
    return body;
  };
  
  render(){
    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
          // ...styles[row.sector.toLowerCase()],
        }}
        // onClick={()=> alert(JSON.stringify(row))}
      />
    );

    const onSelectBoxChange = (data) => {
      console.log("onSelectBoxChange", data)
      switch(data){
        case "cpu":
          console.log("cpu")
          // this.setState({rows:this.props.rowData.cpu});

          this.callApi()
          .then((res) => {
            this.setState({ rows: res.project_usage_top5.cpu });
          })
          .catch((err) => console.log(err));

          break;
        case "memory":
          console.log("memory")
          // this.setState({rows:this.props.rowData.memory});

          this.callApi()
          .then((res) => {
            this.setState({ rows: res.project_usage_top5.memory });
          })
          .catch((err) => console.log(err));

          break;
        default:
          this.setState(this.props.rowData.cpu);
      }
    }

    const selectBoxData = [{name:"cpu", value:"cpu"},{name:"memory", value:"memory"}];
    return (
      <div className="content-box col-sep-2">
        <div className="cb-header">
          Project Usage Top5
          <SelectBox rows={selectBoxData} onSelectBoxChange={onSelectBoxChange} defaultValue=""></SelectBox>
        </div>
        
        <div className="cb-body table-style">
          {this.state.aaa}
          <Grid
            rows = {this.state.rows}
            columns = {this.state.columns}>

            {/* Sorting */}
            <SortingState
            defaultSorting={[{ columnName: 'usage', direction: 'desc' }]}
            />
            <IntegratedSorting />

            <Table/>
            <TableHeaderRow showSortingControls rowComponent={HeaderRow}/>
          </Grid>
        </div>
      </div>
    );
  }
}

class NodeUsageTop5 extends Component {
  state = {
    columns: [
      { name: "name", title: "Name" },
      { name: "usage", title: "Usage" },
    ],
    rows : this.props.rowData.cpu,
    unit : " core",
  }

  callApi = async () => {
    const response = await fetch(`/clusters/${apiParams}/overview`);
    const body = await response.json();
    return body;
  };
  
  render(){
    const Cell = (props) => {
      const { column } = props;
      // console.log("cell : ", props);
      if (column.name === "usage") {
        return (
          <Table.Cell {...props} style={{ cursor: "pointer" }}>
            {props.value + this.state.unit}
          </Table.Cell>
        );
      }
      return <Table.Cell {...props} />;
    };

    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
          // ...styles[row.sector.toLowerCase()],
        }}
        // onClick={()=> alert(JSON.stringify(row))}
      />
    );

    const onSelectBoxChange = (data) => {
      switch(data){
        case "cpu":
          this.setState({
            rows:this.props.rowData.cpu,
            unit : " core"
          });

          // this.callApi()
          // .then((res) => {
          //   this.setState({ rows: res.node_usage_top5.cpu });
          // })
          // .catch((err) => console.log(err));

          break;
        case "memory":
          this.setState({
            rows:this.props.rowData.memory,
            unit : " Gi"
          });

          // this.callApi()
          // .then((res) => {
          //   this.setState({ rows: res.node_usage_top5.memory });
          // })
          // .catch((err) => console.log(err));

          break;
        default:
          this.setState(this.props.rowData.cpu);
      }
    }

    const selectBoxData = [{name:"cpu", value:"cpu"},{name:"memory", value:"memory"}];
    return (
      <div className="content-box col-sep-2">
        <div className="cb-header">
          Node Usage Top5
          <SelectBox rows={selectBoxData} onSelectBoxChange={onSelectBoxChange}
          defaultValue=""></SelectBox>
        </div>
        
        <div className="cb-body table-style">
          {this.state.aaa}
          <Grid
            rows = {this.state.rows}
            columns = {this.state.columns}>

            {/* Sorting */}
            <SortingState
              defaultSorting={[{ columnName: 'usage', direction: 'desc' }]}
            />
            <IntegratedSorting />

            <Table cellComponent={Cell}/>
            <TableHeaderRow showSortingControls rowComponent={HeaderRow}/>
          </Grid>
        </div>
      </div>
    );
  }
}





// 갱신전

let normal = {
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 1.2 },
        { name: "Total", value: 72 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 42.3 },
        { name: "Total", value: 197.3 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 71.3 },
        { name: "Total", value: 1724.6 }
      ]
    }
  }
  

  let stress50 = {
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 35.3 },
        { name: "Total", value: 72 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 95.3 },
        { name: "Total", value: 197.3 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 71.3 },
        { name: "Total", value: 1724.6 }
      ]
    }
  }

    
  let stress70 = {
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 55.1 },
        { name: "Total", value: 72 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 148.1 },
        { name: "Total", value: 197.3 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 71.3 },
        { name: "Total", value: 1724.6 }
      ]
    }
  }

    

  let stress80 = {
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 63.2 },
        { name: "Total", value: 72 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 175.4 },
        { name: "Total", value: 197.3 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 71.3 },
        { name: "Total", value: 1724.6 }
      ]
    }
  }

  let stress0 = {
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 52.2 },
        { name: "Total", value: 72 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 101.4 },
        { name: "Total", value: 197.3 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 71.3 },
        { name: "Total", value: 1724.6 }
      ]
    }
  }

  let eks ={
    cpu: {
      counts: 3,
      unit: "core",
      status: [
        { name: "Used", value: 3.5 },
        { name: "Total", value: 6 }
      ]
    },
    memory: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 13.4 },
        { name: "Total", value: 24.2 }
      ]
    },
    storage: {
      counts: 3,
      unit: "Gi",
      status: [
        { name: "Used", value: 31.4 },
        { name: "Total", value: 75.0 }
      ]
    }
  }
    
    let gke2 ={
      cpu: {
        counts: 3,
        unit: "core",
        status: [
          { name: "Used", value: 1.6 },
          { name: "Total", value: 8 }
        ]
      },
      memory: {
        counts: 3,
        unit: "Gi",
        status: [
          { name: "Used", value: 12.5 },
          { name: "Total", value: 24.2 }
        ]
      },
      storage: {
        counts: 3,
        unit: "Gi",
        status: [
          { name: "Used", value: 19.3 },
          { name: "Total", value: 197.0 }
        ]
      }
}

  let colors3 = [
    "#0088FE", //파랑
    "#ecf0f5",
  ];

  let colors2 = [
    "#ff8042",// 주황
    "#ecf0f5",
  ];

  let count = 1;
class ClusterResourceUsage extends Component {
  state = {
    rows : this.props.rowData,
    colors : [
      "#0088FE",
      "#ecf0f5"
    ],
    unhColors : [
      "#0088FE",
      "#ecf0f5"
    ]
  }

  
  angle = {
    full : {
      startAngle : 0,
      endAngle : 360
    },
    half : {
      startAngle : 180,
      endAngle : 0
    }  
  }

  componentWillUpdate(prevProps, prevState){
    if (this.props.rowData !== prevProps.rowData) {
        this.setState({
          rows: prevProps.rowData,
        });
      }
  }

  onClick = () => {
    if(count === 1){
      this.setState({
        rows:stress50,
        unhColors:colors3
      })
    } else if (count === 2){
      this.setState({
        rows:stress70, 
        unhColors:colors3
      })
    } else if (count === 3 ){
      this.setState({
        rows:stress80, 
        unhColors:colors2
      })
    }  else if (count === 4 ){
      this.setState({
        rows:stress0, 
        unhColors:colors3
      })
    }  else {
      count = 1
      this.setState({
        rows:normal, 
        unhColors:colors3
      })
    }
    count++
    
    // this.props.onRefresh();
  }

  onSelectBoxChange = (data) => {
    count=1
    // console.log(eks-cluster1)
    if(data === "cluster1"){
      this.setState({
        rows:normal,
        unhColors:colors3
      })} else if(data === "eks-cluster1") {
        this.setState({
          rows:eks,
          unhColors:colors3
        })
      } else {
      this.setState({
        rows:gke2,
        unhColors:colors3
      })
    }
    // } 
    // else if(data === "eks-cluster1"){
    //   this.setState({
    //     rows:warning, //addAfter
    //     unhColors:colors2
    //   })
    // } else if(data === "cluster1"){
    //   this.setState({
    //     rows:normal,
    //     unhColors:colors3
    //   })
    // } else if(data === "gke-cluster1"){
    //   this.setState({
    //     rows:gke,
    //     unhColors:colors3
    //   })
    // }else {
    //   this.setState({
    //     rows:gke,
    //     unhColors:colors3
    //   })
    // }


    // this.setState({
    //   rows:warning,
    //   unhColors:colors2
    // })

    // this.setState({
    //   rows:healthy,
    //   unhColors:colors3
    // })
    
  }

  
  
  render(){
    // const colors = [
    //   "#0088FE",
    //   "#ecf0f5",
    // ];

    // const colors2 = [
    //   "#ff8042", //#0088FE, #ff8042;
    //   "#ecf0f5",
    // ];

    // const selectBoxData = [
    //   {name:"cluster1", value:"cluster1"},
    //   {name:"cluster2", value:"cluster2"},
    //   {name:"cluster3", value:"cluster2"},
    //   {name:"eks-cluster1", value:"eks-cluster1"},
    //   {name:"gke-cluster1", value:"gke-cluster1"},
    // ];

    

    return (
      <div className="content-box">
        <div className="cb-header">
          
          <span style={{cursor:"cluster1"}} onClick={this.onClick} >
            Cluster Resource Usage
          </span>
          {this.props.clusterInfo.provider === "aks" ? 
          <ChangeAKSReource clusterInfo={this.props.clusterInfo} /> : "" }
        </div>
        <div className="cb-body flex">
          <div className="cb-body-content pie-chart">
            <div className="cb-sub-title">CPU</div>
            <PieReChart2 data={this.state.rows.cpu} angle={this.angle.half} unit={this.state.rows.cpu.unit} colors={this.state.unhColors}></PieReChart2>
          </div>
          <div className="cb-body-content pie-chart">
            <div className="cb-sub-title">Memory</div>
            <PieReChart2 data={this.state.rows.memory} angle={this.angle.half} unit={this.state.rows.memory.unit} colors={this.state.unhColors}></PieReChart2>
          </div>
          <div className="cb-body-content pie-chart">
            <div className="cb-sub-title">Storage</div>
            <PieReChart2 data={this.state.rows.storage} angle={this.angle.half} unit={this.state.rows.storage.unit} colors={this.state.colors}></PieReChart2>
          </div>
        </div>
      </div>
    );
  }
}

class KubernetesStatus extends Component {
  state = {
    rows : this.props.rowData
  }
  render(){
    
    return(
      <div className="content-box cb-kube-status">
        <div className="cb-header">Kubernetes Status</div>
        <div className="cb-body flex">
          {this.state.rows.map((item)=>{
            return (
          <div className={"cb-body-content "+item.status}>
            <div>{item.name}</div>
            <div>({item.status})</div>
          </div>)
          })}
        </div>
      </div>
    );
  };
};

class Events extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "project", title: "Project" },
        { name: "type", title: "Type" },
        { name: "reason", title: "Reason" },
        { name: "object", title: "Object" },
        { name: "message", title: "Message" },
        { name: "time", title: "Time" },
      ],
      defaultColumnWidths: [
        { columnName: "project", width: 150 },
        { columnName: "type", width: 150 },
        { columnName: "reason", width: 150 },
        { columnName: "object", width: 240 },
        { columnName: "message", width: 440 },
        { columnName: "time", width: 180 },
      ],
      rows: this.props.rowData,

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10, 
      pageSizes: [5, 10, 15, 0],

      completed: 0,
    };
  }

  componentWillMount() {
    // this.props.onSelectMenu(false, "");
  }

  

  // callApi = async () => {
  //   const response = await fetch("/clusters");
  //   const body = await response.json();
  //   return body;
  // };

  // progress = () => {
  //   const { completed } = this.state;
  //   this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  // };

  // //컴포넌트가 모두 마운트가 되었을때 실행된다.
  // componentDidMount() {
  //   //데이터가 들어오기 전까지 프로그래스바를 보여준다.
  //   this.timer = setInterval(this.progress, 20);
  //   this.callApi()
  //     .then((res) => {
  //       this.setState({ rows: res });
  //       clearInterval(this.timer);
  //     })
  //     .catch((err) => console.log(err));
  // };

  render() {
    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
          // ...styles[row.sector.toLowerCase()],
        }}
        // onClick={()=> alert(JSON.stringify(row))}
      />
    );
    const Row = (props) => {
      // console.log("row!!!!!! : ",props);
      return <Table.Row {...props} key={props.tableRow.key}/>;
    };

    return (
      <div className="content-box">
        <div className="cb-header">Events</div>
        <div className="cb-body">
        <Paper>
            {this.state.rows ? (
              [
                <Grid
                  rows={this.state.rows}
                  columns={this.state.columns}
                >
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState defaultValue="" />
                  <IntegratedFiltering />
                  <SearchPanel style={{ marginLeft: 0 }} />

                  {/* Sorting */}
                  <SortingState
                    defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
                  />
                  <IntegratedSorting />

                  {/* 페이징 */}
                  <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* 테이블 */}
                  <Table rowComponent={Row} />
                  <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                </Grid>,
              ]
            ) : (
              <CircularProgress
                variant="determinate"
                value={this.state.completed}
                style={{ position: "absolute", left: "50%", marginTop: "20px" }}
              ></CircularProgress>
            )}
          </Paper>
        </div>
      </div>
    );
  };
};

export default CsOverview;

