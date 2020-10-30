// 각 컨텐츠의 왼쪽 고정 매뉴바
import React, { Component } from "react";
import "../../css/style.css";
// import {ArrowBackIos, NavigateNext} from '@material-ui/icons';
import { NavLink } from 'react-router-dom';
import * as fnMenuList from './LeftMenuData.js';
// import NavigateNextIcon from '@material-ui/icons/NavigateNext';

class LeftMenu extends Component {
  constructor(props){
    super(props);
    
    this.state = {
      params : this.props.menu, //프로젝트에 따라서 수정되야함
    }
  }
  render() {
    const menuList = fnMenuList.getMenu(this.props.title);
    // console.log(menuList[this.props.menu]);
    // console.log(ad, this.props.menu);
    // console.log("leftsubcomp : ", this.props.pathParam)

    const lists = [];
    menuList[this.props.menu].map((item) => {
      if(item.type === "single"){
        lists.push(
          <li className="" >
            <NavLink to={item.path} activeClassName="active">
              <i className="fa fa-dashboard"></i>
              <span>{item.title}</span>
            </NavLink>
          </li>
        )
      } else {
        lists.push(
          <li className="treeview">
            <NavLink to={item.path} activeClassName="active">
              <i className="fa fa-dashboard"></i>
              <span>{item.title}</span>
              <span className="pull-right-container">
              </span>
            </NavLink>

            <ul className="treeview-menu">
              {
                item.sub.map((subItem)=>{
                  return(
                    <li>
                      <NavLink to={subItem.path} activeClassName="active">
                        <i className="fa fa-circle-o"></i> <span>{subItem.title}</span>
                      </NavLink>
                    </li>
                  );
                })
              }
            </ul>
          </li>
        )
      }
    });
    return (
      <div>
        {this.props.title !== undefined ? 
        <aside className="main-sidebar">
        <section className="sidebar">
          <div className="user-panel">
            <div className="pull-left image">
            </div>
            <div className="pull-left info">
              <p>{this.props.title}</p>
              {/* <a href="/">{this.state.createDate}</a> */}
            </div>
          </div>
          
          <ul className="sidebar-menu tree" data-widget="tree">
            {lists}
          </ul>
        </section>
      </aside>
         : ""}
      </div>
    );
  }
}

export default LeftMenu;
