import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import { TextField } from "@material-ui/core";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
import axios from "axios";
import * as utilLog from "../../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
// import Confirm2 from "../../../modules/Confirm2";

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

class EditKVMAuth extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: this.props.open,
      seq: "",
      cluster : "",
      agentURL : "",
      mClusterName : "",
      mClusterPwd : "",
    };
  }

  componentWillMount() {
    
  }
  
  componentDidUpdate(prevProps, prevState) {
    if (this.props.open !== prevProps.open) {
      if(this.props.new === false){
        this.setState({
          ...this.state,
          open: this.props.open,
          seq : this.props.data.seq,
          cluster : this.props.data.cluster,
          agentURL : this.props.data.agentURL,
          mClusterName : this.props.data.mClusterName,
          mClusterPwd : this.props.data.mClusterPwd,
        });
      } else {
        this.setState({
          ...this.state,
          open: this.props.open,
          cluster : "",
          agentURL : "",
          mClusterName : "",
          mClusterPwd : "",
        });
      }
    }
  }

  onChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value,
    });
  };

  handleClose = () => {
    this.props.callBackClosed()
  };

  handleSave = (e) => {
    if (this.state.cluster === ""){
      alert("Please enter Cluster Name");
      return;
    } else if (this.state.agentURL === ""){
      alert("Please enter Agent URL");
      return;
    } else if (this.state.mClusterName === ""){
      alert("Please enter Master Cluster VM Name");
      return;
    } else if (this.state.mClusterPwd === ""){
      alert("Please enter Master Cluster Password");
      return;
    } 

    //post 호출
    const url = `/settings/config/pca/kvm`;
    if(this.props.new){
      const data = {
        seq : this.state.seq,
        cluster : this.state.cluster,
        agentURL : this.state.agentURL,
        mClusterName : this.state.mClusterName,
        mClusterPwd : this.state.mClusterPwd,
      };

      axios.post(url, data)
      .then((res) => {
        this.props.callBackClosed()
      })
      .catch((err) => {
        //close modal popup
        console.log("Error : ",err);
      });
    } else {
      const data = {
        seq : this.props.data.seq,
        cluster : this.state.cluster,
        agentURL : this.state.agentURL,
        mClusterName : this.state.mClusterName,
        mClusterPwd : this.state.mClusterPwd,
      };

      axios.put(url, data)
      .then((res) => {
        this.props.callBackClosed()
      })
      .catch((err) => {
        //close modal popup
        console.log("Error : ",err);
      });
    }

      

    // alert(this.state.cluster+","+ this.state.secretKey+","+this.state.accessKey)

    // loging deployment migration
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PJ-MD01");

    
  };

  callApi = async (uri) => {
    // const response = await fetch("/aws/clusters");
    const response = await fetch(uri);
    const body = await response.json();
    return body;
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
          // onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="lg"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            {this.props.title}
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body" style={{minWidth:"500px"}}>
              <section className="md-content">
                <div className="props">
                  <p>Cluster</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="cluster name"
                    variant="outlined"
                    value={this.state.cluster}
                    fullWidth={true}
                    name="cluster"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Agent URL</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="agentURL"
                    variant="outlined"
                    value={this.state.agentURL}
                    fullWidth={true}
                    name="agentURL"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Master Cluster VM Name</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="master cluster VM Name"
                    variant="outlined"
                    value={this.state.mClusterName}
                    fullWidth={true}
                    name="mClusterName"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Master Cluster Password</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="master cluster VM Password"
                    variant="outlined"
                    value={this.state.mClusterPwd}
                    fullWidth={true}
                    name="mClusterPwd"
                    onChange={this.onChange}
                  />
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

export default EditKVMAuth;
