import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
import { Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
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
// import { NavigateNext} from '@material-ui/icons';
import * as utilLog from '../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';

// import Editor from "./../modules/Editor";

// let apiParams = "";
class VPA extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "namespace", title: "Project" },
        { name: "reference", title: "Reference" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 200 },
        { columnName: "namespace", width: 130 },
        { columnName: "reference", width: 250 },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10, 
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      editorContext : ``,
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }


  

  callApi = async () => {
    const response = await fetch(`/vpa`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  //컴포넌트가 모두 마운트가 되었을때 실행된다.
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
    utilLog.fn_insertPLogs(userId, 'log-PD-VW01');
  };

  render() {

    // 셀 데이터 스타일 변경
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
            pathname: `/pods/${props.value}/overview`,
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
      // background:"rgba(0, 0, 0, 0.77)", 
      // background:"rgba(0, 0, 0, 0.5)", 
      <div className="sub-content-wrapper">
        {/* <div className="loading" style={{height:"109vh", width:"100%", position:"absolute", background:"rgba(0, 0, 0, 0.5)", zIndex:"100", top:"0", left: "0", margin:"auto" }}></div> */}
        {/* 컨텐츠 헤더 */}
        {/* <section className="content-header">
          <h1>
            VPA
            <small>(Vertical Pod Autoscaler)</small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">
              <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
              Pods
            </li>
          </ol>
        </section> */}
        <section className="content" style={{ position: "relative" }}>
        {/* <div className="VPA-TEMP">
            VPA 
            <small> (Vertical Pod Autoscaler)</small>
          </div> */}
          <Paper>
            {this.state.rows ? (
              [
                // <Editor title="create" context={this.state.editorContext}/>,
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
                  <Table cellComponent={Cell} rowComponent={Row} />
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
        </section>
      </div>
    );
  }
}

export default VPA;
