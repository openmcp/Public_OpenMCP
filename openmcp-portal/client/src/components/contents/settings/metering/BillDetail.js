import React, { Component } from "react";
import { NavLink } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import { NavigateNext } from "@material-ui/icons";

class BillDetail extends Component {
  constructor(props) {
    super(props);
    this.state = {
      rows: "",
      completed: 0,
      reRender: "",
      propsRow: "",
    };
  }

  // componentWillMount() {
  //   // this.props.menuData("none");
  //   if(this.props.location.state !== undefined){
  //     this.setState({propsRow:this.props.location.state.data})
  //   }
  // }

  // componentDidMount() {
  //   //데이터가 들어오기 전까지 프로그래스바를 보여준다.
  //   this.timer = setInterval(this.progress, 20);
  //   this.callApi()
  //     .then((res) => {
  //       console.log(res);
  //       if(res == null){
  //         this.setState({ rows: [] });
  //       } else {
  //         this.setState({ rows: res });
  //       }
  //       clearInterval(this.timer);
  //     })
  //     .catch((err) => console.log(err));
  //   let userId = null;
  //   AsyncStorage.getItem("userName",(err, result) => {
  //     userId= result;
  //   })
  //   utilLog.fn_insertPLogs(userId, 'log-ND-VW02');
  // }

  // callApi = async () => {
  //   var param = this.props.match.params;
  //   const response = await fetch(`/nodes/${param.node}${this.props.location.search}`);
  //   const body = await response.json();
  //   return body;
  // };

  // progress = () => {
  //   const { completed } = this.state;
  //   this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  // };

  // onUpdateData = () => {
  //   console.log("onUpdateData={this.props.onUpdateData}")
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
    // console.log("CsOverview_Render : ",this.state.rows.basic_info);
    return (
      <div>
        <div className="content-wrapper node-detail fulled">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
              Bill Detail
              {/* <small>2021-10</small> */}
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext
                  style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
                />
                <NavLink to="/nodes">Settings</NavLink>
              </li>
              <li className="active">
                <NavigateNext
                  style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
                />
                Metering
              </li>
              <li className="active">
                <NavigateNext
                  style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
                />
                Bill
              </li>
            </ol>
          </section>

          {/* 내용부분 */}
          <section className="content">
                <BasicInfo />
                <DataTransferBill />
                <ResourceUsageBill />
            {/* {this.state.rows ? (
              [
                <BasicInfo />,
                <DataTransferPrice />,
                <ResourcePrice />,
                // <Events rowData={this.state.rows.events}/>
              ]
            ) : (
              <CircularProgress
                variant="determinate"
                value={this.state.completed}
                style={{ position: "absolute", left: "50%", marginTop: "20px" }}
              ></CircularProgress>
            )} */}
          </section>
        </div>
      </div>
    );
  }
}

class BasicInfo extends Component {
  render() {
    return (
      <div className="content-box">
        <div className="cb-header">
          <span>Basic Info</span>
        </div>
        <div className="cb-body">
          {/* <div>
            <span>Date : </span>
            <strong>2021-10</strong>
          </div> */}
          <div style={{ display: "flex" }}>
            <div className="cb-body-left">
              <div>
                <span>Date : </span>
                <span>2021-10</span>
              </div>
            </div>
            <div className="cb-body-right">
              <div>
                <span>Total Bill : </span>
                <span>$ 100,000</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

class DataTransferBill extends Component {
  render() {
    return (
      <div className="content-box">
        <div className="cb-header">
          <span>Data Transfer Bill : </span>
          <span>$ 2.6</span>
          {/* <NdTaintConfig name={this.props.rowData.name} taint={this.props.rowData.taint}/> */}
        </div>
        <div className="cb-body">
          <div style={{ display: "flex", margin:"0px 0 0 10px"}}>
            <div style={{minWidth:"200px"}}>
              <div>
              <strong>Date Transfer In</strong>
              </div>
            </div>
            <div style={{minWidth:"200px"}}>
              <div>
                <span>Usage : </span>
                <span>1.2GB</span>
              </div>
            </div>
            <div className="cb-body-right">
              <div>
                <span>Bill : </span>
                <span>$ 1.3</span>
              </div>
            </div>
          </div>
        </div>
        
        <br/>
        <div className="cb-body">
          <div style={{ display: "flex", margin:"0px 0 0 10px" }}>
            <div style={{minWidth:"200px"}}>
              <div>
                <strong>Date Transfer Out</strong>
              </div>
            </div>
            <div style={{minWidth:"200px"}}>
              <div>
                <span>Usage : </span>
                <span>1.2GB</span>
              </div>
            </div>
            <div className="cb-body-right">
              <div>
                <span>Bill : </span>
                <span>$ 1.3</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

class ResourceUsageBill extends Component {
  render() {
    return (
      <div className="content-box">
        <div className="cb-header">
          <span>Resource Usage Bill : </span>
          <span>$ 100.12</span>
          {/* <NdTaintConfig name={this.props.rowData.name} taint={this.props.rowData.taint}/> */}
        </div>
        <div className="cb-body">
          <div style={{ display: "flex", margin:"0px 0 0 10px" }}>
            <div style={{minWidth:"100px"}}>
              <div>
              <strong>CPU</strong>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Set : </span>
                <span>2 vCPU</span>
              </div>
            </div>
            <div style={{minWidth:"200px"}}>
              <div>
                <span>Usage : </span>
                <span>720 hours</span>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Bill : </span>
                <span>$ 32.43</span>
              </div>
            </div>
          </div>
        </div>
        <br/>

        <div className="cb-body">
          <div style={{ display: "flex", margin:"0px 0 0 10px" }}>
          <div style={{minWidth:"100px"}}>
              <div>
              <strong>Memory</strong>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Set : </span>
                <span>4 GiB</span>
              </div>
            </div>
            <div style={{minWidth:"200px"}}>
              <div>
                <span>Usage : </span>
                <span>720 hours</span>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Bill : </span>
                <span>$ 32.43</span>
              </div>
            </div>
          </div>
        </div>
        <br/>

        <div className="cb-body">
          <div style={{ display: "flex", margin:"0px 0 0 10px" }}>
            <div style={{minWidth:"100px"}}>
              <div>
              <strong>Disk</strong>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Set : </span>
                <span>28GB</span>
              </div>
            </div>
            <div style={{minWidth:"200px"}}>
              <div>
                <span>Usage : </span>
                <span>720 hours</span>
              </div>
            </div>
            <div style={{minWidth:"150px"}}>
              <div>
                <span>Bill : </span>
                <span>$ 32.43</span>
              </div>
            </div>
          </div>
        </div>

        
      </div>
    );
  }
}

export default BillDetail;
