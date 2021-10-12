import React, { Component } from "react";
import PieReChart from "../../modules/PieReChart";
import CircularProgress from "@material-ui/core/CircularProgress";
import { Link } from 'react-router-dom';
import { NavigateNext} from '@material-ui/icons';
import TreeView from '../../modules/TreeView';
import TreeView2 from '../../modules/TreeView2';
import * as utilLog from '../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
import { RiDashboardFill } from "react-icons/ri";
import CustomDynamicView from "./CustomDynamicView";
import DashboardSelectModule from "../modal/dashboard/DashboardSelectModule";
// import RefreshButton from './../modules/RefreshButton';
// import PieHalfReChart from './../modules/PieHalfReChart';

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      rows: "",
      completed: 0,
      reRender: "",
      masterCluster: "",
      componentList: []
    };
  }

  componentWillMount() {
    this.props.menuData("none");
  }
  
  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: "" });
        } else {
          console.log(res)
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId = result;
    });
    
    utilLog.fn_insertPLogs(userId, 'log-DS-VW01');
  }


  callApi = async () => {
    const response = await fetch(`/dashboard`);
    const body = await response.json();
    return body;
  };

  

  onRefresh = () => {
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: "" });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };
  angle = {
    full : {
      startAngle : 0,
      endAngle : 360
    },
    half : {
      startAngle : 0,
      endAngle : 180
    }  
  }

  render() {
    // let classNam = 'content-wrapper';
    // console.log(this.state.rows);
    
    
    return (
      <div className="content-wrapper dashboard-main">
        {/* 컨텐츠 헤더 */}
        <section className="content-header">
          <h1>
            <i><RiDashboardFill/></i>
            <span onClick={this.onRefresh} style={{cursor:"pointer"}}>Dashboard</span>
            <span><DashboardSelectModule/></span>
            <small></small>
            
          </h1>
          <ol className="breadcrumb">
            <li>
              <a href="/dashboard">
                <i className="fa fa-dashboard"></i> Home
              </a>
            </li>
            <li className="active">
              <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
              Dashboard
            </li>
          </ol>
        </section>

        {/* 컨텐츠 내용 */}
        <section className="content" style={{ minWidth: 1160 }}>
          {this.state.rows ? (
            [
              // <div onClick={this.onRefresh}><RefreshButton ></RefreshButton></div>
              // ,
              <div style={{ display: "flex" }}>
                <DashboardCard01
                  title="Clusters"
                  width="24%"
                  data={this.state.rows.clusters}
                  path="/clusters"
                  angle={this.angle.full}
                ></DashboardCard01>
                <DashboardCard01
                  title="Nodes"
                  width="24%"
                  data={this.state.rows.nodes}
                  path="/nodes"
                  angle={this.angle.full}
                ></DashboardCard01>
                <DashboardCard01
                  title="Pods"
                  width="24%"
                  data={this.state.rows.pods}
                  path="/pods"
                  angle={this.angle.full}
                ></DashboardCard01>
                <DashboardCard01
                  title="Projects"
                  width="24%"
                  data={this.state.rows.projects}
                  path="/projects"
                  angle={this.angle.full}
                ></DashboardCard01>
              </div>,
              <div style={{ display: "flex" }}>
                <DashboardCard04
                  title="Clusters Groups"
                  width="100%"
                  data={this.state.rows.regions}
                ></DashboardCard04>
              </div>,
              <div style={{ display: "flex" }}>
                <DashboardCard03
                  title="OpenMCP Management Clusters"
                  width="100%"
                  data={this.state.rows.joined_clusters}
                  // width="472px"
                  // data={this.state.masterCluster}
                ></DashboardCard03>
              </div>
            ]
          ) : (
            <CircularProgress
              variant="determinate"
              value={this.state.completed}
              style={{ position: "absolute", left: "50%", marginTop: "20px" }}
            ></CircularProgress>
          )}
          
        </section>
        <section className="content" style={{ minWidth: 1160 }}>
          <CustomDynamicView componentList={this.state.componentList}/>
        </section>
      </div>
    );
  }
}

class DashboardCard01 extends Component {
  onClickRf = () =>{
    this.props.onClickRefresh()
  }
  render() {
    const colors = [
      "#0088FE",
      // "#00C49F",
      "#ff8042",
      "#FFBB28",
      "#cccccc",
      "#00C49F",
      "#FFBB28",
      "#00C49F",
      "#FFBB28",
    ];
    return (
      <div className="content-box" style={{ width: this.props.width }}>
        <div className="cb-header">
          <div onClick={this.onClickRf}></div>
          <span>{this.props.title}</span>
          <span> : {this.props.data.counts}</span>
          
          <Link to={this.props.path}>
            <div className="cb-btn">
              detail
            </div>
          </Link>
          
        </div>
        <div
          className="cb-body"
          style={{ position: "relative", width: "100%" }}
        >
          <PieReChart data={this.props.data} angle={this.props.angle} colors={colors} unit=""></PieReChart>
        </div>
      </div>
    );
  }
}

// class DashboardCard02 extends Component {
//   render() {
//     return (
//       <div className="content-box" style={{ width: this.props.width }}>
//       <div className="cb-header">
//         <span>{this.props.title}</span>
//         {/* <div className="cb-btn">
//           <Link to={this.props.path}>detail</Link>
//         </div> */}
//       </div>
//       <div
//         className="cb-body"
//         style={{ position: "relative", width: "100%", display:"flex"}}
//       >
//         <PieHalfReChart data={this.props.data.cpu} angle={this.props.angle}></PieHalfReChart>
//         <PieHalfReChart data={this.props.data.memory} angle={this.props.angle}></PieHalfReChart>
//         <PieHalfReChart data={this.props.data.storage} angle={this.props.angle}></PieHalfReChart>
//       </div>
//     </div>
//     );
//   }
// }


class DashboardCard03 extends Component {
  // constructor(props){
  //   super(props)
  //   this.state = {
  //     rows:""
  //   }
  // }
  

  // componentDidMount() {
  //   //데이터가 들어오기 전까지 프로그래스바를 보여준다.
  //   this.timer = setInterval(this.progress, 20);
  //   this.callApi()
  //     .then((res) => {
  //       if(res == null){
  //         this.setState({ rows: [] });
  //       } else {
  //         console.log(res)
  //         this.setState({ rows: res });
  //       }
  //       clearInterval(this.timer);
  //     })
  //     .catch((err) => console.log(err));
  // }
  
  // callApi = async () => {
  //   const response = await fetch(`/dashboard-master-cluster`);
  //   const body = await response.json();
  //   return body;
  // };

  // onRefresh = () => {
  //   this.callApi()
  //     .then((res) => {
  //       if(res == null){
  //         this.setState({ rows: [] });
  //       } else {
  //         this.setState({ rows: res });
  //       }
  //       clearInterval(this.timer);
  //     })
  //     .catch((err) => console.log(err));
  // };
  

  render() {
    return (
      <div className="content-box" style={{ width: this.props.width }}>
      {/* <div className="cb-header" onClick={this.onRefresh} style={{paddingTop: "140px", fontSize:"24px"}}> */}
      <div className="cb-header" onClick={this.onRefresh}>
        <span style={{cursor:"pointer"}} >{this.props.title}</span>
        {/* <div className="cb-btn">
          <Link to={this.props.path}>detail</Link>
        </div> */}
      </div>
      <div
        className="cb-body"
        style={{ position: "relative", width: "100%", display:"flex"}}
        >
          <TreeView data={this.props.data}/>
          {/* <TreeView data=[{this.props.data}]/> */}
        {/* {this.state.rows ? (
        <TreeView data={this.state.rows}/>
        ) : (
          <CircularProgress
            variant="determinate"
            value={this.state.completed}
            style={{ position: "absolute", left: "50%", marginTop: "20px" }}
          ></CircularProgress>
        )} */}
      </div>
    </div>
    
    );
  }
}

class DashboardCard04 extends Component {
  render() {
    // console.log("BasicInfo:", this.props.data)

    return (
      <div className="content-box" style={{ width: this.props.width }}>
      <div className="cb-header">
        <span>{this.props.title}</span>
        {/* <div className="cb-btn">
          <Link to={this.props.path}>detail</Link>
        </div> */}
      </div>
      <div
        className="cb-body"
        style={{ position: "relative", width: "100%", display:"flex"}}
      >
        <TreeView2 data={this.props.data}/>
      </div>
    </div>
    );
  }
}

export default Dashboard;
