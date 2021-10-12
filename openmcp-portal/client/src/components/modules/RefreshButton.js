import React, { Component } from 'react';
import RefreshRoundedIcon from '@material-ui/icons/RefreshRounded';

class RefreshButton extends Component {
  render() {
    return (
      <div>
        <RefreshRoundedIcon onClick={this.onRefresh}></RefreshRoundedIcon>
      </div>
    );
  }
}

export default RefreshButton;