import React, { Component } from "react";
// import { Settings } from "@material-ui/icons";
import AccountCircleIcon from "@material-ui/icons/AccountCircle";

import { Link } from "react-router-dom";
// //import ClickAwayListener from '@material-ui/core/ClickAwayListener';
// import Grow from '@material-ui/core/Grow';
// import Paper from '@material-ui/core/Paper';
// import Popper from '@material-ui/core/Popper';
// import MenuItem from '@material-ui/core/MenuItem';
// import MenuList from '@material-ui/core/MenuList';
import * as utilLog from "./../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
import LanguageSwitch from "../modules/LanguageSwitch.js";
import LanguageListMenu from "../modules/LanguageListMenu.js";

class Head extends Component {
  constructor(props) {
    super(props);
    this.state = {
      // anchorEl:null,
      open: false,
      selectedMenu: "dashboard",
    };
    this.anchorRef = React.createRef();
    this.prevOpen = React.createRef(this.state.open);
  }

  componentWillUpdate(prevProps, prevState) {
    if (this.props.path !== prevProps.path) {
      // if (prevProps.path.indexOf('/dashboard') >= 0 ){
      //   this.setState({selectedMenu:'dashboard'})
      // } else if (prevProps.path.indexOf('/clusters') >= 0 ) {
      //   this.setState({selectedMenu:'clusters'})
      // } else if (prevProps.path.indexOf('/nodes') >= 0 ) {
      //   this.setState({selectedMenu:'nodes'})
      // } else if (prevProps.path.indexOf('/projects') >= 0 ) {
      //   this.setState({selectedMenu:'projects'})
      // } else if (prevProps.path.indexOf('/deployments') >= 0 ) {
      //   this.setState({selectedMenu:'deployments'})
      // } else if (prevProps.path.indexOf('/pods') >= 0 ) {
      //   this.setState({selectedMenu:'pods'})
      // } else if (prevProps.path.indexOf('/network') >= 0 ) {
      //   this.setState({selectedMenu:'network'})
      // } else if (prevProps.path.indexOf('/settings') >= 0 ){
      //   this.setState({selectedMenu:'settings'})
      // }
      this.selectionMenu(prevProps.path);
    }
  }

  componentWillMount() {
    // var menu = window.location.pathname
    // // console.log(menu)
    // if (menu.indexOf('/dashboard') >= 0 ){
    //   this.setState({selectedMenu:'dashboard'})
    // } else if (menu.indexOf('/clusters') >= 0 ) {
    //   this.setState({selectedMenu:'clusters'})
    // } else if (menu.indexOf('/nodes') >= 0 ) {
    //   this.setState({selectedMenu:'nodes'})
    // } else if (menu.indexOf('/projects') >= 0 ) {
    //   this.setState({selectedMenu:'projects'})
    // } else if (menu.indexOf('/deployments') >= 0 ) {
    //   this.setState({selectedMenu:'deployments'})
    // } else if (menu.indexOf('/pods') >= 0 ) {
    //   this.setState({selectedMenu:'pods'})
    // } else if (menu.indexOf('/network') >= 0 ) {
    //   this.setState({selectedMenu:'network'})
    // } else if (menu.indexOf('/settings') >= 0 ){
    //   this.setState({selectedMenu:'settings'})
    // }
    this.selectionMenu(this.props.path);
  }

  selectionMenu = (path) => {
    if (path.indexOf("/dashboard") >= 0) {
      this.setState({ selectedMenu: "dashboard" });
    } else if (path.indexOf("/clusters") >= 0) {
      this.setState({ selectedMenu: "clusters" });
    } else if (path.indexOf("/nodes") >= 0) {
      this.setState({ selectedMenu: "nodes" });
    } else if (path.indexOf("/projects") >= 0) {
      this.setState({ selectedMenu: "projects" });
    } else if (path.indexOf("/deployments") >= 0) {
      this.setState({ selectedMenu: "deployments" });
    } else if (path.indexOf("/pods") >= 0) {
      this.setState({ selectedMenu: "pods" });
    } else if (path.indexOf("/network") >= 0) {
      this.setState({ selectedMenu: "network" });
    } else if (path.indexOf("/maintenance") >= 0) {
      this.setState({ selectedMenu: "maintenance" });
    } else if (path.indexOf("/settings") >= 0) {
      this.setState({ selectedMenu: "settings" });
    }
  };

  onLogout = (e) => {
    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-LG-LG02");

    // localStorage.removeItem("token");
    // localStorage.removeItem("userName");
    // localStorage.removeItem("roles");

    AsyncStorage.setItem("token", null);
    AsyncStorage.setItem("userName", null);
    AsyncStorage.setItem("roles", null);
    AsyncStorage.setItem("projects", null);
  };

  handleToggle = () => {
    this.setState({ open: !this.prevOpen.current });
  };

  handleClose = (event) => {
    this.setState({ open: false });
  };

  onSelectMenu = (e) => {
    this.setState({ selectedMenu: e.currentTarget.id });
  };

  onClick = (e) => {
    e.preventDefault();
  };

  componentDidUpdate() {}

  render() {
    // const handleListKeyDown = (event) => {
    //   if (event.key === 'Tab') {
    //     event.preventDefault();
    //     this.setState({open:false});
    //   }
    // }

    let userName = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userName = result;
    });
    return (
      <header className="main-header">
        <Link to="/dashboard" className="logo">
          <span className="logo-lg">
            <b>OpenMCP</b>
          </span>
        </Link>

        <nav className="navbar navbar-static-top">
          <div className="top-menu navbar-right">
            <div>
              <LanguageListMenu />
            </div>
            <div
              className={"main-menu " + this.state.selectedMenu}
              id="accounts"
              style={{ position: "relative", textAlign: "left" }}
            >
              <Link to="/" onClick={this.onClick}>
                <AccountCircleIcon />
                <div
                  style={{
                    position: "absolute",
                    display: "inline-block",
                    right: "15px",
                    top: "12px",
                  }}
                >
                  {userName}
                </div>
              </Link>
              <div className="sub-menu accounts">
                <Link to="/login" onClick={this.onLogout}>
                  Logout
                </Link>
                {/* <LanguageSwitch/> */}
              </div>
            </div>

            {/* <div className={"main-menu " + this.state.selectedMenu} id="settings" onClick={this.onSelectMenu}>
              <Link to="/settings/accounts" onClick={this.onSelectMenu}>
                <div style={{ fontSize: 20}}><Settings></Settings></div>
              </Link>
              <div className="sub-menu settings">
                <Link to="/settings/accounts" onClick={this.onSelectMenu}>Accounts</Link>
                <Link to="/settings/group-role" onClick={this.onSelectMenu}>Group Role</Link>
                <Link to="/settings/policy" onClick={this.onSelectMenu}>Policy</Link>
                <Link to="/settings/alert" onClick={this.onSelectMenu}>Alert</Link>
                <Link to="/settings/config" onClick={this.onSelectMenu}>Config</Link>
              </div>
            </div> */}
          </div>
        </nav>
      </header>
    );
  }
}

export default Head;
