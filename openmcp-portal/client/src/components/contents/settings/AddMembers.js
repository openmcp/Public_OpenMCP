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

class AddMembers extends Component {
  constructor(props) {
    super(props);
    this.state = {
      userid: "",
      password: "",
      retype_password: "",
      open: false,
    };
    this.onChanged = this.onChanged.bind(this);
  }

  componentWillMount() {
  }

  onChanged(e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  handleClickOpen = () => {
    this.setState({ open: true });
  };

  handleClose = () => {
    this.setState({ open: false });
    this.props.menuClose();
  };

  handleRegister = (e) => {
    e.preventDefault();
    const { password, retype_password } = this.state;
    if (password !== retype_password) {
      alert("Password confirmation does not match");
    } else {
      const url = `/create_account`;
      const data = {
        userid:this.state.userid,
        password:this.state.password,
        role:"{omcp_monitor}",
      };
      axios.post(url, data)
      .then((res) => {
          alert(res.data.message);
          this.setState({ open: false });
          this.props.menuClose();
      })
      .catch((err) => {
          alert(err);
      });

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
            textTransform: "initial",
          }}
        >
          Create an Account
        </div>
        <Dialog
          // onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          // maxWidth={false}
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Create Account
          </DialogTitle>
          <DialogContent dividers>
            <div className="signup">
              <form onSubmit={this.submitForm}>
                <input
                  type="text"
                  placeholder="Userid"
                  name="userid"
                  value={this.state.userid}
                  onChange={this.onChanged}
                />
                <input
                  type="password"
                  placeholder="Password"
                  name="password"
                  value={this.state.password}
                  onChange={this.onChanged}
                />
                <input
                  type="password"
                  placeholder="Retype password"
                  name="retype_password"
                  value={this.state.retype_password}
                  onChange={this.onChanged}
                />
              </form>
            </div>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleRegister} color="primary">
              register
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

export default withStyles(styles)(AddMembers);
