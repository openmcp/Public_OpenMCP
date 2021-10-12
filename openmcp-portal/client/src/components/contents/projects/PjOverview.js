import React, { Component } from "react";
import { NavLink } from 'react-router-dom';
import CircularProgress from "@material-ui/core/CircularProgress";
// import line_chart_sample from './../../../json/line_chart_sample.json'
import {NavigateNext} from '@material-ui/icons';

import {
  // SearchState,
  // IntegratedFiltering,
  // PagingState,
  // IntegratedPaging,
  SortingState,
  IntegratedSorting,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  // Toolbar,
  // SearchPanel,
  TableHeaderRow,
  TableColumnResizing,
  // PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import LineReChart from './../../modules/LineReChart';
import SelectBox from './../../modules/SelectBox';
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';


let apiParams = "";
class PjOverview extends Component {
  state = {
    rows:"",
    completed: 0,
    reRender : "",
  }

  componentWillMount() {
    //왼쪽 메뉴쪽에 타이틀 데이터 전달
    const result = {
      menu : "projects",
      title : this.props.match.params.project,
      pathParams : {
        searchString : this.props.location.search,
        project : this.props.match.params.project,
        // state : this.props.location.state.data
      },
    }
    this.props.menuData(result);
    apiParams = this.props.match.params.project;
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
    utilLog.fn_insertPLogs(userId, 'log-PJ-VW02');
  }  

  callApi = async () => {
    var param = this.props.match.params.project;
    const response = await fetch(`/projects/${param}/overview${this.props.location.search}`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  render() {
    // console.log("PjOverview_Render : ",this.state.rows.basic_info);
    return (
      <div>
        <div className="content-wrapper">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
            Overview
              <small>{this.props.match.params.project}</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                <NavLink to="/projects">Projects</NavLink>
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
            <BasicInfo rowData={this.state.rows.basic_info}/>,
            <div style={{display:"flex"}}>
              <ProjectResources rowData={this.state.rows.project_resource}/>
              <UsageTop5 rowData={this.state.rows.usage_top5} query ={this.props.location.search}/>
            </div>,
            <PhysicalResources rowData={this.state.rows.physical_resources}/>
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
    // console.log("BasicInfo:", this.props.rowData.name)
    
    return (
      <div className="content-box">
        <div className="cb-header">Basic Info</div>
        <div className="cb-body">
        <div style={{display:"flex"}}>
            <div className="cb-body-left">
              <div>
                <span>Name : </span>
                <strong>{this.props.rowData.name}</strong>
              </div>
              <div>
                <span>Cluster : </span>
                {this.props.rowData.cluster}
              </div>
              <div>
                  <span>Labels : </span>
                  <div style={{margin : "-25px 0px 0px 66px"}}>
                    {
                      Object.keys(this.props.rowData.labels).length > 0 ?
                        (
                          Object.entries(this.props.rowData.labels).map(i=>{
                          return (<div>{i.join(" : ")}</div>)
                        })
                        ) : 
                        "-"
                    }
                  </div>
                </div>
            </div>
            <div className="cb-body-right">
              <div>
                <span>Status : </span>
                {this.props.rowData.status}
              </div>
              <div>
                <span>UID : </span>
                {this.props.rowData.uid}
              </div>
              <div>
                <span>Created Time : </span>
                {this.props.rowData.created_time}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

class ProjectResources extends Component {
  state = {
    columns: [
      { name: "resource", title: "Resource" },
      { name: "total", title: "Total" },
      { name: "abnormal", title: "Abnormal" },
    ],
  }
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
    return (
      <div className="content-box col-sep-2 ">
        <div className="cb-header">Project Resources</div>
        <div className="cb-body table-style">
          <Grid
            rows = {this.props.rowData}
            columns = {this.state.columns}
            >

            {/* Sorting */}
            <SortingState
            // defaultSorting={[{ columnName: 'city', direction: 'desc' }]}
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




class UsageTop5 extends Component {
  state = {
    columns: [
      { name: "name", title: "Name" },
      // { name: "type", title: "Type" },
      { name: "usage", title: "Usage" },
    ],
    defaultColumnWidths: [
      { columnName: "name", width: 350 },
      { columnName: "usage", width: 200 },
    ],
    rows : this.props.rowData.cpu,
  }

  callApi = async () => {
    const response = await fetch(`/projects/${apiParams}/overview${this.props.query}`);
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
      switch(data){
        case "cpu":
          this.setState({rows:this.props.rowData.cpu});

          // this.callApi()
          // .then((res) => {
          //   this.setState({ rows: res.usage_top5.cpu });
          // })
          // .catch((err) => console.log(err));

          break;
        case "memory":
          this.setState({rows:this.props.rowData.memory});

          // this.callApi()
          // .then((res) => {
          //   console.log(res.usage_top5.memory)
          //   this.setState({ rows: res.usage_top5.memory });
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
          Usage Top5
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
            // defaultSorting={[{ columnName: 'city', direction: 'desc' }]}
            />
            <IntegratedSorting />

            <Table/>
            <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
            <TableHeaderRow showSortingControls rowComponent={HeaderRow}/>
          </Grid>
        </div>
      </div>
    );
  }
}

// class PhysicalResources extends Component {
//   render(){
//     return (
//       <div className="content-box">
//         <div className="cb-header">Physical Resources</div>
//         <div className="cb-body">
//           <div className="cb-body-content" style={{height:"250px"}}>
//             <LineChart data={line_chart_sample[0].cpu} ></LineChart>
//           </div>
//           <div className="cb-body-content" style={{height:"250px"}}>
//             <LineChart data={line_chart_sample[1].memory} ></LineChart>
//           </div>
//           <div className="cb-body-content" style={{height:"250px"}}>
//             <LineChart data={line_chart_sample[2].network} ></LineChart>
//           </div>
//         </div>
//       </div>
//     );
//   }
// }


class PhysicalResources extends Component {
  render() {
    const network_title = ["in", "out"];
    return (
      <div className="content-box line-chart">
        <div className="cb-header">Physical Resources</div>
        <div className="cb-body">
          <div className="cb-body-content">
            <LineReChart
              rowData={this.props.rowData.cpu}
              unit="m"
              name="cpu"
              title="CPU"
              cardinal={false}
            ></LineReChart>
          </div>
          <div className="cb-body-content">
            <LineReChart
              rowData={this.props.rowData.memory}
              unit="mib"
              name="memory"
              title="Memory"
              cardinal={false}
            ></LineReChart>
          </div>
          <div className="cb-body-content">
            <LineReChart
              rowData={this.props.rowData.network}
              unit="Bps"
              name={network_title}
              title="Network"
              cardinal={true}
            ></LineReChart>
          </div>
        </div>
      </div>
    );
  }
}




export default PjOverview;

