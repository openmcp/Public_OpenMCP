import React, { Component } from "react";
import "../../css/style.css";

class LeftSlideMenu extends Component {
  render() {
    return (
      <aside className="main-sidebar">
        <section className="sidebar">
          <div className="user-panel">
            <div className="pull-left image">
              <img
                src="https://s.gravatar.com/avatar/065d2006dd053e2bf0a44e5c90e4bfc6?s=80"
                className="img-circle"
                alt="User Image"
              />
            </div>
            <div className="pull-left info">
              <p>kh lee</p>
              <a href="#">
                <i className="fa fa-circle text-success"></i> Logged In
              </a>
            </div>
          </div>
          <ul className="sidebar-menu tree" data-widget="tree">
            <li className="header">USER ACTIONS</li>
            <li className="active">
              <a href="/dashboard/">
                <i className="fa fa-dashboard"></i> <span>Dashboard</span>
              </a>
            </li>

            <li className="">
              <a href="/domain/add">
                <i className="fa fa-plus"></i> <span>New Domain</span>
              </a>
            </li>

            <li className="header">ADMINISTRATION</li>
            <li className="">
              <a href="/admin/pdns">
                <i className="fa fa-info-circle"></i> <span>PDNS</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/global-search">
                <i className="fa fa-search"></i> <span>Global Search</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/history">
                <i className="fa fa-calendar"></i> <span>History</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/templates/list">
                <i className="fa fa-clone"></i> <span>Domain Templates</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/manage-account">
                <i className="fa fa-industry"></i> <span>Accounts</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/manage-user">
                <i className="fa fa-users"></i> <span>Users</span>
              </a>
            </li>
            <li className="">
              <a href="/admin/manage-keys">
                <i className="fa fa-key"></i> <span>API Keys</span>
              </a>
            </li>
            <li className="treeview">
              <a href="#">
                <i className="fa fa-cog"></i> <span>Settings</span>
                <span className="pull-right-container">
                  <i className="fa fa-angle-left pull-right"></i>
                </span>
              </a>
              <ul className="treeview-menu">
                <li>
                  <a href="/admin/setting/basic">
                    <i className="fa fa-circle-o"></i> <span>Basic</span>
                  </a>
                </li>
                <li>
                  <a href="/admin/setting/dns-records">
                    <i className="fa fa-circle-o"></i> <span>Records</span>
                  </a>
                </li>

                <li>
                  <a href="/admin/setting/pdns">
                    <i className="fa fa-circle-o"></i> <span>PDNS</span>
                  </a>
                </li>
                <li>
                  <a href="/admin/setting/authentication">
                    <i className="fa fa-circle-o"></i>{" "}
                    <span>Authentication</span>
                  </a>
                </li>
              </ul>
            </li>
          </ul>
        </section>
      </aside>
    );
  }
}

export default LeftSlideMenu;
