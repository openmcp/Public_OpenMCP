import React, { Component } from "react";
import { Link, Route, Switch } from "react-router-dom";
import { withStyles } from "@material-ui/core/styles";

import PropTypes from "prop-types";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Box from "@material-ui/core/Box";
import { Container } from "@material-ui/core";
// import { NavigateNext } from '@material-ui/icons';
import ConfigEKS from './ConfigEKS';
import ConfigGKE from './ConfigGKE';
import ConfigAKS from './ConfigAKS';
import ConfigKVM from './ConfigKVM';

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

class PublicCloud extends Component {
  state = {
    // rows: "",
    // completed: 0,
    reRender: "",
    value: 0,
    tabHeader: [
      { label: "EKS", index: 0, param:"eks" },
      { label: "GKE", index: 1, param:"gke" },
      { label: "AKS", index: 2, param:"aks" },
      { label: "KVM", index: 3, param:"kvm" },
    // { label: "DaemonSets", index: 3 },
    ],
  };

  componentWillMount(){
    var menu = window.location.pathname
    if (menu.indexOf('/eks') >= 0 ){
      this.setState({value: 0})
    } else if (menu.indexOf('/gke') >= 0 ) {
      this.setState({value: 1})
    } else if (menu.indexOf('/aks') >= 0 ) {
      this.setState({value: 2})
    } else if (menu.indexOf('/kvm') >= 0 ) {
      this.setState({value: 3})
    } 
  }
  

  render() {
    const handleChange = (event, newValue) => {
      this.setState({ value: newValue });
    };
    const { classes } = this.props;
    return (
      <div>
        <div className="sub-content-wrapper">
          {/* 내용부분 */}
          <section class="pca-tab">
            {/* 탭매뉴가 들어간다. */}
            <div className={classes.root}>
              <AppBar position="static" className="app-bar">
                <Tabs
                  value={this.state.value}
                  onChange={handleChange}
                  aria-label="simple tabs example"
                  // style={{ backgroundColor: "#257790" }}
                  style={{ backgroundColor: "#ecf0f5",
                    padding: "14px 12px 0px 12px",
                    minHeight:"35px"
                  }}
                  // indicatorColor="primary"
                  // indicatorColor="primary"
                  TabIndicatorProps ={{ style:{backgroundColor:"#9ccee6"}}}
                >
                  {this.state.tabHeader.map((i) => {
                    return (
                    <Tab label={i.label} {...a11yProps(i.index)}
                          component={Link}
                          to={{
                            pathname: `/settings/config/public-cloud/${i
                              .param}`
                          }}
                          style={{minHeight:"35px", fontSize: "13px", minWidth:"100px"  }}
                          
                    />
                    );
                  })}
                </Tabs>
              </AppBar>
              <TabPanel className="tab-panel" value={this.state.value} index={0}>
                <Switch>
                  <Route path="/settings/config/public-cloud/eks"
                    render={({match,location}) => <ConfigEKS  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={1}>
                <Switch>
                  <Route path="/settings/config/public-cloud/gke"
                    render={({match,location}) => <ConfigGKE  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={2}>
                <Switch>
                  <Route path="/settings/config/public-cloud/aks"
                    render={({match,location}) => <ConfigAKS  match={match} location={location} menuData={this.onMenuData}/>} >
                  </Route>
                </Switch>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={3}>
                <Switch>
                  <Route path="/settings/config/public-cloud/kvm"
                    render={({match,location}) => <ConfigKVM  match={match} location={location} menuData={this.onMenuData}/>} >
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

export default withStyles(styles)(PublicCloud);
