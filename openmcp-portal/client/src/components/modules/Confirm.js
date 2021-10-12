import React, { Component } from "react";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogTitle from "@material-ui/core/DialogTitle";

/* use example

// state
confirmInfo : {
  title :"Cluster Join Confrim",
  context :"Are you sure you want to Join the Cluster?",
  button : {
    open : "JOIN",
    yes : "JOIN",
    no : "CANCEL",
  }
},
confrimTarget : "false" // "false" : must have target, "" : 

//component
<Confirm confirmInfo={this.state.confirmInfo} confrimTarget ={this.state.confrimTarget} confirmed={this.confirmed}/>

//callback
confirmed = (result) => {
  if(result) {
    //Unjoin proceed
    console.log("confirmed")
  } else {
    console.log("cancel")
  }
}
*/

class Confirm extends Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false,
      title: this.props.confirmInfo.title,
      context: this.props.confirmInfo.context,
      confrimTarget: "",
      button: {
        open: this.props.confirmInfo.button.open,
        yes: this.props.confirmInfo.button.yes,
        no: this.props.confirmInfo.button.no,
      },
    };
  }

  componentDidMount() {}

  render() {
    const handleClickOpen = () => {
      if (this.props.confrimTarget === "false") {
        alert("Please select target");
        this.setState({ open: false });
        return;
      }

      this.setState({
        open: true,
        confrimTarget: this.props.confrimTarget,
      });
    };
    const handleYes = () => {
      this.props.confirmed(true);
      this.setState({ open: false });
      this.props.menuClose();
    };

    const handleNo = () => {
      this.props.confirmed(false);
      this.setState({ open: false });
      this.props.menuClose();
    };

    return (
      <div>
        {/* <Button
          variant="outlined"
          color="primary"
          onClick={handleClickOpen}
          style={{
            position: "absolute",
            right: "26px",
            marginTop: "15px",
            zIndex: "10",
            width: "148px",
            height: "31px",
            textTransform: "capitalize",
          }}
        >
          {this.state.button.open}
        </Button> */}
        <div
          variant="outlined"
          color="primary"
          onClick={handleClickOpen}
          style={{
            // position: "absolute",
            // right: "26px",
            // marginTop: "15px",
            zIndex: "10",
            width: "148px",
            // height: "31px",
            textTransform: "capitalize",
          }}
        >
          {this.state.button.open}
        </div>
        <Dialog
          open={this.state.open}
          // onClose={handleNo}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">{this.state.title}</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              <div>{this.state.context}</div>
              {this.state.confrimTarget !== "" ? (
                <div style={{ fontSize: "16px" }}>
                  <strong> Target : {this.state.confrimTarget}</strong>
                </div>
              ) : (
                ""
              )}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleYes} color="primary">
              {this.state.button.yes}
            </Button>
            <Button onClick={handleNo} color="primary" autoFocus>
              {this.state.button.no}
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default Confirm;
