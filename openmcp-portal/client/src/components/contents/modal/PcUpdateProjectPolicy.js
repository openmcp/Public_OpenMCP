import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  TextField,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
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

class PcUpdateProjectPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {

      project:"",
      cluster:"",
      podsCPURate:"",
      podsMemRate:"",
      clusterCPURate:"",
      clusterMemRate:"",
      open: this.props.onOpen,

      rows: "",
      rowData : "",
    };
  }

  componentWillMount() {
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.onOpen !== prevProps.onOpen) {
      this.setState({
        ...this.state,
        project: this.props.rowData.project,
        cluster: this.props.rowData.cluster,
        clusterCPURate: this.props.rowData.cls_cpu_trh_r,
        clusterMemRate: this.props.rowData.cls_mem_trh_r,
        podsCPURate: this.props.rowData.pod_cpu_trh_r,
        podsMemRate: this.props.rowData.pod_mem_trh_r,
        open: this.props.onOpen,
      });
    }
  }

  initState = () => {
    this.setState({
      podsCPURate:"",
      podsMemRate:"",
      clusterCPURate:"",
      clusterMemRate:"",
    });
  }

  callApi = async () => {
    const response = await fetch("/projects");
    const body = await response.json();
    return body;
  };

  handleClose = () => {
    this.initState();
    this.setState({
      open: false,
      rowData : "",
    });
    this.props.onCloseUpdatePolicy(false);
  };

  onChange = (e) =>{
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  handleSave = (e) => {
    if (this.state.clusterCPURate==="" || this.state.clusterMemRate==="" ||this.state.podsCPURate===""||this.state.podsMemRate===""){
      alert("Please insert threshold data");
      return;
    }


    const url = `/settings/policy/project-policy`;
    const data = {
      project:this.state.project,
      cluster:this.state.cluster,
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

    
    // loging Edit Project Policy
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PO-MD02");

    //close modal popup
    this.props.onCloseUpdatePolicy(false);
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
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Edit Project Policy
          </DialogTitle>
          <DialogContent dividers>

          
            <div className="md-contents-body add-node">
            <section className="md-content">
                {/* deployment informations */}
                <p>Target Project</p>
                <div id="md-content-info">
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>Cluster : </strong></span>
                      <span>{this.state.cluster}</span>
                    </div>
                  </div>
                  <div class="md-partition">
                    <div class="md-item">
                      <span><strong>Project : </strong></span>
                      <span>{this.state.project}</span>
                    </div>
                  </div>
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
                      value = {this.state.clusterCPURate}
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
                      value = {this.state.clusterMemRate}
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
                      value = {this.state.podsCPURate}
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
                      value = {this.state.podsMemRate}
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
              update
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

export default PcUpdateProjectPolicy;