import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
// import { NavLink } from "react-router-dom";
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
  TableColumnVisibility
} from "@devexpress/dx-react-grid-material-ui";
// import { NavigateNext} from '@material-ui/icons';
import * as utilLog from '../../../util/UtLogs.js';
import { AsyncStorage } from 'AsyncStorage';
// import AddMembers from "./AddMembers";
// import Editor from "../../modules/Editor";
import PcSetOMCPPolicy from '../../modal/setOMCPPolicy/PcSetOMCPPolicy';;

class OpenMCPPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "name", title: "Policy"},
        { name: "status", title: "Status"},
        { name: "value", title: "Policy Value"},
        
      ],
      defaultColumnWidths: [
        { columnName: "name", width: 300 },
        { columnName: "status", width: 100 },
        { columnName: "value", width: 300 },
      ],
      defaultHiddenColumnNames :[
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
      selectedPolicyName : "analytic-metrics-weight"
    };
  }

  componentWillMount() {
    // this.props.menuData("none");
  }

  callApi = async () => {
    const response = await fetch(`/settings/policy/openmcp-policy`);
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
        let result = [];
        if(res == null){
          this.setState({ rows: [] });
        } else {
          res.forEach((item) => {
            if(item.name != "has-target-cluster"){
              console.log(item.name);
              result.push(item);
            }
          });
          this.setState({ rows: result });
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
        let result = [];
        if(res == null){
          this.setState({ rows: [] });
        } else {
          res.forEach((item) => {
            if(item.name != "has-target-cluster"){
              console.log(item.name);
              result.push(item);
            }
          });
          this.setState({ rows: result });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  };

  render() {

    const Cell = (props) => {
      // const { column, row } = props;
      const { column } = props;

      const fnEnterCheck = () => {
        if(props.value === undefined){
          return ""
        } else {
          return (
            props.value.indexOf("|") > 0 ? 
              props.value.split("|").map( item => {
                return (
                  <p>{item}</p>
              )}) : 
                props.value
          )
        }
      }

      // const arrayToString = () => {
      //   const stringData = props.value.reduce((result, item, index, arr) => {
      //     if (index+1 === arr.length){
      //       return `${result}${item}`
      //     } else {
      //       return `${result}${item}, `
      //     }
      //   }, "")

      //   return stringData
      // }

      if (column.name === "value") {
        return <Table.Cell>{fnEnterCheck()}</Table.Cell>;
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

    const onSelectionChange = (selection) => {
      // console.log(this.state.rows[selection[0]])
      if (selection.length > 1) selection.splice(0, 1);
      this.setState({ 
        selection: selection, 
        selectedPolicyName : this.state.rows[selection[0]] ? this.state.rows[selection[0]].name : ""});
      this.setState({ selectedRow: this.state.rows[selection[0]] ? this.state.rows[selection[0]] : {} });
    };

    return (
      <div className="sub-content-wrapper fulled">
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                <PcSetOMCPPolicy policy={this.state.selectedRow} onUpdateData={this.onUpdateData} policyName={this.state.selectedPolicyName}/>, 
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
                    defaultSorting={[{ columnName: 'policy_id', direction: 'asc' }]}
                  />

                  <SelectionState
                    selection={this.state.selection}
                    onSelectionChange={onSelectionChange}
                  />

                  <IntegratedFiltering />
                  <IntegratedSelection />
                  <IntegratedSorting />
                  <IntegratedPaging />

                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableColumnResizing defaultColumnWidths={this.state.defaultColumnWidths} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                  <TableColumnVisibility
                    defaultHiddenColumnNames={this.state.defaultHiddenColumnNames}
                  />
                  
                  <TableSelection
                    // selectByRowClick
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

export default OpenMCPPolicy;