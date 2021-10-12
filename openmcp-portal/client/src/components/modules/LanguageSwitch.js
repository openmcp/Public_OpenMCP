import React, { Component } from "react";
import { alpha, styled } from "@material-ui/core/styles";
import { pink } from "@material-ui/core/colors";
import Switch from "@material-ui/core/Switch";

class LanguageSwitch extends Component {
  constructor(props) {
    super(props);
    this.state = {
      checked: false,
      language : "En"
    };
  }
  render() {
    // const GreenSwitch = styled(Switch)(({ theme }) => ({
    //   '& .MuiSwitch-switchBase.Mui-checked': {
    //     color: "#0088fe",
    //     // '&:hover': {
    //     //   backgroundColor: "#BECFDDBB",
    //     // },
    //   },
    //   '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': {
    //     backgroundColor: "#FFFFFF",
    //   },
    // }));

    const GreenSwitch = styled(Switch)(({ theme }) => ({
      "& .MuiSwitch-switchBase.Mui-checked": {
        color: pink[600],
        "&:hover": {
          backgroundColor: "#BECFDDBB",
        },
      },
      "& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track": {
        backgroundColor: pink[600],
      },
    }));

    const label = { inputProps: { "aria-label": "Switch demo" } };

    const handleChange = (event) => {
      this.setState({ checked: event.target.checked, language: event.target.checked ? "한글" : "EN"});
      
    };

    return (
      <div style={{fontSize:"12px"}}>
        <Switch
          {...label}
          checked={this.state.checked}
          onChange={handleChange}
          color="default"
        />
        <span>{this.state.language}</span>
      </div>
    );
  }
}

export default LanguageSwitch;
