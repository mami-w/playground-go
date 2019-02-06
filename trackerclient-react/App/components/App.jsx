import React from 'react'
import Users from './Users'
import Entries from './Entries'

export default class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedUser: null
        };
        console.log("app constructor")
    }

    componentWillMount() {
        console.log("app will mount")
    }

    componentWillUnmount() {
        console.log("app will unmount")
    }

    render() {
        const selectedUser = this.state.selectedUser;

         return (
             <div>
                 <Users selectUser={this.selectUser} selectedUser={selectedUser}/>
                 <Entries selectedUser={selectedUser}/>
             </div>
         )
    }

    selectUser = (userId) => {
        console.log("called selectUser")
        this.setState( {
            selectedUser : userId
        })
    }
}