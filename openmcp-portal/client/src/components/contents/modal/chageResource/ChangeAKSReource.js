import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import * as utilLog from "../../../util/UtLogs.js";
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
  PagingState,
  SortingState,
  SelectionState,
  IntegratedFiltering,
  IntegratedPaging,
  IntegratedSorting,
  IntegratedSelection,
  // SearchState,
  // EditingState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
  // SearchPanel,
  // Toolbar,
  // TableEditRow,
  // TableEditColumn,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
import Confirm2 from './../../../modules/Confirm2';
import ProgressTemp from './../../../modules/ProgressTemp';
import axios from 'axios';

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

class ChangeAKSReource extends Component {
  constructor(props) {
    super(props);
    this.state = {

      columns: [
        { name: "code", title: "Type" },
        { name: "etc", title: "Resources" },
        { name: "description", title: "Tier" },
      ],
      defaultColumnWidths: [
        { columnName: "code", width: 200 },
        { columnName: "etc", width: 200 },
        { columnName: "description", width: 300 },
      ],


      aPoolColumns: [
        { name: "name", title: "Pool" },
        { name: "vmssname", title: "VmssName" },
        { name: "nodecount", title: "NodeCount" },
      ],
      defaultColumnWidths2: [
        { columnName: "name", width: 200 },
        { columnName: "vmssname", width: 300 },
        { columnName: "nodecount", width: 150 },
      ],

      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 3,
      pageSizes: [3, 6, 12, 0],

      instTypes: [],
      clusterInfo:[],
      selection: [],
      selectedRow: "",
      selection2: [],
      selectedRow2: "",

      open: false,

      confirmOpen: false,
      confirmInfo : {
        title :"Change AKS Pool Resource",
        context :"Are you sure you want to Change AKS Pool Resources?",
        button : {
          open : "",
          yes : "CONFIRM",
          no : "CANCEL",
        }
      },
      confrimTarget : "",
      confirmTargetKeyname:""
    };
  }

  componentWillMount() {
    this.callApi()
      .then((res) => {
        this.setState({ instTypes: res });
      })
      .catch((err) => console.log(err));
    

    this.callApi2()
      .then((res) => {
        this.setState({ clusterInfo: res });
      })
      .catch((err) => console.log(err));
  }

  initState = () => {
    this.setState({
      selection: [],
      selectedRow: "",
      selection2: [],
      selectedRow2: "",
    });
  };

  callApi = async () => {
    const response = await fetch("/azure/aks-type");
    const body = await response.json();
    return body;
  };

  callApi2 = async () => {
    // const response = await fetch(`/azure/pool/${this.props.clusterInfo.name}`);
    const response = await fetch(`/azure/pool/aks-cluster-01`);
    const body = await response.json();
    return body;
  }


  handleClickOpen = () => {
    this.setState({
      open: true,
    });
  };

  handleClose = () => {
    this.initState();
    this.setState({
      open: false,
    });
  };

  handleSaveClick = () => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select Instance Type");
      return;
    } else if(Object.keys(this.state.selectedRow2).length  === 0) {
      alert("Please select Agent Pool");
      return;
    }
    else {
      this.setState({
        confirmOpen: true,
      })
    }
  };

  confirmed = (result) => {
    this.setState({confirmOpen:false});

    //show progress loading...
    this.setState({openProgress:true});

    if(result) {
      const url = `/clusters/aks/change`;
      const data = {
        // cluster : this.props.clusterInfo.name,
        cluster : "aks-cluster-01",
        type :  this.state.selectedRow.code,
        tier :  this.state.selectedRow.description,
        poolName : this.state.selectedRow2.name,
      };

      axios.post(url, data)
        .then((res) => {
          if(res.data.error){
            alert(res.data.message)
            return
          }
          this.setState({open: false, openProgress:false})
        })
        .catch((err) => {
        });
        
      // loging Add Node
      let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
      utilLog.fn_insertPLogs(userId, "log-ND-MD02");
    } else {
      this.setState({open: false, openProgress:false})
    }
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

  onSelectionChange = (selection) => {
    if (selection.length > 1) selection.splice(0, 1);
    this.setState({ selection: selection });
    this.setState({
      selectedRow: this.state.instTypes[selection[0]]
        ? this.state.instTypes[selection[0]]
        : {},
    });
  };

  onSelectionChange2 = (selection) => {
    if (selection.length > 1) selection.splice(0, 1);
    this.setState({ selection2: selection });
    this.setState({
      selectedRow2: this.state.clusterInfo.agentpools[selection[0]]
        ? this.state.clusterInfo.agentpools[selection[0]]
        : {},
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

    return (
      <div>
        <Button variant="outlined" color="primary" onClick={this.handleClickOpen} style={{position:"absolute", right:"0px", top:"0px", zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
          Resource Config
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
           Pool Resource Configration
          </DialogTitle>
          <DialogContent dividers>
            <div>
              {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""}

              <Confirm2
                confirmInfo={this.state.confirmInfo} 
                confrimTarget ={this.state.confrimTarget} 
                confirmTargetKeyname = {this.state.confirmTargetKeyname}
                confirmed={this.confirmed}
                confirmOpen={this.state.confirmOpen}/>
                
              <div className="md-contents-body">
                <section className="md-content">
                  {/* deployment informations */}
                  <p>AKS Cluster Info</p>
                  <div id="md-content-info">
                    <div class="md-partition">
                      <div class="md-item">
                        <span><strong>Name : </strong></span>
                        <span>{this.props.clusterInfo.name}</span>
                      </div>
                      <div class="md-item">
                        <span><strong>Provider : </strong></span>
                        <span>{this.props.clusterInfo.provider}</span>
                      </div>
                      <div class="md-item">
                        <span><strong>Resource Group : </strong></span>
                        <span>{this.state.clusterInfo.resourcegroup}</span>
                      </div>
                      <div class="md-item">
                        <span><strong>Location : </strong></span>
                        <span>{this.state.clusterInfo.location}</span>
                      </div>
                    </div>
                    <div class="md-partition">
                      <div class="md-item">
                        <span><strong>Status : </strong></span>
                        <span>{this.props.clusterInfo.status}</span>
                      </div>
                      <div class="md-item">
                        <span><strong>Region : </strong></span>
                        <span>{this.props.clusterInfo.region}</span>
                      </div>
                      <div class="md-item">
                        <span><strong>Node Resource Group : </strong></span>
                        <span>{this.state.clusterInfo.resourcegroup}</span>
                      </div>
                    </div>
                  </div>
                </section>
                
                <section className="md-content">
                  <div>
                    <p>Agents Pool</p>
                    {/* cluster selector */}
                    <Paper>
                      <Grid
                        rows={this.state.clusterInfo.agentpools}
                        columns={this.state.aPoolColumns}
                      >
                        {/* Sorting */}
                        <SortingState
                          defaultSorting={[
                            { columnName: "name", direction: "asc" },
                          ]}
                        />

                        {/* 페이징 */}
                        <PagingState
                          defaultCurrentPage={0}
                          defaultPageSize={this.state.pageSize}
                        />
                        <PagingPanel pageSizes={this.state.pageSizes} />
                        <SelectionState
                          selection={this.state.selection2}
                          onSelectionChange={this.onSelectionChange2}
                        />

                        <IntegratedFiltering />
                        <IntegratedSorting />
                        <IntegratedSelection />
                        <IntegratedPaging />

                        {/* 테이블 */}
                        <Table />
                        <TableColumnResizing
                          defaultColumnWidths={
                            this.state.defaultColumnWidths2
                          }
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
                      </Grid>
                    </Paper>
                  </div>
                </section>
                <section className="md-content">
                  <div>
                    <p>Instance Type</p>
                    {/* cluster selector */}
                    <Paper>
                      <Grid
                        rows={this.state.instTypes}
                        columns={this.state.columns}
                      >
                        {/* Sorting */}
                        <SortingState
                          defaultSorting={[
                            { columnName: "code", direction: "asc" },
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
                          onSelectionChange={this.onSelectionChange}
                        />

                        <IntegratedFiltering />
                        <IntegratedSorting />
                        <IntegratedSelection />
                        <IntegratedPaging />

                        {/* 테이블 */}
                        <Table />
                        <TableColumnResizing
                          defaultColumnWidths={
                            this.state.defaultColumnWidths
                          }
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
                      </Grid>
                    </Paper>
                  </div>
                </section>
              </div>
            </div>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleSaveClick} color="primary">
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

export default ChangeAKSReource;