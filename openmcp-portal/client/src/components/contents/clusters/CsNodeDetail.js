import React, { Component } from "react";
// import { NavLink} from 'react-router-dom';
// import CircularProgress from "@material-ui/core/CircularProgress";
// import { NavigateNext} from '@material-ui/icons';
// import Paper from "@material-ui/core/Paper";
// import {
//   SearchState,
//   IntegratedFiltering,
//   PagingState,
//   IntegratedPaging,
//   SortingState,
//   IntegratedSorting,
// } from "@devexpress/dx-react-grid";
// import {
//   Grid,
//   Table,
//   Toolbar,
//   SearchPanel,
//   TableColumnResizing,
//   TableHeaderRow,
//   PagingPanel,
// } from "@devexpress/dx-react-grid-material-ui";
// import PieReChart2 from '../../modules/PieReChart2';
// import NdTaintConfig from './../modal/NdTaintConfig';
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
// import NdResourceConfig from './../modal/NdResourceConfig';
// import Confirm2 from './../../modules/Confirm2';
// import Button from "@material-ui/core/Button";
// import ProgressTemp from './../../modules/ProgressTemp';
// import axios from "axios";
import NdNodeDetail from "../nodes/NdNodeDetail";


class CsNodeDetail extends Component {
  constructor(props){
    super(props);
    this.state = {
      rows:"",
      completed: 0,
      reRender : "",
      propsRow : ""
    }
  }

  componentWillMount() {
    const result = {
      menu : "clusters",
      title : this.props.match.params.cluster,
      pathParams : {
        cluster : this.props.match.params.cluster
      }
    }
    this.props.menuData(result);
    if(this.props.location.state !== undefined){
      this.setState({propsRow:this.props.location.state.data})
    }
    // apiParams = this.props.match.params.cluster;
  }

  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        console.log(res);
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
    utilLog.fn_insertPLogs(userId, 'log-ND-VW02');
  }  

  callApi = async () => {
    var param = this.props.match.params;
    const response = await fetch(`/nodes/${param.node}${this.props.location.search}`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  onUpdateData = () => {
    console.log("onUpdateData={this.props.onUpdateData}")
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
    // console.log("CsOverview_Render : ",this.state.rows.basic_info);
    return (
      <div>
        
        <NdNodeDetail match={this.props.match} location={this.props.location} menuData={this.props.menuData}/>
      </div>
    );
  }
}

// class BasicInfo extends Component {
//   render(){
//     return (
//       <div className="content-box">
//         <div className="cb-header">
//           <span>
//             Basic Info
//           </span>
//             <NdTaintConfig name={this.props.rowData.name} taint={this.props.rowData.taint}/>
//         </div>
//         <div className="cb-body">
//           <div>
//             <span>Name : </span>
//             <strong>{this.props.rowData.name}</strong>
//           </div>
//           <div style={{display:"flex"}}>

//             <div className="cb-body-left">
//               <div>
//                 <span>Status : </span>
//                 <span
//                   style={{
//                     color:
//                     this.props.rowData.status === "Healthy" ? "#1ab726"
//                       : this.props.rowData.status === "Unhealthy" ? "red"
//                         : this.props.rowData.status === "Unknown" ? "#b5b5b5"
//                           : this.props.rowData.status === "Warning" ? "#ff8042" : "black",
//                   }}
//                 >{this.props.rowData.status}</span>
//               </div>
//               <div>
//                 <span>Role : </span>
//                 {this.props.rowData.role}
//               </div>
//               <div>
//                 <span>Cluster : </span>
//                 {this.props.rowData.cluster}
//               </div>
//               <div>
//                 <span>Kubernetes : </span>
//                 {this.props.rowData.kubernetes}
//               </div>
//               <div>
//                 <span>Kubernetes Proxy : </span>
//                 {this.props.rowData.kubernetes_proxy}
//               </div>
//             </div>
//             <div className="cb-body-right">
//               <div>
//                   <span>IP : </span>
//                   {this.props.rowData.ip}
//                 </div>
//                 <div>
//                   <span>OS : </span>
//                   {this.props.rowData.os}
//                 </div>
//                 <div>
//                   <span>Docker : </span>
//                   {this.props.rowData.docker}
//                 </div>
//                 <div>
//                   <span>Created Time : </span>
//                   {this.props.rowData.created_time}
//                 </div>
//                 <div>
//                   <span>Provider : </span>
//                   {this.props.propsRow.provider}
//                 </div>
//             </div>
//           </div>
          
//         </div>
//       </div>
//     );
//   }
// }

// class NodeResourceUsage extends Component {
//   state = {
//     rows : this.props.rowData,
//     nodeData : this.props.nodeData
//   }
//   angle = {
//     full : {
//       startAngle : 0,
//       endAngle : 360
//     },
//     half : {
//       startAngle : 180,
//       endAngle : 0
//     }  
//   }
//   render(){
//     const colors = [
//       "#0088FE",
//       "#ecf0f5",
//     ];
//     return (
//       <div className="content-box">
//         <div className="cb-header">
//         <span>
//           Node Resource Usage
//           </span>
//           {this.props.propsRow.provider === "eks" || this.props.propsRow.provider === "kvm" ? 
//             <NdResourceConfig 
//               rows = {this.state.rows}
//               rowData={this.props.rowData}
//               nodeData={this.props.nodeData}
//               onUpdateData={this.props.onUpdateData}
//               propsRow={this.props.propsRow}/> : ""}
//         </div>
//         <div className="cb-body flex">
//           <div className="cb-body-content pie-chart">
//             <div className="cb-sub-title">CPU</div>
//             <PieReChart2 data={this.state.rows.cpu} angle={this.angle.half} unit={this.state.rows.cpu.unit} colors={colors}></PieReChart2>
//           </div>
//           <div className="cb-body-content pie-chart">
//             <div className="cb-sub-title">Memory</div>
//             <PieReChart2 data={this.state.rows.memory} angle={this.angle.half} unit={this.state.rows.memory.unit} colors={colors}></PieReChart2>
//           </div>
//           <div className="cb-body-content pie-chart">
//             <div className="cb-sub-title">Storage</div>
//             <PieReChart2 data={this.state.rows.storage} angle={this.angle.half} unit={this.state.rows.storage.unit} colors={colors}></PieReChart2>
//           </div>
//           <div className="cb-body-content pie-chart">
//             <div className="cb-sub-title">Pod</div>
//             <PieReChart2 data={this.state.rows.pods} angle={this.angle.half} unit={this.state.rows.pods.unit} colors={colors}></PieReChart2>
//           </div>
//         </div>
//       </div>
//     );
//   }
// }

// class KubernetesStatus extends Component {
//   constructor(props){
//     super(props);
//     this.state={
//       rows : this.props.rowData,
//       confirmType : "",
//       confirmOpen: false,
//       confirmInfo : {
//         title :"Confirm Stop Node",
//         context :"Are you sure you want to stop Node?",
//         button : {
//           open : "",
//           yes : "CONFIRM",
//           no : "CANCEL",
//         }
//       },
//       confrimTarget : "",
//       confirmTargetKeyname:"",
//       powerflag:"on",
//     }
//   }

//   handleClickStart = () => {
//     this.setState({
//       confirmType: "power",
//       confirmOpen: true,
//       powerFlag : "on",
//       confirmInfo : {
//         title :"Confirm Start Node",
//         context :"Are you sure you want to Start Node?",
//         button : {
//           open : "",
//           yes : "CONFIRM",
//           no : "CANCEL",
//         }
//       }
//     })
//   }

//   handleClickStop = () => {
//     this.setState({
//       confirmType: "power",
//       confirmOpen: true,
//       powerFlag : "off",
//       confirmInfo : {
//         title :"Confirm Stop Node",
//         context :"Are you sure you want to stop Node?",
//         button : {
//           open : "",
//           yes : "CONFIRM",
//           no : "CANCEL",
//         }
//       }
//     })
//   }

//   handleClickDelete = () => {
//     this.setState({
//       confirmType: "delete",
//       confirmOpen: true,
//       confirmInfo : {
//         title :"Confirm Delete KVM VM",
//         context :"Are you sure you want to delete vm(node)?",
//         button : {
//           open : "",
//           yes : "CONFIRM",
//           no : "CANCEL",
//         }
//       }
//     })
//   }

//   //callback
//   confirmed = (result) => {
//     this.setState({confirmOpen:false});

//     //show progress loading...
//     this.setState({openProgress:true});
//     const provider = this.props.propsRow.provider;
//     let data = {};
//     let url = "";

//     if(result) {
//       if(this.state.confirmType === "power"){
//         if(this.state.powerFlag === "on"){
//           console.log("poweron")
//           url = `/nodes/${provider}/start`;
//           // utilLog.fn_insertPLogs(userId, "log-ND-PO01"); //poweron log
//         } else if (this.state.powerFlag === "off"){
//           console.log("poweroff")
//           url = `/nodes/${provider}/stop`;
//           // utilLog.fn_insertPLogs(userId, "log-ND-PO02"); //poweroff log
//         }
       
//         if(provider === "eks"){
//           //eks
//           data = {
//             region: this.props.propsRow.region,
//             node: this.props.propsRow.name,
//             cluster : this.props.propsRow.cluster
//           };
//         } else if (provider === "aks"){
//           data = {
//             cluster : this.props.propsRow.cluster,
//             node : this.props.propsRow.name,
//           };
//         } else if (provider === "kvm"){
//           data = {
//             cluster : this.props.propsRow.cluster,
//             node : this.props.propsRow.name,
//           };
//         } else {
//           alert(provider + " is not supported Type");
//           this.setState({openProgress:false})
//           return;
//         }
//       } else if (this.state.confirmType === "delete"){
//         url = `/nodes/delete/kvm`;
//         data = {
//           cluster : this.props.propsRow.cluster,
//           node: this.props.propsRow.name,
//         };
//       }
      

//       axios.post(url, data)
//       .then((res) => {
//         if(res.data.error){
//           alert(res.data.message)
//         }
//       })
//       .catch((err) => {
//           console.log(err)
//       });

//       this.setState({openProgress:false})
//     } else {
//       this.setState({openProgress:false})
//     }
//   }

//   render(){
//     return(
//       <div className="content-box cb-kube-status">
//         {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""}
//         <Confirm2
//           confirmInfo={this.state.confirmInfo} 
//           confrimTarget ={this.state.confrimTarget} 
//           confirmTargetKeyname = {this.state.confirmTargetKeyname}
//           confirmed={this.confirmed}
//           confirmOpen={this.state.confirmOpen}/>

//         <div className="cb-header">
//           <span>Kubernetes Node Status</span>
//           {this.props.propsRow.provider !== "gke" ?
//             <div style={{position:"absolute", top:"0px", right:"0px"}}>
//               <Button variant="outlined" color="primary" onClick={this.handleClickStart} style={{marginRight:"10px",zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
//                 Start Node
//               </Button>

//               <Button variant="outlined" color="primary" onClick={this.handleClickStop} style={{zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
//                 Stop Node
//               </Button>

//               {this.props.propsRow.provider === "kvm" ? 
//                 <Button variant="outlined" color="primary" onClick={this.handleClickDelete} style={{marginLeft:"10px", zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
//                   Delete Node
//                 </Button> : ""
//               }
//             </div> : ""
//           }
//         </div>
//         <div className="cb-body flex">
//           {this.state.rows.map((item)=>{
//             return (
//           <div className={"cb-body-content "+item.status}>
//             <div>{item.name}</div>
//             <div>({item.status})</div>
//           </div>)
//           })}
//         </div>
//       </div>
//     );
//   };
// };

// class Events extends Component {
//   constructor(props) {
//     super(props);
//     this.state = {
//       columns: [
//         { name: "project", title: "Project" },
//         { name: "type", title: "Type" },
//         { name: "reason", title: "Reason" },
//         { name: "object", title: "Object" },
//         { name: "message", title: "Message" },
//         { name: "time", title: "Time" },
//       ],
//       defaultColumnWidths: [
//         { columnName: "project", width: 150 },
//         { columnName: "type", width: 150 },
//         { columnName: "reason", width: 150 },
//         { columnName: "object", width: 240 },
//         { columnName: "message", width: 440 },
//         { columnName: "time", width: 180 },
//       ],
//       rows: this.props.rowData,

//       // Paging Settings
//       currentPage: 0,
//       setCurrentPage: 0,
//       pageSize: 10, 
//       pageSizes: [5, 10, 15, 0],

//       completed: 0,
//     };
//   }

//   componentWillMount() {
//     // this.props.onSelectMenu(false, "");
//   }

  

//   // callApi = async () => {
//   //   const response = await fetch("/clusters");
//   //   const body = await response.json();
//   //   return body;
//   // };

//   // progress = () => {
//   //   const { completed } = this.state;
//   //   this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
//   // };

//   // //컴포넌트가 모두 마운트가 되었을때 실행된다.
//   // componentDidMount() {
//   //   //데이터가 들어오기 전까지 프로그래스바를 보여준다.
//   //   this.timer = setInterval(this.progress, 20);
//   //   this.callApi()
//   //     .then((res) => {
//   //       this.setState({ rows: res });
//   //       clearInterval(this.timer);
//   //     })
//   //     .catch((err) => console.log(err));
//   // };

//   render() {
//     const HeaderRow = ({ row, ...restProps }) => (
//       <Table.Row
//         {...restProps}
//         style={{
//           cursor: "pointer",
//           backgroundColor: "whitesmoke",
//           // ...styles[row.sector.toLowerCase()],
//         }}
//         // onClick={()=> alert(JSON.stringify(row))}
//       />
//     );
//     const Row = (props) => {
//       // console.log("row!!!!!! : ",props);
//       return <Table.Row {...props} key={props.tableRow.key}/>;
//     };

//     return (
//       <div className="content-box">
//         <div className="cb-header">Events</div>
//         <div className="cb-body">
//         <Paper>
//             {this.state.rows ? (
//               [
//                 <Grid
//                   rows={this.state.rows}
//                   columns={this.state.columns}
//                 >
//                   <Toolbar />
//                   {/* 검색 */}
//                   <SearchState defaultValue="" />
//                   <IntegratedFiltering />
//                   <SearchPanel style={{ marginLeft: 0 }} />

//                   {/* Sorting */}
//                   <SortingState
//                     defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
//                   />
//                   <IntegratedSorting />

//                   {/* 페이징 */}
//                   <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
//                   <IntegratedPaging />
//                   <PagingPanel pageSizes={this.state.pageSizes} />

//                   {/* 테이블 */}
//                   <Table rowComponent={Row} />
//                   <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
//                   <TableHeaderRow
//                     showSortingControls
//                     rowComponent={HeaderRow}
//                   />
//                 </Grid>,
//               ]
//             ) : (
//               <CircularProgress
//                 variant="determinate"
//                 value={this.state.completed}
//                 style={{ position: "absolute", left: "50%", marginTop: "20px" }}
//               ></CircularProgress>
//             )}
//           </Paper>
//         </div>
//       </div>
//     );
//   };
// };
export default CsNodeDetail;