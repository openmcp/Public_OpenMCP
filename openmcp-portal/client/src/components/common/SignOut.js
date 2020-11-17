import React, { Component } from 'react';
import { Link } from 'react-router-dom';

class SignOut extends Component {
    constructor(props){
        super(props)
        // localStorage.removeItem("token")
        sessionStorage.removeItem("token")
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