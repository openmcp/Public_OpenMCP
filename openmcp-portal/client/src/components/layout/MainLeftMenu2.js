import React, { Component } from "react";
import { NavLink } from "react-router-dom";
import { AiFillDashboard, AiOutlineDeploymentUnit, AiOutlineUser, AiFillAlert, AiOutlineSetting, AiOutlineAreaChart} from "react-icons/ai";
import { RiDashboardFill } from "react-icons/ri";
import { GrBundle } from "react-icons/gr";
import { FaBuffer} from "react-icons/fa";
import { CgServer } from "react-icons/cg";
import { FiBox } from "react-icons/fi";
import { BiNetworkChart } from "react-icons/bi";
import { HiOutlineDuplicate,HiOutlineCamera } from "react-icons/hi";
import { RiInboxUnarchiveLine } from "react-icons/ri";
import { IoKeyOutline } from "react-icons/io5";
import { SiGraphql } from "react-icons/si";





class LeftMenu2 extends Component {
  // constructor(props) {
  //   super(props);

  //   // this.state = {
  //   //   params: this.props.menu,
  //   // };
  // }

  // shouldComponentUpdate(prevProps, prevState) {
  //   if (
  //     this.props.menu !== prevProps.menu ||
  //     this.props.title !== prevProps.title
  //   ) {
  //     return true;
  //   } else {
  //     if (
  //       this.props.pathParams.hasOwnProperty("searchString") &&
  //       this.props.pathParams.searchString !== prevProps.pathParams.searchString
  //     ) {
  //       return true;
  //     }
  //     return false;
  //   }
  // }

  render() {
    const menuList = [
      {
        type: "multi",
        title: "Monitor",
        path: "/dashboard",
        icon : <AiFillDashboard className="leftMenu-main-icon"/>,
        sub: [
          {
            title: "Dashboard",
            path: "/dashboard",
            icon : <RiDashboardFill className="leftMenu-sub-icon"/>,
          },
          // {
          //   title: "Multiple Metrics",
          //   path: "/dashboard",
          //   icon : <AiOutlineDashboard className="leftMenu-sub-icon"/>,
          // },
        ],
      },
      {
        type: "multi",
        title: "Multiple Clusters",
        path: "/clusters",
        icon : <GrBundle className="leftMenu-main-icon"/>,
        sub: [
          {
            title: "Clusters",
            path: "/clusters/joined",
            icon : <FaBuffer className="leftMenu-sub-icon"/>,
          },
          {
            title: "Nodes",
            path: "/nodes",
            icon : <CgServer className="leftMenu-sub-icon"/>,
          },
          
        ],
      },
      {
        type: "multi",
        title: "Workloads",
        path: "/projects",
        icon : <HiOutlineDuplicate className="leftMenu-main-icon"/>,
        sub: [
          {
            title: "Projects",
            path: "/projects",
            icon : <HiOutlineDuplicate className="leftMenu-sub-icon"/>,
          },
          {
            title: "Deployments",
            path: "/deployments",
            icon : <AiOutlineDeploymentUnit className="leftMenu-sub-icon"/>,
          },
          {
            title: "Pods",
            path: "/pods",
            icon : <FiBox className="leftMenu-sub-icon"/>,
          },
          {
            title: "Network",
            path: "/network/dns",
            icon : <BiNetworkChart className="leftMenu-sub-icon"/>,
          },
        ],
      },
      {
        type: "multi",
        title: "Advenced",
        path: "/maintenance/Migrations",
        icon : <AiFillDashboard className="leftMenu-main-icon"/>,
        sub: [
          {
            title: "Migrations",
            path: "/maintenance/migration",
            icon : <RiInboxUnarchiveLine className="leftMenu-sub-icon"/>,
          },
          {
            title: "Snapshots",
            path: "/maintenance/snapshot",
            icon : <HiOutlineCamera className="leftMenu-sub-icon"/>,
          },
        ],
      },
      {
        type: "multi",
        title: "Settings",
        path: "/settings",
        icon : <AiFillDashboard className="leftMenu-main-icon"/>,
        sub: [
          {
            title: "Accounts",
            path: "/settings/accounts",
            icon : <AiOutlineUser className="leftMenu-sub-icon"/>,
          },
          {
            title: "Group Role",
            path: "/settings/group-role",
            icon : <SiGraphql className="leftMenu-sub-icon"/>,
          },
          {
            title: "Policy",
            path: "/settings/policy",
            icon : <IoKeyOutline className="leftMenu-sub-icon"/>,
          },
          {
            title: "Alert",
            path: "/settings/alert",
            icon : <AiFillAlert className="leftMenu-sub-icon"/>,
          },
          {
            title: "Meterings",
            path: "/settings/meterings",
            icon : <AiOutlineAreaChart className="leftMenu-sub-icon"/>,
          },
          {
            title: "Config",
            path: "/settings/config",
            icon : <AiOutlineSetting className="leftMenu-sub-icon"/>,
          },
        ],
      },
    ];
    const lists = [];
    menuList.forEach((item) => {
      lists.push(
        <li className="treeview left-main-menu">
          <div className="sidebar-main-title">
            
            <span>{item.title}</span>
          </div>
          <ul className="treeview-menu sidebar-sub-title">
            {item.sub.map((subItem) => {
              return (
                <li>
                  <NavLink
                    to={{
                      pathname: `${subItem.path}`,
                    }}
                    activeClassName="active"
                  >
                    {subItem.icon}
                    <span>{subItem.title}</span>
                  </NavLink>
                </li>
              );
            })}
          </ul>
        </li>
      );
    });
    // console.log("this.props.title", this.props.title)
    return (
      <div>
        <aside className="main-sidebar">
          <section className="sidebar">
            <ul className="sidebar-menu tree" data-widget="tree">
              {lists}
            </ul>
          </section>
        </aside>
      </div>
    );
  }
}

export default LeftMenu2;
