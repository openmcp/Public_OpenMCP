import React, { Component } from "react";
import { NavLink, Link, Route, Switch } from "react-router-dom";
// import CircularProgress from "@material-ui/core/CircularProgress";
// import {
//   SearchState,
//   IntegratedFiltering,
//   PagingState,
//   IntegratedPaging,
//   SortingState,
//   IntegratedSorting,
// } from "@devexpress/dx-react-grid";
// import {
//   Grid,
//   Table,
//   Toolbar,
//   SearchPanel,
//   TableHeaderRow,
//   PagingPanel,
// } from "@devexpress/dx-react-grid-material-ui";
import { withStyles } from "@material-ui/core/styles";

import PropTypes from "prop-types";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
// import Typography from "@material-ui/core/Typography";
import Box from "@material-ui/core/Box";
// import Paper from "@material-ui/core/Paper";
// import Editor from "../../../modules/Editor";
import { Container } from "@material-ui/core";
import PjwDeployments from './PjwDeployments';
import PjwDeploymentDetail from './PjwDeploymentDetail';
import { NavigateNext } from '@material-ui/icons';
import PjwStatefulsets from './PjwStatefulsets';
import PjwStatefulSetDetail from './PjwStatefulSetDetail';

const styles = (theme) => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
  },
  // indicator: {
  //   display: 'flex',
  //   justifyContent: 'center',
  //   backgroundColor: 'transparent',
  //   '& > span': {
  //     maxWidth: 40,
  //     width: '100%',
  //     backgroundColor: '#635ee7',
  //   },
  // },
});

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Container>
          <Box>
            {children}
          </Box>
        </Container>
      )}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

class PjWorkloads extends Component {
  state = {
    // rows: "",
    // completed: 0,
    reRender: "",
    value: 0,
    tabHeader: [
      { label: "Deployments", index: 1, param:"deployments" },
      { label: "StatefulSets", index: 2, param:"statefulsets" },
    // { label: "DaemonSets", index: 3 },
    // { label: "CronJobs", index: 4 },
    // { label: "Jobs", index: 5 }
    ],
  };

  componentWillMount() {
    //왼쪽 메뉴쪽에 타이틀 데이터 전달

    const result = {
      menu : "projects",
      title : this.props.match.params.project,
      pathParams : {
        searchString : this.props.location.search,
        project : this.props.match.params.project
      }
    }
    this.props.menuData(result);
     if(this.props.match.url.indexOf("statefulsets") > 0 ){
       this.setState({ value: 1 });
     } else {
      this.setState({ value: 0 });
     }
    // apiParams = this.props.match.params.project;
  }

  render() {
    const handleChange = (event, newValue) => {
      this.setState({ value: newValue });
    };

    // const StyledTabs = withStyles({
    //   indicator: {
    //     display: 'flex',
    //     justifyContent: 'center',
    //     backgroundColor: 'transparent',
    //     '& > span': {
    //       maxWidth: 40,
    //       width: '100%',
    //       backgroundColor: '#635ee7',
    //     },
    //   },
    // })((props) => <Tabs {...props} TabIndicatorProps={{ children: <span /> }} />);

    // console.log("PjWorkloads_Render : ", this.state.rows.basic_info);
    const { classes } = this.props;
    return (
      <div>
        <div className="content-wrapper">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
              Workloads
              <small>{this.props.match.params.project}</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                <NavLink to="/Projects">Projects</NavLink>
              </li>
              
              <li className="active">
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                Resources
              </li>
            </ol>
          </section>

          {/* 내용부분 */}
          <section>
            {/* 탭매뉴가 들어간다. */}
            <div className={classes.root}>
              <AppBar position="static" className="app-bar">
                <Tabs
                  value={this.state.value}
                  onChange={handleChange}
                  aria-label="simple tabs example"
                  style={{
                    backgroundColor: "#3c8dbc",
                    minHeight: "42px",
                    position: "fixed",
                    width: "100%",
                    zIndex: "990",
                  }}
                  TabIndicatorProps ={{ style:{backgroundColor:"#00d0ff"}}}
                >
                  {this.state.tabHeader.map((i) => {
                    return (
                    <Tab label={i.label} {...a11yProps(i.index)}
                          component={Link}
                          to={{
                            pathname: `/projects/${this.props.match.params.project}/resources/workloads/${i.param}`,
                            search : this.props.location.search
                            // state: {
                            //   data : row
                            // }
                          }}
                          style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}
                    />
                    );
                  })}
                  {/* <Tab label="Item One" {...a11yProps(0)} />
                  <Tab label="Item Two" {...a11yProps(1)} />
                  <Tab label="Item Three" {...a11yProps(2)} /> */}
                </Tabs>
              </AppBar>
              {/* {this.props.rows.map((i) => {
                    return (
                      <Tab label={i.lable} {...a11yProps(i.index)} />
                      <TabPanel value={this.state.value} index={0}></TabPanel>
                      );
                  })} */}
              <TabPanel className="tab-panel" value={this.state.value} index={0}>
                <Switch>
                  <Route path="/projects/:project/resources/workloads/deployments/:deployment" 
                    render={({match,location}) => <PjwDeploymentDetail  match={match} location={location}/>} >
                  </Route>
                  <Route path="/projects/:project/resources/workloads/deployments" 
                    render={({match,location}) => <PjwDeployments  match={match} location={location}/>} >
                  </Route>
            
                  {/* <Route path="/projects/:project/resources/workloads/deployments/:deployment/pods/:pod" 
                    render={({match,location}) => <PjwDeployment_PodDetail  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                  <Route path="/projects/:project/resources/workloads/deployments/:deployment/containers/:container" 
                    render={({match,location}) => <PjwDeployment_ContainerDetail  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route> */}
                </Switch>
                {/* <DeploymentsTab pathParam={this.props.match.params.project}></DeploymentsTab> */}
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={1}>
                <Switch>
                  <Route path="/projects/:project/resources/workloads/statefulsets/:statefulset" 
                    render={({match,location}) => <PjwStatefulSetDetail  match={match} location={location}/>} >
                  </Route>
                  <Route path="/projects/:project/resources/workloads/statefulsets" 
                        render={({match,location}) => <PjwStatefulsets  match={match} location={location}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              {/* <TabPanel  className="tab-panel"value={this.state.value} index={2}>
                Item Three
              </TabPanel> */}
            </div>
          </section>
        </div>
      </div>
    );
  }
}



export default withStyles(styles)(PjWorkloads);
