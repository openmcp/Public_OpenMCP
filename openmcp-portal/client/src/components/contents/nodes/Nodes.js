import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
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
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
// import Editor from "./../modules/Editor";
import { NavigateNext} from '@material-ui/icons';
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
import NdAddNode from './../modal/NdAddNode';
import FiberManualRecordSharpIcon from '@material-ui/icons/FiberManualRecordSharp';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import LinearProgressBar from './../../modules/LinearProgressBar';
import { CgServer} from "react-icons/cg";


class Nodes extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Node" },
        { name: "status", title: "Status" },
        { name: "cluster", title: "Cluster"},
        { name: "role", title: "Role" },
        { name: "system_version", title: "System Version" },
        { name: "cpu", title: "CPU" },
        { name: "memory", title: "Memory" },
        { name: "pods", title: "Pods" },
        { name: "provider", title: "Provider" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 250 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 100},
        { columnName: "role", width: 75 },
        { columnName: "system_version", width: 200 },
        { columnName: "cpu", width: 155 },
        { columnName: "memory", width: 155 },
        { columnName: "pods", width: 155 },
        { columnName: "provider", width: 100 },
        { columnName: "region", width: 90 },
        { columnName: "zone", width: 80 },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10, 
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      anchorEl:null,
      infoWindowOpen:null,
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }

  callApi = async () => {
    const response = await fetch(`/nodes`);
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
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));

    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, 'log-ND-VW01');
  };

  onUpdateData = () => {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));

    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PJ-VW03");
  };

  onRefresh = () => {
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  };

      //셀
  infoWindowOpen = (e) => {
    document.getElementById("info-window-"+e.currentTarget.id).style.visibility='visible';
  };

  infoWindowClose = (e) => {
    // this.setState({infoWindowOpen:null});
    document.getElementById("info-window-"+e.currentTarget.id).style.visibility='hidden';
  };

  render() {
    // 셀 데이터 스타일 변경
    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell>
        <span
          style={{
            color:
            value === "Healthy" ? "#1ab726"
              : value === "Unhealthy" ? "red"
                : value === "Unknown" ? "#b5b5b5" : "black",
          }}
        >
          <FiberManualRecordSharpIcon style={{fontSize:12, marginRight:4,
          backgroundColor: 
          value === "Healthy" ? "rgba(85,188,138,.1)"
            : value === "Unhealthy" ? "rgb(152 13 13 / 10%)"
              : value === "Unknown" ? "rgb(255 255 255 / 10%)" : "white",
          boxShadow: 
          value === "Healthy" ? "0 0px 5px 0 rgb(85 188 138 / 36%)"
            : value === "Unhealthy" ? "rgb(188 85 85 / 36%) 0px 0px 5px 0px"
              : value === "Unknown" ? "rgb(255 255 255 / 10%)" : "white",
          borderRadius: "20px",
          // WebkitBoxShadow: "0 0px 1px 0 rgb(85 188 138 / 36%)",
          }}></FiberManualRecordSharpIcon>
        </span>
        <span
          style={{
            color:
            value === "Healthy" ? "#1ab726" : 
              value === "Unhealthy" ? "red" : 
                value === "Unknown" ? "#b5b5b5" : "black"
          }}>
          {value}
        </span>
      </Table.Cell>
    );

    const Cell = (props) => {
      const { column, row } = props;
      const fnEnterCheck = () => {
        let data = []
        if(props.value === undefined){
          return ""
        } else {
          props.value.indexOf("|") > 0 ? 
            props.value.split("|").map( (item,index) => 
                data[index] = item
            ) : data[0] = props.value
          return data.length > 1 ? <p>{data[1] + " ("+data[0]+")" }</p> : props.value
        }
      }

      const fn_linearProgressBar = () =>{
        var data = [];
        if(props.value.indexOf("|") > -1) {
          props.value.split("|").forEach( item => {
            if(item.indexOf(" ") > -1) {
              item.split(" ").map((i, index) => data[index] = i);
            }
          });
        } else {
          data = [];
        }
        
        return (
          <p style={{marginTop:"5px"}}>
            <LinearProgressBar value={data[0]} total={data[2]}/>
          </p>
        )
      }

      if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "name") {
        return (
          <Table.Cell
            {...props}
            style={{ cursor: "pointer" }}
            aria-haspopup="true"
          >
            <InfoWindow data={row} rowId={props.tableRow.rowId}/>
            <Link  
            id={props.tableRow.rowId}
            onMouseEnter={this.infoWindowOpen}
            onMouseLeave={this.infoWindowClose}
            to={{
              pathname: `/nodes/${props.value}`,
              // search:`clustername=${row.cluster}&provider=${row.provider}`,
              search:`clustername=${row.cluster}`,
              state: {
                data : row
              }
            }}>
              {fnEnterCheck()}
            </Link>
          </Table.Cell>
        ) 
      } else if (column.name === "cpu" || column.name === "memory" || column.name === "pods"){
        return <Table.Cell>
          {fnEnterCheck()}
          {fn_linearProgressBar()}
          </Table.Cell>
        // 
      };
      return <Table.Cell>{fnEnterCheck()}</Table.Cell>;
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
      return <Table.Row {...props} key={props.tableRow.key}/>;
    };

    const handleClick = (event) => {
      if(this.state.anchorEl === null){
        this.setState({anchorEl : event.currentTarget});
      } else {
        this.setState({anchorEl : null});
      }
    };

    const handleClose = () => {
      this.setState({anchorEl : null});
    };

    const open = Boolean(this.state.anchorEl);

    return (
      <div className="content-wrapper nodes">
        {/* 컨텐츠 헤더 */}
        <section className="content-header" onClick={this.onRefresh}>
          <h1>
            {/* <span  style={{cursor:"pointer"}}> */}
            <i><CgServer/></i>
            <span>
              Nodes
            </span>
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">
              <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
              Nodes
            </li>
          </ol>
        </section>
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                // <Editor title="add node"/>,
                <div style={{
                  position: "absolute",
                  right: "21px",
                  top: "20px",
                  zIndex: "10",
                  textTransform: "capitalize",
                }}>
                  <IconButton
                    aria-label="more"
                    aria-controls="long-menu"
                    aria-haspopup="true"
                    onClick={handleClick}
                  >
                    <MoreVertIcon />
                  </IconButton>

                  <Popper open={open} anchorEl={this.state.anchorEl} role={undefined} transition disablePortal placement={'bottom-end'}>
                    {({ TransitionProps, placement }) => (
                      <Grow
                      {...TransitionProps}
                      style={{ transformOrigin: placement === 'bottom' ? 'center top' : 'center top' }}
                      >
                        <Paper>
                          <MenuList autoFocusItem={open} id="menu-list-grow">
                              <MenuItem style={{ textAlign: "center", display: "block", fontSize:"14px"}}>
                                <NdAddNode onUpdateData = {this.onUpdateData}
                                 menuClose={handleClose}/>
                              </MenuItem>
                            </MenuList>
                          </Paper>
                      </Grow>
                    )}
                  </Popper>
                </div>,
                <Grid
                  rows={this.state.rows}
                  columns={this.state.columns}
                >
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState defaultValue="" />
                  <IntegratedFiltering />
                  <SearchPanel style={{ marginLeft: 0 }} />

                  {/* Sorting */}
                  <SortingState
                    defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
                  />
                  <IntegratedSorting />

                  {/* 페이징 */}
                  <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* 테이블 */}
                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
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

class InfoWindow extends Component {
  constructor(props) {
    super(props);
    this.state = {
      row : this.props.data
    }
  }
  render(){
    return(
      
      <div className="info-window" id={"info-window-" + this.props.rowId}
        style={{visibility: "hidden" }}
        >
        <div className="info">
          <div>{this.state.row.name}</div> 
        </div>
      </div>
      
    )
  }
}

export default Nodes;
