import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import Typography from "@material-ui/core/Typography";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import * as utilLog from './../../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';


// import axios from 'axios';


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



function valuetext(value) {
  return `${value}°C`;
}

class MtSetMetering extends Component {
  constructor(props){
    super(props)
    this.cpu_max = 5000;
    this.memory_max = 10000;
    this.state={
      title : this.props.name,
      open : false,
      cpu_req : "No Request",
      cpu_limit : "Limitless",
      mem_req : "No Request",
      mem_limit : "Limitless",
      // cpu_value : [this.props.resources.cpu_min.replace('m',''),this.props.resources.cpu_max.replace('m','')],
      // mem_value : [this.props.resources.memory_min.replace('Mi',''),this.props.resources.memory_max.replace('Mi','')],
      cpu_value : [0, 0],
       mem_value : [0, 0]
    }
    this.onChange = this.onChange.bind(this);
  }
  componentWillMount() {
  }

  onChange(e) {
    // console.log("onChangedd");
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  handleClickOpen = () => {
    this.setState({
       open: true,
       cpu_value : [0, 0],
       mem_value : [0, 0]
      });
    
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  handleSave = (e) => {
    //Save modification data (Resource Changed)

    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, 'log-PD-MD01');
    this.setState({open:false});
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

    const handleChangeCpu = (e, newValue) => {
      this.setState({
        cpu_value: newValue,
      });
    };
    
    const handleChangeMem = (e, newValue) => {
      this.setState({
        mem_value: newValue,
      });
    };

    return (
      <div>
        {/* <div
          className="btn-join"
          onClick={this.handleClickOpen}
          style={{
            position: "absolute",
            right: "12px",
            top: "0px",
            zIndex: "10",
          }}
        >
          
        </div> */}
        <Button variant="outlined" color="primary" onClick={this.handleClickOpen} style={{position:"absolute", right:"31px", top:"28px", zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
        Set Metering
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth={false}
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
          Metering Setting
          </DialogTitle>
          <DialogContent dividers>
            <div className="mt-set-metering">
              <div className="res">
                <Typography id="range-slider" gutterBottom>
                  Resource Price
                </Typography>
                <div className="txt">
                  <span>CPU : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                          style ={{textAlign: "right",
                            padding: "0"}}
                    />
                  </div>
                  <span className="per">vCPU/hour</span>

                  <span className="price">Price : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="cpu_req" 
                      // value={
                      //     this.state.cpu_value[0]  === 0 ? "No Request" : this.state.cpu_value[0]
                      // }
                      // onChange={this.onChange}
                       />
                    <span>$</span>
                  </div>
                </div>
                <div className="txt">
                  <span>Memory : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                          style ={{textAlign: "right",
                            padding: "0"}}
                    />
                  </div>
                  <span className="per">GiB/hour</span>

                  <span className="price">Price : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="cpu_req" 
                      // value={
                      //   this.state.cpu_value[0]  === 0 ? "No Request" : this.state.cpu_value[0]
                      // }
                      // onChange={this.onChange} 
                      />
                    <span>$</span>
                  </div>
                </div>
                <div className="txt">
                  <span>Disk : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                          style ={{textAlign: "right",
                            padding: "0"}}
                    />
                  </div>
                  <span className="per">GB/hour</span>

                  <span className="price">Price : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="cpu_req" 
                      // value={
                      //   this.state.cpu_value[0]  === 0 ? "No Request" : this.state.cpu_value[0]
                      // }
                      // onChange={this.onChange} 
                      />
                    <span>$</span>
                  </div>
                </div>
              </div>
              <div className="res transfer-price">
                <Typography id="range-slider" gutterBottom>
                    Data Transfer Price
                </Typography>
                <div className="txt">
                  <span>Data Transfer Out : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                          style ={{textAlign: "right",
                            padding: "0"}}
                    />
                  </div>
                  <span className="per">GB/Month</span>

                  <span className="price">Price : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                    />
                    <span>$</span>
                  </div>

                </div>
                <div className="txt">
                  <span>Data Transfer In : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                          style ={{textAlign: "right",
                            padding: "0"}}
                    />
                  </div>
                  <span className="per">GB/Month</span>

                  <span className="price">Price : </span>
                  <div className="input-txt">
                    <input type="number" placeholder="0" name="mem_req" 
                          // value={
                          //   this.state.mem_value[0]  === 0 ? "No Request" : this.state.mem_value[0]
                          // }
                          // onChange={this.onChange} 
                    />
                    <span>$</span>
                  </div>

                </div>
              </div>
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

export default MtSetMetering;