import { ResponsiveLine } from "@nivo/line";
import React, { Component } from "react";
// make sure parent container have a defined height when using
// responsive component, otherwise height will be 0 and
// no chart will be rendered.
// website examples showcase many properties,
// you'll often use just a few of them.

class LineChart extends Component {
  constructor(props) {
    super(props);
    this.state = {
      data: this.props.data,
    };
  }
  render() {
    //this.props.data를 통해서 데이터를 받아온다.
    // console.log("MyResponsiveLine", this.props.data)
    return (
      <ResponsiveLine
        data={this.state.data}
        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
        xScale={{ type: "point" }}
        yScale={{
          type: "linear",
          min: "auto",
          max: "auto",
          stacked: false,
          reverse: false,
        }}
        yFormat=" >-.2f"
        axisTop={null}
        axisRight={null}
        axisBottom={{
          orient: "bottom",
          tickSize: 5,
          tickPadding: 5,
          tickRotation: 0,
          legend: "transportation",
          legendOffset: 36,
          legendPosition: "middle",
        }}
        axisLeft={{
          orient: "left",
          tickSize: 5,
          tickPadding: 5,
          tickRotation: 0,
          legend: "count",
          legendOffset: -40,
          legendPosition: "middle",
        }}
        enablePoints={false}
        pointSize={2}
        pointColor={{ theme: "background" }}
        pointBorderWidth={4}
        pointBorderColor={{ from: "serieColor" }}
        enablePointLabel={true}
        pointLabelYOffset={-12}
        enableArea={true}
        debugSlices={false}
        crosshairType="x"
        useMesh={true}
        enableSlices="x"
        sliceTooltip={({ slice }) => {
          return (
            <div
              style={{
                background: "white",
                padding: "9px 12px",
                border: "1px solid #ccc",
              }}
            >
              <div>x: {slice.id}</div>
              {slice.points.map((point) => (
                <div
                  key={point.id}
                  style={{
                    color: point.serieColor,
                    padding: "3px 0",
                  }}
                >
                  <strong>{point.serieId}</strong> [{point.data.yFormatted}]
                  <strong>{}</strong> [{point.data.yFormatted}]
                  <strong>{point.serieId}</strong> [{point.data.yFormatted}]
                  <strong>{point.serieId}</strong> [{point.data.yFormatted}]
                  <strong>{point.serieId}</strong> [{point.data.yFormatted}]
                </div>
              ))}
            </div>
          );
        }}
        // axisLeft={{
        //   format: (value) =>
        //     Number(value).toLocaleString("ru-RU", {
        //       minimumFractionDigits: 2,
        //     }),
        // }}
        // tooltip={({ id, color, data }) => (
        //   <strong style={{ color }}>
        //     {id}: {data}
        //   </strong>
        // )}
        // theme={{
        //   tooltip: {
        //     container: {
        //       background: "#333",
        //     },
        //   },
        // }}
        legends={[
          {
            anchor: "bottom-right",
            direction: "column",
            justify: true,
            translateX: 100,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: "left-to-right",
            itemWidth: 80,
            itemHeight: 20,
            itemOpacity: 0.75,
            symbolSize: 12,
            symbolShape: "circle",
            symbolBorderColor: "rgba(0, 0, 0, .5)",
            effects: [
              {
                on: "hover",
                style: {
                  itemBackground: "rgba(0, 0, 0, .03)",
                  itemOpacity: 1,
                },
              },
            ],
          },
        ]}
        motionConfig={{
          mass: 1,
          tension: 500,
          friction: 36,
          clamp: true,
          precision: 0.01,
          velocity: 0,
        }}
      />
    );
  }
}

export default LineChart;
