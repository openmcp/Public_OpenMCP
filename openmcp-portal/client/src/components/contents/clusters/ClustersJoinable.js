import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
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
import * as utilLog from '../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
import Confirm from './../../modules/Confirm';
import FiberManualRecordSharpIcon from '@material-ui/icons/FiberManualRecordSharp';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';

import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
import axios from "axios";
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';


class ClustersJoinable extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "endpoint", title: "Endpoint" },
        { name: "provider", title: "Provider" },
        { name: "region", title: "Region" },
        { name: "zone", title: "Zone" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 180 },
        { columnName: "endpoint", width: 200},
        { columnName: "provider", width: 130 },
        { columnName: "region", width: 130 },
        { columnName: "zone", width: 130 },
      ],
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10, 
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      selection: [],
      selectedRow: "",

      confirmInfo : {
        title :"Cluster Join Confrim",
        context :"Are you sure you want to Join the Cluster?",
        button : {
          open : "JOIN",
          yes : "JOIN",
          no : "CANCEL",
        }  
      },
      confrimTarget : "false",
      editorContext : ``,
      anchorEl: null,
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }
  

  callApi = async () => {
    const response = await fetch("/clusters-joinable");
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
    utilLog.fn_insertPLogs(userId, 'log-CL-VW01');
  };

  confirmed = (result) => {
    if(result) {
      //Unjoin proceed
      console.log("confirmed")
      
      const url = `/cluster/join`;
      const data = {
        clusterName: this.state.selectedRow.name,
        clusterAddress : this.state.selectedRow.endpoint
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


  render() {

    // 셀 데이터 스타일 변경
    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell
        {...restProps}
        style={{
          // backgroundColor:
          //   value === "Healthy" ? "white" : value === "Unhealthy" ? "white" : undefined,
          // cursor: "pointer",
          ...style,
        }}>
        <span
          style={{
            color:
            value === "Healthy" ? "#1ab726"
              : value === "Unhealthy" ? "red"
                : value === "Unknown" ? "#b5b5b5"
                  : value === "Warning" ? "#ff8042" : "black",
          }}
        >
          <FiberManualRecordSharpIcon style={{fontSize:12, marginRight:4,
          backgroundColor: 
          value === "Healthy" ? "rgba(85,188,138,.1)"
            : value === "Unhealthy" ? "rgb(152 13 13 / 10%)"
              : value === "Unknown" ? "rgb(255 255 255 / 10%)"
                : value === "Warning" ? "rgb(109 31 7 / 10%)" : "white",
          boxShadow: 
          value === "Healthy" ? "0 0px 5px 0 rgb(85 188 138 / 36%)"
            : value === "Unhealthy" ? "rgb(188 85 85 / 36%) 0px 0px 5px 0px"
              : value === "Unknown" ? "rgb(255 255 255 / 10%)"
                : value === "Warning" ? "rgb(188 114 85 / 36%) 0px 0px 5px 0px" : "white",
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

    //셀
    const Cell = (props) => {
      const { column, row } = props;
      // console.log("cell : ", props);
      if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } else if  (column.name === "provider") {
        if(row.provider === ""){
          return <Table.Cell>-</Table.Cell>
        }
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
      // console.log("row!!!!!! : ",props);
      return <Table.Row {...props} key={props.tableRow.key}/>;
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
      this.setState({anchorEl : null});
    };

    const open = Boolean(this.state.anchorEl);
    return (
      <div className="sub-content-wrapper fulled">
        {/* 컨텐츠 헤더 */}
        {/* <section className="content-header" onClick={this.onRefresh}>
          <h1>
            <span onClick={this.onRefresh} style={{cursor:"pointer"}}>
              Joinable Clusters
            </span>
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">
              <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
              Clusters
            </li>
          </ol>
        </section> */}
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
                {/* <Editor title="create" context={this.state.editorContext}/> */}
        {this.state.rows ? (
                    <Grid
                      rows={this.state.rows}
                      columns={this.state.columns}
                    >
                      <Toolbar />
                      {/* 검색 */}
                      <SearchState defaultValue="" />
                      <SearchPanel style={{ marginLeft: 0 }} />

                      {/* Sorting */}
                      <SortingState
                        defaultSorting={[{ columnName: 'status', direction: 'desc' }]}
                      />

                      {/* 페이징 */}
                      <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
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
                      <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
                      <TableHeaderRow
                        showSortingControls
                        rowComponent={HeaderRow}
                      />
                      <TableSelection
                        selectByRowClick
                        highlightRow
                      />
                    </Grid>
        
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

export default ClustersJoinable;
