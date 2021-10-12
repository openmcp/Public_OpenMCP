import React, { Component } from "react";
import WorldMapClusterStatus from "./WorldMapClusterStatus";

class CustomDynamicView extends Component {
  constructor(props) {
    super(props);
    this.state = {
      componentList: this.props.componentList,
    };
  }

  render() {
    
    return (
      <div>
        <WorldMapClusterStatus />
      </div>
    );
  }
}

export default CustomDynamicView;
