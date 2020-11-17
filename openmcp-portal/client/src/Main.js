import React, { Component } from "react";
import Head from "./components/layout/Head";
import Contents from "./components/layout/Contents";
import { Redirect } from "react-router-dom";
import LeftMenu from './components/layout/LeftMenu';

class Main extends Component {
  constructor(props) {
    super(props);
    // const token = localStorage.getItem("token");
    const token = sessionStorage.getItem("token");

    let loggedIn = true;
    if (token == null) {
      loggedIn = false;
    }

    this.state = {
      isLeftMenuOn: false,
      isLogined: true,
      loggedIn,
      windowHeight: undefined,
      windowWidth: undefined
    };
  }

//   componentWillMount(){
//     console.log("WINDOW : ",window);
//     this.setState({height: window.innerHeight + 'px',width:window.innerWidth+'px'});
// }

  handleResize = () => this.setState({
    windowHeight: window.innerHeight,
    windowWidth: window.innerWidth
  });

  componentDidMount() {
    this.handleResize();
    window.addEventListener('resize', this.handleResize)
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize)
  }

  render() {
    if (!this.state.loggedIn) {
      return <Redirect to="/login"></Redirect>;
    }
    return (
      <div className="wrapper" style={{minHeight:this.state.windowHeight}}>
        <Head onSelectMenu={this.onLeftMenu} />
        {/* <LeftMenu /> */}
        <Contents path={this.props.location.pathname} onSelectMenu={this.onLeftMenu} info={this.props}/>
      </div>
    );
  }
}

export default Main;
