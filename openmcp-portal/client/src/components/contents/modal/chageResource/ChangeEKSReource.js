import React, { Component } from "react";
import {
  PagingState,
  SortingState,
  SelectionState,
  IntegratedFiltering,
  IntegratedPaging,
  IntegratedSorting,
  IntegratedSelection,
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
import ProgressTemp from './../../../modules/ProgressTemp';
import Confirm2 from './../../../modules/Confirm2';
import * as utilLog from "./../../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';

// const styles = (theme) => ({
//   root: {
//     margin: 0,
//     padding: theme.spacing(2),
//   },
//   closeButton: {
//     position: "absolute",
//     right: theme.spacing(1),
//     top: theme.spacing(1),
//     color: theme.palette.grey[500],
//   },
// });

class ChangeEKSReource extends Component {
  constructor(props) {
    super(props);
    this.state = {

      columns: [
        { name: "code", title: "Type" },
        { name: "etc", title: "Resources" },
        { name: "description", title: "Description" },
      ],
      defaultColumnWidths: [
        { columnName: "code", width: 100 },
        { columnName: "etc", width: 200 },
        { columnName: "description", width: 300 },
      ],

      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 3,
      pageSizes: [3, 6, 12, 0],

      instTypes: [],
      selection: [],
      selectedRow: "",

      open: false,

      confirmOpen: false,
      confirmInfo : {
        title :"Change EKS Node Resource",
        context :"Are you sure you want to Change Node Resources?",
        button : {
          open : "",
          yes : "CONFIRM",
          no : "CANCEL",
        }
      },
      confrimTarget : "",
      confirmTargetKeyname:""
    };
    // this.onChange = this.onChange.bind(this);
  }

  componentWillMount() {
    // console.log("Migration will mount");
    this.callApi()
      .then((res) => {
        this.setState({ instTypes: res });
      })
      .catch((err) => console.log(err));
  }

  initState = () => {
    this.setState({
      selection: [],
      selectedRow: "",
    });
  };

  callApi = async () => {
    const response = await fetch("/aws/eks-type");
    const body = await response.json();
    return body;
  };

  // onChange = (e) => {
  //   this.setState({
  //     [e.target.name]: e.target.value,
  //   });
  // };

  handleClose = () => {
    this.initState();
    this.setState({
      open: false,
    });
  };

  handleSaveClick = (e) => {
    if (Object.keys(this.state.selectedRow).length  === 0) {
      alert("Please select Instance Type");
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
      const url = `/nodes/eks/change`;
      const data = {
        cluster : "eks-cluster1",
        region : "ap-northeast-2",
        node : "ip-172-31-0-123.ap-northeast-2.compute.internal",
        // cluster : this.props.nodeData.cluster,
        type :  this.state.selectedRow.code,
        // region : this.props.propsRow.region,
        // node : this.props.nodeData.name,
      };

      axios.post(url, data)
        .then((res) => {
          if(res.data.error){
            alert(res.data.message)
            return
          }
        })
        .catch((err) => {
        });
        
        this.props.handleClose()
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

  onSelectionChange = (selection) => {
    if (selection.length > 1) selection.splice(0, 1);
    this.setState({ selection: selection });
    this.setState({
      selectedRow: this.state.instTypes[selection[0]]
        ? this.state.instTypes[selection[0]]
        : {},
    });
  };

  render() {
    return (
      <div>
        {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""}

        <Confirm2
          confirmInfo={this.state.confirmInfo} 
          confrimTarget ={this.state.confrimTarget} 
          confirmTargetKeyname = {this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}/>
          
        <div className="md-contents-body">
          <section className="md-content">
            {/* deployment informations */}
            <p>EKS Node Info</p>
            <div id="md-content-info">
              <div class="md-partition">
                <div class="md-item">
                  <span><strong>Name : </strong></span>
                  <span>{this.props.nodeData.name}</span>
                </div>
                <div class="md-item">
                  <span><strong>OS : </strong></span>
                  <span>{this.props.nodeData.os}</span>
                </div>
              </div>
              <div class="md-partition">
                <div class="md-item">
                  <span><strong>Status : </strong></span>
                  <span>{this.props.nodeData.status}</span>
                </div>
                <div class="md-item">
                  <span><strong>IP : </strong></span>
                  <span>{this.props.nodeData.ip}</span>
                </div>
              </div>
            </div>
          </section>
          <section className="md-content">
            <div>
              <p>Instance Type</p>
              {/* cluster selector */}
              <Paper>
                <Grid
                  rows={this.state.instTypes}
                  columns={this.state.columns}
                >
                  {/* Sorting */}
                  <SortingState
                    defaultSorting={[
                      { columnName: "code", direction: "asc" },
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
                    onSelectionChange={this.onSelectionChange}
                  />

                  <IntegratedFiltering />
                  <IntegratedSorting />
                  <IntegratedSelection />
                  <IntegratedPaging />

                  {/* 테이블 */}
                  <Table />
                  <TableColumnResizing
                    defaultColumnWidths={
                      this.state.defaultColumnWidths
                    }
                  />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={this.HeaderRow}
                  />
                  <TableSelection
                    selectByRowClick
                    highlightRow
                    // showSelectionColumn={false}
                  />
                </Grid>
              </Paper>
            </div>
          </section>
        </div>
      </div>
    );
  }
}

export default ChangeEKSReource;
