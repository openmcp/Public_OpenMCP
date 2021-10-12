import React, { Component } from "react";
import { PieChart, Pie, Sector, Cell, Legend, ResponsiveContainer  } from "recharts";

class PieReChart extends Component {
  constructor(props) {
    super(props);
    // console.log("pierechar constructor");
    this.state = {
      activeIndex: 0,
      rows: this.props.data.status,
    };
  }

  componentDidUpdate(prevProps, prevState) {
    if (this.props.data.status !== prevProps.data.status) {
      this.setState({
        ...this.state,
        rows: this.props.data.status,
      });
    }
  }

  onPieEnter = (data, index) => {
    this.setState({
      activeIndex: index,
    });
  };

  render() {
    // console.log("render", this.props.data.status);

    // const {PieChart, Pie, Sector} = Recharts;
    // const data = [
    //   { name: "healthy", value: 400 },
    //   { name: "partially_healthy", value: 100 },
    //   { name: "converging", value: 10 },
    //   { name: "unhealthy", value: 20 },
    // ];
    // const COLORS = [
    //   "#0088FE",
    //   "#00C49F",
    //   "#FFBB28",
    //   "#FF8042",
    //   "#00C49F",
    //   "#FFBB28",
    //   "#00C49F",
    //   "#FFBB28",
    // ];
    // console.log("color", this.props.color);
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
        payload,
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
          <text x={cx} y={cy} dy={0} textAnchor="middle" fill={fill}>
            {payload.name}
          </text>
          <text x={cx} y={cy} dy={15} textAnchor="middle" fill={fill}>
            {`${(percent * 100).toFixed(0)}%`}
          </text>

          <Sector
            cx={cx}
            cy={cy}
            innerRadius={innerRadius}
            outerRadius={outerRadius}
            startAngle={startAngle}
            endAngle={endAngle}
            fill={fill}
          />
          {/* <Sector
              cx={cx}
              cy={cy}
              startAngle={startAngle}
              endAngle={endAngle}
              innerRadius={outerRadius + 6}
              outerRadius={outerRadius + 10}
              fill={fill}
            />
            <path d={`M${sx},${sy}L${mx},${my}L${ex},${ey}`} stroke={fill} fill="none" />
            <circle cx={ex} cy={ey} r={2} fill={fill} stroke="none" />
            <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} textAnchor={textAnchor} fill="#333">{`PV ${value}`}</text>
            <text x={ex + (cos >= 0 ? 1 : -1) * 12} y={ey} dy={18} textAnchor={textAnchor} fill="#999">
              {`(Rate ${(percent * 100).toFixed(2)}%)`}
            </text> */}
        </g>
      );
    };
    const style = {
      // top: 48,
      // left: 200,
      // position: "relative",
      right:"8px",
      lineHeight: "25px",
      fontSize:"0.9vw",
    };
    return (
      <div style={{ position: "relative", height: "150px"}} className="pie-chart">
        <ResponsiveContainer  width="100%" height="100%">
          <PieChart>
            <Pie
              activeIndex={this.state.activeIndex}
              activeShape={renderActiveShape}
              data={this.state.rows}
              cx={70}
              cy={70}
              startAngle={this.props.angle.startAngle}
              endAngle={this.props.angle.endAngle}
              innerRadius={35}
              outerRadius={60}
              fill="#367fa9"
              dataKey="value"
              paddingAngle={0}
              onMouseEnter={this.onPieEnter}
            >
              {this.state.rows.map((entry, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <Legend
              iconSize={10}
              // width={180}
              // height={140}
              align= "right"
              layout="vertical"
              verticalAlign="middle"
              wrapperStyle={style}
              payload={this.state.rows.map((item, index) => ({
                id: item.name,
                type: "square",
                value: `${item.name} (${item.value}${this.props.unit})`,
                color: COLORS[index % COLORS.length],
              }))}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>
    );
  }
}

export default PieReChart;
