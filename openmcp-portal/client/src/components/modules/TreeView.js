import React from "react";
import Tree from "react-d3-tree";

const myTreeData = [
  {
    name: "Top Level",
    attributes: {
      keyA: "val A",
      keyB: "val B",
      keyC: "val C",
    },
    children: [
      {
        name: "Level 2: A",
        attributes: {
          keyA: "val A",
          keyB: "val B",
          keyC: "val C",
        },
      },
      {
        name: "Level 2: B",
      },
    ],
  },
];

class NodeLabel extends React.PureComponent {
  render() {
    const { className, nodeData } = this.props;
    return (
      <div className={className}>
        <h2>{nodeData.name}</h2>
        {nodeData._children && (
          <button>{nodeData._collapsed ? "Expand" : "Collapse"}</button>
        )}
      </div>
    );
  }
}

class MyComponent extends React.Component {
  state = {};
  componentDidMount() {
    const dimensions = this.treeContainer.getBoundingClientRect();
    this.setState({
      translate: {
        x: dimensions.width / 2,
        y: dimensions.height / 2
      }
    });
  }

  render() {
    const svgSquare = {
      shape: "rect",
      shapeProps: {
        width: 20,
        height: 20,
        x: -10,
        y: -10,
      },
    };



    
    // const styles={
    //   links: <svgStyleObject>,
    //   nodes: {
    //     node: {
    //       circle: <svgStyleObject>,
    //       name: <svgStyleObject>,
    //       attributes: <svgStyleObject>,
    //     },
    //     leafNode: {
    //       circle: <svgStyleObject>,
    //       name: <svgStyleObject>,
    //       attributes: <svgStyleObject>,
    //     },
    //   },
    // };

    const containerStyles = {
      width: "100%",
      height: "50vh"
    };

    return (
      /* <Tree /> will fill width/height of its container; in this case `#treeWrapper` */
      // <div id="treeWrapper" style={{ width: "50em", height: "20em" }}>
        <div style={containerStyles} ref={(tc) => (this.treeContainer = tc)}>
        <Tree
          data={myTreeData}
          nodeSvgShape={svgSquare}
          orientation="vertical" //horizontal
          translate={this.state.translate}
          // styles = {styles}
        />

        <Tree
          data={myTreeData}
          allowForeignObjects
          translate={this.state.translate}
          nodeLabelComponent={{
            render: <NodeLabel className="myLabelComponentInSvg" />,
            foreignObjectWrapper: {
              y: 24,
            },
          }}
          orientation="vertical" //horizontal
        />
      </div>
    );
  }
}

export default MyComponent;
