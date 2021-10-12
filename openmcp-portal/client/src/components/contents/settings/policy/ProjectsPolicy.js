import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
// import { Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
  // IntegratedSelection,
  // SelectionState,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableColumnResizing,
  TableHeaderRow,
  // TableSelection,
  PagingPanel,
  // TableColumnVisibility
} from "@devexpress/dx-react-grid-material-ui";
import * as utilLog from '../../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
import PcAddProjectPolicy from './../../modal/PcAddProjectPolicy';
import PcUpdateProjectPolicy from './../../modal/PcUpdateProjectPolicy';


class ProjectsPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "project", title: "Project"},
        { name: "cluster", title: "cluster"},
        { name: "cls_cpu_trh_r", title: "Cluster CPU"},
        { name: "cls_mem_trh_r", title: "Cluster Memory"},
        { name: "pod_cpu_trh_r", title: "Pods CPU"},
        { name: "pod_mem_trh_r", title: "Pods Memory"},
        { name: "updated_time", title: "Updated Time"},
      ],
      defaultColumnWidths: [
        { columnName: "project", width: 300 },
        { columnName: "cluster", width: 200 },
        { columnName: "cls_cpu_trh_r", width: 150 },
        { columnName: "cls_mem_trh_r", width: 150 },
        { columnName: "pod_cpu_trh_r", width: 150 },
        { columnName: "pod_mem_trh_r", width: 150 },
        { columnName: "updated_time", width: 200 },
      ],
      // defaultHiddenColumnNames :[
      //   "rate", "period", "policy_id"
      // ],
      rows: "",
      selectedRowData:"",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 10,
      pageSizes: [5, 10, 15, 0],

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
    const response = await fetch(`/settings/policy/project-policy`);
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
      const { column, row } = props;

      if (column.name === "project") {
        // // console.log("name", props.value);
        // console.log("this.props.match.params", this.props)
        return (
          <Table.Cell {...props} style={{ cursor: "pointer", color:"#3c8dbc"}}>
            <div onClick={()=>this.onClickUpdatePolicy(row)}>
              {props.value}
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
                // <PcSetOMCPPolicy rowData={this.state.selectedRow} onUpdateData={this.onUpdateData}/>,
                // <AcChangeRole rowData={this.state.selectedRow} onUpdateData={this.onUpdateData}/>,
                <PcAddProjectPolicy onUpdateData={this.onUpdateData}/>,
                <PcUpdateProjectPolicy onUpdateData={this.onUpdateData} onOpen={this.state.onClickUpdatePolicy} onCloseUpdatePolicy={this.onCloseUpdatePolicy} rowData={this.state.selectedRowData}/>,
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

export default ProjectsPolicy;