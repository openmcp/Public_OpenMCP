import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
  TextField,
} from "@material-ui/core";
import {
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  SelectionState,
  IntegratedSelection,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
import axios from "axios";
import LensIcon from "@material-ui/icons/Lens";

import Stepper from "@material-ui/core/Stepper";
import Step from "@material-ui/core/Step";
import StepLabel from "@material-ui/core/StepLabel";

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
  backButton: {
    marginRight: theme.spacing(1),
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
});

class ThCreateThreshold extends Component {
  constructor(props) {
    super(props);
    this.state = {
      // clusterName: "",
      // nodeName: "",
      cpuWarn: 0,
      cpuDanger: 0,
      ramWarn: 0,
      ramDanger: 0,
      storageWarn: 0,
      stroageDanger: 0,
      open: false,

      preRows:this.props.rowDatas,

      rows: [],

      // groupName : "",
      // description : "",

      selectedCluster: [],
      clusterSelectionId: [],

      selectedNode: [],
      NodeSelectionId: [],

      // selectedRoleIds: [],
      // roleSelectionId: [],

      // selectedUserIds : [],
      // userSelectionId : [],

      // selectedProjects : [],
      // projectSelectionId : [],

      activeStep: 0,
    };
    // this.onChange = this.onChange.bind(this);
  }

  onChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value,
    });
  };

  componentWillMount() {}

  handleClickOpen = () => {
    this.setState({
      open: true,
    });
  };

  handleClose = () => {
    this.setState({
      clusterName: "",
      nodeName: "",
      cpuWarn: 0,
      cpuDanger: 0,
      ramWarn: 0,
      ramDanger: 0,
      storageWarn: 0,
      stroageDanger: 0,
      open: false,
    });
    this.props.menuClose();
  };

  handleSave = (e) => {
    if (Object.keys(this.state.selectedCluster).length === 0) {
      alert("Please select 'Cluster'");
      return;
    } else if (Object.keys(this.state.selectedNode).length === 0) {
      alert("Please select 'Node'");
      return;
    } else if (this.state.cpuWarn === 0) {
      alert("Please set 'cpu warning threshold (%)'");
      return;
    } else if (this.state.cpuDanger === 0) {
      alert("Please set 'cpu danger threshold (%)'");
      return;
    } else if (this.state.ramWarn === 0) {
      alert("Please set 'ram warning threshold (%)'");
      return;
    } else if (this.state.ramDanger === 0) {
      alert("Please set 'ram danger threshold (%)'");
      return;
    } else if (this.state.stroageWarn === 0) {
      alert("Please set 'storage warning threshold (%)'");
      return;
    } else if (this.state.stroageDanger === 0) {
      alert("Please set 'storage danger threshold (%)'");
      return;
    }

    // insert host threshold
    const url = `/settings/threshold`;
    const data = {
      nodeName: this.state.selectedNode[0],
      clusterName: this.state.selectedCluster[0],
      cpuWarn: this.state.cpuWarn,
      cpuDanger: this.state.cpuDanger,
      ramWarn: this.state.ramWarn,
      ramDanger: this.state.ramDanger,
      storageWarn: this.state.storageWarn,
      stroageDanger: this.state.stroageDanger,
    };
    axios
      .post(url, data)
      .then((res) => {
        alert(res.data.message);
        this.setState({ open: false });
        this.props.menuClose();
        this.props.onUpdateData();
      })
      .catch((err) => {
        alert(err);
      });

    // loging 호스트 임계설정
    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-TH-TH01"); //호스트 임계설정(insert)

    //close modal popup
    this.setState({ open: false });
  };

  onSelectCluster = (rows, selectionId) => {
    this.setState({
      selectedCluster: [rows.length > 0 ? rows[0].name : ""],
      clusterSelectionId: selectionId,
    });
  };

  onSelectNode = (rows, selectionId) => {
    this.setState({
      selectedNode: [rows.length > 0 ? rows[0].name : ""],
      NodeSelectionId: selectionId,
    });
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

    const steps = ["Select Cluster", "Select Node(host)", "Set Threshold"];
    const handleNext = () => {
      switch (this.state.activeStep) {
        case 0:
          if (Object.keys(this.state.selectedCluster).length === 0) {
            alert("Please select 'Cluster'");
            return;
          } else {
            this.setState({ activeStep: this.state.activeStep + 1 });
            return;
          }
        case 1:
          if (Object.keys(this.state.selectedNode).length === 0) {
            alert("Please select 'Node'");
            return;
          } else {
            this.setState({ activeStep: this.state.activeStep + 1 });
            return;
          }
        case 2:
          if (this.state.cpuWarn === 0) {
            alert("Please set 'cpu warning threshold (%)'");
            return;
          } else if (this.state.cpuDanger === 0) {
            alert("Please set 'cpu danger threshold (%)'");
            return;
          } else if (this.state.ramWarn === 0) {
            alert("Please set 'ram warning threshold (%)'");
            return;
          } else if (this.state.ramDanger === 0) {
            alert("Please set 'ram danger threshold (%)'");
            return;
          } else if (this.state.stroageWarn === 0) {
            alert("Please set 'storage warning threshold (%)'");
            return;
          } else if (this.state.stroageDanger === 0) {
            alert("Please set 'storage danger threshold (%)'");
            return;
          } else {
            this.setState({ activeStep: this.state.activeStep + 1 });
            return;
          }
        default:
          return;
      }
    };

    const handleBack = () => {
      this.setState({ activeStep: this.state.activeStep - 1 });
    };


    return (
      <div>
        <div
          onClick={this.handleClickOpen}
          style={{
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Create Threshold
        </div>
        <Dialog
          // onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Create Host Threshold
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body small-grid">
              <Stepper activeStep={this.state.activeStep} alternativeLabel>
                {steps.map((label) => (
                  <Step key={label}>
                    <StepLabel>{label}</StepLabel>
                  </Step>
                ))}
              </Stepper>
              <div>
                <Typography>
                  {this.state.activeStep === 0 ? (
                    <section className="md-content">
                      <ThCluseter
                        selection={this.state.clusterSelectionId}
                        onSelectCluster={this.onSelectCluster}
                      />
                    </section>
                  ) : this.state.activeStep === 1 ? (
                    <section className="md-content">
                      <ThNode
                        preRows = {this.state.preRows}
                        selectedCluster={this.state.selectedCluster}
                        selection={this.state.NodeSelectionId}
                        onSelectNode={this.onSelectNode}
                      />
                    </section>
                  ) : (
                    <div>
                      <section className="md-content">
                        <div style={{ display: "flex", marginTop: "10px" }}>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%", marginRight: "10px" }}
                          >
                            <p>CPU Threshold (Warning)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.cpuWarn}
                              fullWidth={true}
                              name="cpuWarn"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%" }}
                          >
                            <p>CPU Threshold (Danger)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.cpuDanger}
                              fullWidth={true}
                              name="cpuDanger"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                        </div>
                      </section>
                      <section className="md-content">
                        <div style={{ display: "flex", marginTop: "10px" }}>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%", marginRight: "10px" }}
                          >
                            <p>Memory Threshold (Warning)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.ramWarn}
                              fullWidth={true}
                              name="ramWarn"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%" }}
                          >
                            <p>Memory Threshold (Danger)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.ramDanger}
                              fullWidth={true}
                              name="ramDanger"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                        </div>
                      </section>
                      <section className="md-content">
                        <div style={{ display: "flex", marginTop: "10px" }}>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%", marginRight: "10px" }}
                          >
                            <p>Storage Threshold (Warning)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.storageWarn}
                              fullWidth={true}
                              name="storageWarn"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                          <div
                            className="props pj-pc-textfield"
                            style={{ width: "50%" }}
                          >
                            <p>Storage Threshold (Danger)</p>
                            <TextField
                              id="outlined-multiline-static"
                              rows={1}
                              type="number"
                              placeholder="threshold rate"
                              variant="outlined"
                              value={this.state.stroageDanger}
                              fullWidth={true}
                              name="stroageDanger"
                              onChange={this.onChange}
                            />
                            <span style={{bottom: "8px"}}>%</span>
                          </div>
                        </div>
                      </section>
                    </div>
                  )}
                </Typography>
              </div>
            </div>
          </DialogContent>
          <DialogActions>
            <div>
              <Button
                disabled={this.state.activeStep === 0}
                onClick={handleBack}
              >
                Back
              </Button>
              {this.state.activeStep === steps.length - 1 ? (
                <Button onClick={this.handleSave} color="primary">
                  save
                </Button>
              ) : (
                <Button color="primary" onClick={handleNext}>
                  next
                </Button>
              )}
            </div>
            {/* <Button onClick={this.handleSave} color="primary">
              save
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

class ThCluseter extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
        { name: "nodes", title: "Nodes" },
        { name: "cpu", title: "CPU" },
        { name: "ram", title: "Memory" },
        { name: "provider", title: "Provider" },
        // { name: "disk", title: "Disk" },
        // { name: "network", title: "Network" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 100 },
        { columnName: "status", width: 100 },
        { columnName: "region", width: 100 },
        { columnName: "zone", width: 80 },
        { columnName: "nodes", width: 100 },
        { columnName: "cpu", width: 150 },
        { columnName: "ram", width: 150 },
        { columnName: "provider", width: 150 },
        // { columnName: "disk", width: 150 },
        // { columnName: "network", width: 150 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],
      rows: "",
      selection: this.props.selection,
      selectedRow: [],
      completed: 0,
    };
  }

  callApi = async () => {
    const response = await fetch("/clusters");
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  componentWillMount() {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        this.setState({ rows: res });
        let selectedRows = [];
        this.props.selection.forEach((id) => {
          selectedRows.push(res[id]);
        });
        this.setState({ selectedRow: selectedRows });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

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

    // const onSelectionChange = (selection) => {
    //   if (selection.length > 1) selection.splice(0, 1);
    //   this.setState({ selection: selection });
    //   this.setState({
    //     selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {},
    //     confrimTarget : this.state.rows[selection[0]] ? this.state.rows[selection[0]].group_name : "false" ,
    //   });

    // };
    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      let selectedRows = [];

      selection.forEach((id) => {
        selectedRows.push(this.state.rows[id]);
      });
      this.setState({ selectedRow: selectedRows });
      this.setState({ selection: selection });

      this.props.onSelectCluster(selectedRows, selection);
    };

    return (
      <div>
        <p>Selected Clusters</p>
        <div
          id="md-content-info"
          style={{ display: "block", minHeight: "40px", marginBottom: "10px" }}
        >
          {this.state.selectedRow.length > 0 ? (
            this.state.selectedRow.map((row) => {
              return (
                <div>
                  <span>
                    <LensIcon style={{ fontSize: "8px", marginRight: "5px" }} />
                    {row.name}
                  </span>
                  {/* <div>
                    <span>{row.cpu}</span><span>{row.ram}</span>
                  </div> */}
                </div>
              );
            })
          ) : (
            <div
              style={{
                color: "#9a9a9a",
                textAlign: "center",
              }}
            >
              Please Select Cluster
            </div>
          )}
        </div>
        {/* <p>Select Role</p> */}
        <Paper>
          {this.state.rows ? (
            [
              <Grid rows={this.state.rows} columns={this.state.columns}>
                {/* <Toolbar /> */}
                {/* 검색 */}
                {/* <SearchState defaultValue="" />
          <SearchPanel style={{ marginLeft: 0 }} /> */}

                {/* Sorting */}
                <SortingState
                  defaultSorting={[{ columnName: "status", direction: "asc" }]}
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
                <Table />
                <TableColumnResizing
                  defaultColumnWidths={this.state.defaultColumnWidths}
                />
                <TableHeaderRow showSortingControls rowComponent={HeaderRow} />
                <TableSelection selectByRowClick highlightRow />
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
    );
  }
}

class ThNode extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Node" },
        { name: "status", title: "Status" },
        { name: "cluster", title: "Cluster" },
        { name: "role", title: "Role" },
        { name: "system_version", title: "System Version" },
        { name: "cpu", title: "CPU" },
        { name: "memory", title: "Memory" },
        { name: "pods", title: "Pods" },
        { name: "provider", title: "Provider" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 250 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 100 },
        { columnName: "role", width: 75 },
        { columnName: "system_version", width: 200 },
        { columnName: "cpu", width: 155 },
        { columnName: "memory", width: 155 },
        { columnName: "pods", width: 155 },
        { columnName: "provider", width: 100 },
        { columnName: "region", width: 90 },
        { columnName: "zone", width: 80 },
      ],
      rows: "",
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],
      selection: this.props.selection,
      selectedRow: [],
      completed: 0,
    };
  }

  callApi = async () => {
    const response = await fetch(
      `/clusters/${this.props.selectedCluster[0]}/nodes`
    );
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  componentWillMount() {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((result) => {
        var res = [];
        result.forEach((item)=>{
          var isExsits = false;
          
          this.props.preRows.forEach((pre)=>{
            if(item.cluster === pre.cluster_name && item.name === pre.node_name){
              isExsits = true;
            }
          })
          if (!isExsits){
            res.push(item);
          }
        })


        this.setState({ rows: res });
        let selectedRows = [];
        this.props.selection.forEach((index) => {
          selectedRows.push(res[index]);
        });
        this.setState({ selectedRow: selectedRows });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

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

    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      let selectedRows = [];

      selection.forEach((id) => {
        selectedRows.push(this.state.rows[id]);
      });
      this.setState({ selectedRow: selectedRows });
      this.setState({ selection: selection });

      this.props.onSelectNode(selectedRows, selection);
    };

    return (
      <div>
        <p>Selected Node</p>
        <div
          id="md-content-info"
          style={{ display: "block", minHeight: "40px", marginBottom: "10px" }}
        >
          {this.state.selectedRow.length > 0 ? (
            this.state.selectedRow.map((row) => {
              return (
                <span>
                  <LensIcon style={{ fontSize: "8px", marginRight: "5px" }} />
                  {row.name}
                </span>
              );
            })
          ) : (
            <div
              style={{
                color: "#9a9a9a",
                textAlign: "center",
              }}
            >
              Please Select Node
            </div>
          )}
        </div>
        {/* <p>Select Role</p> */}
        <Paper>
          {this.state.rows ? (
            [
              <Grid rows={this.state.rows} columns={this.state.columns}>
                {/* <Toolbar /> */}
                {/* 검색 */}
                {/* <SearchState defaultValue="" />
            <SearchPanel style={{ marginLeft: 0 }} /> */}

                {/* Sorting */}
                <SortingState
                  defaultSorting={[{ columnName: "status", direction: "asc" }]}
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
                <Table />
                <TableColumnResizing
                  defaultColumnWidths={this.state.defaultColumnWidths}
                />
                <TableHeaderRow showSortingControls rowComponent={HeaderRow} />
                <TableSelection selectByRowClick highlightRow />
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
    );
  }
}


export default ThCreateThreshold;
