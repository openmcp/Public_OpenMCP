import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import CloseIcon from "@material-ui/icons/Close";
// import SelectBox from "../../modules/SelectBox";
import * as utilLog from "../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import CircularProgress from "@material-ui/core/CircularProgress";
import FiberManualRecordSharpIcon from '@material-ui/icons/FiberManualRecordSharp';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  IconButton,
  Typography,
  TextField,
} from "@material-ui/core";
import {
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  SelectionState,
  IntegratedSelection,
  TableColumnVisibility,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  TableColumnResizing,
  TableHeaderRow,
  PagingPanel,
  TableSelection,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
import axios from 'axios';
import LensIcon from '@material-ui/icons/Lens';

import Stepper from '@material-ui/core/Stepper';
import Step from '@material-ui/core/Step';
import StepLabel from '@material-ui/core/StepLabel';

// import Typography from "@material-ui/core/Typography";
// import DialogActions from "@material-ui/core/DialogActions";
// import DialogContent from "@material-ui/core/DialogContent";
// import Button from "@material-ui/core/Button";
// import Dialog from "@material-ui/core/Dialog";
// import IconButton from "@material-ui/core/IconButton";
// import axios from 'axios';
// import { ContactlessOutlined } from "@material-ui/icons";

const styles = (theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
  backButton: {
    marginRight: theme.spacing(1),
  },
  instructions: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
});
let preRoleSelection=[];
let preRoleSelectedRow=[];
let preUserSelection=[];
let preUserSelectedRow=[];
let preProjectSelection=[];
let preProjectSelectedRow=[];
class GrEditGroup extends Component {
  constructor(props) {
    super(props);
    this.state = {
      groupId : "",
      groupName : "",
      description : "",
      open: this.props.open,

      rows : [],

      selectedRoleIds : [], //선택된 row의 role_id
      roleSelectionId : [], //컬럼 index

      selectedUserIds : [],
      userSelectionId : [],

      selectedProjects : [],
      projectSelectionId : [],

      activeStep : 0,
    };
    // this.onChange = this.onChange.bind(this);
  }
  componentWillMount(){
    preRoleSelection=[];
    preRoleSelectedRow=[];
    preUserSelection=[];
    preUserSelectedRow=[];
    preProjectSelection=[];
    preProjectSelectedRow=[];
  }

  onChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value,
    });
  };

  handleClickOpen = () => {
    if (Object.keys(this.props.rowData).length  === 0) {
      alert("Please select a Group Role");
      this.setState({ open: false });
      return;
    }

    this.setState({ 
      open: true,
      groupId : this.props.rowData.group_id,
      groupName : this.props.rowData.group_name,
      description : this.props.rowData.description,
      selectedRoleIds : [], //선택된 row의 role_id
      roleSelectionId : [], //컬럼 index

      selectedUserIds : [],
      userSelectionId : [],

      selectedProjects : [],
      projectSelectionId : [],
    });
  };

  handleClose = () => {
    this.setState({
      groupName : "",
      description : "",
      rows : [],
      selectedRoleIds : [],
      roleSelectionId : [],
      selectedUserIds : [],
      userSelectionId : [],
      selectedProjects : [],
      projectSelectionId : [],
      
      activeStep : 0,
      open: false,
    });
    preRoleSelection = [];
    preRoleSelectedRow = [];
    preUserSelection = [];
    preUserSelectedRow = [];
    this.props.menuClose();
  };

  handleUpdate = (e) => {
    if (this.state.groupName === "") {
      alert("Please insert 'group name' data");
      return;
    } else if (this.state.description === "") {
      alert("Please insert 'description' data");
      return;
    } else if (Object.keys(this.state.selectedRoleIds).length === 0) {
      alert("Please select roles");
      return;
    } else if (Object.keys(this.state.selectedUserIds).length === 0) {
      alert("Please select users");
      return;
    }

    // Update user role
    const url = `/settings/group-role`;
    const data = {
      group_id : this.state.groupId,
      groupName:this.state.groupName,
      description:this.state.description,
      role_id:this.state.selectedRoleIds,
      user_id:this.state.selectedUserIds,
      projects:this.state.selectedProjects,
    };
    axios.put(url, data)
    .then((res) => {
        alert(res.data.message);
        this.setState({ open: false });
        this.props.menuClose();
        this.props.onUpdateData();
    })
    .catch((err) => {
        alert(err);
      });


    // loging deployment migration
    let userId = null;
    AsyncStorage.getItem("userName",(err, result) => { 
      userId= result;
    })
    utilLog.fn_insertPLogs(userId, "log-PJ-MD01");

    //close modal popup
    this.setState({ open: false });
  };

  onSelectRoles = (rows, selectionId) => {
    let roleIds = [];
    rows.forEach((role) => {
      roleIds.push(role.role_id);
    });

    this.setState({
      selectedRoleIds : roleIds,
      roleSelectionId : selectionId
    })
  }

  onSelectUsers = (rows, selectionId) => {
    let userIds = [];
    rows.forEach((user) => {
      userIds.push(user.user_id);
    });

    this.setState({
      selectedUserIds : userIds,
      userSelectionId : selectionId
    })
  }

  onSelectProjects = (rows, selectionId) => {
    let projects = [];
    rows.forEach((project) => {
      projects.push(project.name);
    });

    this.setState({
      selectedProjects : projects,
      projectSelectionId : selectionId
    })
  }

  render() {
    const DialogTitle = withStyles(styles)((props) => {
      const { children, classes, onClose, ...other } = props;
      return (
        <MuiDialogTitle disableTypography className={classes.root} {...other}>
          <Typography variant="h6">{children}</Typography>
          {onClose ? (
            <IconButton
              aria-label="close"
              className={classes.closeButton}
              onClick={onClose}
            >
              <CloseIcon/>
            </IconButton>
          ) : null}
        </MuiDialogTitle>
      );
    });

    const steps = ['Set Group Informations', 'Select Group Roles','Select Projects', 'Select Group Users'];
    const handleNext = () => {
      switch (this.state.activeStep){
        case 0 :
          if (this.state.groupName === "") {
            alert("Please insert 'group name' data");
            return;
          } else if (this.state.description === "") {
            alert("Please insert 'description' data");
            return;
          } else {
            this.setState({activeStep : this.state.activeStep + 1});
            return;
          }
        case 1:
          if (Object.keys(this.state.selectedRoleIds).length === 0) {
            alert("Please select roles");
            return;
          } else {
            this.setState({activeStep : this.state.activeStep + 1});
            return;
          }
        case 2:
          if (Object.keys(this.state.selectedProjects).length === 0) {
            alert("Please select projects");
            return;
          } else {
            this.setState({activeStep : this.state.activeStep + 1});
            return;
          }
        default :
          return;
      }
    };
  
    const handleBack = () => {
      this.setState({activeStep : this.state.activeStep - 1});
    };
    return (
      <div>
        <div
          onClick={this.handleClickOpen}
          style={{
            zIndex: "10",
            width: "148px",
            textTransform: "capitalize",
          }}
        >
          Edit Group
        </div>
        <Dialog
          // onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth="md"
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Edit Group Role
          </DialogTitle>
          <DialogContent dividers>
            <div className="md-contents-body small-grid">
            <Stepper activeStep={this.state.activeStep} alternativeLabel>
              {steps.map((label) => (
                <Step key={label}>
                  <StepLabel>{label}</StepLabel>
                </Step>
              ))}
            </Stepper>
            <div>
            <Typography>
              {this.state.activeStep === 0 ? (
                <div>
                  <section className="md-content">
                    <p>Group Role Name </p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      placeholder="group role name"
                      variant="outlined"
                      value={this.state.groupName}
                      fullWidth={true}
                      name="groupName"
                      onChange={this.onChange}
                    />
                  </section>
                  <section className="md-content">
                    <p>Description</p>
                    <TextField
                      id="outlined-multiline-static"
                      rows={1}
                      placeholder="group information"
                      variant="outlined"
                      value={this.state.description}
                      fullWidth={true}
                      name="description"
                      onChange={this.onChange}
                    />
                  </section>
                </div>
              ) : this.state.activeStep === 1 ? (
                <section className="md-content">
                  <GrRoles 
                    selection={this.state.roleSelectionId}
                    onSelectedRoles={this.onSelectRoles}
                    preSelection={this.props.rowData.role}
                  />
                </section>
              ) : this.state.activeStep === 2 ? (
                <section className="md-content">
                  <GrProjects 
                    selection={this.state.projectSelectionId}
                    onSelectedProjects={this.onSelectProjects}
                    preSelection={this.props.rowData.projects}
                  />
                </section>
              ) : (
                <section className="md-content">
                  <GrUsers 
                    selection={this.state.userSelectionId}
                    onSelectedUsers={this.onSelectUsers}
                    preSelection={this.props.rowData.member}
                  />
                </section>
              )}
            </Typography>
          </div>
              
             
            </div>
          </DialogContent>
          <DialogActions>
            <div>
              <Button
                disabled={this.state.activeStep === 0}
                onClick={handleBack}
              >
                Back
              </Button>
              {this.state.activeStep === steps.length - 1 ? (
                <Button onClick={this.handleUpdate} color="primary">
                  update
                </Button>
              ) : (
                <Button color="primary" onClick={handleNext}>
                  next
                </Button>
              )}
              
            </div>
            <Button onClick={this.handleClose} color="primary">
              cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

class GrRoles extends Component{
  constructor(props){
    super(props);
    this.state = {
      columns: [
        { name: "role_name", title: "Role" },
        { name: "description", title: "Description" },
        { name: "role_id", title: "Role id" },
      ],
      defaultColumnWidths: [
        { columnName: "role_name", width: 200 },
        { columnName: "description", width: 480 },
        { columnName: "role_id", width: 0 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 0,
      pageSizes: [5, 10, 15, 0],
      rows:[],
      selection: this.props.selection,
      selectedRow : [],

      preSelection : this.props.preSelection,
    }
  }

  callApi = async () => {
    const response = await fetch("/account-roles");
    const body = await response.json();
    return body;
  };

  componentWillMount(){
    this.callApi()
    .then((res) => {
      this.setState({ rows: res});
      let selectedRows = [];
      this.props.selection.forEach((id) => {
        selectedRows.push(res[id]);
      });
      this.setState({ selectedRow: selectedRows});
      })
      .catch((err) => console.log(err));
  }

  
  Cell = (props) => {
    const { column } = props;

    if (column.name === "role_name") {
      
      if(this.props.preSelection.includes(props.value)){
        if(!preRoleSelection.includes(props.tableRow.rowId)){
          preRoleSelection.push(props.tableRow.rowId);
          preRoleSelectedRow.push(props.row)
          this.setState({
            selection: preRoleSelection,
          });

          this.props.onSelectedRoles(preRoleSelectedRow, preRoleSelection);
        }
      }
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

  Row = (props) => {
    console.log("Row")
    return <Table.Row {...props} key={props.tableRow.key}/>;
  };

  HeaderRow = ({ row, ...restProps }) => (
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


  render(){
    console.log("rende");
    const onSelectionChange = (selection) => {
      let selectedRows = [];
      selection.forEach((id) => {
        selectedRows.push(this.state.rows[id]);
      });
      this.setState({ selectedRow: selectedRows});
      this.setState({ selection: selection });

      this.props.onSelectedRoles(selectedRows, selection);
    };

    return(
      <div>
        <p>Select Group Roles</p>
          <div id="md-content-info" style={{display:"block", minHeight:"95px",marginBottom:"10px"}}>
              {this.state.selectedRow.length > 0 
                ? this.state.selectedRow.map((row)=>{
                  return (
                    <span>
                      <LensIcon style={{fontSize:"8px", marginRight:"5px"}}/>
                      {row.role_name}
                    </span>
                  );
                }) 
                : <div style={{
                  color:"#9a9a9a",
                  textAlign: "center",
                  paddingTop: "30px"}}>
                    Please Select Roles
                  </div>}
          </div>
        {/* <p>Select Role</p> */}
        <Paper>
          <Grid rows={this.state.rows} columns={this.state.columns}>
            {/* <Toolbar /> */}
            {/* 검색 */}
            {/* <SearchState defaultValue="" />
            <SearchPanel style={{ marginLeft: 0 }} /> */}

            {/* Sorting */}
            <SortingState
              defaultSorting={[{ columnName: "status", direction: "asc" }]}
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
            <Table cellComponent={this.Cell} rowComponent={this.Row}/>
            <TableColumnResizing
              defaultColumnWidths={this.state.defaultColumnWidths}
            />
            <TableHeaderRow
              showSortingControls
              rowComponent={this.HeaderRow}
            />
            <TableColumnVisibility
              defaultHiddenColumnNames={['role_id']}
            />
            <TableSelection
              selectByRowClick
              highlightRow
            />
          </Grid>
        </Paper>
      </div>
    );
  }
}

class GrUsers extends Component{
  constructor(props){
    super(props);
    this.state = {
      columns: [
        { name: "user_id", title: "User ID" },
        { name: "role", title: "Roles" },
      ],
      defaultColumnWidths: [
        { columnName: "user_id", width: 200 },
        { columnName: "role", width: 500 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],
      rows:[],
      selection: this.props.selection,
      selectedRow : [],

      preSelection : this.props.preSelection,
    }
  }

  callApi = async () => {
    const response = await fetch("/settings/accounts");
    const body = await response.json();
    return body;
  };

  componentWillMount(){
    this.callApi()
      .then((res) => {
        this.setState({ rows: res});
        let selectedRows = [];
        this.props.selection.forEach((index) => {
          selectedRows.push(res[index]);
        });
        this.setState({ selectedRow: selectedRows});
        })
      .catch((err) => console.log(err));
  }

  Cell = (props) => {
    const { column } = props;

    if (column.name === "user_id") {
      
      if(this.props.preSelection.includes(props.value)){
        if(!preUserSelection.includes(props.tableRow.rowId)){
          preUserSelection.push(props.tableRow.rowId);
          preUserSelectedRow.push(props.row)
          this.setState({
            selection: preUserSelection,
          });

          this.props.onSelectedUsers(preUserSelectedRow, preUserSelection);
        }
      }
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };


  render(){
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

    const onSelectionChange = (selection) => {
      let selectedRows = [];
      selection.forEach((id) => {
        selectedRows.push(this.state.rows[id]);
      });
      this.setState({ selectedRow: selectedRows});
      this.setState({ selection: selection });

      this.props.onSelectedUsers(selectedRows, selection);
    };

    return(
      <div>
        <p>Select Users</p>
          <div id="md-content-info" style={{display:"block", minHeight:"95px",marginBottom:"10px"}}>
              {this.state.selectedRow.length > 0 
                ? this.state.selectedRow.map((row)=>{
                  return (
                    <span>
                      <LensIcon style={{fontSize:"8px", marginRight:"5px"}}/>
                      {row.user_id}
                    </span>
                  );
                }) 
                : <div style={{
                  color:"#9a9a9a",
                  textAlign: "center",
                  paddingTop: "30px"}}>
                    Please Select Roles
                  </div>}
          </div>
        {/* <p>Select Role</p> */}
        <Paper>
          <Grid rows={this.state.rows} columns={this.state.columns}>
            {/* <Toolbar /> */}
            {/* 검색 */}
            {/* <SearchState defaultValue="" />
            <SearchPanel style={{ marginLeft: 0 }} /> */}

            {/* Sorting */}
            <SortingState
              defaultSorting={[{ columnName: "status", direction: "asc" }]}
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
            <Table cellComponent={this.Cell}/>
            <TableColumnResizing
              defaultColumnWidths={this.state.defaultColumnWidths}
            />
            <TableHeaderRow
              showSortingControls
              rowComponent={HeaderRow}
            />
            <TableColumnVisibility
              defaultHiddenColumnNames={['role_id']}
            />
            <TableSelection
              selectByRowClick
              highlightRow
            />
          </Grid>
        </Paper>
      </div>
    );
  }
}

class GrProjects extends Component{
  constructor(props){
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Name" },
        { name: "status", title: "Status" },
        { name: "cluster", title: "Cluster" },
        { name: "created_time", title: "Created Time" },
        { name: "labels", title: "Labels" },
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 200 },
        { columnName: "status", width: 100 },
        { columnName: "cluster", width: 100 },
        { columnName: "created_time", width: 180 },
        { columnName: "labels", width: 180 },
      ],
      rows:[],
      
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],
      
      selection: this.props.selection,
      selectedRow : [],
      completed: 0,

      preSelection : this.props.preSelection,
    }
  }

  callApi = async () => {
    const response = await fetch("/projects");
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  componentDidMount(){
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
          let selectedRows = [];
          this.props.selection.forEach((index) => {
            selectedRows.push(res[index]);
          });
          this.setState({ selectedRow: selectedRows});
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

  render(){

    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell>
        <span
          style={{
            color:
            value === "Active" ? "#1ab726"
              : value === "Deactive" ? "red" : "black",
          }}
        >
          <FiberManualRecordSharpIcon style={{fontSize:12, marginRight:4,
          backgroundColor: 
          value === "Active" ? "rgba(85,188,138,.1)"
            : value === "Deactive" ? "rgb(152 13 13 / 10%)" : "white",
          boxShadow: 
          value === "Active" ? "0 0px 5px 0 rgb(85 188 138 / 36%)"
            : value === "Deactive" ? "rgb(188 85 85 / 36%) 0px 0px 5px 0px" : "white",
          borderRadius: "20px",
          // WebkitBoxShadow: "0 0px 1px 0 rgb(85 188 138 / 36%)",
          }}></FiberManualRecordSharpIcon>
        </span>
        <span
          style={{
            color:
              value === "Active" ? "#1ab726" 
                : value === "Deactive" ? "red" : undefined,
          }}
        >
          {value}
        </span>
      </Table.Cell>
    );

    const Cell = (props) => {

      const fnEnterCheck = (prop) => {
        var arr = [];
        var i;
        for(i=0; i < Object.keys(prop.value).length; i++){
          const str = Object.keys(prop.value)[i] + " : " + Object.values(prop.value)[i]
          arr.push(str)
        }
        return (
         arr.map(item => {
           return (
             <p>{item}</p>
           )
         })
        )
        // return (
          // props.value.indexOf("|") > 0 ? 
          //   props.value.split("|").map( item => {
          //     return (
          //       <p>{item}</p>
          //   )}) : 
          //     props.value
        // )
      }

      const { column } = props;
      // console.log("cell : ", props);


      if (column.name === "name") {
        if(this.props.preSelection.includes(props.value)){
          if(!preProjectSelection.includes(props.tableRow.rowId)){
            preProjectSelection.push(props.tableRow.rowId);
            preProjectSelectedRow.push(props.row)
            this.setState({
              selection: preProjectSelection,
            });
  
            this.props.onSelectedProjects(preProjectSelectedRow, preProjectSelection);
          }
        }
      } else if (column.name === "status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "labels"){
        return (
        <Table.Cell>{fnEnterCheck(props)}</Table.Cell>
        )
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

    const onSelectionChange = (selection) => {
      let selectedRows = [];
      selection.forEach((id) => {
        selectedRows.push(this.state.rows[id]);
      });
      this.setState({ selectedRow: selectedRows});
      this.setState({ selection: selection });

      this.props.onSelectedProjects(selectedRows, selection);
    };

    return(
      <div>
        <p>Select Projects</p>
          <div id="md-content-info" style={{display:"block", minHeight:"95px",marginBottom:"10px"}}>
              {this.state.selectedRow.length > 0 
                ? this.state.selectedRow.map((row)=>{
                  return (
                    <span>
                      <LensIcon style={{fontSize:"8px", marginRight:"5px"}}/>
                      {row.name}
                    </span>
                  );
                }) 
                : <div style={{
                  color:"#9a9a9a",
                  textAlign: "center",
                  paddingTop: "30px"}}>
                    Please Select Projects
                  </div>}
          </div>
        {/* <p>Select Role</p> */}
        <Paper>
           {this.state.rows.length > 0 ? (
              [
          <Grid rows={this.state.rows} columns={this.state.columns}>
            {/* <Toolbar /> */}
            {/* 검색 */}
            {/* <SearchState defaultValue="" />
            <SearchPanel style={{ marginLeft: 0 }} /> */}

            {/* Sorting */}
            <SortingState
              defaultSorting={[{ columnName: "status", direction: "asc" }]}
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
            <Table  cellComponent={Cell}/>
            <TableColumnResizing
              defaultColumnWidths={this.state.defaultColumnWidths}
            />
            <TableHeaderRow
              showSortingControls
              rowComponent={HeaderRow}
            />
            <TableColumnVisibility
              // defaultHiddenColumnNames={['role_id']}
            />
            <TableSelection
              selectByRowClick
              highlightRow
            />
          </Grid>]
            ) : (
              <CircularProgress
                variant="determinate"
                value={this.state.completed}
                style={{ position: "absolute", left: "50%", marginTop: "20px" }}
              ></CircularProgress>
            )}
        </Paper>
      </div>
    );
  }
}

export default GrEditGroup;
