// import React, { Component } from "react";
// import Paper from "@material-ui/core/Paper";
// import { NavLink, Link } from "react-router-dom";
// import CircularProgress from "@material-ui/core/CircularProgress";
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
//   TableHeaderRow,
//   TableColumnResizing,
//   PagingPanel,
// } from "@devexpress/dx-react-grid-material-ui";
// import Editor from "./../../modules/Editor";
// import { NavigateNext} from '@material-ui/icons';
// import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
// import PjCreateProject from '../modal/PjCreateProject';





// class Projects extends Component {
//   constructor(props) {
//     super(props);
//     this.state = {
//       columns: [
//         { name: "name", title: "Name" },
//         { name: "status", title: "Status" },
//         { name: "cluster", title: "Cluster" },
//         { name: "labels", title: "Labels" },
//         { name: "cluster_cpu", title: "CPU/Cluster" },
//         { name: "cluster_mem", title: "Memory/Cluster" },
//         { name: "cluster_pods", title: "Pods/Cluster" },
//         { name: "pod_cpu", title: "CPU/Pod" },
//         { name: "pod_mem", title: "Memory/Pod" },
//         { name: "created_time", title: "Created Time" },
//       ],
//       defaultColumnWidths: [
//         { columnName: "name", width: 200 },
//         { columnName: "status", width: 100 },
//         { columnName: "cluster", width: 100 },
//         { columnName: "labels", width: 180 },
//         { columnName: "cluster_cpu", width: 130 },
//         { columnName: "cluster_mem", width: 150 },
//         { columnName: "cluster_pods", width: 130 },
//         { columnName: "pod_cpu", width: 110 },
//         { columnName: "pod_mem", width: 130 },
//         { columnName: "created_time", width: 180 },
//       ],
//       rows: "",

//       // Paging Settings
//       currentPage: 0,
//       setCurrentPage: 0,
//       pageSize: 10,
//       pageSizes: [5, 10, 15, 0],

//       completed: 0,
//     };
//   }

//   componentWillMount() {
//     this.props.menuData("none");
//   }

  

//   callApi = async () => {
//     const response = await fetch("/projects");
//     const body = await response.json();
//     return body;
//   };

//   progress = () => {
//     const { completed } = this.state;
//     this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
//   };

//   //컴포넌트가 모두 마운트가 되었을때 실행된다.
//   componentDidMount() {
//     //데이터가 들어오기 전까지 프로그래스바를 보여준다.
//     this.timer = setInterval(this.progress, 20);
//     this.callApi()
//       .then((res) => {
//         if(res == null){
//           this.setState({ rows: [] });
//         } else {
//           this.setState({ rows: res });
//         }
//         clearInterval(this.timer);
//       })
//       .catch((err) => console.log(err));



//     let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
//     utilLog.fn_insertPLogs(userId, 'log-PJ-VW01');

//   };

//   onRefresh = () => {
//     this.callApi()
//       .then((res) => {
//         if(res == null){
//           this.setState({ rows: [] });
//         } else {
//           this.setState({ rows: res });
//         }
//         clearInterval(this.timer);
//       })
//       .catch((err) => console.log(err));
//   };

//   render() {

//     // 셀 데이터 스타일 변경
//     const HighlightedCell = ({ value, style, row, ...restProps }) => (
//       <Table.Cell>
//         <span
//           style={{
//             color:
//               value === "Active" ? "#1ab726" : value === "Deactive" ? "red" : undefined,
//           }}
//         >
//           {value}
//         </span>
//       </Table.Cell>
//     );


    
    
//     const Cell = (props) => {

//       const fnEnterCheck = (prop) => {
//         var arr = [];
//         var i;
//         for(i=0; i < Object.keys(prop.value).length; i++){
//           const str = Object.keys(prop.value)[i] + " : " + Object.values(prop.value)[i]
//           arr.push(str)
//         }
//         return (
//          arr.map(item => {
//            return (
//              <p>{item}</p>
//            )
//          })
//         )
//         // return (
//           // props.value.indexOf("|") > 0 ? 
//           //   props.value.split("|").map( item => {
//           //     return (
//           //       <p>{item}</p>
//           //   )}) : 
//           //     props.value
//         // )
//       }

//       const { column, row } = props;
//       // console.log("cell : ", props);
//       if (column.name === "status") {
//         return <HighlightedCell {...props} />;
//       } else if (column.name === "name") {
//         return (
//           <Table.Cell
//             {...props}
//             style={{ cursor: "pointer" }}
//           ><Link to={{
//             pathname: `/projects/${props.value}/overview`,
//             search: "cluster="+row.cluster,
//             state: {
//               data : row
//             }
//           }}>{props.value}</Link></Table.Cell>
//         );
//       } else if (column.name === "labels"){
//         return (
//         <Table.Cell>{fnEnterCheck(props)}</Table.Cell>
//         )
//       }
//       return <Table.Cell {...props} />;
//     };

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
//       <div className="content-wrapper fulled">
//         {/* 컨텐츠 헤더 */}
//         <section className="content-header">
//           <h1>
//             <span onClick={this.onRefresh} style={{cursor:"pointer"}}>
//               Projects
//             </span>
//             <small></small>
//           </h1>
//           <ol className="breadcrumb">
//             <li>
//               <NavLink to="/dashboard">Home</NavLink>
//             </li>
//             <li className="active">
//               <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
//               Projects
//             </li>
//           </ol>
//         </section>
//         <section className="content" style={{ position: "relative" }}>
//           <Paper>
//             {this.state.rows ? (
//               [
//                 <PjCreateProject/>,
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
//                     defaultSorting={[{ columnName: 'created_time', direction: 'desc' }]}
//                   />
//                   <IntegratedSorting />

//                   {/* 페이징 */}
//                   <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
//                   <IntegratedPaging />
//                   <PagingPanel pageSizes={this.state.pageSizes} />

                  

//                   {/* 테이블 */}
//                   <Table cellComponent={Cell} rowComponent={Row} />
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
//         </section>
//       </div>
//     );
//   }
// }

// export default Projects;
