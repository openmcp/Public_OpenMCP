import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
// import SelectBox from "../../modules/SelectBox";
// import { Link } from "react-router-dom";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
import {
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  SelectionState,
  IntegratedSelection,
  TableColumnVisibility,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  // TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
  // TableFixedColumns,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
// import axios from "axios";
// import Typography from "@material-ui/core/Typography";
// import DialogActions from "@material-ui/core/DialogActions";
// import DialogContent from "@material-ui/core/DialogContent";
// import Button from "@material-ui/core/Button";
// import Dialog from "@material-ui/core/Dialog";
// import IconButton from "@material-ui/core/IconButton";
// import axios from 'axios';
// import { ContactlessOutlined } from "@material-ui/icons";
// import Confirm from './../../modules/Confirm';
import Confirm2 from "./../../modules/Confirm2";

const styles = (theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
});

class SnapShotControl extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "SnapShot" },
        { name: "created_time", title: "Created Time" },
        { name: "control", title: " " },
      ],
      tableColumnExtensions: [
        { columnName: "name", width: "45%" },
        { columnName: "created_time", width: "30%" },
        { columnName: "control", align: "center" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 400 },
        { columnName: "created_time", width: 300 },
        { columnName: "control", width: 100 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],

      open: false,
      account: "",
      account_role: "",
      rows: [],

      selection: [],
      selectedRow: "",
      rightColumns: ["control"],

      confirmOpen: false,
      confirmInfo: {
        title: "Cluster Join Confrim",
        context: "Are you sure you want to Join the Cluster?",
        button: {
          open: "",
          yes: "JOIN",
          no: "CANCEL",
        },
      },
      confrimTarget: "false",
      confirmTargetKeyname: "snapshot",
    };
    // this.onChange = this.onChange.bind(this);
  }

  callApi = async () => {
    const response = await fetch(`/snapshots`);
    const body = await response.json();
    return body;
  };

  componentWillMount() {
    // this.callApi()
    //   .then((res) => {
    //     if(res == null) {
    //       this.setState({ rows: [] });
    //     } else {
    //       this.setState({ rows: res });
    //     }
    //   })
    //   .catch((err) => console.log(err));
  }

  // onChange(e) {
  //   console.log("onChangedd");
  //   this.setState({
  //     [e.target.name]: e.target.value,
  //   });
  // }

  //snapshot open 버튼
  handleClickOpen = () => {
    if (Object.keys(this.props.rowData).length  === 0) {
      alert("Please select deployement");
      this.setState({ open: false });
      return;
    }

    this.setState({ open: true });

    this.callApi()
      .then((res) => {
        console.log(res);
        this.setState({ rows: res });
      })
      .catch((err) => console.log(err));
  };

  handleClose = () => {
    this.setState({
      account: "",
      role_id: "",
      open: false,
    });
    this.props.menuClose();
  };

  handleSave = (e) => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select snapshot");
      return;
    }

    // loging deployment migration
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PJ-MD01");

    //close modal popup
    this.setState({ open: false });
    this.props.menuClose();
  };

  onSnapshotDelete = (data) => {
    console.log("Delete snapshot", data);
    // alert("Delete snapshot", data)
    this.setState({
      confirmOpen: true,
      confirmInfo: {
        title: "Snapshot Delete",
        context: "Are you sure you want to Snapshot Delete?",
        button: {
          open: "",
          yes: "Delete",
          no: "Cancel",
        },
      },
      confrimTarget: data,
      confirmTargetKeyname: "snapshot",
    });
  };

  onSnapshotRevert = (data) => {
    // console.log("Revert snapshot",data)

    this.setState({
      confirmOpen: true,
      confirmInfo: {
        title: "Snapshot Revert",
        context: "Are you sure you want to Revert?",
        button: {
          open: "",
          yes: "Revert",
          no: "Cancel",
        },
      },
      confrimTarget: data,
      confirmTargetKeyname: "snapshot",
    });
    // alert("Revert snapshot", data)
  };

  Cell = (props) => {
    const { column, row } = props;
    if (column.name === "control") {
      return (
        <Table.Cell
          {...props}
          style={{
            borderRight: "1px solid #e0e0e0",
            borderLeft: "1px solid #e0e0e0",
            textAlign: "center",
            background: "whitesmoke",
          }}
        >
          <div className="snapshot">
            <span
              className="revert"
              style={{ cursor: "pointer", display: "inline-block" }}
              onClick={() => this.onSnapshotRevert(row.name)}
            >
              Revert
            </span>
            <span style={{ margin: "0 5px" }}> | </span>
            <span
              className="delete"
              style={{ cursor: "pointer", display: "inline-block" }}
              onClick={() => this.onSnapshotDelete(row.name)}
            >
              Delete
            </span>
          </div>
        </Table.Cell>
      );
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

  // callback function
  confirmed = (result) => {
    if (result) {
      //Unjoin proceed
      console.log("confirmed");
      let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
      utilLog.fn_insertPLogs(userId, "log-CL-MO03");
    } else {
      console.log("cancel");
    }
    this.setState({ confirmOpen: false, open:false });
  };

  render() {
    const DialogTitle = withStyles(styles)((props) => {
      const { children, classes, onClose, ...other } = props;
      return (
        <MuiDialogTitle disableTypography className={classes.root} {...other}>
          <Typography variant="h6">{children}</Typography>
          {onClose ? (
            <IconButton
              aria-label="close"
              className={classes.closeButton}
              onClick={onClose}
            >
              <CloseIcon />
            </IconButton>
          ) : null}
        </MuiDialogTitle>
      );
    });

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
      return <Table.Row {...props} key={props.tableRow.key} />;
    };

    //셀
    // const Cell = (props) => {
    //   const { column, row } = props;
    //   if (column.name === "control") {
    //     return (
    //       <Table.Cell
    //         {...props}
    //         style={{ borderRight:"1px solid #e0e0e0", borderLeft:"1px solid #e0e0e0", textAlign:"center",background:"whitesmoke"}}
    //       >
    //         <div onClick={this.onSnapshotRevert(props)}>
    //           <span style={{cursor:"pointer", display:"inline-block"}} onClick={this.onSnapshotRevert(props)}>Revert</span>
    //           <span> | </span>
    //           <span style={{cursor:"pointer", display:"inline-block"}} onClick={this.onSnapshotDelete(props)}>Delete</span>
    //         </div>
    //       </Table.Cell>
    //     );
    //   }
    //   return <Table.Cell>{props.value}</Table.Cell>;
    // };
    const onSelectionChange = (selection) => {
      // console.log(selection);
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({
        selectedRow: this.state.rows[selection[0]]
          ? this.state.rows[selection[0]]
          : {},
      });
    };

    return (
      <div>
        <div
          // variant="outlined"
          // color="primary"
          onClick={this.handleClickOpen}
          style={{
            // position: "absolute",
            // right: "272px",
            // top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          SnapShot
        </div>
        {/* <Button
          variant="outlined"
          color="primary"
          onClick={this.handleClickOpen}
          style={{
            // position: "absolute",
            // right: "272px",
            // top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          SnapShot
        </Button> */}
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Snapshots
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body">
              {/* <section className="md-content">
                <p>User Info</p>
                <div id="md-content-info">
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>UserID : </strong></span>
                      <span>{this.state.account}</span>
                    </div>
                  </div>
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>Current Role : </strong></span>
                      <span>{this.state.account_role}</span>
                    </div>
                  </div>
                </div>
              </section> */}
              <section className="md-content">
                <p>Snapshot List</p>
                {/* cluster selector */}
                <Paper>
                  <Confirm2
                    confirmInfo={this.state.confirmInfo}
                    confrimTarget={this.state.confrimTarget}
                    confirmTargetKeyname={this.state.confirmTargetKeyname}
                    confirmed={this.confirmed}
                    confirmOpen={this.state.confirmOpen}
                  />
                  <Grid rows={this.state.rows} columns={this.state.columns}>
                    {/* <Toolbar /> */}
                    {/* 검색 */}
                    {/* <SearchState defaultValue="" />
                  <SearchPanel style={{ marginLeft: 0 }} /> */}

                    {/* Sorting */}
                    <SortingState
                      defaultSorting={[
                        { columnName: "status", direction: "asc" },
                      ]}
                    />

                    {/* 페이징 */}
                    <PagingState
                      defaultCurrentPage={0}
                      defaultPageSize={this.state.pageSize}
                    />
                    <PagingPanel pageSizes={this.state.pageSizes} />
                    <SelectionState
                      selection={this.state.selection}
                      onSelectionChange={onSelectionChange}
                    />

                    <IntegratedFiltering />
                    <IntegratedSorting />
                    <IntegratedSelection />
                    <IntegratedPaging />

                    {/* 테이블 */}
                    <Table
                      cellComponent={this.Cell}
                      rowComponent={Row}
                      columnExtensions={this.state.tableColumnExtensions}
                    />
                    {/* <TableColumnResizing
                    // defaultColumnWidths={this.state.defaultColumnWidths}
                  /> */}
                    <TableHeaderRow
                      showSortingControls
                      rowComponent={HeaderRow}
                    />
                    <TableColumnVisibility
                      defaultHiddenColumnNames={["role_id"]}
                    />
                    <TableSelection
                      // selectByRowClick
                      highlightRow
                      // showSelectionColumn={false}
                    />
                    {/* <TableFixedColumns
                    cellComponent={Cell}
                    rightColumns={this.state.rightColumns}
                  /> */}
                  </Grid>
                </Paper>
              </section>
            </div>
            {/* <div className="pj-create">
              <div className="create-content">
                <p>Deployment</p>
                <TextField
                  id="outlined-multiline-static"
                  label="name"
                  rows={1}
                  variant="outlined"
                  defaultValue = {this.props.rowData.name}
                  fullWidth	={true}
                  name="deployment_name"
                  InputProps={{
                    readOnly: true,
                  }}
                />
                <p className="pj-cluster">Cluster</p>
                <SelectBox className="selectbox" rows={this.state.selectBoxData} onSelectBoxChange={onSelectBoxChange}  defaultValue={this.state.cluster}></SelectBox>
              </div>
            </div> */}
          </DialogContent>
          <DialogActions>
            {/* <Button onClick={this.handleSave} color="primary">
              Take a Snapshot
            </Button> */}
            <Button onClick={this.handleClose} color="primary">
              cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default SnapShotControl;
