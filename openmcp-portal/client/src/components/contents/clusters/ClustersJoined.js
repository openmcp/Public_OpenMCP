import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
import { Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  SelectionState,
  IntegratedSelection,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
} from "@devexpress/dx-react-grid-material-ui";
// import Editor from "./../modules/Editor";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import Confirm from './../../modules/Confirm';
// import ProgressTemp from './../../modules/ProgressTemp';
import PieReChartMini from '../../modules/PieReChartMini';
import FiberManualRecordSharpIcon from '@material-ui/icons/FiberManualRecordSharp';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';

import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
import LinearProgressBar from "../../modules/LinearProgressBar.js";
import axios from "axios";
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';

class ClustersJoined extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
        { name: "nodes", title: "Nodes" },
        { name: "cpu", title: "CPU" },
        { name: "ram", title: "Memory" },
        { name: "provider", title: "Provider" },
        // { name: "disk", title: "Disk" },
        // { name: "network", title: "Network" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 100 },
        { columnName: "status", width: 100 },
        { columnName: "region", width: 100 },
        { columnName: "zone", width: 80 },
        { columnName: "nodes", width: 100 },
        { columnName: "cpu", width: 150 },
        { columnName: "ram", width: 150 },
        { columnName: "provider", width: 150 },
        // { columnName: "disk", width: 150 },
        // { columnName: "network", width: 150 },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      selection: [],
      selectedRow: "",

      confirmInfo : {
        title :"Cluster Unjoin Confrim",
        context :"Are you sure you want to cancel the Cluster Join?",
        button : {
          open : "UNJOIN",
          yes : "UNJOIN",
          no : "CANCEL",
        }  
      },
      confrimTarget : "false",
      openProgress : false,
      anchorEl: null,
    };
  }
  componentWillMount() {
    // this.props.menuData("none");
    
  }

  callApi = async () => {
    const response = await fetch("/clusters");
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
    utilLog.fn_insertPLogs(userId, "log-CL-VW01");
  }

  confirmed = (result) => {
    if(result) {
      //Unjoin proceed
      console.log("confirmed");

      const url = `/cluster/unjoin`;
      const data = {
        clusterName: this.state.selectedRow.name,
      };
  
      axios
        .post(url, data)
        .then((res) => {
          // alert(res.data.message);
          this.setState({ open: false });
          this.onRefresh();
        })
        .catch((err) => {
          alert(err);
        });


      let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
      utilLog.fn_insertPLogs(userId, "log-CL-MO03");
    } else {
      console.log("cancel")
    }
  }
 
  onRefresh = () => {
    // this.timer = setInterval(this.progress2, 20);
    // clearInterval(this.timer);
   
    if(this.state.openProgress){
      this.setState({openProgress:false})
    } else {
      this.setState({openProgress:true})
    }

    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
      })
      .catch((err) => console.log(err));
  };

  closeProgress = () => {
    this.setState({openProgress:false})
  }
  
  render() {
    // 셀 데이터 스타일 변경
    const HighlightedCell = ({ value, style, row, ...restProps }) => {
      var cpuPct = parseFloat(row.cpu.split("/")[0])/parseFloat(row.cpu.split("/")[1].split(" Core")[0]) * 100
      var memPct = parseFloat(row.ram.split("/")[0])/parseFloat(row.ram.split("/")[1].split(" Gi")[0]) * 100
      // console.log(cpuPct, memPct)
      var status = cpuPct >= 90 || memPct >= 90 ? "Warning" : value
      return (
      <Table.Cell
        {...restProps}
        style={{
          // backgroundColor:
          //   value === "Healthy" ? "white" : value === "Unhealthy" ? "white" : undefined,
          // cursor: "pointer",
          ...style,
        }}
      >
        <span
          style={{
            color:
            status === "Healthy" ? "#1ab726"
              : status === "Unhealthy" ? "red"
                : status === "Unknown" ? "#b5b5b5"
                  : status === "Warning" ? "#ff8042" : "black",
          }}
        >
          <FiberManualRecordSharpIcon style={{fontSize:12, marginRight:4,
          backgroundColor: 
            status === "Healthy" ? "rgba(85,188,138,.1)"
                : status === "Unhealthy" ? "rgb(152 13 13 / 10%)"
                  : status === "Unknown" ? "rgb(255 255 255 / 10%)"
                    : status === "Warning" ? "rgb(109 31 7 / 10%)" : "white",
          boxShadow: 
            status === "Healthy" ? "0 0px 5px 0 rgb(85 188 138 / 36%)"
                : status === "Unhealthy" ? "rgb(188 85 85 / 36%) 0px 0px 5px 0px"
                  : status === "Unknown" ? "rgb(255 255 255 / 10%)"
                    : status === "Warning" ? "rgb(188 114 85 / 36%) 0px 0px 5px 0px" : "white",
          borderRadius: "20px",
          // WebkitBoxShadow: "0 0px 1px 0 rgb(85 188 138 / 36%)",
          }}></FiberManualRecordSharpIcon>
        </span>
          
        <span
          style={{
            color:
            status === "Healthy"
                ? "#1ab726"
                : status === "Unhealthy"
                ? "red"
                : status === "Unknown"
                ? "#b5b5b5"
                : status === "Warning"
                ? "#ff8042"
                : "black",
          }}
        >
          {status}
        </span>
      </Table.Cell>
    )};

    //셀
    const Cell = (props) => {
      const { column, row } = props;

      const fn_linearProgressBar = () =>{
        var data = [];
        if(props.value.indexOf(" ") > -1) {
          props.value.split(" ").forEach( (item) => {
            if(item.indexOf("/") > -1) {
              item.split("/").map((i, index) => data[index] = i);
            }
          });
        } else {
          data = [];
        }

        var percent = data[0]/data[1]*100;
        
        return (
          <div>
            <p>{props.value+ " (" + percent.toFixed(1) + "%)"}</p>
            <p style={{marginTop:"5px"}}>
              <LinearProgressBar value={data[0]} total={data[1]}/>
            </p>
          </div>
        )
      }

      // console.log("cell : ", props);
      if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "name") {
        return (
          <Table.Cell {...props} style={{ cursor: "pointer" }}>
            <Link
              to={{
                pathname: `/clusters/${props.value}/overview`,
                state: {
                  data: row,
                },
              }}
            >
              {props.value}
            </Link>
          </Table.Cell>
        );
      } else if (column.name === "cpu" || column.name === "ram"){
        return <Table.Cell>
          {/* <p>{props.value}</p> */}
          {fn_linearProgressBar()}
          </Table.Cell>
        // 
      };
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
      // console.log("row!!!!!! : ",props);
      return <Table.Row {...props} key={props.tableRow.key} />;
    };

    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({
        selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {} ,
        confrimTarget : this.state.rows[selection[0]] ? this.state.rows[selection[0]].name : "false" ,
      });
    };

    const handleClick = (event) => {
      if(this.state.anchorEl === null){
        this.setState({anchorEl : event.currentTarget});
      } else {
        this.setState({anchorEl : null});
      }
    };

    const handleClose = () => {
      this.setState({anchorEl : false});
    };

    const open = Boolean(this.state.anchorEl);

    return ([
      <div className="sub-content-wrapper cluster-list">
      {/* {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""} */}
        {/* 컨텐츠 헤더 */}
        {/* <section className="content-header" onClick={this.onRefresh} >
          <h1>
            <i><FaBuffer/></i>
            <span>
              Joined Clusters
            </span>
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">
              <NavigateNext
                style={{ fontSize: 12, margin: "-2px 2px", color: "#444" }}
              />
              Clusters
            </li>
          </ol>
        </section> */}
        {this.state.rows ? (
          [
            <section className="content" style={{ position: "relative" }}>
              
              <Paper>
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
                                <Confirm confirmInfo={this.state.confirmInfo} confrimTarget ={this.state.confrimTarget} confirmed={this.confirmed}
                                menuClose={handleClose}
                                />
                              </MenuItem>
                            </MenuList>
                          </Paper>
                      </Grow>
                    )}
                  </Popper>
                </div>
                <Grid rows={this.state.rows} columns={this.state.columns}>
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState className="search-Satste" defaultValue="" />
                  {/* <Confirm confirmInfo={this.state.confirmInfo} confrimTarget ={this.state.confrimTarget} confirmed={this.confirmed}/> */}
                  <SearchPanel className="search-Satste" style={{ marginLeft: 200, backgroundColor:"#000000" }} />

                  {/* Sorting */}
                  <SortingState
                    defaultSorting={[
                      { columnName: "status", direction: "desc" },
                    ]}
                  />

                  {/* 페이징 */}
                  <PagingState
                    defaultCurrentPage={0}
                    defaultPageSize={this.state.pageSize}
                  />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  />

                  <IntegratedFiltering />
                  <IntegratedSorting />
                  <IntegratedSelection />
                  <IntegratedPaging />

                  {/* 테이블 */}
                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing
                    defaultColumnWidths={this.state.defaultColumnWidths}
                  />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                  <TableSelection
                    selectByRowClick
                    highlightRow
                  />
                </Grid>
              </Paper>
            </section>,
            <section
              className="content"
              style={{ position: "relative" }}
            ><ResourceStatus data={this.state.rows}/>
            </section>
          ]
        ) : (
          <CircularProgress
            variant="determinate"
            value={this.state.completed}
            style={{ position: "absolute", left: "50%", marginTop: "20px" }}
          ></CircularProgress>
        )}
      </div>
      ]
    );
  }
}

class ResourceStatus extends Component{
  constructor(props){
    super(props);
    this.state = {
      data : this.props.data
    }
  }

  angle = {
    full : {
      startAngle : 0,
      endAngle : 360
    },
    half : {
      startAngle : 180,
      endAngle : 0
    }  
  }

  componentWillUpdate(prevProps, prevState){
    if (this.props.data !== prevProps.data) {
        this.setState({
          data: prevProps.data,
        });
      }
  }

  render(){
    const colors = [
      "#51BFFF",
      "#ecf0f5",
    ];
    return(
      <div>
      <div className="rs-title">Cluster Resource Status</div>
      {/* <div className="rap-contents" style={{display: "flex",overflow: "auto"}}> */}
      <div className="rap-contents" style={{overflow: "auto"}}>
        {this.state.data.map((item)=>{
          return(
            
            <Paper className="rs-status">
              <div className="status-title">
                <div>{item.name}</div>
              </div>

              <div class="status-content">
                <PieReChartMini name="CPU" data={item.resourceUsage.cpu} angle={this.angle.full} unit={item.resourceUsage.cpu.unit} colors={colors}></PieReChartMini>
                <div class="sub-content">
                  <div>{item.resourceUsage.cpu.status[0].name}</div>
                  <div>{item.resourceUsage.cpu.status[0].value.toFixed(1)+" "+item.resourceUsage.cpu.unit }</div>
                </div>
                <div class="sub-content">
                  <div> {item.resourceUsage.cpu.status[1].name}</div>
                  <div>{item.resourceUsage.cpu.status[1].value.toFixed(1)+" "+item.resourceUsage.cpu.unit }</div>
                </div>
              </div>

              <div class="status-content">
                <PieReChartMini name="Memory" data={item.resourceUsage.memory} angle={this.angle.full} unit={item.resourceUsage.memory.unit} colors={colors}></PieReChartMini>
                <div class="sub-content">
                  <div>{item.resourceUsage.memory.status[0].name}</div>
                  <div>{item.resourceUsage.memory.status[0].value.toFixed(1)+" "+item.resourceUsage.memory.unit }</div>
                </div>
                <div class="sub-content">
                  <div> {item.resourceUsage.memory.status[1].name}</div>
                  <div>{item.resourceUsage.memory.status[1].value.toFixed(1)+" "+item.resourceUsage.memory.unit }</div>
                </div>
              </div>

              <div class="status-content">
                <PieReChartMini name="Storage" data={item.resourceUsage.storage} angle={this.angle.full} unit={item.resourceUsage.cpu.unit} colors={colors}></PieReChartMini>
                <div class="sub-content">
                  <div>{item.resourceUsage.storage.status[0].name}</div>
                  <div>{item.resourceUsage.storage.status[0].value.toFixed(1)+" "+item.resourceUsage.storage.unit }</div>
                </div>
                <div class="sub-content">
                  <div> {item.resourceUsage.storage.status[1].name}</div>
                  <div>{item.resourceUsage.storage.status[1].value.toFixed(1)+" "+item.resourceUsage.storage.unit }</div>
                </div>
              </div>
              
              
              
              

              {/* <div className="status-title">
                <div>{item.name}</div>
              </div>
              <div className="status-content">
                <div style={{color:(item.status === "Warning" ? "#ff8042" : "#000000")}}>
                  <span>CPU: </span>
                  <span>{item.cpu}</span>
                </div>
                <div style={{color:(item.status === "Warning" ? "#ff8042" : "#000000")}}>
                  <span>Memory: </span>
                  <span>{item.ram}</span>
                </div>
                <div>
                  <span>Disk: </span>
                  <span>{item.disk}</span>
                </div>
                <div>
                  <span>Network: </span>
                  <span>{item.network}</span>
                </div>

              </div> */}
            </Paper>
          )
        })}
        </div>
      </div>
    )
  }
}



export default ClustersJoined;
