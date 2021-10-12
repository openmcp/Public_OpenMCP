import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import CircularProgress from "@material-ui/core/CircularProgress";
// import { Link } from "react-router-dom";
import {
  TextField,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
import {
  PagingState,
  SortingState,
  SelectionState,
  IntegratedFiltering,
  IntegratedPaging,
  IntegratedSorting,
  IntegratedSelection,
  // RowDetailState,
  SearchState,
  // EditingState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
  // TableRowDetail,
  SearchPanel,
  Toolbar,
  // TableEditRow,
  // TableEditColumn,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
// import AppBar from "@material-ui/core/AppBar";
// import Tabs from "@material-ui/core/Tabs";
// import Tab from "@material-ui/core/Tab";
// import { Container } from "@material-ui/core";
// import Box from "@material-ui/core/Box";
// import PropTypes from "prop-types";
import axios from 'axios';
// import SelectBox from "../../modules/SelectBox";

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

class PcAddProjectPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {

      podsCPURate:"",
      podsMemRate:"",
      clusterCPURate:"",
      clusterMemRate:"",
      

      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status" },
        { name: "cluster", title: "Cluster" },
        { name: "created_time", title: "Created Time" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 200 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 100 },
        { columnName: "created_time", width: 180 },
      ],

      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 3,
      pageSizes: [3, 6, 9, 0],

      open: false,

      rows: [],

      selection: [],
      selectedRow : "",

      value: 0,
    };
    // this.onChange = this.onChange.bind(this);
  }

  componentWillMount() {
  }

  initState = () => {
    this.setState({
      podsCPURate:"",
      podsMemRate:"",
      clusterCPURate:"",
      clusterMemRate:"",
      selection : [],
      selectedRow:"",
    });
  }

  callApi = async () => {
    const response = await fetch("/projects");
    const body = await response.json();
    return body;
  };

  onChange = (e) =>{
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  handleClickOpen = () => {
    this.initState();
    this.setState({ 
      open: true,
    });
    this.callApi()
    .then((res) => {
      this.setState({ rows: res });
    })
    .catch((err) => console.log(err));
  };

  handleClose = () => {
    this.initState();
    this.setState({
      open: false,
    });
  };

  handleSave = (e) => {
    if (Object.keys(this.state.selectedRow).length === 0) {
      alert("Please select project");
      return;
    } else if (this.state.clusterCPURate === "" || this.state.clusterMemRate === "" ||this.state.podsCPURate === ""||this.state.podsMemRate === ""){
      alert("Please insert threshold data");
      return;
    }


    const url = `/settings/policy/project-policy`;
    const data = {
      project:this.state.selectedRow.name,
      cluster:this.state.selectedRow.cluster,
      cls_cpu_trh_r:this.state.clusterCPURate,
      cls_mem_trh_r:this.state.clusterMemRate,
      pod_cpu_trh_r:this.state.podsCPURate,
      pod_mem_trh_r: this.state.podsMemRate,
    };
    axios.put(url, data)
    .then((res) => {
      this.props.onUpdateData();
    })
    .catch((err) => {
        alert(err);
    });

    
    // loging Add Project Policy
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PO-CR01");

    //close modal popup
    this.setState({ open: false });
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
        }}
      >
        <span
          style={{
            color:
              value === "Warning"
                ? "orange"
                : value === "Unschedulable"
                ? "red"
                : value === "Stop"
                ? "red"
                : value === "Running"
                ? "#1ab726"
                : "black",
          }}
        >
          {value}
        </span>
      </Table.Cell>
    );

    //셀
    const Cell = (props) => {
      const { column } = props;

      if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } 

      return <Table.Cell>{props.value}</Table.Cell>;
    };
    
    const Row = (props) => {
      // console.log("row!!!!!! : ",props);
      return <Table.Row {...props} key={props.tableRow.key} />;
    };

    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({ selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {} });
    };

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



    return (
      <div>
        <Button
          variant="outlined"
          color="primary"
          onClick={this.handleClickOpen}
          style={{
            position: "absolute",
            right: "26px",
            top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Add Policy
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Add Project Policy
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body add-node">
              <section className="md-content">
                <div className="outer-table">
                  <p>Select Project</p>
                  {/* cluster selector */}
                  <Paper>
                    {this.state.rows ? (
                      [
                        <Grid rows={this.state.rows} columns={this.state.columns}>
                          <Toolbar />
                          {/* 검색 */}
                          <SearchState defaultValue="" />

                          <SearchPanel style={{ marginLeft: 0 }} />

                          {/* Sorting */}
                          <SortingState
                          defaultSorting={[{ columnName: 'created_time', direction: 'desc' }]}
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

                          <IntegratedFiltering />
                          <IntegratedSorting />
                          <IntegratedSelection />
                          <IntegratedPaging />

                          {/* 테이블 */}
                          <Table cellComponent={Cell} rowComponent={Row} />
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
                </div>
              </section>
              <section className="md-content">
                <div style={{display:"flex", marginTop:"10px"}}>
                  <div className="props pj-pc-textfield" style={{width:"50%", marginRight:"10px"}}>
                    <p>Cluster CPU Threshold</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      type="number"
                      placeholder="cpu threshold rate"
                      variant="outlined"
                      value = {this.state.secretKey}
                      fullWidth	={true}
                      name="clusterCPURate"
                      onChange = {this.onChange}
                    />
                    <span>%</span>
                  </div>
                  <div className="props pj-pc-textfield" style={{width:"50%"}}>
                    <p>Cluster Memory Threshold</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      type="number"
                      placeholder="memory threshold rate"
                      variant="outlined"
                      value = {this.state.accessKey}
                      fullWidth	={true}
                      name="clusterMemRate"
                      onChange = {this.onChange}
                    />
                    <span>%</span>
                  </div>
                </div>
              </section>
              <section className="md-content">
                <div style={{display:"flex", marginTop:"10px"}}>
                  <div className="props pj-pc-textfield" style={{width:"50%", marginRight:"10px"}}>
                    <p>Pods CPU Threshold</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      type="number"
                      placeholder="cpu threshold rate"
                      variant="outlined"
                      value = {this.state.secretKey}
                      fullWidth	={true}
                      name="podsCPURate"
                      onChange = {this.onChange}
                    />
                    <span>%</span>
                  </div>
                  <div className="props pj-pc-textfield" style={{width:"51%"}}>
                    <p>Pods Memory Threshold</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      type="number"
                      placeholder="memory threshold rate"
                      variant="outlined"
                      value = {this.state.accessKey}
                      fullWidth	={true}
                      name="podsMemRate"
                      onChange = {this.onChange}
                    />
                    <span>%</span>
                  </div>
                </div>
              </section>
            </div>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleSave} color="primary">
              save
            </Button>
            <Button onClick={this.handleClose} color="primary">
              cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default PcAddProjectPolicy;