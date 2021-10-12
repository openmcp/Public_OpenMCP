import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { AsyncStorage } from 'AsyncStorage';

class SignOut extends Component {
    
    constructor(props){
        super(props)
        // localStorage.removeItem("token");
        // localStorage.removeItem("userName");
        // localStorage.removeItem("roles");
        AsyncStorage.removeItem("token");
        AsyncStorage.removeItem("userName");
        AsyncStorage.removeItem("roles");
        AsyncStorage.removeItem("projects");
    }
    render() {
        return (
            <div>
                <h1>You have been logged out!!!</h1>
                <Link to="/">Login Again</Link>
            </div>
        );
    }
}

export default SignOut;