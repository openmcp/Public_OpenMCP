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
import * as utilLog from '../../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
import WarningRoundedIcon from "@material-ui/icons/WarningRounded";


class AlertLog extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "node_name", title: "Node"},
        { name: "cluster_name", title: "Cluster"},
        { name: "status", title: "Status"},
        { name: "message", title: "Message"},
        { name: "resource", title: "Resource"},
        { name: "created_time", title: "Created Time"},
      ],
      defaultColumnWidths: [
        { columnName: "node_name", width: 250 },
        { columnName: "cluster_name", width: 130 },
        { columnName: "status", width: 100 },
        { columnName: "message", width: 500 },
        { columnName: "resource", width: 120 },
        { columnName: "created_time", width: 200 },
      ],
      // defaultHiddenColumnNames :[
      //   "rate", "period", "policy_id"
      // ],
      rows: "",
      selectedRowData:"",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 30,
      pageSizes: [30, 40, 50, 0],

      completed: 0,
      onClickUpdatePolicy: false,
      // selection: [],
      // selectedRow: "",
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }

  callApi = async () => {
    const response = await fetch(`/settings/threshold/log`);
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
    this.setState({
      selection : [],
      selectedRow:"",
    })
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

  onClickUpdatePolicy = (rowData) => {
    this.setState({
      onClickUpdatePolicy: true,
      selectedRowData : rowData
    })
  }

  onCloseUpdatePolicy = (value) => {
    this.setState({onClickUpdatePolicy:value})
  }

  render() {
    const Cell = (props) => {
      const { column } = props;

      if (
        column.name === "status"
      ) {
        return (
          <Table.Cell
            {...props}
            // style={{ cursor: "pointer" }}
            aria-haspopup="true"
          >
            <div style={{ position: "relative", top: "-3px" }}>
              
              <WarningRoundedIcon
                style={{
                  fontSize: "19px",
                  marginRight: "5px",
                  position: "relative",
                  top: "5px",
                  color: props.value === "warn" ? "#efac17" : "#dc0505",
                }}
              />
              <span>{props.value}</span>
            </div>
          </Table.Cell>
        );
      }
      return <Table.Cell>{props.value}</Table.Cell>;
    };

    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
        }}
      />
    );
    const Row = (props) => {
      return <Table.Row {...props} key={props.tableRow.key}/>;
    };

    // const onSelectionChange = (selection) => {
    //   // console.log(this.state.rows[selection[0]])
    //   if (selection.length > 1) selection.splice(0, 1);
    //   this.setState({ selection: selection });
    //   this.setState({ selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {} });
    // };

    return (
      <div className="sub-content-wrapper fulled">
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
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
                    defaultSorting={[{ columnName: 'updated_time', direction: 'desc' }]}
                  />

                  {/* <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  /> */}

                  <IntegratedFiltering />
                  {/* <IntegratedSelection /> */}
                  <IntegratedSorting />
                  <IntegratedPaging />

                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                  {/* <TableColumnVisibility
                    defaultHiddenColumnNames={this.state.defaultHiddenColumnNames}
                  /> */}
                  
                  {/* <TableSelection
                    selectByRowClick
                    highlightRow
                    // showSelectionColumn={false}
                  /> */}
                  
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

export default AlertLog;