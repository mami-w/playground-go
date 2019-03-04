import React from 'react'
import Users from './Users'
import Entries from './Entries'
import Auth from './Auth';

export default class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedUser: null,
            isAuthorized: false
        };
        this.auth = new Auth();
    }

    componentDidMount() {
        this.auth.parseHash();
        const idToken = localStorage.getItem('id_token');
        if (idToken) { this.setAuthorized(true) } else { this.setAuthorized(false)};
    }

    render() {
        const selectedUser = this.state.selectedUser;

        if (!this.state.isAuthorized) {
            return (
                <div>
                    <button onClick={this.login}>Login</button>
                </div>
            );
        }
        else {
            return (
                <div>
                    <button onClick={this.logout}>Logout</button>
                    <Users selectUser={this.selectUser} selectedUser={selectedUser} setUnauthorized={this.setUnauthorized}/>
                    <Entries selectedUser={selectedUser} setUnauthorized={this.setUnauthorized}/>
                </div>
            );
        }
    }

    selectUser = (userId) => {
        this.setState( {
            selectedUser : userId
        })
    }

    setUnauthorized = () => {
         this.setAuthorized(false)
    }

    // example of how to do login from client
    login = () => {
         this.auth.login();
    }

    logout = () => {
         this.auth.logout();
    }

    // Set user login state
    setAuthorized(value){
        if(this.state.isAuthorized != value){
            this.setState({ isAuthorized: value });
        }
    }
}