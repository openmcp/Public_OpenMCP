// import React, {  Component } from "react";
// import { PieChart, Pie, Sector, Cell, Legend } from "recharts";

// class PieHalfReChart extends Component {
//   constructor(props) {
//     super(props);
//     console.log("pierechar constructor");
//     this.state = {
//       activeIndex: 0,
//       rows: this.props.data.status,
//     };
//   }

//   componentDidUpdate(prevProps, prevState) {
//     if (this.props.data.status !== prevProps.data.status) {
//       this.setState({
//         ...this.state,
//         rows: this.props.data.status,
//       });
//     }
//   }

//   // onPieEnter = (data, index) => {
//   //   this.setState({
//   //     activeIndex: index,
//   //   });
//   // };

//   render() {
//     const COLORS = this.props.colors;
//     const renderActiveShape = (props) => {
//       // const RADIAN = Math.PI / 180;
//       const {
//         cx,
//         cy,
//         // midAngle,
//         innerRadius,
//         outerRadius,
//         startAngle,
//         endAngle,
//         fill,
//         payload,
//         percent,
//         // value,
//       } = props;
//       // const sin = Math.sin(-RADIAN * midAngle);
//       // const cos = Math.cos(-RADIAN * midAngle);
//       // const sx = cx + (outerRadius + 10) * cos;
//       // const sy = cy + (outerRadius + 10) * sin;
//       // const mx = cx + (outerRadius + 30) * cos;
//       // const my = cy + (outerRadius + 30) * sin;
//       // const ex = mx + (cos >= 0 ? 1 : -1) * 22;
//       // const ey = my;
//       // const textAnchor = cos >= 0 ? "start" : "end";

//       return (
//         <g style={{fontSize:"14px"}}>
//           <text x={cx} y={cy} dy={3} textAnchor="middle" fill={fill}>
//             {payload.name}
//           </text>
//           <text x={cx} y={cy} dy={20} textAnchor="middle" fill={fill}>
//             {`${(percent * 100).toFixed(0)}%`}
//           </text>

//           <Sector
//             cx={cx}
//             cy={cy}
//             innerRadius={innerRadius}
//             outerRadius={outerRadius}
//             startAngle={startAngle}
//             endAngle={endAngle}
//             fill={fill}
//           />
          
//         </g>
//       );
//     };
//     const style = {
//       top: 48,
//       left: 200,
//       lineHeight: "25px",
//       fontSize:"14px",
//     };
//     return (
//       <div style={{ position: "relative" }} className="pie-chart">
//         <PieChart width={200} height={200}>
//           <Pie
//             activeIndex={this.state.activeIndex}
//             activeShape={renderActiveShape}
//             data={this.state.rows}
//             cx={95}
//             cy={95}
//             startAngle={this.props.angle.startAngle}
//             endAngle={this.props.angle.endAngle}
//             innerRadius={50}
//             outerRadius={80}
//             fill="#367fa9"
//             dataKey="value"
//             paddingAngle={0}
//             onMouseEnter={this.onPieEnter}
//           >
//             {this.state.rows.map((entry, index) => (
//               <Cell
//                 key={`cell-${index}`}
//                 fill={COLORS[index % COLORS.length]}
//               />
//             ))}
//           </Pie>
//           <Legend
//             iconSize={10}
//             width={180}
//             height={140}
//             layout="vertical"
//             verticalAlign="middle"
//             wrapperStyle={style}
//             payload={this.state.rows.map((item, index) => ({
//               id: item.name,
//               type: "square",
//               value: `${item.name} (${item.value} ${this.props.unit})`,
//               color: COLORS[index % COLORS.length],
//             }))}
//           />
//         </PieChart>
//       </div>
//     );
//   }
// }

// export default PieHalfReChart;
