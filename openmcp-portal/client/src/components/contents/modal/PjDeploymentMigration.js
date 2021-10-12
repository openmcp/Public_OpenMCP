import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
// import SelectBox from "../../modules/SelectBox";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  // TextField,
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
  // SearchState,
  // EditingState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  // Toolbar,
  // SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  // TableEditRow,
  // TableEditColumn,
  TableSelection,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
// import Typography from "@material-ui/core/Typography";
// import DialogActions from "@material-ui/core/DialogActions";
// import DialogContent from "@material-ui/core/DialogContent";
// import Button from "@material-ui/core/Button";
// import Dialog from "@material-ui/core/Dialog";
// import IconButton from "@material-ui/core/IconButton";
import axios from 'axios';
// import { ContactlessOutlined } from "@material-ui/icons";

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

class PjDeploymentMigration extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status" },
        { name: "nodes", title: "nodes" },
        { name: "cpu", title: "CPU(%)" },
        { name: "ram", title: "Memory(%)" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
        
        // { name: "edit", title: "edit" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 130 },
        { columnName: "status", width: 130 },
        { columnName: "nodes", width: 130 },
        { columnName: "cpu", width: 100 },
        { columnName: "ram", width: 120 },
        { columnName: "region", width: 130 },
        { columnName: "zone", width: 130 },
        // { columnName: "edit", width: 170 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],

      open: false,
      dpName : "",
      dpStatus : "",
      dpImage : "",
      dpCluster : "",
      clusters: [],

      selection: [],
      selectedRow : "",
      YAML : `apiVersion: openmcp.k8s.io/v1alpha1
      kind: Migration
      metadata:
        name: migrations8
      spec:
        MigrationServiceSource:
        - SourceCluster: cluster1
          TargetCluster: cluster2
          NameSpace: testmig
          ServiceName: testim
          MigrationSource:
          - ResourceName: testim-dp
            ResourceType: Deployment
          - ResourceName: testim-sv
            ResourceType: Service
          - ResourceName: testim-pv
            ResourceType: PersistentVolume
          - ResourceName: testim-pvc
            ResourceType: PersistentVolumeClaim
      `
    };
    // this.onChange = this.onChange.bind(this);
  }

  callApi = async () => {
    const response = await fetch("/clusters");
    const body = await response.json();
    return body;
  };

  handleClickOpen = () => {
    // console.log("count:",Object.keys(this.props.rowData).length, this.props.rowData);
    if (Object.keys(this.props.rowData).length === 0) {
      alert("Please select deployment");
      this.setState({ open: false });
      return;
    }

    this.setState({ 
      open: true,
      dpName : this.props.rowData.name,
      dpStatus : this.props.rowData.status,
      dpImage : this.props.rowData.image,
      dpCluster : this.props.rowData.cluster,
      selection : [],
    });

    this.callApi()
      .then((res) => {
        this.setState({ clusters: res });
        // console.log(res[0])
        // this.setState({ cluster: res[0], firstValue: res[0] });
      })
      .catch((err) => console.log(err));
  };

  componentWillMount() {
    // console.log("Migration will mount");
    // cluster list를 가져오는 api 호출
    // this.callApi()
    //   .then((res) => {
    //     this.setState({ clusters: res });
    //     // console.log(res[0])
    //     // this.setState({ cluster: res[0], firstValue: res[0] });
    //   })
    //   .catch((err) => console.log(err));
  }

  // onChange(e) {
  //   console.log("onChangedd");
  //   this.setState({
  //     [e.target.name]: e.target.value,
  //   });
  // }

  handleClose = () => {
    this.setState({
      project_name: "",
      project_description: "",
      cluster: this.state.firstValue,
      open: false,
    });
    this.props.menuClose();
  };

  handleSave = (e) => {
    if (Object.keys(this.state.selectedRow).length === 0) {
      alert("Please select target cluster");
      return;
    
    } 
    // let YAML =`apiVersion: openmcp.k8s.io/v1alpha1
    // kind: Migration
    // metadata:
    //   name: migrations8
    // spec:
    //   MigrationServiceSource:
    //   - SourceCluster: cluster1
    //     TargetCluster: cluster2
    //     NameSpace: testmig
    //     ServiceName: testim
    //     MigrationSource:
    //     - ResourceName: testim-dp
    //       ResourceType: Deployment
    //     - ResourceName: testim-sv
    //       ResourceType: Service
    //     - ResourceName: testim-pv
    //       ResourceType: PersistentVolume
    //     - ResourceName: testim-pvc
    //       ResourceType: PersistentVolumeClaim`;
    // let YAML2 = `apiVersion: openmcp.k8s.io/v1alpha1
    // kind: Migration
    // metadata:
    //   name: migrations8
    // spec:
    //   MigrationServiceSource:
    //   - SourceCluster: ${this.state.dpCluster}
    //     TargetCluster: ${this.state.selectedRow.name}
    //     NameSpace: ${this.state.dpName}
    //     ServiceName: testim
    //     MigrationSource:
    //     - ResourceName: testim-dp
    //       ResourceType: Deployment
    //     - ResourceName: testim-sv
    //       ResourceType: Service
    //     - ResourceName: testim-pv
    //       ResourceType: PersistentVolume
    //     - ResourceName: testim-pvc
    //       ResourceType: PersistentVolumeClaim`;
          
    const url = `/deployments/migration`;
//     const data = {
//       yaml:`apiVersion: openmcp.k8s.io/v1alpha1
// kind: Migration
// metadata:
//   name: migrations9
// spec:
//   MigrationServiceSource:
//   - SourceCluster: cluster1
//     TargetCluster: cluster2
//     NameSpace: testmig
//     ServiceName: testim
//     MigrationSource:
//     - ResourceName: testim-dp
//       ResourceType: Deployment
//     - ResourceName: testim-sv
//       ResourceType: Service
//     - ResourceName: testim-pv
//       ResourceType: PersistentVolume
//     - ResourceName: testim-pvc
//       ResourceType: PersistentVolumeClaim`
//     };
    const data = {
      yaml:`apiVersion: openmcp.k8s.io/v1alpha1
kind: Migration
metadata:
  name: migrations1
spec:
  MigrationServiceSource:
  - SourceCluster: cluster1
    TargetCluster: cluster2
    NameSpace: default
    ServiceName: iotservice
    MigrationSource:
    - ResourceName: iot-gateway
      ResourceType: Deployment
    - ResourceName: iot-gateway-svc
      ResourceType: Service`
    };
    axios.post(url, data)
    .then((res) => {
        alert(res.data.message);
        this.setState({ open: false });
        this.props.onUpdateData();
        this.props.menuClose();
    })
    .catch((err) => {
        alert(err);
    });

    // implement migration workflow
    // ......
    this.props.onUpdateData();
    
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

    const onSelectionChange = (selection) => {
      // console.log(selection);
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({ selectedRow: this.state.clusters[selection[0]] ? this.state.clusters[selection[0]] : {} });
    };

    return (
      <div>
        {/* <Button
          variant="outlined"
          color="primary"
          onClick={this.handleClickOpen}
          style={{
            // position: "absolute",
            // right: "115px",
            // top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Migration
        </Button> */}
        <div
          // variant="outlined"
          // color="primary"
          onClick={this.handleClickOpen}
          style={{
            // position: "absolute",
            // right: "115px",
            // top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Migration
        </div>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Deployment Migration
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body">
              <section className="md-content">
                {/* deployment informations */}
                <p>Target Deployment</p>
                <div id="md-content-info">
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>Name : </strong></span>
                      <span>{this.state.dpName}</span>
                    </div>
                    <div class="md-item">
                      <span><strong>Image : </strong></span>
                      <span>{this.state.dpImage}</span>
                    </div>
                  </div>
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>Current Cluster : </strong></span>
                      <span>{this.state.dpCluster}</span>
                    </div>
                  </div>
                </div>
              </section>
              <section className="md-content">
                <p>Select Cluster</p>
                {/* cluster selector */}
                <Paper>
                <Grid rows={this.state.clusters} columns={this.state.columns}>
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
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                  <TableSelection
                    selectByRowClick
                    highlightRow
                    // showSelectionColumn={false}
                  />
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
            <Button onClick={this.handleSave} color="primary">
              excution
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

export default PjDeploymentMigration;
