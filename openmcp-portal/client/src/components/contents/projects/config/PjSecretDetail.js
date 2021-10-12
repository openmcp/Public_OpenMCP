import React, { Component } from "react";
import { NavLink } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import { NavigateNext } from "@material-ui/icons";
import Paper from "@material-ui/core/Paper";
// import LineChart from '../../../modules/LineChart';
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
} from "@devexpress/dx-react-grid";
// import LineReChart from '../../../modules/LineReChart';
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import * as utilLog from "./../../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';

// let apiParams = "";
class PjSecretDetail extends Component {
  state = {
    rows: "",
    completed: 0,
    reRender: "",
  };

  componentWillMount() {
    const result = {
      menu: "projects",
      title: this.props.match.params.project,
      pathParams : {
        searchString : this.props.location.search,
        project : this.props.match.params.project
      }
    };
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
    utilLog.fn_insertPLogs(userId, "log-PJ-VW16");
  }

  callApi = async () => {
    var param = this.props.match.params;
    const response = await fetch(
      `/projects/${param.project}/config/secrets/${param.secret}${this.props.location.search}`
    );
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
        <div className="content-wrapper pod-detail">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
              {this.props.match.params.secret}
              <small>Secret Overview</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext
                  style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
                />
                <NavLink to="/projects">Projects</NavLink>
              </li>
              <li className="active">
                <NavigateNext
                  style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
                />
                Config
              </li>
            </ol>
          </section>

          {/* 내용부분 */}
          <section className="content">
            {this.state.rows ? (
              [
                <BasicInfo rowData={this.state.rows.basic_info} />,
                <Data rowData={this.state.rows.data} />,
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
  render() {
    return (
      <div className="content-box">
        <div className="cb-header">Basic Info</div>
        <div className="cb-body">
          <div style={{ display: "flex" }}>
            <div className="cb-body-left">
              <div>
                <span>Name : </span>
                <strong>{this.props.rowData.name}</strong>
              </div>
              <div>
                <span>Project : </span>
                {this.props.rowData.project}
              </div>
              
            </div>
            <div className="cb-body-right">
              <div>
                <span>Type : </span>
                {this.props.rowData.type}
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
class Data extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "key", title: "Key" },
        { name: "value", title: "Value" },
      ],
      defaultColumnWidths: [
        { columnName: "key", width: 150 },
        { columnName: "value", width: 800 },
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
      return <Table.Row {...props} key={props.tableRow.key} />;
    };

    const Cell = (props) => {
      const { column } = props;
      if (column.name === "value") {
        // console.log("cell : ", props);
        return (
          <Table.Cell
            {...props}
            style={{
              cursor: "pointer",
              wordBreak: "break-all",
              whiteSpace: "inherit",
              lineHeight: "1.5",
            }}
          >
            {props.value}
          </Table.Cell>
        );
      }
      return <Table.Cell {...props} />;
    };

    return (
      <div className="content-box">
        <div className="cb-header">Data</div>
        <div className="cb-body">
          <Paper>
            {this.state.rows ? (
              [
                <Grid rows={this.state.rows} columns={this.state.columns}>
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
                  <PagingState
                    defaultCurrentPage={0}
                    defaultPageSize={this.state.pageSize}
                  />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* 테이블 */}
                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing
                    defaultColumnWidths={this.state.defaultColumnWidths}
                  />
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
  }
}

export default PjSecretDetail;
