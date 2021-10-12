import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
import { NavLink } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  IntegratedSelection,
  SelectionState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  TableSelection,
  PagingPanel,
  // TableColumnVisibility
} from "@devexpress/dx-react-grid-material-ui";
import { NavigateNext} from '@material-ui/icons';
import * as utilLog from './../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
// import AddMembers from "./AddMembers";
// import Editor from "../../modules/Editor";
// import AcChangeRole from './../modal/AcChangeRole';
import IconButton from "@material-ui/core/IconButton";
import MenuItem from "@material-ui/core/MenuItem";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';
import GrCreateGroup from './../modal/GrCreateGroup';
import GrEditGroup from './../modal/GrEditGroup';
import axios from 'axios';
import Confirm2 from './../../modules/Confirm2';
import { SiGraphql} from "react-icons/si";


class GroupRole extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "group_name", title: "Group Name"},
        { name: "description", title: "Description" },
        { name: "member", title: "Member" },
        { name: "role", title: "Roles" },
        { name: "projects", title: "Project" },
        { name: "group_id", title: "Group ID" },
      ],
      defaultColumnWidths: [
        { columnName: "group_name", width: 150 },
        { columnName: "description", width: 200 },
        { columnName: "member", width: 200 },
        { columnName: "role", width: 400 },
        { columnName: "projects", width: 400 },
        { columnName: "group_id", width: 0 },
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
      anchorEl: null,

      grEditOpen: false,

      confirmOpen: false,
      confirmInfo : {
        title :"Delete Group Role",
        context :"Are you sure you want to Delete Group Role?",
        button : {
          open : "",
          yes : "CONFIRM",
          no : "CANCEL",
        }
      },
      confrimTarget : "",
      confirmTargetKeyname:"Group Name"

    };
  }

  componentWillMount() {
    this.props.menuData("none");
  }

  callApi = async () => {
    const response = await fetch(`/settings/group-role`);
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  componentDidMount() {
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
    utilLog.fn_insertPLogs(userId, 'log-AC-VW01');
  };

  onUpdateData = () => {
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        this.setState({ 
          selection : [],
          selectedRow : "",
          rows: res 
        });

        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  };

   Cell = (props) => {
    const { column } = props;
    // const { column } = props;
    let stringData = []

    const arrayToString = () => {
      if (props.value === null){
        stringData = ["-"]
      } else {
        stringData = props.value.reduce((result, item, index, arr) => {
          if (index+1 === arr.length){
            return `${result}${item}`
          } else {
            return `${result}${item} | `
          }
        }, "")
      }

      return stringData
    }

    if (column.name === "role" || column.name === "member" || column.name === "projects") {
      return (
        <Table.Cell {...props}>
          {arrayToString()}
        </Table.Cell>
      );
    } 
    // else if (column.name == "group_name") {
    //   return (
    //     <Table.Cell {...props} 
    //       style={{ cursor: "pointer", color:"#3c8dbc", marginLeft:"20px"}} 
    //       // onClick = {this.editDialogOpen}
    //       onClick= {e => this.editDialogOpen(props.row)}
    //     >
    //         {props.value}
    //     </Table.Cell>
    //   )
    // }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

   HeaderRow = ({ row, ...restProps }) => (
    <Table.Row
      {...restProps}
      style={{
        cursor: "pointer",
        backgroundColor: "whitesmoke",
      }}
    />
  );
  Row = (props) => {
    return <Table.Row {...props} key={props.tableRow.key}/>;
  };

  // editDialogOpen = (row) => {
  //   this.setState({
  //     grEditOpen : true,
  //     selectedRow : row
  //   });
  // }
  handleDeleteClick = (e) => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select a Group Role");
      return;
    } else {
      this.setState({
        confirmOpen: true,
      })
    }
  }

  confirmed = (result) => {
    this.setState({confirmOpen:false});

    //show progress loading...
    this.setState({openProgress:true});

    if(result) {
      const url = `/settings/group-role`;

      const data = {
        group_id : this.state.selectedRow.group_id,
      };

      axios.delete(url, {data:data})
        .then((res) => {
          alert(res.data.message);
          this.setState({ open: false });
          this.handleClose();
          this.onUpdateData();
        })
        .catch((err) => {
        });
      
      this.setState({openProgress:false})
      
      // loging Add Node
      let userId = null;
      AsyncStorage.getItem("userName",(err, result) => { 
        userId= result;
      })
        utilLog.fn_insertPLogs(userId, "log-ND-MD02");
    } else {
      this.setState({openProgress:false})
    }
  };

  handleClick = (event) => {
    if(this.state.anchorEl === null){
      this.setState({anchorEl : event.currentTarget});
    } else {
      this.setState({anchorEl : null});
    }
  };

  handleClose = () => {
    this.setState({ 
      anchorEl: null,
      selection: [],
      selectedRow: ""});
  };

  render() {
    const onSelectionChange = (selection) => {
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({ 
        selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {},
        confrimTarget : this.state.rows[selection[0]] ? this.state.rows[selection[0]].group_name : "false" ,
      });


    };
    

    const open = Boolean(this.state.anchorEl);

    return (
      <div className="content-wrapper fulled">
        <Confirm2
          confirmInfo={this.state.confirmInfo} 
          confrimTarget ={this.state.confrimTarget} 
          confirmTargetKeyname = {this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}/>

        <section className="content-header">
          <h1>
          <i><SiGraphql/></i>
          <span>Group Role</span>
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">
              <NavigateNext style={{fontSize:12, margin: "-2px 2px", color: "#444"}}/>
              Settings
            </li>
          </ol>
        </section>
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                
                <div
                style={{
                  position: "absolute",
                  right: "21px",
                  top: "20px",
                  zIndex: "10",
                  textTransform: "capitalize",
                }}
              >
                <IconButton
                  aria-label="more"
                  aria-controls="long-menu"
                  aria-haspopup="true"
                  onClick={this.handleClick}
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
                            <MenuItem
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                            >
                              <GrCreateGroup onUpdateData={this.onUpdateData} menuClose={this.handleClose}/>
                            </MenuItem>
                            <MenuItem
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                              >
                              <GrEditGroup rowData={this.state.selectedRow} onUpdateData={this.onUpdateData} menuClose={this.handleClose}
                              />
                            </MenuItem>
                            <MenuItem
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                              >
                              <div onClick={this.handleDeleteClick}>Delete Group</div>
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
                  <SearchState defaultValue="" />
                  <SearchPanel style={{ marginLeft: 0 }} />


                  <PagingState defaultCurrentPage={0} defaultPageSize={this.state.pageSize} />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  <SortingState
                    defaultSorting={[{ columnName: 'user_id', direction: 'asc' }]}
                  />

                  <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  />

                  <IntegratedFiltering />
                  <IntegratedSelection />
                  <IntegratedSorting />
                  <IntegratedPaging />

                  <Table cellComponent={this.Cell} rowComponent={this.Row} />
                  <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={this.HeaderRow}
                  />
                  
                  <TableSelection
                    selectByRowClick
                    highlightRow
                    // showSelectionColumn={false}
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

export default GroupRole;