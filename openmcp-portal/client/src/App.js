import React, { Component } from "react";
import "./css/style.css";
import SignIn from "./components/common/SignIn";
import Main from './Main';
import { Switch, Route } from "react-router-dom";

class App extends Component {

  render() {
    return (
      <Switch>
        <Route exact path="/login" component={SignIn} />
        <Route path="/" component={Main} />
      </Switch>
    );
  }
}

export default App;
