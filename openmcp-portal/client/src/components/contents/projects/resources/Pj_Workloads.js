import React, { Component } from "react";
import { NavLink, Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableHeaderRow,
  PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import { withStyles } from "@material-ui/core/styles";

import PropTypes from "prop-types";
import AppBar from "@material-ui/core/AppBar";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Typography from "@material-ui/core/Typography";
import Box from "@material-ui/core/Box";
import Paper from "@material-ui/core/Paper";
import Editor from "../../../common/Editor";
import { Container } from "@material-ui/core";

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

let apiParams = "";
class Pj_Workloads extends Component {
  state = {
    rows: "",
    completed: 0,
    reRender: "",
    value: 0,
    tabHeader: [{ label: "Deployments", index: 1 },
    { label: "SatetfulSets", index: 2 },
    // { label: "DaemonSets", index: 3 },
    // { label: "CronJobs", index: 4 },
    // { label: "Jobs", index: 5 }
    ],
  };

  componentWillMount() {
    //왼쪽 메뉴쪽에 타이틀 데이터 전달
    const result = {
      menu : "clusters",
      title : this.props.match.params.name
    }
    this.props.menuData(result);
    apiParams = this.props.match.params;
  }
  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        this.setState({ rows: res });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

  callApi = async () => {
    var param = this.props.match.params.name;
    const response = await fetch(`/projects/${param}/overview`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

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

    // console.log("Pj_Workloads_Render : ", this.state.rows.basic_info);
    const { classes } = this.props;
    return (
      <div>
        <div className="content-wrapper">
          {/* 컨텐츠 헤더 */}
          <section className="content-header">
            <h1>
              Workloads
              <small>{this.props.match.params.name}</small>
            </h1>
            <ol className="breadcrumb">
              <li>
                <NavLink to="/dashboard">Home</NavLink>
              </li>
              <li>
                <NavLink to="/Projects">Projects</NavLink>
              </li>
              <li className="active">Resources</li>
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
                  style={{ backgroundColor: "#16586c" }}
                  indicatorColor="primary"
                  // indicator={{backgroundColor:"#ffffff"}}
                  // TabIndicatorProps ={{ ind:"#635ee7"}}
                >
                  {this.state.tabHeader.map((i) => {
                    return <Tab label={i.label} {...a11yProps(i.index)} />;
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
                <DeploymentsTab pathParam={this.props.match.params.name}></DeploymentsTab>
              </TabPanel>
              <TabPanel className="tab-panel" value={this.state.value} index={1}>
                Item Two
              </TabPanel>
              <TabPanel  className="tab-panel"value={this.state.value} index={2}>
                Item Three
              </TabPanel>
            </div>
          </section>
        </div>
      </div>
    );
  }
}


class DeploymentsTab extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "project_name", title: "Name" },
        { name: "project_status", title: "Status" },
        { name: "project_creator", title: "Createor" },
        { name: "project_create_time", title: "Created Time" },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      setPageSize: 5,
      pageSizes: [5, 10, 15, 0],

      completed: 0,
    };
  }

  componentWillMount() {
    // this.props.onSelectMenu(false, "");
  }

  

  callApi = async () => {
    const response = await fetch("/api/projects");
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  //컴포넌트가 모두 마운트가 되었을때 실행된다.
  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        this.setState({ rows: res });
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

  render() {

    // 셀 데이터 스타일 변경
    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell
        {...restProps}
        style={{
          backgroundColor:
            value === "Healthy" ? "white" : value === "Unhealthy" ? "white" : undefined,
          cursor: "pointer",
          ...style,
        }}
      >
        <span
          style={{
            color:
              value === "Healthy" ? "green" : value === "Unhealthy" ? "red" : undefined,
          }}
        >
          {value}
        </span>
      </Table.Cell>
    );

    //셀
    const Cell = (props) => {
      const { column, row } = props;
      if (column.name === "project_status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "project_name") {
        return (
          // <Table.Cell
          //   component={Link}
          //   to={{
          //     pathname: `/projects/${this.props.pathParam}/resources/workloads/deployments/${props.value}`,
          //     state: {
          //       data : row
          //     }
          //   }}
          //   {...props}
          //   style={{ cursor: "pointer" }}
          // ></Table.Cell>
          <Table.Cell
            {...props}
            style={{ cursor: "pointer" }}
          ><Link to={{
            pathname: `/projects/${this.props.pathParam}/resources/workloads/deployments/${props.value}`,
            state: {
              data : row
            }
          }}>{props.value}</Link></Table.Cell>
        );
      }
      return <Table.Cell {...props} />;
    };

    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
          // ...styles[row.sector.toLowerCase()],
        }}
        // onClick={()=> alert(JSON.stringify(row))}
      />
    );
    const Row = (props) => {
      return <Table.Row {...props} />;
    };

    return (
      <div className="content-wrapper full">
        {/* 컨텐츠 헤더 */}
       
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                <Editor />,
                <Grid
                  rows={this.state.rows}
                  columns={this.state.columns}
                >
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState defaultValue="" />
                  <IntegratedFiltering />
                  <SearchPanel style={{ marginLeft: 0 }} />

                  {/* 페이징 */}
                  <PagingState defaultCurrentPage={0} defaultPageSize={5} />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* Sorting */}
                  <SortingState
                  // defaultSorting={[{ columnName: 'city', direction: 'desc' }]}
                  />
                  <IntegratedSorting />

                  {/* 테이블 */}
                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                </Grid>,
              ]
            ) : (
              <CircularProgress
                variant="determinate"
                value={this.state.completed}
                style={{ position: "absolute", left: "50%", marginTop: "20px" }}
              ></CircularProgress>
            )}
          </Paper>
        </section>
      </div>
    );
  }
}


export default withStyles(styles)(Pj_Workloads);
