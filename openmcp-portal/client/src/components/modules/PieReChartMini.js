import React, { Component } from "react";
import { PieChart, Pie, Sector, Cell,
  //  Legend
} from "recharts";

class PieReChartMini extends Component {
  constructor(props) {
    super(props);
    // console.log("pierechar constructor");
    const result = []
    result.push(this.props.data.status[0])
    result.push({
      name : "left",
      value : this.props.data.status[1].value -  this.props.data.status[0].value
    })
    this.state = {
      activeIndex: 0,
      rows: this.props.data.status,
      calResult : result
    };
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.data.status !== prevProps.data.status) {
      const result = []
      result.push(this.props.data.status[0])
      result.push({
        name : "left",
        value : this.props.data.status[1].value -  this.props.data.status[0].value
      })
      this.setState({
        ...this.state,
        rows: this.props.data.status,
        calResult : result
      });
    }
  }

  // onPieEnter = (data, index) => {
  //   this.setState({
  //     activeIndex: index,
  //   });
  // };

  render() {
    const COLORS = this.props.colors;
    const renderActiveShape = (props) => {
      // const RADIAN = Math.PI / 180;
      const {
        cx,
        cy,
        // midAngle,
        innerRadius,
        outerRadius,
        startAngle,
        endAngle,
        fill,
        // payload,
        percent,
        // value,
      } = props;
      // const sin = Math.sin(-RADIAN * midAngle);
      // const cos = Math.cos(-RADIAN * midAngle);
      // const sx = cx + (outerRadius + 10) * cos;
      // const sy = cy + (outerRadius + 10) * sin;
      // const mx = cx + (outerRadius + 30) * cos;
      // const my = cy + (outerRadius + 30) * sin;
      // const ex = mx + (cos >= 0 ? 1 : -1) * 22;
      // const ey = my;
      // const textAnchor = cos >= 0 ? "start" : "end";

      return (
        <g style={{fontSize:"14px"}}>
          <g>
            <text x={90} y={28} dy={3} textAnchor="left" fill={"#89caff"} style={{fontSize:"13px"}}>
              {this.props.name}
            </text>
            <text x={90} y={28} dy={20} textAnchor="left" fill={"#ffffff"}>
              {`${(percent * 100).toFixed(1)}%`}
              {/* {payload.name} */}
              {/* {`${(percent * 100).toFixed(0)}%`} */}
            </text>
          </g>
          <Sector
            cx={cx}
            cy={cy}
            innerRadius={innerRadius}
            outerRadius={outerRadius}
            startAngle={startAngle}
            endAngle={endAngle}
            fill={fill}
          />
        </g>
      );
    };
    // const style = {
    //   top: 48,
    //   left: 200,
    //   // lineHeight: "25px",
    //   fontSize:"14px",
    // };
    return (
      <div style={{ position: "relative" }} className="pie-chart">
        <PieChart width={140} height={70}>
          <Pie
            data={[this.state.rows[1]]}
            cx={30}
            cy={30}
            startAngle={this.props.angle.startAngle}
            endAngle={this.props.angle.endAngle}
            innerRadius={12}
            outerRadius={20}
            fill="#ecf0f5"
            dataKey="value"
            paddingAngle={0}
            onMouseEnter={this.onPieEnter}
            isAnimationActive={false}
          >
          </Pie>
          <Pie
            activeIndex={this.state.activeIndex}
            activeShape={renderActiveShape}
            data={this.state.calResult}
            cx={30}
            cy={30}
            startAngle={this.props.angle.startAngle}
            endAngle={this.props.angle.endAngle}
            innerRadius={12}
            outerRadius={20}
            fill="#36A1A9"
            dataKey="value"
            paddingAngle={0}
            onMouseEnter={this.onPieEnter}
          >
            {this.state.calResult.map((entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
              />
            ))}
          </Pie>
          {/* <Legend
            iconSize={10}
            width={180}
            height={140}
            layout="vertical"
            verticalAlign="middle"
            wrapperStyle={style}
            payload={this.state.rows.map((item, index) => ({
              id: item.name,
              type: "square",
              value: `${item.name} (${item.value} ${this.props.unit})`,
              color: COLORS[index % COLORS.length],
            }))}
          /> */}
          
        </PieChart>
      </div>
    );
  }
}

export default PieReChartMini;
