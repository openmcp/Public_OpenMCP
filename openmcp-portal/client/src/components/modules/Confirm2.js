import React, { Component } from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

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

class Confirm2 extends Component {
  constructor(props){
    super(props);
    this.state = {
      open : this.props.confirmOpen,
      title : this.props.confirmInfo.title,
      context : this.props.confirmInfo.context,
      confrimTarget : "",
      confirmTargetKeyname : "",
      button : {
        open : this.props.confirmInfo.button.open,
        yes : this.props.confirmInfo.button.yes,
        no : this.props.confirmInfo.button.no
      }
    }
  }

  componentDidMount(){
  } 

  // componentDidUpdate(prevProps, prevState) {
  //   console.log("pods update");
  //   if (this.props.rowData !== prevProps.rowData) {
  //     this.setState({
  //       ...this.state,
  //       rows: this.props.rowData,
  //     });
  //   }
  // }

  componentWillUpdate(prevProps, prevState){
    // console.log("aaaaaaaaaaaaaaaaaaaaaaaa", prevProps.confirmOpen, this.props.confirmOpen)
    if(prevProps.confirmOpen !== this.props.confirmOpen || prevProps.confirmInfo.title !== this.props.confirmInfo.title){
      // console.log("componentWillUpdate", prevProps.confirmOpen !== this.props.confirmOpen)
      // console.log("componentWillUpdate", prevProps.confirmOpen,this.props.confirmOpen)
      // console.log("componentWillUpdate", prevProps.confirmInfo.title,this.props.confirmInfo.title)
      this.setState({
        open:prevProps.confirmOpen,
        confrimTarget:prevProps.confrimTarget,
        title : prevProps.confirmInfo.title,
        context : prevProps.confirmInfo.context,
        button : {
          open : prevProps.confirmInfo.button.open,
          yes : prevProps.confirmInfo.button.yes,
          no : prevProps.confirmInfo.button.no
        },
        confirmTargetKeyname : prevProps.confirmTargetKeyname
      })
    }
  }

  handleYes = () => {
    // this.setState({open:false});
    this.props.confirmed(true)
  };

  handleNo = () => {
    // this.setState({open:false});
    this.props.confirmed(false)
  };
  render() {


    return (
      <div>
      <Dialog
        open={this.state.open}
        // onClose={this.handleNo}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">
          {this.state.title}
        </DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            <div>{this.state.context}</div>
            {this.state.confrimTarget !== ""
             ? <div style={{fontSize:"16px"}}>
                <strong> {this.state.confirmTargetKeyname} : {this.state.confrimTarget}</strong>
               </div>
             : ""}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.handleYes} color="primary">
            {this.state.button.yes}
          </Button>
          <Button onClick={this.handleNo} color="primary" autoFocus>
            {this.state.button.no}
          </Button>
        </DialogActions>
      </Dialog>
    </div>
    );
  }
}

export default Confirm2;