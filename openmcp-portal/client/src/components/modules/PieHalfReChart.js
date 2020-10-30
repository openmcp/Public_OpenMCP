import React, { PureComponent, Component } from "react";
import { PieChart, Pie, Sector, Cell, Legend } from "recharts";

class PieHalfReChart extends Component {
  constructor(props) {
    super(props);
    
    const result = this.calRowData(this.props.data.status[0])
    console.log("pierechar constructor", result);


    this.state = {
      activeIndex: 0,
      rows: result,
    };
  }

  calRowData = (object) => {

    Number(object.value)
    const rows = []
    const remineValue = 100 - object.value;
    console.log(remineValue);
    rows.push(object);
    rows.push({name:"none", value:remineValue});
    return rows
  }

  componentDidUpdate(prevProps, prevState) {
    // console.log("props.data.status", this.props.data.status);
    // console.log("prevProps.status", prevProps.data.status);
    if (this.props.data !== prevProps.data) {
    //   console.log("업뎃?");
      const result = this.calRowData(this.props.data.status[0])
      this.setState({
        ...this.state,
        rows: result,
      });
    }
  }
  //   componentWillMount(){
  //     console.log("componentWillUnmount",this.props.data.status);
  //     this.setState({data:this.props.data.status})
  //   }

  onPieEnter = (data, index) => {
    // this.setState({
    //   activeIndex: index,
    // });
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
    const COLORS = [
      "#0088FE",
      "#ecf0f5",
    ];
    const renderActiveShape = (props) => {
      const RADIAN = Math.PI / 180;
      const {
        cx,
        cy,
        midAngle,
        innerRadius,
        outerRadius,
        startAngle,
        endAngle,
        fill,
        payload,
        percent,
        value,
      } = props;
      const sin = Math.sin(-RADIAN * midAngle);
      const cos = Math.cos(-RADIAN * midAngle);
      const sx = cx + (outerRadius + 10) * cos;
      const sy = cy + (outerRadius + 10) * sin;
      const mx = cx + (outerRadius + 30) * cos;
      const my = cy + (outerRadius + 30) * sin;
      const ex = mx + (cos >= 0 ? 1 : -1) * 22;
      const ey = my;
      const textAnchor = cos >= 0 ? "start" : "end";

      return (
        <g>
          <text x={cx} y={cy} dy={3} textAnchor="middle" fill={fill}>
            {payload.name}
          </text>
          <text x={cx} y={cy} dy={20} textAnchor="middle" fill={fill}>
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
        </g>
      );
    };
    const style = {
      top: 29,
      left: 238,
      lineHeight: "33px",
    };
    return (
      <div style={{ position: "relative" }} className="pieChart half">
        <PieChart width={216} height={136}>
          <Pie
            activeIndex={this.state.activeIndex}
            activeShape={renderActiveShape}
            data={this.state.rows}
            cx={95}
            cy={95}
            startAngle={180}
            endAngle={0}
            innerRadius={60}
            outerRadius={100}
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
              value: `${item.name} (${item.value})`,
              color: COLORS[index % COLORS.length],
            }))}
          /> */}
        </PieChart>
      </div>
    );
  }
}

export default PieHalfReChart;
