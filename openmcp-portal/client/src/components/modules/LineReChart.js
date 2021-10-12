import React, { Component } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";

class LineReChart extends Component {
  state = {
    rows: this.props.rowData,
  };

  componentDidUpdate(prevProps, prevState) {
    console.log("componentDidUpdate")
    if (this.props.rowData !== prevProps.rowData) {
      this.setState({
        ...this.state,
        rows: this.props.rowData,
      });
    }
  }
  render() {
    // console.log(this.state.rows, this.props.name);
    const color = ["#367fa9", "#3cb0bc"];
    const avg = [];
    const key = this.props.name;
    if (this.props.cardinal) {
      this.props.name.forEach((i) => {
        let sum = this.state.rows.reduce(function (prev, current) {
          return prev + +current[i];
        }, 0);
        avg.push((sum / this.state.rows.length).toFixed(2));
      });
    } else {
      let sum = this.state.rows.reduce(function (prev, current) {
        return prev + +current[key];
      }, 0);
      avg.push((sum / this.state.rows.length).toFixed(2));
    }
    return (
      // color, name, unit
      <div>
        <h4>
            {this.props.title}
            <span style={{fontSize:"13px", fontWeight:"normal", paddingLeft:"5px"}}>
                {this.props.cardinal ? 
                this.props.name.map((i, index) => {
                    return <span>{index === this.props.name.length-1 ? "("+avg[index] +" "+ this.props.unit +" "+ i +")" : "("+avg[index] +" "+ this.props.unit +" "+ i+") /"}</span>;
                }) : 
                <span>{"("+avg[0] + this.props.unit+")"}</span>
                }
            </span>
        </h4>
        <AreaChart
          width={800}
          height={200}
          data={this.state.rows}
          margin={{
            top: 24,
            right: 30,
            left: 30,
            bottom: 0,
          }}
          style={{ fontSize: "12px" }}
          stackOffset="expand"
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip
            // labelStyle={{color:"#ffffff"}}
            // itemStyle = {{color:"#ffffff"}}
            contentStyle={{
              fontSize: "14px",
              color: "#ffffff",
              background: "#222d32",
              borderRadius: "4px",
            }}
            // wrapperStyle={{fontSize: "14px"}}
          />
          {this.props.cardinal
            ? this.state.rows.map((i, index) => {
                // console.log(
                //   this.props.name[index],
                //   color[index % this.props.name.length]
                // );
                return [
                  <defs>
                    <linearGradient
                      id={"colorNetwork" + index}
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop
                        offset="5%"
                        stopColor={color[index % this.props.name.length]}
                        stopOpacity={0.8}
                      />
                      <stop
                        offset="95%"
                        stopColor={color[index % this.props.name.length]}
                        stopOpacity={0}
                      />
                    </linearGradient>
                  </defs>,
                  <Area
                    type="linear" //basis, linear, natural, monotone, step
                    dataKey={this.props.name[index]}
                    stroke={color[index % this.props.name.length]}
                    fill={"url(#colorNetwork" + index + ")"}
                    unit={this.props.unit}
                    name={this.props.name[index]}
                  />,
                ];
              })
            : [
                <defs>
                  <linearGradient
                    id="colorCpuMemory"
                    x1="0"
                    y1="0"
                    x2="0"
                    y2="1"
                  >
                    <stop offset="5%" stopColor="#367fa9" stopOpacity={0.8} />
                    <stop offset="95%" stopColor="#367fa9" stopOpacity={0} />
                  </linearGradient>
                </defs>,
                <Area
                  type="linear" //basis, linear, natural, monotone, step
                  dataKey={this.props.name}
                  stroke="#367fa9"
                  fill="url(#colorCpuMemory)"
                  unit={this.props.unit}
                  name={this.props.name}
                />,
              ]}
        </AreaChart>
      </div>
    );
  }
}

export default LineReChart;
