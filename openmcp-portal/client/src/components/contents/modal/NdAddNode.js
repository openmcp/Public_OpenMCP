import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import { Container } from "@material-ui/core";
import Box from "@material-ui/core/Box";
import PropTypes from "prop-types";
import AddEKSNode from "./addnode/AddEKSNode"
import AddAKSNode from "./addnode/AddAKSNode"
import AddGKENode from "./addnode/AddGKENode"
import AddKVMNode from "./addnode/AddKVMNode"

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

function TabPanel(props) {
  const { children, value, index, ...other } = props;
  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Container>
          <Box>
            {children}
          </Box>
        </Container>
      )}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

const childAddNodeRef = React.createRef();
class NdAddNode extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: 0
    };
  }

  componentWillMount() {
  }

  handleClickOpen = () => {
    this.setState({ 
      open: true,
    });
  };

  handleClose = () => {
    this.setState({
      open: false,
      value: 0
    });
    this.props.menuClose();
  };

  handleChange = (event, newValue) => {
    this.setState({ 
      value: newValue
     });
  };
  
  handleSaveClick = () => {
    if (childAddNodeRef.current) {
      childAddNodeRef.current.handleSaveClick();
    }
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
          // variant="outlined"
          // color="primary"
          onClick={this.handleClickOpen}
          style={{
            // position: "absolute",
            // right: "26px",
            // top: "26px",
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Add Node
        </div>
        <Dialog
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Add Node
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body add-node">
              <AppBar position="static" className="app-bar">
                <Tabs
                  value={this.state.value}
                  onChange={this.handleChange}
                  aria-label="simple tabs example"
                  style={{ backgroundColor: "#3c8dbc", minHeight:"42px"}}
                  TabIndicatorProps ={{ style:{backgroundColor:"#00d0ff"}}}
                >
                  <Tab label="EKS" {...a11yProps(0)} style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}/>
                  <Tab label="GKE" {...a11yProps(1)} style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}/>
                  <Tab label="AKS" {...a11yProps(2)} style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}/>
                  <Tab label="KVM" {...a11yProps(3)} style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}/>
                </Tabs>
              </AppBar>
              <TabPanel className="tab-panel" value={this.state.value} index={0}>
                {/* <AddEKSNode/> */}
                {/* <AddEKSNode handleSaveData={this.handleSaveDialog}/> */}
                <AddEKSNode ref={childAddNodeRef} handleClose={this.handleClose}/>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={1}>
                <AddGKENode ref={childAddNodeRef} handleClose={this.handleClose}/>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={2}>
                <AddAKSNode ref={childAddNodeRef} handleClose={this.handleClose}/>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={3}>
                <AddKVMNode ref={childAddNodeRef} handleClose={this.handleClose}/>
              </TabPanel>
            </div>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleSaveClick} color="primary">
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

export default NdAddNode;