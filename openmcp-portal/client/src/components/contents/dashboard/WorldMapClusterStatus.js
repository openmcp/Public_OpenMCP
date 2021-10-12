import React, { Component } from "react";
import WorldMap from "react-svg-worldmap";
// import './App.css';

class WorldMapClusterStatus extends Component {
  constructor(props) {
    super(props);
    this.state = {
      mapSize: Math.min(window.innerHeight, window.innerWidth),
    };
  }

  onSizeUp = () => {
    let mapSize = this.state.mapSize;
    this.setState({ mapSize: mapSize + 100 });
  };

  onSizeDown = () => {
    let mapSize = this.state.mapSize;
    this.setState({ mapSize: mapSize - 100 });
  };

  render() {
    const data = [
      { country: "cn", value: 1 }, // china
      { country: "in", value: 2 }, // india
      { country: "us", value: 3 }, // united states
      { country: "id", value: 4 }, // indonesia
      { country: "pk", value: 5 }, // pakistan
      { country: "br", value: 6 }, // brazil
      { country: "ng", value: 7 }, // nigeria
      { country: "bd", value: 8 }, // bangladesh
      { country: "ru", value: 9 }, // russia
      { country: "kr", value: 10 }, // mexico
    ];
    return (
      <div className="content-box" style={{ width: this.props.width, display: "inline-block"}}>
        <div className="cb-header">
          <span>World Cluster Status</span>
        </div>
        <div
          className="cb-body"
          style={{ position: "relative", width: "100%"}}
        >
          <button class="btn-worldmap-size" onClick={this.onSizeUp}>
            +
          </button>
          <button class="btn-worldmap-size" onClick={this.onSizeDown}>
            -
          </button>
          <div  style={{textAlign:"center" }}>
            <WorldMap
              color="#0088fe"
              // title="Top 10 Populous Countries"
              value-suffix="people"
              // size="responsive"
              size={this.state.mapSize}
              // frame={true}
              data={data}
            />
          </div>
          
        </div>
      </div>
    );
  }
}

export default WorldMapClusterStatus;
