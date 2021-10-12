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
import Slider from "@material-ui/core/Slider";
import * as utilLog from "../../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
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

class PcSetNumericPolicy extends Component {
  constructor(props) {
    super(props);
    // this.g_rate_max = 10;
    // this.period_max = 10;
    this.state = {
      open: false,
      policyData: [],
    };
    this.onChange = this.onChange.bind(this);
  }

  componentWillMount() {}

  onChange(e, newValue) {
    if(e.target.id !== ""){
      let tempData = this.state.policyData;
      tempData[e.target.id].value = newValue;
    }
  }

  handleClickOpen = () => {
    if (Object.keys(this.props.policy).length === 0) {
      alert("Please Select Policy");
      this.setState({ open: false });
      return;
    }

    let policyData = [];
    this.props.policy.value
      .slice(0, -1)
      .split("|")
      .forEach((item) => {
        let itemSplit = item.split(" : ");
        policyData.push({
          key: itemSplit[0],
          value: parseFloat(itemSplit[1]),
        });
      });
    console.log(policyData);

    this.setState({
      open: true,
      policyData: policyData,
    });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  handleSave = (e) => {
    let valueData = [];
    this.state.policyData.forEach((item, index)=>{
      let object = {
        op: "replace",
        path: "/spec/template/spec/policies/"+index.toString()+"/value/0",
        value: item.value.toString(),
      }

      valueData.push(object);
    })

    // Save modification data (Policy Changed)
    const url = `/settings/policy/openmcp-policy`;
    const data = {
      policyName: this.props.policyName,
      values: valueData,
    };
    axios
      .post(url, data)
      .then((res) => {
        if (res.data.length > 0) {
          if (res.data[0].code == 200) {
            this.props.onUpdateData();

            let userId = null;
            AsyncStorage.getItem("userName", (err, result) => {
              userId = result;
            });
            utilLog.fn_insertPLogs(userId, "log-PO-MD01");
          }
          alert(res.data[0].text);
        }
      })
      .catch((err) => {
        alert(err);
      });
    this.setState({ open: false });
  };

  g_float_marks = [
    { value: 0, label: "0" },
    { value: 0.1, label: "0.1" },
    { value: 0.2, label: "0.2" },
    { value: 0.3, label: "0.3" },
    { value: 0.4, label: "0.4" },
    { value: 0.5, label: "0.5" },
    { value: 0.6, label: "0.6" },
    { value: 0.7, label: "0.7" },
    { value: 0.8, label: "0.8" },
    { value: 0.9, label: "0.9" },
    { value: 1, label: "1" },
  ];

  
  g_int_marks = [
    {value: 0,label: "0"},
    {value: 1,label: "1"},
    {value: 2,label: "2"},
    {value: 3,label: "3"},
    {value: 4,label: "4"},
    {value: 5,label: "5"},
    {value: 6,label: "6"},
    {value: 7,label: "7"},
    {value: 8,label: "8"},
    {value: 9,label: "9"},
    {value: 10,label: "10"},
  ];

  render() {
    const DialogTitle = withStyles(styles)((props) => {
      const { children, classes, onClose, ...other } = props;
      return (
        <div>
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
        </div>
      );
    });

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
        <Button
          variant="outlined"
          color="primary"
          onClick={this.handleClickOpen}
          style={{
            position: "absolute",
            right: "26px",
            top: "26px",
            zIndex: "10",
            width: "148px",
            height: "31px",
            textTransform: "capitalize",
          }}
        >
          edit policy
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth={false}
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            {this.props.policyName}
          </DialogTitle>
          <DialogContent dividers>
            <div className="pd-resource-config">
              {this.state.policyData.map((item, index) => {
                return (
                  <div className="res">
                    <Typography id="range-slider" gutterBottom>
                      {item.key}
                    </Typography>
                    <Slider
                      id={index}
                      className="sl"
                      name="policyData"
                      defaultValue={item.value}
                      onChange={this.onChange}
                      valueLabelDisplay="auto"
                      aria-labelledby="range-slider"
                      // getAriaValueText={valuetext}
                      step={null}
                      min={0}
                      max={this.props.isFloat ? 1 : 10}
                      marks={this.props.isFloat ? this.g_float_marks : this.g_int_marks}
                    />
                    <div className="txt"></div>
                  </div>
                );
              })}
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

export default PcSetNumericPolicy;
