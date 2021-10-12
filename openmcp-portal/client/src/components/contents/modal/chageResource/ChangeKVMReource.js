import React, { Component } from "react";
// import {
//   PagingState,
//   SortingState,
//   SelectionState,
//   IntegratedFiltering,
//   IntegratedPaging,
//   IntegratedSorting,
//   IntegratedSelection,
// } from "@devexpress/dx-react-grid";
import {
  // Grid,
  Table,
  // TableColumnResizing,
  // TableHeaderRow,
  // PagingPanel,
  // TableSelection,
} from "@devexpress/dx-react-grid-material-ui";
// import Paper from "@material-ui/core/Paper";
import axios from 'axios';
import ProgressTemp from './../../../modules/ProgressTemp';
import Confirm2 from './../../../modules/Confirm2';
// import SelectBox from "../../../modules/SelectBox";
// import InputLabel from '@material-ui/core/InputLabel';
import Slider from "@material-ui/core/Slider";
import Typography from "@material-ui/core/Typography";
// import { withStyles, makeStyles } from '@material-ui/core/styles';
import * as utilLog from "./../../../util/UtLogs.js";
import { AsyncStorage } from 'AsyncStorage';

// const styles = (theme) => ({
//   root: {
//     color: '#52af77',
//     height: 8,
//   },
//   thumb: {
//     height: 24,
//     width: 24,
//     backgroundColor: '#fff',
//     border: '2px solid currentColor',
//     marginTop: -8,
//     marginLeft: -12,
//     '&:focus, &:hover, &$active': {
//       boxShadow: 'inherit',
//     },
//   },
//   active: {},
//   valueLabel: {
//     left: 'calc(-50% + 4px)',
//   },
//   track: {
//     height: 8,
//     borderRadius: 4,
//   },
//   rail: {
//     height: 8,
//     borderRadius: 4,
//   },
// });

class ChangeKVMReource extends Component {
  constructor(props) {
    super(props);
    this.state = {

      cpu : "",
      memory : "",

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
        title :"Change KVM Node Resource",
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
    // console.log(typeof this.props.rows.cpu.status[1].value)
    // console.log(this.props.rows.cpu.status[1].value)
    // console.log(this.props.rows.memory.status[1].value)
    this.setState({
      cpu: parseInt(this.props.rows.cpu.status[1].value),
      memory: parseInt(this.props.rows.memory.status[1].value),
    });

    
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

  onChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value,
    });
  };

  handleClose = () => {
    this.initState();
    this.setState({
      open: false,
    });
  };

  handleSaveClick = (e) => {
    // if (Object.keys(this.state.selectedRow).length  === 0) {
    //   alert("Please select Instance Type");
    //   return;
    // } else {
      this.setState({
        confirmOpen: true,
      })
    // }
  }

  confirmed = (result) => {
    this.setState({confirmOpen:false});

    //show progress loading...
    this.setState({openProgress:true});

    if(result) {
      const url = `/nodes/kvm/change`;
      const data = {
        cluster : this.props.nodeData.cluster,
        node : this.props.nodeData.name,
        cpu : this.state.cpu.toString(),
        memory : this.state.memory.toString(),
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

  // onSelectionChange = (selection) => {
  //   if (selection.length > 1) selection.splice(0, 1);
  //   this.setState({ selection: selection });
  //   this.setState({
  //     selectedRow: this.state.instTypes[selection[0]]
  //       ? this.state.instTypes[selection[0]]
  //       : {},
  //   });
  // };

  handleChangeCpu = (e, newValue) => {
    // console.log("handleChangeCpu",newValue)
    this.setState({
      cpu: newValue,
    });
  };

  handleChangeMemory = (e, newValue) => {
    // console.log("handleChangeMemory",newValue)
    this.setState({
      memory: newValue,
    });
  };


  render() {
    // const onSelectCpuChange = (data) => {
    //   switch(data){
    //     case "cpu":

    //       break;
    //     case "memory":
    //       break;
    //     default:
    //   }
    // }

    // const onSelectMemoryChange = (data) => {
    //   switch(data){
    //     case "cpu":

    //       break;
    //     case "memory":
    //       break;
    //     default:
    //   }
    // };

    // const selectBoxCpu = [
    //   {name:"1", value:"1"},
    //   {name:"2", value:"2"}
    // ];
    // const selectBoxMemory = [
    //   {name:"1", value:"1"},
    //   {name:"2", value:"2"},
    //   {name:"4", value:"4"}
    // ];

    const cpu_marks = [
      {
        value: 1,
        label: "1",
      },
      {
        value: 2,
        label: "2",
      },
      {
        value: 4,
        label: "4",
      },
    ];
  
    const memory_marks = [
      {
        value: 1,
        label: "1",
      },
      {
        value: 2,
        label: "2",
      },
      {
        value: 4,
        label: "4",
      },
      {
        value: 8,
        label: "8",
      },
    ];

    return (
      <div>
        {this.state.openProgress ? <ProgressTemp openProgress={this.state.openProgress} closeProgress={this.closeProgress}/> : ""}

        <Confirm2
          confirmInfo={this.state.confirmInfo} 
          confrimTarget ={this.state.confrimTarget} 
          confirmTargetKeyname = {this.state.confirmTargetKeyname}
          confirmed={this.confirmed}
          confirmOpen={this.state.confirmOpen}/>
          
        <div className="md-contents-body" style={{minWidth:"700px"}}>
          <section className="md-content">
            {/* deployment informations */}
            <p>KVM Node Info</p>
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
            <p>Resources</p>
            <div className="pd-resource-config">
              <div className="res">
                <Typography id="range-slider" gutterBottom>
                  CPU
                </Typography>
                <Slider
                  className="sl"
                  name="cpu"
                  // value={this.state.cpu}
                  onChange={this.handleChangeCpu}
                  valueLabelDisplay="auto"
                  aria-labelledby="pretto-slider"
                  // getAriaValueText={this.handleChangeCpu}
                  defaultValue={this.state.cpu}
                  step={null}
                  min={1}
                  max={4}
                  marks={cpu_marks}
                />
              </div>
              <div className="res">
                <Typography id="range-slider" gutterBottom>
                  Memory
                </Typography>
                <Slider
                  className="sl"
                  name="memory"
                  // value={this.state.memory}
                  onChange={this.handleChangeMemory}
                  valueLabelDisplay="auto"
                  aria-labelledby="pretto-slider"
                  // getAriaValueText={this.handleChangeMemory}
                  step={null}
                  defaultValue={this.state.memory}
                  min={1}
                  max={8}
                  marks={memory_marks}
                />
              </div>
            </div>
          </section>
        </div>
      </div>
    );
  }
}

export default ChangeKVMReource;
