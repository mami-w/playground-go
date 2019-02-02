import React from 'react'
import Users from './Users'
import Entries from './Entries'

export default class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedUser: null
        };
    }

    componentWillMount() {
        console.log("app will mount")
    }

    componentWillUnmount() {
        console.log("app will unmount")
    }

    render() {
         return (
                <div>
                    <Users
                        selectUser={this.selectUser}
                        selectedUser={this.state.selectedUser}
                    />
                    <Entries selectedUser={this.state.selectedUser}/>
                </div>
         )

    }

    selectUser = (id) => {
                        console.log("called selectUser")
        this.setState( {
            selectedUser : id
        })
    }
}