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
  SelectionState,
  IntegratedSelection,
  // FilteringState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
  // TableFilterRow,
} from "@devexpress/dx-react-grid-material-ui";
// import {  Button,} from "@material-ui/core";
import Editor from "./../../modules/Editor";
import * as utilLog from "./../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
import PjDeploymentMigration from "./../modal/PjDeploymentMigration";
import { NavigateNext } from "@material-ui/icons";
import axios from "axios";
// import ProgressTemp from './../../modules/ProgressTemp';
import SnapShotControl from "./../modal/SnapShotControl";

import IconButton from "@material-ui/core/IconButton";
import MenuItem from "@material-ui/core/MenuItem";
import MoreVertIcon from "@material-ui/icons/MoreVert";

import Popper from "@material-ui/core/Popper";
import MenuList from "@material-ui/core/MenuList";
import Grow from "@material-ui/core/Grow";
import { AiOutlineDeploymentUnit } from "react-icons/ai";
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';

// let apiParams = "";
class Deployments extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Ready" },
        { name: "cluster", title: "Cluster" },
        { name: "project", title: "Project" },
        { name: "image", title: "Image" },
        { name: "created_time", title: "Created Time" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 250 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 130 },
        { columnName: "project", width: 200 },
        { columnName: "image", width: 370 },
        { columnName: "created_time", width: 170 },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      selection: [],
      selectedRow: "",
      clusterName: "",
      editorContext: `apiVersion: openmcp.k8s.io/v1alpha1
kind: OpenMCPDeployment
metadata:
  name: openmcp-deployment2
  namespace: openmcp
spec:
  replicas: 3
  labels:
      app: openmcp-nginx
  template:
    spec:
      template:
        spec:
          containers:
          - image: nginx
            name: nginx`,
      openProgress: false,
      anchorEl: null,
      projects: "",
    };
  }

  componentWillMount() {
    var projects = "";

    AsyncStorage.getItem("projects", (err, result) => {
      projects = result;
    });

    this.setState({
      projects: projects,
    });
    this.props.menuData("none");
  }

  callApi = async () => {
    // var param = this.props.match.params.cluster;
    // queryString = queryString.parse(this.props.location.search).cluster
    // console.log(query);
    const response = await fetch(`/deployments`);
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
      .then((result) => {
        // var res = [];
        // if (result == null) {
        //   this.setState({ rows: [] });
        // } else {
        //   result.map((item) => {
        //     if (this.state.projects.indexOf(item.project) < 0) {
        //       res.push(item);
        //     }
        //   });

        //   // this.setState({ rows: res });
        // }
        this.setState({ rows: result });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));

    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-PJ-VW03");
  }

  onUpdateData = () => {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        this.setState({
          selection: [],
          selectedRow: "",
          rows: res,
        });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));

    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-PJ-VW03");
  };

  excuteScript = (context) => {
    if (this.state.openProgress) {
      this.setState({ openProgress: false });
    } else {
      this.setState({ openProgress: true });
    }

    const url = `/deployments/create`;
    const data = {
      yaml: context,
    };
    // console.log(context)
    axios
      .post(url, data)
      .then((res) => {
        // alert(res.data.message);
        this.setState({ open: false });
        this.onUpdateData();
      })
      .catch((err) => {
        alert(err);
      });
  };

  onRefresh = () => {
    if (this.state.openProgress) {
      this.setState({ openProgress: false });
    } else {
      this.setState({ openProgress: true });
    }
    this.callApi()
      .then((res) => {
        this.setState({
          // selection : [],
          // selectedRow : "",
          rows: res,
        });
      })
      .catch((err) => console.log(err));
  };

  closeProgress = () => {
    this.setState({ openProgress: false });
  };

  //셀
  Cell = (props) => {
    console.log("CEll");
    const { column, row } = props;

    if (column.name === "name") {
      // // console.log("name", props.value);
      // console.log("this.props.match.params", this.props)
      return (
        <Table.Cell {...props} style={{ cursor: "pointer" }}>
          <Link
            to={{
              pathname: `/deployments/${props.value}`,
              search: `cluster=${row.cluster}&project=${row.project}`,
              state: {
                data: row,
              },
            }}
          >
            {props.value}
          </Link>
        </Table.Cell>
      );
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

  HeaderRow = ({ row, ...restProps }) => (
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

  Row = (props) => {
    // if(this.state.projects.indexOf(props.tableRow.row.project) > -1){
    //   return <Table.Row {...props} key={props.tableRow.key} />;
    // } else {
    //   return null;
    // }

    return <Table.Row {...props} key={props.tableRow.key} />;
  };
  render() {
    const onSelectionChange = (selection) => {
      // console.log(this.state.rows[selection[0]])
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({
        selectedRow: this.state.rows[selection[0]]
          ? this.state.rows[selection[0]]
          : {},
      });
    };

    const handleClick = (event) => {
      if (this.state.anchorEl === null) {
        this.setState({ anchorEl: event.currentTarget });
      } else {
        this.setState({ anchorEl: null });
      }
    };

    const handleClose = () => {
      this.setState({ anchorEl: null });
    };

    const open = Boolean(this.state.anchorEl);

    return (
      <div className="content-wrapper fulled">
        {/* {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""} */}
        {this.state.clusterName}
        {/* 컨텐츠 헤더 */}
        <section className="content-header" onClick={this.onRefresh}>
          <h1>
            <i><AiOutlineDeploymentUnit/></i>
            <span>Deployments</span>
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <Link to="/dashboard">Home</Link>
            </li>
            <li className="active">
              <NavigateNext
                style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
              />
              Deployments
            </li>
          </ol>
        </section>
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                <div
                  style={{
                    position: "absolute",
                    right: "21px",
                    top: "20px",
                    zIndex: "10",
                    textTransform: "capitalize",
                  }}
                >
                  <IconButton
                    aria-label="more"
                    aria-controls="long-menu"
                    aria-haspopup="true"
                    onClick={handleClick}
                  >
                    <MoreVertIcon />
                  </IconButton>

                  <Popper
                    open={open}
                    anchorEl={this.state.anchorEl}
                    role={undefined}
                    transition
                    disablePortal
                    placement={"bottom-end"}
                  >
                    {({ TransitionProps, placement }) => (
                      <Grow
                        {...TransitionProps}
                        style={{
                          transformOrigin:
                            placement === "bottom"
                              ? "center top"
                              : "center top",
                        }}
                      >
                        <Paper>
                          <MenuList autoFocusItem={open} id="menu-list-grow">
                            <MenuItem
                              style={{
                                textAlign: "center",
                                display: "block",
                                fontSize: "14px",
                              }}
                            >
                              <SnapShotControl
                                title="create snapshot"
                                rowData={this.state.selectedRow}
                                onUpdateData={this.onUpdateData}
                                menuClose={handleClose}
                              />
                            </MenuItem>
                            <MenuItem
                              style={{
                                textAlign: "center",
                                display: "block",
                                fontSize: "14px",
                              }}
                            >
                              <PjDeploymentMigration
                                title="pod migration"
                                rowData={this.state.selectedRow}
                                onUpdateData={this.onUpdateData}
                                menuClose={handleClose}
                              />
                            </MenuItem>
                            <MenuItem
                              // onClick={handleClose}
                              style={{
                                textAlign: "center",
                                display: "block",
                                fontSize: "14px",
                              }}
                            >
                              <Editor
                                btTitle="create"
                                title="Create Deployment"
                                context={this.state.editorContext}
                                excuteScript={this.excuteScript}
                                menuClose={handleClose}
                              />
                            </MenuItem>
                          </MenuList>
                        </Paper>
                      </Grow>
                    )}
                  </Popper>
                </div>,
                <Grid rows={this.state.rows} columns={this.state.columns}>
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState defaultValue="" />

                  <SearchPanel style={{ marginLeft: 0 }} />

                  {/* Sorting */}
                  <SortingState
                    defaultSorting={[
                      { columnName: "created_time", direction: "desc" },
                    ]}
                  />

                  {/* 페이징 */}
                  <PagingState
                    defaultCurrentPage={0}
                    defaultPageSize={this.state.pageSize}
                  />

                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* <EditingState
                    onCommitChanges={commitChanges}
                  /> */}
                  <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  />
                  {/* <FilteringState/> */}

                  <IntegratedFiltering />
                  <IntegratedSorting />
                  <IntegratedSelection />
                  <IntegratedPaging />

                  {/* 테이블 */}
                  <Table cellComponent={this.Cell} />
                  <TableColumnResizing
                    defaultColumnWidths={this.state.defaultColumnWidths}
                  />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={this.HeaderRow}
                  />
                  <TableSelection
                    selectByRowClick
                    highlightRow
                    rowComponent={this.Row}
                    // showSelectionColumn={false}
                  />

                  {/* <TableFilterRow showFilterSelector={true}/> */}
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

export default Deployments;
