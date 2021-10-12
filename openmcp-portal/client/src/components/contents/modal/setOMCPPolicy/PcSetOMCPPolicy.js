import React, { Component } from "react";
import PcSetNumericPolicy from "./PcSetNumericPolicy.js";
import PcSetTextValuePolicy from "./PcSetTextValuePolicy.js";
// import PcSetAnalyticMetricsWeight from "./PcSetNumericPolicy.js";
// import PcSetHpaMinMaxDistributionMode from "./PcSetHpaMinMaxDistributionMode.js";
// import PcSetLoadBalancingControllerPolicy from "./PcSetLoadBalancingControllerPolicy.js";
// import PcSetLogLevel from "./PcSetLogLevel.js";
// import PcSetMetricCollectorPeriod from "./PcSetMetricCollectorPeriod.js";

// function valuetext(value) {
//   return `${value}°C`;
// }

class PcSetOMCPPolicy extends Component {
  constructor(props) {
    super(props);
    this.state = {
      policyName : this.props.policyName
    }
  }

  componentWillMount() {
  }
  
  componentDidUpdate(prevProps, prevState) {
    
    if (this.props.policyName !== prevProps.policyName) {
      this.setState({
        ...this.state,
        policyName: this.props.policyName,
      });
    }
  }

  render() {
    const PolicyDialog = () => {
      switch(this.state.policyName){
        case "metric-collector-period":
          return <PcSetNumericPolicy isFloat= {false} policyName={this.state.policyName} policy={this.props.policy} onUpdateData={this.props.onUpdateData}/>
        case "log-level":
          return <PcSetNumericPolicy isFloat= {false} policyName={this.state.policyName} policy={this.props.policy} onUpdateData={this.props.onUpdateData}/>
        case "loadbalancing-controller-policy":
          return <PcSetNumericPolicy isFloat= {true} policyName={this.state.policyName} policy={this.props.policy} onUpdateData={this.props.onUpdateData}/>
        case "analytic-metrics-weight":
          return <PcSetNumericPolicy isFloat= {true} policyName={this.state.policyName} policy={this.props.policy} onUpdateData={this.props.onUpdateData}/>
        default : //hpa-minmax-distribution-mode
         return <PcSetTextValuePolicy isFloat= {false} policyName={this.state.policyName} policy={this.props.policy} onUpdateData={this.props.onUpdateData}/>

      }
    }
    return (
      <div>
        <PolicyDialog/>
      </div>
    );
  }
}

export default PcSetOMCPPolicy;
