import React, { Component } from "react";
import { Redirect } from "react-router-dom";

class SignIn extends Component {
  constructor(props) {
    super(props);
    // const token = localStorage.getItem("token")
    const token = sessionStorage.getItem("token");

    let loggedIn = true;
    console.log(token);
    if (token == null) {
      loggedIn = false;
    }

    this.state = {
      username: "",
      password: "",
      loggedIn,
    };
    this.onChange = this.onChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }
  onChange(e) {
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  submitForm(e) {
    e.preventDefault();
    const { username, password } = this.state;
    //login magic
    if (
      (username === "admin" && password === "") ||
      (username === "demo" && password === "")
    ) {
      
      sessionStorage.setItem("userName", username);

      this.setState({
        loggedIn: true,
      });
    } else {
      alert("Please check username & password !!");
    }
  }
  render() {
    if (this.state.loggedIn) {
      console.log("SignIn", this.state.loggedIn);
      return <Redirect to="/dashboard"></Redirect>;
    }
    return (
      <div className="login">
        <div className="login-form">
          <h1>OpenMCP-Portal</h1>
          <form onSubmit={this.submitForm}>
            <input
              type="text"
              placeholder="Username"
              name="username"
              value={this.state.username}
              onChange={this.onChange}
            />
            <input
              type="password"
              placeholder="Password"
              name="password"
              value={this.state.password}
              onChange={this.onChange}
            />
            <input className="btn-signIn" type="submit" value="Sign In"/>
          </form>
        </div>
      </div>
    );
  }
}

export default SignIn;
