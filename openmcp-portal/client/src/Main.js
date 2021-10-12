import React, { Component } from "react";
import Head from "./components/layout/Head";
import Contents from "./components/layout/Contents";
import { Redirect } from "react-router-dom";
import { AsyncStorage } from 'AsyncStorage';
import MainLeftMenu2 from './components/layout/MainLeftMenu2';

class Main extends Component {
  constructor(props) {
    super(props);
    let token = null;
    AsyncStorage.getItem("token", (err, result) => { 
      token = result;
    })

    let loggedIn = true;
    if (token === null || token === "null" || token==="") {
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
        <Head onSelectMenu={this.onLeftMenu} path={this.props.location.pathname}/>
        <MainLeftMenu2 info={this.props}/>

        <Contents path={this.props.location.pathname} onSelectMenu={this.onLeftMenu} info={this.props}/>
      </div>
    );
  }
}

export default Main;
