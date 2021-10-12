import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
// import SelectBox from "../../modules/SelectBox";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
  TextField,
} from "@material-ui/core";

// import Paper from "@material-ui/core/Paper";
import axios from "axios";

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

class ThEditThreshold extends Component {
  constructor(props) {
    super(props);
    this.state = {
      cluster:"",
      node:"",
      cpuWarn: 0,
      cpuDanger: 0,
      ramWarn: 0,
      ramDanger: 0,
      storageWarn: 0,
      storageDanger: 0,
      open: false,
    };
  }
  componentWillMount() {}

  onChange = (e) =>{
    this.setState({
      [e.target.name]: e.target.value,
    });
  }
  handleClickOpen = () => {
    if (Object.keys(this.props.rowData).length === 0) {
      alert("Please select a Host Threshold");
      this.setState({ open: false });
      return;
    }

    this.setState({
      cluster :this.props.rowData.cluster_name,
      node:this.props.rowData.node_name,
      cpuWarn: this.props.rowData.cpu_warn,
      cpuDanger: this.props.rowData.cpu_danger,
      ramWarn: this.props.rowData.ram_warn,
      ramDanger: this.props.rowData.ram_danger,
      storageWarn: this.props.rowData.storage_warn,
      storageDanger: this.props.rowData.storage_danger,
      open: true,
    });
  };

  handleClose = () => {
    this.setState({
      cluster :"",
      node:"",
      cpuWarn: 0,
      cpuDanger: 0,
      ramWarn: 0,
      ramDanger: 0,
      storageWarn: 0,
      storageDanger: 0,
      rows: [],
      
      open: false,
    });
   
    this.props.menuClose();
  };

  handleUpdate = (e) => {
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
    } else if (this.state.storageWarn === 0) {
      alert("Please set 'storage warning threshold (%)'");
      return;
    } else if (this.state.storageDanger === 0) {
      alert("Please set 'storage danger threshold (%)'");
      return;
    }

    // Update user role
    const url = `/settings/threshold`;
    const data = {
      clusterName : this.state.cluster,
      nodeName : this.state.node,
      cpuWarn: this.state.cpuWarn,
      cpuDanger: this.state.cpuDanger,
      ramWarn: this.state.ramWarn,
      ramDanger: this.state.ramDanger,
      storageWarn: this.state.storageWarn,
      storageDanger: this.state.storageDanger,
    };
    axios
      .put(url, data)
      .then((res) => {
        alert(res.data.message);
        this.setState({ open: false });
        this.props.menuClose();
        this.props.onUpdateData();
      })
      .catch((err) => {
        alert(err);
      });

    // loging deployment migration
    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-PJ-MD01");

    //close modal popup
    this.setState({ open: false });
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
        <div
          onClick={this.handleClickOpen}
          style={{
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Edit Threshold
        </div>
        <Dialog
          // onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Edit Threshold
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body small-grid">
              <div>
                <Typography>
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
                          <span style={{ bottom: "8px" }}>%</span>
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
                          <span style={{ bottom: "8px" }}>%</span>
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
                          <span style={{ bottom: "8px" }}>%</span>
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
                          <span style={{ bottom: "8px" }}>%</span>
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
                          <span style={{ bottom: "8px" }}>%</span>
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
                            value={this.state.storageDanger}
                            fullWidth={true}
                            name="storageDanger"
                            onChange={this.onChange}
                          />
                          <span style={{ bottom: "8px" }}>%</span>
                        </div>
                      </div>
                    </section>
                  </div>
                </Typography>
              </div>
            </div>
          </DialogContent>
          <DialogActions>
            <div>
              <Button onClick={this.handleUpdate} color="primary">
                update
              </Button>
            </div>
            <Button onClick={this.handleClose} color="primary">
              cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default ThEditThreshold;
