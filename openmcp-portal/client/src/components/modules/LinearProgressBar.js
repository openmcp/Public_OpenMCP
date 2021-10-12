import React, { Component } from "react";
import LinearProgress from "@material-ui/core/LinearProgress";
import { withStyles } from "@material-ui/core/styles";
// import { debug } from "request";

const styles = (props) => ({
  normalColor: {
    backgroundColor: "#3f51b5",
  },
  normalBaseColor: {
    backgroundColor: "#b6bce2",
  },
  warnColor: {
    backgroundColor: "#FFA228",
  },
  warnBaseColor: {
    backgroundColor: "#FDEFCF",
  },
  dangerColor: {
    backgroundColor: "#FF3628",
  },
  dangerBaseColor: {
    backgroundColor: "#FFD5C4",
  },
});

class LinearProgressBar extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: this.props.value,
      total: this.props.total,
      color: this.props.classes.normalColor,
      baseColor: this.props.classes.normalBaseColor,
    };
  }

  componentWillMount() {
    if ((this.state.value / this.state.total) * 100 > 90) {
      this.setState({
        color: this.props.classes.dangerColor,
        baseColor: this.props.classes.dangerBaseColor,
      });
    } else if ((this.state.value / this.state.total) * 100 > 80) {
      this.setState({
        color: this.props.classes.warnColor,
        baseColor: this.props.classes.warnBaseColor,
      });
    } else {
      this.setState({
        color: this.props.classes.normalColor,
        baseColor: this.props.classes.normalBaseColor,
      });
    }
    // if(this.state.percentage === null){
    //   this.setState({progress: this.state.value / this.state.total * 100})
    // }
  }
  render() {
    // const { classes } = this.props;
    return (
      <div className="linear-progress">
        <LinearProgress
          {...this.props}
          variant="determinate"
          value={(this.state.value / this.state.total) * 100}
          classes={{
            colorPrimary: this.state.baseColor,
            barColorPrimary: this.state.color,
          }}
        />
        {/* .VolumeBar > * { background-color:green; }
.VolumeBar{background-color:gray ;} */}
      </div>
    );
  }
}

export default withStyles(styles)(LinearProgressBar);
