import React, { Component } from "react";
import "../../css/style.css";
import { Settings } from "@material-ui/icons";

// import styled from "styled-components";
// import oc from "open-color";
import { Link, NavLink } from 'react-router-dom';


// const Mainheader = styled.div`
// /* 레이아웃 */
//     display: flex;
//     position: fixed;
//     align-items: center;
//     justify-content: center;
//     height: 60px;
//     width: 100%;
//     top: 0px;
//     z-index: 5;

//     /* 색상 */
//     background: ${oc.indigo[6]};
//     color: white;
//     border-bottom: 1px solid ${oc.indigo[7]};
//     box-shadow: 0 3px 6px rgba(0,0,0,0.10), 0 3px 6px rgba(0,0,0,0.20);

//     /* 폰트 */
//     font-size: 2.5rem;
// `;

class Head extends Component {
  // onSelectTopMenu = (e) => {
  //   this.props.onSelectMenu(false,"");
  // };

  // onSelectLeftMenu = (e) => {
  //   this.props.onSelectMenu(true,"");
  // };

  // const userName = sessionStorage.getItem("username");
  // componentWillMount(){
  //   userName = 
  // }
  
  onLogout = (e) => {
    // debugger;
    // localStorage.removeItem("token");
    sessionStorage.removeItem("token");
  }
  render() {
    const userName = sessionStorage.getItem("userName");
    // console.log(sessionStorage.getItem("userName"));
    return (
      <header className="main-header">
        <NavLink to="/dashboard" className="logo">
          <span className="logo-lg">
            <b>OpenMCP</b>
          </span>
        </NavLink>

        <nav className="navbar navbar-static-top">
          <div className="top-menu navbar-left">
            <NavLink to="/dashboard" activeClassName="active">
              <span>Dashboard</span>
            </NavLink>
            <NavLink to="/clusters" activeClassName="active" >
              <span>Clusters</span>
            </NavLink>
            <NavLink to="/nodes" activeClassName="active" >
              <span>Nodes</span>
            </NavLink>
            <NavLink to="/projects" activeClassName="active" >
              <span>Projects</span>
            </NavLink>
            <NavLink to="/pods" activeClassName="active" >
              <span>Pods</span>
            </NavLink>
          </div>
          {/* <a
            href="/"
            className="sidebar-toggle"
            data-toggle="push-menu"
            role="button"
          >
          </a> */}

          <div className="navbar-custom-menu">
            <ul className="nav navbar-nav">
              <li className="dropdown user user-menu">
                <Link to="/login" className="dropdown-toggle" data-toggle="dropdown" onClick={this.onLogout}>
                  <span className="hidden-xs">{userName}</span>
                </Link>
              </li>
              <li className="dropdown user user-menu">
                <a href="/" className="dropdown-toggle" data-toggle="dropdown">
                  {/* <span className="hidden-xs">설정아이콘</span> */}
                  <Settings style={{ fontSize: 20 }}></Settings>
                </a>
              </li>
            </ul>
          </div>
        </nav>
      </header>
    );
  }
}

export default Head;
