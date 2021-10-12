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

class EditAKSAuth extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: this.props.open,
      seq: "",
      cluster : "",
      clientId:"",
      clientSec:"",
      tenantId:"",
      subId:"",
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
          clientId:this.props.data.clientId,
          clientSec:this.props.data.clientSec,
          tenantId:this.props.data.tenantId,
          subId:this.props.data.subId,
        });
      } else {
        this.setState({
          ...this.state,
          open: this.props.open,
          cluster : "",
          clientId:"",
          clientSec:"",
          tenantId:"",
          subId:"",
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
    } else if (this.state.clientId === ""){
      alert("Please enter Client ID");
      return;
    } else if (this.state.clientSec === ""){
      alert("Please enter Client Sec");
      return;
    } else if (this.state.tenantId === ""){
      alert("Please enter Tenant ID");
      return;
    } else if (this.state.subId === ""){
      alert("Please enter Sub ID");
      return;
    } 

    //post 호출
    const url = `/settings/config/pca/aks`;
  
    if(this.props.new){
      const data = {
        cluster : this.state.cluster,
        clientId: this.state.clientId,
        clientSec: this.state.clientSec,
        tenantId: this.state.tenantId,
        subId: this.state.subId,
      };

      // clientID = "1edadbd7-d466-43b1-ad73-15a2ee9080ff"
      // clientSec = "07.Tx2r7GobBf.Suq7quNRhO_642z-p~6a"
      // tenantID = "bc231a1b-ab45-4865-bdba-7724c2893f1c"
      // subID := "dc80d3cf-4e1a-4b9a-8785-65c4b739e8d2"

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
        clientId: this.state.clientId,
        clientSec: this.state.clientSec,
        tenantId: this.state.tenantId,
        subId: this.state.subId,
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
                  <p>Client ID</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="clientId"
                    variant="outlined"
                    value={this.state.clientId}
                    fullWidth={true}
                    name="clientId"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Client Sec</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="clientSec"
                    variant="outlined"
                    value={this.state.clientSec}
                    fullWidth={true}
                    name="clientSec"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Tenant ID</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="tenant Id"
                    variant="outlined"
                    value={this.state.tenantId}
                    fullWidth={true}
                    name="tenantId"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Sub ID</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="sub Id"
                    variant="outlined"
                    value={this.state.subId}
                    fullWidth={true}
                    name="subId"
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

export default EditAKSAuth;
