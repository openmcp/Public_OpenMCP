import React, { Component } from 'react';
import {Link, Redirect} from 'react-router-dom';

class Admin extends Component {

    constructor(props){
        super(props)
        // const token = localStorage.getItem("token")
        const token = sessionStorage.getItem("token")
        

        let loggedIn = true
        if(token == null){
            loggedIn = false
        }

        this.state = {
            loggedIn
        }
    }
    render() {
        if(this.state.loggedIn === false){
            return <Redirect to="/"></Redirect>
        }
        return (
            <div>
                <h1>This is an Adminpage. only Auth people can see this.</h1>
                <Link to="/logout">Logout</Link>
            </div>
        );
    }
}

export default Admin;