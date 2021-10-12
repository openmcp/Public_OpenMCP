import React, { Component } from "react";
import { NavLink, Link, Route, Switch } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

import PropTypes from "prop-types";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Box from "@material-ui/core/Box";
import { Container } from "@material-ui/core";
import { NavigateNext } from '@material-ui/icons';
import Pods from './Pods';
import HPA from './HPA';
import VPA from './VPA';
import { FiBox } from "react-icons/fi";


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

class PodMenu extends Component {
  state = {
    // rows: "",
    // completed: 0,
    reRender: "",
    value: 0,
    tabHeader: [
      { label: "pod", index: 1, param:"pod" },
      { label: "hpa", index: 2, param:"hpa" },
      { label: "vpa", index: 3, param:"vpa" },
    ],
  };

  componentWillMount() {
     if(this.props.match.url.indexOf("hpa") > 0 ){
       this.setState({ value: 1 });
     } else if(this.props.match.url.indexOf("vpa") > 0){
      this.setState({ value: 2 });
     } else {
      this.setState({ value: 0 });
     }
     this.props.menuData("none");
  }

  render() {
    const handleChange = (event, newValue) => {
      this.setState({ value: newValue });
    };
    const { classes } = this.props;


    


    return (
      <div>
        <div className="content-wrapper fulled">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
            
            <i><FiBox/></i>
              <span>Pods</span>
              <small>{this.props.match.params.project}</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
                <NavLink to="/clusters">Pods</NavLink>
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
                            pathname: `/pods/${i.param}`
                          }}
                          style={{minHeight:"42px", fontSize: "13px", minWidth:"100px"  }}
                    />
                    );
                  })}
                </Tabs>
              </AppBar>
              <TabPanel className="tab-panel" value={this.state.value} index={0}>
                <Switch>
                <Route path="/pods/pod"
                    render={({match,location}) => <Pods  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={1}>
               <Switch>
                  <Route path="/pods/hpa"
                    render={({match,location}) => <HPA  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={2}>
               <Switch>
                  <Route path="/pods/vpa"
                    render={({match,location}) => <VPA  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
             
            </div>
          </section>
        </div>
      </div>
    );
  }
}

export default withStyles(styles)(PodMenu);
