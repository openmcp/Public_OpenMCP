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

class EditEKSAuth extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: this.props.open,
      seq: "",
      cluster : "",
      accessKey : "",
      secretKey : "",
      region : "",
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
          accessKey : this.props.data.accessKey,
          secretKey : this.props.data.secretKey,
          region : this.props.data.region
        });
      } else {
        this.setState({
          ...this.state,
          open: this.props.open,
          cluster : "",
          accessKey : "",
          secretKey : "",
          region : "",
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
    } else if (this.state.secretKey === ""){
      alert("Please enter Secret Key");
      return;
    } else if (this.state.accessKey === ""){
      alert("Please enter Access Key");
      return;
    } else if (this.state.region === ""){
      alert("Please enter Region");
      return;
    }

    //post 호출
    const url = `/settings/config/pca/eks`;
  
    if(this.props.new){
      const data = {
        cluster : this.state.cluster,
        accessKey : this.state.accessKey,
        secretKey : this.state.secretKey,
        region : this.state.region,
      };

      // cluster := cluster1
      // secretkey := "QnD+TaxAwJme1krSz7tGRgrI5ORiv0aCiZ95t1XK"
      // accessKey := "AKIAJGFO6OXHRN2H6DSA"

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
        accessKey : this.state.accessKey,
        secretKey : this.state.secretKey,
        region : this.state.region,
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
                  <p>Secret key</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="secret key"
                    variant="outlined"
                    value={this.state.secretKey}
                    fullWidth={true}
                    name="secretKey"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Access key</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="access key"
                    variant="outlined"
                    value={this.state.accessKey}
                    fullWidth={true}
                    name="accessKey"
                    onChange={this.onChange}
                  />
                </div>
                <div className="props">
                  <p>Region</p>
                  <TextField
                    id="outlined-multiline-static"
                    rows={1}
                    placeholder="region"
                    variant="outlined"
                    value={this.state.region}
                    fullWidth={true}
                    name="region"
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

export default EditEKSAuth;
