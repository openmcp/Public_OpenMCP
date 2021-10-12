import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
import { NavLink, Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  IntegratedSelection,
  SelectionState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  TableSelection,
  PagingPanel,
  // TableColumnVisibility
} from "@devexpress/dx-react-grid-material-ui";
// import { NavigateNext } from "@material-ui/icons";
import * as utilLog from "./../../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
// import AddMembers from "./../AddMembers";
// import Editor from "../../modules/Editor";
// import AcChangeRole from "./../../modal/AcChangeRole";
import IconButton from "@material-ui/core/IconButton";
import MenuItem from "@material-ui/core/MenuItem";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import Popper from "@material-ui/core/Popper";
import MenuList from "@material-ui/core/MenuList";
import Grow from "@material-ui/core/Grow";
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';
// import GrCreateGroup from "./../../modal/GrCreateGroup";
// import GrEditGroup from "./../../modal/GrEditGroup";
import axios from "axios";
import Confirm2 from "./../../../modules/Confirm2";
import ThCreateThreshold from "../../modal/ThCreateThreshold.js";
import ThEditThreshold from "../../modal/ThEditThreshold.js";
import WarningRoundedIcon from "@material-ui/icons/WarningRounded";

class BillList extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "date", title: "Date" },
        { name: "total", title: "Total Bill" },
      ],
      defaultColumnWidths: [
        { columnName: "date", width: 150 },
        { columnName: "total", width: 150 },
      ],
      defaultHiddenColumnNames: [],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10,
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      selection: [],
      selectedRow: "",
      anchorEl: null,

      grEditOpen: false,

      confirmOpen: false,
      confirmInfo: {
        title: "Delete Threshold",
        context: "Are you sure you want to Delete Host Threshold?",
        button: {
          open: "",
          yes: "CONFIRM",
          no: "CANCEL",
        },
      },
      confrimTarget: "",
      confirmTargetKeyname: "Threshold Name",
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }

  callApi = async () => {
    const response = await fetch(`/apis/metering/bill`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  componentDidMount() {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if (res == null) {
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));

    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-AC-VW01");
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
  };

  // editDialogOpen = (row) => {
  //   this.setState({
  //     grEditOpen : true,
  //     selectedRow : row
  //   });
  // }
  handleDeleteClick = (e) => {
    if (Object.keys(this.state.selectedRow).length === 0) {
      alert("Please select a Host Threshold");
      return;
    } else {
      this.setState({
        confirmOpen: true,
      });
    }
  };

  confirmed = (result) => {
    this.setState({ confirmOpen: false });

    //show progress loading...
    this.setState({ openProgress: true });

    if (result) {
      const url = `/settings/threshold`;

      const data = {
        cluster: this.state.selectedRow.cluster_name,
        node: this.state.selectedRow.node_name,
      };

      axios
        .delete(url, { data: data })
        .then((res) => {
          alert(res.data.message);
          this.setState({ open: false });
          this.handleClose();
          this.onUpdateData();
        })
        .catch((err) => {});

      this.setState({ openProgress: false });

      // loging Add Node
      let userId = null;
      AsyncStorage.getItem("userName", (err, result) => {
        userId = result;
      });
      utilLog.fn_insertPLogs(userId, "log-ND-MD02");
    } else {
      this.setState({ openProgress: false });
    }
  };

  handleClick = (event) => {
    if (this.state.anchorEl === null) {
      this.setState({ anchorEl: event.currentTarget });
    } else {
      this.setState({ anchorEl: null });
    }
  };

  handleClose = () => {
    this.setState({
      anchorEl: null,
      selection: [],
      selectedRow: "",
    });
  };

  render() {
    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({
        selectedRow: this.state.rows[selection[0]]
          ? this.state.rows[selection[0]]
          : {},
        confrimTarget: this.state.rows[selection[0]]
          ? this.state.rows[selection[0]].group_name
          : "false",
      });
    };

    const open = Boolean(this.state.anchorEl);

    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
        }}
      />
    );

    const Row = (props) => {
      return <Table.Row {...props} key={props.tableRow.key} />;
    };

    const Cell = (props) => {
      const { column, row } = props;

      // <WarningRoundedIcon style={{ fontSize: "8px", marginRight: "5px" }} />
      if (column.name === "total") {
        return (
          <Table.Cell
            {...props}
            style={{ cursor: "pointer" }}
          ><Link to={{
            pathname: `/settings/meterings/bill/${row.date}`,
            state: {
              data : row
            }
          }}>$ {props.value}</Link></Table.Cell>
        );
      }
      return <Table.Cell>{props.value}</Table.Cell>;
    };

    return (
      <div className="sub-content-wrapper fulled">
        <Confirm2
          confirmInfo={this.state.confirmInfo}
          confrimTarget={this.state.confrimTarget}
          confirmTargetKeyname={this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}
        />
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
                    onClick={this.handleClick}
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
                              <ThCreateThreshold
                                rowDatas={this.state.rows}
                                onUpdateData={this.onUpdateData}
                                menuClose={this.handleClose}
                              />
                            </MenuItem>
                            <MenuItem
                              style={{
                                textAlign: "center",
                                display: "block",
                                fontSize: "14px",
                              }}
                            >
                              <ThEditThreshold
                                rowData={this.state.selectedRow}
                                onUpdateData={this.onUpdateData}
                                menuClose={this.handleClose}
                              />
                            </MenuItem>
                            <MenuItem
                              style={{
                                textAlign: "center",
                                display: "block",
                                fontSize: "14px",
                              }}
                            >
                              <div onClick={this.handleDeleteClick}>
                                Delete Threshold
                              </div>
                            </MenuItem>
                          </MenuList>
                        </Paper>
                      </Grow>
                    )}
                  </Popper>
                </div>,
                <Grid rows={this.state.rows} columns={this.state.columns}>
                  <Toolbar />
                  <SearchState defaultValue="" />
                  <SearchPanel style={{ marginLeft: 0 }} />

                  <PagingState
                    defaultCurrentPage={0}
                    defaultPageSize={this.state.pageSize}
                  />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  <SortingState
                    defaultSorting={[
                      { columnName: "user_id", direction: "asc" },
                    ]}
                  />

                  <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  />

                  <IntegratedFiltering />
                  <IntegratedSelection />
                  <IntegratedSorting />
                  <IntegratedPaging />

                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing
                    defaultColumnWidths={this.state.defaultColumnWidths}
                  />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />

                  <TableSelection
                    selectByRowClick
                    highlightRow
                    // showSelectionColumn={false}
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

export default BillList;