import React, { Component } from "react";
// import {  Button,} from "@material-ui/core";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  SelectionState,
  IntegratedSelection,
  // TableColumnVisibility,
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
import Paper from "@material-ui/core/Paper";
import EditKVMAuth from "../../../modal/public-cloud-auth/EditKVMAuth.js";
import axios from 'axios';
import Confirm2 from './../../../../modules/Confirm2';
import IconButton from "@material-ui/core/IconButton";
import MenuItem from "@material-ui/core/MenuItem";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';

class ConfigKVM extends Component {
  constructor(props) {
    super(props);

    this.state = {
      columns: [
        { name: "seq", title:"No"},
        { name: "cluster", title: "Cluster" },
        { name: "agentURL", title: "Agent URL" },
        { name: "mClusterName", title: "MCluster VM Name" },
        { name: "mClusterPwd", title: "MCluster Passwd" },
      ],
      defaultColumnWidths: [
        { columnName: "seq", width: 100 },
        { columnName: "cluster", width: 200 },
        { columnName: "agentURL", width: 300 },
        { columnName: "mClusterName", width: 300 },
        { columnName: "mClusterPwd", width: 300 },
      ],
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      pageSizes: [5, 10, 15, 0],

      open: false,
      account: "",
      account_role: "",
      rows: [],

      selection: [],
      selectedRow: "",
      popTitle:"",

      confirmOpen: false,
      confirmInfo : {
        title :"Delete KVM PCA Info",
        context :"Are you sure you want to delete KVM PCA config?",
        button : {
          open : "",
          yes : "CONFIRM",
          no : "CANCEL",
        }
      },
      confrimTarget : "",
      confirmTargetKeyname:"",
      anchorEl: null,
    };
  }

  callApi = async () => {
    const response = await fetch(`/settings/config/pca/kvm`);
    const body = await response.json();
    return body;
  };

  componentWillMount() {

    this.callApi()
      .then((res) => {
        this.setState({ rows: res });
      })
      .catch((err) => console.log(err));
  }

  Cell = (props) => {
    const { column } = props;
    if (column.name === "control") {
      return (
        <Table.Cell
          {...props}
          style={{
            borderRight: "1px solid #e0e0e0",
            borderLeft: "1px solid #e0e0e0",
            textAlign: "center",
            background: "whitesmoke",
          }}
        >
        </Table.Cell>
      );
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

  handleClickNew = () => {
    this.setState({ 
      open : true,
      new : true,
      popTitle:"Add KVM Authentication",
      data:{}
    });
  };

  handleClickEdit = () => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select a authentication data row");
      this.setState({ open: false });
      return;
    }

    this.setState({ 
      open: true, 
      new: false, 
      popTitle:"Edit KVM Authentication",
      data:{
        seq : this.state.selectedRow.seq,
        cluster: this.state.selectedRow.cluster,
        agentURL: this.state.selectedRow.agentURL,
        mClusterName: this.state.selectedRow.mClusterName,
        mClusterPwd: this.state.selectedRow.mClusterPwd,
      }
    });
  };

  handleClickDelete = () => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select a authentication data row");
      this.setState({ open: false });
      return;
    } else {
      this.setState({
        confirmOpen: true,
      })
    }
  }

  //callback
  confirmed = (result) => {
    this.setState({confirmOpen:false})

    if(result) {
      const data = {
        seq : this.state.selectedRow.seq,
        cluster : this.state.selectedRow.cluster
      };
  
      const url = `/settings/config/pca/kvm`;
      axios.delete(url, {data:data})
      .then((res) => {
        this.callBackClosed();
      })
      .catch((err) => {
        console.log("Error : ",err);
      });
    } else {
      this.setState({confirmOpen:false})
    }
  }

  callBackClosed = () => {
    this.setState({
      open : false,
      selection: [],
      selectedRow: "",});
    this.callApi()
    .then((res) => {
      this.setState({ rows: res });
    })
    .catch((err) => console.log(err));
  }

  render() {
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
      return <Table.Row {...props} key={props.tableRow.key} />;
    };
    const onSelectionChange = (selection) => {
      // console.log(selection);
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ selection: selection });
      this.setState({
        selectedRow: this.state.rows[selection[0]]
          ? this.state.rows[selection[0]]
          : {},
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
      this.setState({ anchorEl: null });
    };

    const open = Boolean(this.state.anchorEl);

    return (
      <div>

        <Confirm2
          confirmInfo={this.state.confirmInfo} 
          confrimTarget ={this.state.confrimTarget} 
          confirmTargetKeyname = {this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}/>
          
        <EditKVMAuth 
          open={this.state.open}
          new={this.state.new}
          callBackClosed={this.callBackClosed}
          title={this.state.popTitle}
          data={this.state.data}
        />
        
        <div className="md-contents-body">
          <div style={{padding: "8px 15px",
                        fontSize:"13px",
                        backgroundColor: "#bfdcec",
                        boxShadow: "0px 0px 3px 0px #b9b9b9"
                      }}
          > KVM Authentications Configration</div>
          <section className="md-content">
            <Paper>
            <div
                style={{
                  position: "absolute",
                  right: "21px",
                  top: "212px",
                  zIndex: "10",
                  textTransform: "capitalize",
                }}
              >
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
                            <MenuItem
                              onClick={handleClose}
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                            >
                              <div
                                onClick={this.handleClickNew}
                                style={{ width: "148px", textTransform: "capitalize", }}
                              >
                                New </div>
                            </MenuItem>
                            <MenuItem
                              onClick={handleClose}
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                            >
                              <div
                                onClick={this.handleClickEdit}
                                style={{ width: "148px", textTransform: "capitalize", }}
                              >
                                Edit
                              </div>
                            </MenuItem>
                            <MenuItem
                              onClick={handleClose}
                              style={{ textAlign: "center", display: "block", fontSize: "14px"}}
                            >
                              <div
                                onClick={this.handleClickDelete}
                                style={{ width: "148px", textTransform: "capitalize", }}
                              >
                                Delete
                              </div>
                            </MenuItem>
                            </MenuList>
                          </Paper>
                      </Grow>
                    )}
                  </Popper>
              </div>
              <Grid 
                rows={this.state.rows} 
                columns={this.state.columns}
              >
                {/* <div style={{position:"relative"}}>
                  <div style = {{position:"absolute",
                        right: "13px",
                        top: "13px",
                        zIndex: "10",}}>
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={this.handleClickNew}
                      style={{
                        width: "120px",
                        marginRight:"10px",
                        textTransform: "capitalize",
                      }}
                    >
                      New
                    </Button>
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={this.handleClickEdit}
                      style={{
                        width: "120px",
                        marginRight:"10px",
                        textTransform: "capitalize",
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={this.handleClickDelete}
                      style={{
                        width: "120px",
                        textTransform: "capitalize",
                      }}
                    >
                      Delete
                    </Button>
                  </div>
                </div> */}
                <Toolbar />
                {/* 검색 */}
                <SearchState defaultValue="" />
                <SearchPanel style={{ marginLeft: 0 }} />

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
                <Table
                  cellComponent={this.Cell}
                  rowComponent={Row}
                />
                <TableColumnResizing
                    defaultColumnWidths={this.state.defaultColumnWidths}
                />
                <TableHeaderRow showSortingControls rowComponent={HeaderRow} />
                {/* <TableColumnVisibility defaultHiddenColumnNames={["role_id"]} /> */}
                <TableSelection
                  selectByRowClick
                  highlightRow
                  // showSelectionColumn={false}
                />
              </Grid>
            </Paper>
          </section>
        </div>
      </div>
    );
  }
}

export default ConfigKVM;
