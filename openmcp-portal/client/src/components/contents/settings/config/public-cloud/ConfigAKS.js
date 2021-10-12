import React, { Component } from "react";
// import { Button } from "@material-ui/core";
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
  // TableFixedColumns,
} from "@devexpress/dx-react-grid-material-ui";
import Paper from "@material-ui/core/Paper";
import EditAKSAuth from "./../../../modal/public-cloud-auth/EditAKSAuth";
import axios from "axios";
import Confirm2 from "./../../../../modules/Confirm2";
import IconButton from "@material-ui/core/IconButton";
// import Menu from "@material-ui/core/Menu";
import MenuItem from "@material-ui/core/MenuItem";
import MoreVertIcon from "@material-ui/icons/MoreVert";
import Popper from '@material-ui/core/Popper';
import MenuList from '@material-ui/core/MenuList';
import Grow from '@material-ui/core/Grow';
//import ClickAwayListener from '@material-ui/core/ClickAwayListener';

class ConfigAKS extends Component {
  constructor(props) {
    super(props);

    this.state = {
      // tb_auth_eks > seq,cluster,type,accessKey,secretKey
      columns: [
        { name: "seq", title: "No" },
        { name: "cluster", title: "Cluster" },
        { name: "clientId", title: "Client ID" },
        { name: "clientSec", title: "Client Sec" },
        { name: "tenantId", title: "Tenant ID" },
        { name: "subId", title: "Sub ID" },
      ],
      defaultColumnWidths: [
        { columnName: "seq", width: 70 },
        { columnName: "cluster", width: 150 },
        { columnName: "clientId", width: 300 },
        { columnName: "clientSec", width: 300 },
        { columnName: "tenantId", width: 300 },
        { columnName: "subId", width: 300 },
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
      popTitle: "",

      confirmOpen: false,
      confirmInfo: {
        title: "Delete AKS PCA Info",
        context: "Are you sure you want to delete AKS PCA config?",
        button: {
          open: "",
          yes: "CONFIRM",
          no: "CANCEL",
        },
      },
      confrimTarget: "",
      confirmTargetKeyname: "",
      anchorEl: null,
    };
  }

  callApi = async () => {
    const response = await fetch(`/settings/config/pca/aks`);
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
    if (column.name === "privateKey") {
      return (
        <Table.Cell
          {...props}
          style={{
            cursor: "pointer",
            wordBreak: "break-all",
            whiteSpace: "inherit",
            lineHeight: "1.5",
          }}
        ></Table.Cell>
      );
    }
    return <Table.Cell>{props.value}</Table.Cell>;
  };

  handleClickNew = () => {
    this.setState({
      open: true,
      new: true,
      popTitle: "Add AKS Authentication",
      data: {},
    });
  };

  handleClickEdit = () => {
    if (Object.keys(this.state.selectedRow).length === 0) {
      alert("Please select a authentication data row");
      this.setState({ open: false });
      return;
    }

    this.setState({
      open: true,
      new: false,
      popTitle: "Edit AKS Authentication",
      data: {
        seq: this.state.selectedRow.seq,
        cluster: this.state.selectedRow.cluster,
        clientId: this.state.selectedRow.clientId,
        clientSec: this.state.selectedRow.clientSec,
        tenantId: this.state.selectedRow.tenantId,
        subId: this.state.selectedRow.subId,
      },
    });
  };

  handleClickDelete = () => {
    if (Object.keys(this.state.selectedRow).length === 0) {
      alert("Please select a authentication data row");
      this.setState({ open: false });
      return;
    } else {
      this.setState({
        confirmOpen: true,
      });
    }
  };

  //callback
  confirmed = (result) => {
    this.setState({ confirmOpen: false });

    if (result) {
      const data = {
        seq: this.state.selectedRow.seq,
        cluster: this.state.selectedRow.cluster,
      };

      const url = `/settings/config/pca/aks`;
      axios
        .delete(url, { data: data })
        .then((res) => {
          this.callBackClosed();
        })
        .catch((err) => {
          console.log("Error : ", err);
        });
    } else {
      this.setState({ confirmOpen: false });
    }
  };

  callBackClosed = () => {
    this.setState({
      open: false,
      selection: [],
      selectedRow: "",
    });
    this.callApi()
      .then((res) => {
        this.setState({ rows: res });
      })
      .catch((err) => console.log(err));
  };

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
          confrimTarget={this.state.confrimTarget}
          confirmTargetKeyname={this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}
        />

        <EditAKSAuth
          open={this.state.open}
          new={this.state.new}
          callBackClosed={this.callBackClosed}
          title={this.state.popTitle}
          data={this.state.data}
        />

        <div className="md-contents-body">
          <div
            style={{
              padding:  "8px 15px",
              fontSize:"13px",
              backgroundColor: "#bfdcec",
              boxShadow: "0px 0px 3px 0px #b9b9b9",
            }}
          >
            {" "}
            AKS Authentications Configration
          </div>
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
              <Grid rows={this.state.rows} columns={this.state.columns}>
                {/* <div style={{ position: "relative" }}>
                  <div
                    style={{
                      position: "absolute",
                      right: "13px",
                      top: "13px",
                      zIndex: "10",
                    }}
                  >
                    <Button
                      variant="outlined"
                      color="primary"
                      onClick={this.handleClickNew}
                      style={{
                        width: "120px",
                        marginRight: "10px",
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
                        marginRight: "10px",
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
                <Table cellComponent={this.Cell} rowComponent={Row} />
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

export default ConfigAKS;
