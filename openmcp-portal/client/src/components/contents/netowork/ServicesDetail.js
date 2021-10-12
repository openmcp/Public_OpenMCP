import React, { Component } from "react";
import { NavLink, Link} from 'react-router-dom';
import CircularProgress from "@material-ui/core/CircularProgress";
import { NavigateNext} from '@material-ui/icons';
import Paper from "@material-ui/core/Paper";
import {
  SearchState,IntegratedFiltering,PagingState,IntegratedPaging,SortingState,IntegratedSorting,
} from "@devexpress/dx-react-grid";
import {
  Grid,Table,Toolbar,SearchPanel,TableColumnResizing,TableHeaderRow,PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';


// let apiParams = "";
class ServicesDetail extends Component {
  state = {
    rows:"",
    completed: 0,
    reRender : ""
  }

  componentWillMount() {
    const result = {
      menu : "projects",
      title : this.props.match.params.project,
      pathParams : {
        searchString : this.props.location.search,
        project : this.props.match.params.project
      }
    }
    this.props.menuData(result);
    // apiParams = this.props.match.params.project;
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
    utilLog.fn_insertPLogs(userId, 'log-PJ-VW10');
  }  

  callApi = async () => {
    var param = this.props.match.params;
    const response = await fetch(`/projects/${this.props.location.state.data.project}/resources/services/${param.service}${this.props.location.search}`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  render() {
    return (
      <div>
        <div className="content-wrapper fulled">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
            { this.props.match.params.service}
              <small>Service Overview</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                Netowork
              </li>
            </ol>
          </section>

          {/* 내용부분 */}
          <section className="content">
          {this.state.rows ? (
            [
            <BasicInfo rowData={this.state.rows.basic_info}/>,
            // <Workloads rowData={this.state.rows.workloads}/>,
            <Pods rowData={this.state.rows.pods}/>,
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
  // constructor(props){
  //   super(props);
    
  //   this.state = {
  //     selector : this.props.rowData.selector.split(",")
  //   }
  // }
  render(){
    return (
      <div className="content-box">
        <div className="cb-header">Basic Info</div>
        <div className="cb-body">
          <div style={{display:"flex"}}>
            <div className="cb-body-left" style={{width:"50%"}}>
              <div>
                <span>Name : </span>
                <strong>{this.props.rowData.name}</strong>
              </div>
              <div>
                  <span>Project : </span>
                  {this.props.rowData.project}
                </div>
              <div>
                <span>Type : </span>
                {this.props.rowData.type}
              </div>
              <div>
                <span>Session Affinity : </span>
                {this.props.rowData.session_affinity}
              </div>
              <div>
                <span>Selector : </span>
                  {this.props.rowData.selector}
              </div>
              {/* <div>
                <span>Access Type : </span>
                {this.props.rowData.access_type}
              </div> */}
            </div>
            <div className="cb-body-right">
              <div>
                <span>Cluster : </span>
                {this.props.rowData.cluster}
              </div>
              <div>
                <span>Cluster IP : </span>
                {this.props.rowData.cluster_ip}
              </div>
              <div>
                <span>Externer IP : </span>
                {this.props.rowData.external_ip}
              </div>
              <div>
                <span>Endpoints : </span>
                {this.props.rowData.endpoints}
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

// class Workloads extends Component {
//   constructor(props) {
//     super(props);
//     this.state = {
//       columns: [
//         { name: "name", title: "Name" },
//         { name: "status", title: "Status" },
//         { name: "type", title: "Type" },
//       ],
//       defaultColumnWidths: [
//         { columnName: "name", width: 400 },
//         { columnName: "status", width: 150 },
//         { columnName: "type", width: 150 },
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
//         <div className="cb-header">Workloads</div>
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
//                     // defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
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

class Pods extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status"},
        { name: "cluster", title: "Cluster"},
        { name: "project", title: "Project" },
        { name: "pod_ip", title: "Pod IP" },
        { name: "node", title: "Node" },
        { name: "node_ip", title: "Node IP" },
        // { name: "cpu", title: "CPU" },
        // { name: "memory", title: "Memory" },
        { name: "created_time", title: "Created Time" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 330 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 100 },
        { columnName: "project", width: 130 },
        { columnName: "pod_ip", width: 120 },
        { columnName: "node", width: 230 },
        { columnName: "node_ip", width: 130 },
        // { columnName: "cpu", width: 80 },
        // { columnName: "memory", width: 100 },
        { columnName: "created_time", width: 170 },
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

    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell
        {...restProps}
        style={{
          // backgroundColor:
          //   value === "Healthy" ? "white" : value === "Unhealthy" ? "white" : undefined,
          // cursor: "pointer",
          ...style,
        }}>
        <span
          style={{
            color:
              value === "Pending" ? "orange" : 
                value === "Failed" ? "red" : 
                  value === "Unknown" ? "red" : 
                    value === "Succeeded" ? "skyblue" : 
                      value === "Running" ? "#1ab726" : "black"
          }}>
          {value}
        </span>
      </Table.Cell>
    );

    //셀
    const Cell = (props) => {
      const { column, row } = props;
      // console.log("cell : ", props);
      // const values = props.value.split("|");
      // console.log("values", props.value);
      
      // const values = props.value.replace("|","1");
      // console.log("values,values", values)

      const fnEnterCheck = () => {
        if(props.value === undefined){
          return ""
        } else {
          return (
            props.value.indexOf("|") > 0 ? 
              props.value.split("|").map( item => {
                return (
                  <p>{item}</p>
              )}) : 
                props.value
          )
        }
      }


      if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "name") {
        // console.log("name", props.value);
        return (
          <Table.Cell
            {...props}
            style={{ cursor: "pointer" }}
          ><Link to={{
            pathname: `/pods/${props.value}`,
            state: {
              data : row
            }
          }}>{fnEnterCheck()}</Link></Table.Cell>
        );
      }
      return <Table.Cell>{fnEnterCheck()}</Table.Cell>;
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
    const Row = (props) => {
      // console.log("row!!!!!! : ",props);
      return <Table.Row {...props} key={props.tableRow.key}/>;
    };

    return (
      <div className="content-box">
        <div className="cb-header">Pods</div>
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
                    // defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
                  />
                  <IntegratedSorting />

                  {/* 페이징 */}
                  <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* 테이블 */}
                  <Table cellComponent={Cell}  rowComponent={Row} />
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
                    // defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
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


export default ServicesDetail;