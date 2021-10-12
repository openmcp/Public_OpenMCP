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
  
  class EditGKEAuth extends Component {
    constructor(props) {
      super(props);
      this.state = {
        open: this.props.open,
        seq: "",
        cluster : "",
        type:"",
        clientEmail:"",
        projectID:"",
        privateKey:"",
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
            type:this.props.data.type,
            clientEmail:this.props.data.clientEmail,
            projectID:this.props.data.projectID,
            privateKey:this.props.data.privateKey,
          });
        } else {
          this.setState({
            ...this.state,
            open: this.props.open,
            cluster : "",
            type:"",
            clientEmail:"",
            projectID:"",
            privateKey:"",
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
      } else if (this.state.type === ""){
        alert("Please enter Credential Type");
        return;
      } else if (this.state.clientEmail === ""){
        alert("Please enter Client Email");
        return;
      } else if (this.state.projectID === ""){
        alert("Please enter Project ID");
        return;
      } else if (this.state.privateKey === ""){
        alert("Please enter Private Key");
        return;
      } 
  
      //post 호출
      const url = `/settings/config/pca/gke`;
    
      if(this.props.new){
        const data = {
          cluster : this.state.cluster,
          type: this.state.type,
          clientEmail: this.state.clientEmail,
          projectID: this.state.projectID,
          privateKey: this.state.privateKey,
        };

        // credType := "service_account"
        // projectID = "just-advice-302807"
        // clientEmail = "gkeadmin@just-advice-302807.iam.gserviceaccount.com"
        // privateKey = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDWdGxXcM+cRb39\nN6fbCBpibF+EIVFkKGjsuuuGJxEoTQIKp2dnl5FlBFKKSa0cSIz4duwgxc5+25KS\neR5cBB6MjSxBC62qK6VeyNUT2KzyIrQfp/zGmxkBVpXFZ13u0JopiwSH5Kvp4vU1\nOJn4wLA3aLs3QMzUC4rXl6IW0yuyMeClooJLFqxjW7ihry2Y0MjMLuSWeHpqCQCK\n0IntRpqhPoKEkWUjonJnQo7Lem5/iqp8rL80vMDPHuDTPLcQt3pI7Ak6z2qk7etm\ng5jkUS1cVU9Xne2jffEMOjTXPrEgozoHWxN0QLwzrA/7vW6zAt3nfOdO9C6wBzh9\n4GgUeTDDAgMBAAECggEAAlWPaFQ+A5bEE/bVyOM0W6Xk/uyDP50rpzKm+vV/O6UQ\nRKAV1rbQ9PyFuXjxKBb8vHzu4lxvfEn/imtEZ/6o0SF9kyesDZIetq1mRFUIwjSb\n0/cMH/fy3w+GNHkvjeM6ClcNBuhM8WVwWH1JOmqT1caPYxvoHta7/XoVCufkLd2q\nqpFcod8LISW3HN7wSgzB5lpDry+Zk8KoXtxn2bAJyRYeky7tkXQbkCwrE10oUkAs\nivgR27wGF0nowoSvs8KwxWME3zW836fVALyF+dGCBlYVtIMvx6T4cu868dI5JANj\nY6U4H3xjB98MQ/zp7uH6w4kj1/cMxvbfAT7jBTiqAQKBgQD9Q/9bEVxPBc+gKEMo\ncXYCJTCT5XsdAgdw/kXHdR463z70sUbLhvHjt/6xwlNCS5j2jkTbN7InLI6xIwY/\nzdfppXsoW4qyEqrgMHjG1af3AlslEA3GLnkLEIx/VM6zoDKlBlI3uz2PMf4wJiFK\nli3X/5tcpYlyc0pCkIJBQ+o2QQKBgQDYxSf8b2/WW87+L3l6/VlbyWMG9aw5RP2v\nitP0cIqoFj/LkD1pJWtJre0Lnlzgz8JJDcRsbrqDZFuIiWnTc8dy8YM1Pv1kz7xZ\nANvpJGEDr5cZjopOoq+w5zfNDrLf/SPB2g6u9/33Ukds3F0++14901b/f7SjHFN2\nH+OPFwMOAwKBgQDPugrird2Rbwm5qexTaqRI5Cnw1ELjKvvhgJzJGNV/ogXn+tM/\nMeKKTSqYr/NMJ+dBKrVtPERh/xjWTwzcHkBegfz+v/6FSexfT0Jwi2NlpMgPIRi7\nGPjsy1kBQxT6nYWMdx/OWEQIhA+hfFTH8V+OjzbliVyvw8H/0LkVQNgEQQKBgBJr\nhn9T9NvxR0CgRiFmX+6FyW1w+OaQ70G4eVRfL9kist8Yba9+p4RGTEtddKUB4o+U\npOlV63F42LJcguqd/wfMcArZRG0JngauJQHFvpyykhNw4l3WQzm0HDDHm/meqCgz\n4GWL2z/l9P3SJ/ZPI+37BHyHnJDzuj/ia9Lf8LmDAoGBAOm92Sp7qFkrwogzIBfp\nU9PtDc2GeiSj7WJctIakuxQ+bSWtOoPq6CPd8OAWmpgZA8SzCfkWMnBQJhB7A6RQ\nZOA50xvE07ybQ397NLkDKAB56zdQ9hDAYpgkzCFWL1AvIouM8OLU48LLIh3KJLxG\nSUwFrPzKIQz4RKj3em+M+iQP\n-----END PRIVATE KEY-----\n"
  
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
          type: this.state.type,
          clientEmail: this.state.clientEmail,
          projectID: this.state.projectID,
          privateKey: this.state.privateKey,
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
                    <p>Credential Type</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      placeholder="credential type"
                      variant="outlined"
                      value={this.state.type}
                      fullWidth={true}
                      name="type"
                      onChange={this.onChange}
                    />
                  </div>
                  <div className="props">
                    <p>Client Email</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      placeholder="client email"
                      variant="outlined"
                      value={this.state.clientEmail}
                      fullWidth={true}
                      name="clientEmail"
                      onChange={this.onChange}
                    />
                  </div>
                  <div className="props">
                    <p>Project ID</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      placeholder="project id"
                      variant="outlined"
                      value={this.state.projectID}
                      fullWidth={true}
                      name="projectID"
                      onChange={this.onChange}
                    />
                  </div>
                  <div className="props">
                    <p>Private key</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={8}
                      multiline
                      placeholder="private key"
                      variant="outlined"
                      value={this.state.privateKey}
                      fullWidth={true}
                      name="privateKey"
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
  
  export default EditGKEAuth;
  