import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
// import * as utilLog from "../../util/UtLogs.js";
// import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
} from "@material-ui/core";
// import {
//   PagingState,
//   SortingState,
//   SelectionState,
//   IntegratedFiltering,
//   IntegratedPaging,
//   IntegratedSorting,
//   IntegratedSelection,
//   // SearchState,
//   // EditingState,
// } from "@devexpress/dx-react-grid";
// import {
//   Grid,
//   Table,
//   TableColumnResizing,
//   TableHeaderRow,
//   PagingPanel,
//   TableSelection,
//   // SearchPanel,
//   // Toolbar,
//   // TableEditRow,
//   // TableEditColumn,
// } from "@devexpress/dx-react-grid-material-ui";
// import Paper from "@material-ui/core/Paper";
import ChangeEKSReource from "./chageResource/ChangeEKSReource.js";
import ChangeKVMReource from './chageResource/ChangeKVMReource';

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

const childChangeNodeRes = React.createRef();
class NdResourceConfig extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
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
    });
  };

  handleSaveClick = () => {
    if (childChangeNodeRes.current) {
      childChangeNodeRes.current.handleSaveClick();
    }
  };

  handleChange = (event, newValue) => {
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
        <Button variant="outlined" color="primary" onClick={this.handleClickOpen} style={{position:"absolute", right:"0px", top:"0px", zIndex:"10", width:"148px", height:"31px", textTransform: "capitalize"}}>
          Resource Config
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
           Node Resource Configration
          </DialogTitle>
          <DialogContent dividers>
            {this.props.propsRow.provider === "eks" 
            ? <ChangeEKSReource ref={childChangeNodeRes} handleClose={this.handleClose} nodeData={this.props.nodeData} propsRow={this.props.propsRow} />
            : (this.props.propsRow.provider === "kvm" 
                ? <ChangeKVMReource ref={childChangeNodeRes} handleClose={this.handleClose} nodeData={this.props.nodeData} propsRow={this.props.propsRow} rows={this.props.rows}/> 
                : "")
            }
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleSaveClick} color="primary">
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

export default NdResourceConfig;
