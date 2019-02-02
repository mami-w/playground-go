import React from 'react'
import UserReport from './UserReport'
import uuid from 'node-uuid'
import { userEditMode } from './Datastructures'

export default class Users extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            editMode : userEditMode.none,
            users: [],
            needLoadUsers : true
        };
    }

    componentDidMount() {

        console.log("Users did mount")
        this.checkReload();
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        console.log("Users did Update")
        if (prevProps.selectedUser != this.props.selectedUser) {
            this.checkReload();
        }
    }

    componentWillUnmount() {
        console.log("Users will unmount")
    }

    render() {

        console.log("Users render")
        const users = this.state.users;
        const editMode = this.state.editMode;

        const error = this.state.error;
        const needLoadUsers = this.state.needLoadUsers;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (needLoadUsers) {
            return <div>Loading...</div>;
        } else {
            return (
                <div>
                    {(editMode == userEditMode.addUser) ? this.renderNewUser() : null}
                    {(editMode == userEditMode.none) ? <button onClick={this.setAddUser}>New User</button> : null}

                        <div id="users-caption">Users</div>
                        <ul>
                            <li>
                                <div className="user-header">
                                    <div id="user-header-id">Id</div>
                                    <div id="user-header-name">Name</div>
                                </div>
                            </li>
                        </ul>
                        <ul>
                        {users.map(user => {
                            return (
                                <li key={user.id}>
                                <UserReport
                                    user={user}
                                    deleteUser={this.deleteUser.bind(this, user.id)}
                                    finishUpdateUser={this.finishUpdateUser.bind(this)}
                                    selectUser={this.props.selectUser.bind(this, user.id)}
                                    setEditUser={this.setEditUser}
                                    selectedUser={this.props.selectedUser}
                                    editMode={this.state.editMode}
                                />
                                </li>)
                        })}
                        </ul>
                </div>
            )
        }
    }

    renderNewUser () {
        // create new guid
        const id = uuid.v4()
        return (
            <form
                className="inputForm"
                onSubmit={this.finishAddUser}
            >
                <input type="text" name="userId" defaultValue={id} />
                <button type="submit">Create User</button>
            </form>
        )
    }

    checkReload() {
        const users = sessionStorage.getItem("users")
        if (users != null) {
            this.setState({ users : JSON.parse(users) })
        }

        const needLoadUsers = sessionStorage.getItem("needLoadUsers")
        if (needLoadUsers == null || needLoadUsers == true) {
            this.loadUsers();
        }
        else {
            this.setState({needLoadUsers:false})
        }
    }

    loadUsers() {
        const url = "/api/v1.0/tracker/user/1"; // todo: create correct url
        fetch(url)
            .then(
                res => res.json(),
                (error) => {
                    this.setState({needLoadUsers: false, error: error})
                }
            )
            .then(
                (result) => {
                    this.setState({
                        needLoadUsers: false,
                        users: result
                    },
                        () => {
                            sessionStorage.setItem("needLoadUsers", false);
                            sessionStorage.setItem("users", JSON.stringify(this.state.users));
                        });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
                (error) => {
                    this.setState({
                        needLoadUsers: false,
                        error: error
                    });
                }
            )
            .catch(
                (error) => {
                    this.setState({
                        needLoadUsers: true,
                        error: error
                    })
                }
            )
    }

    postUser(user) {
        const url = `/api/v1.0/tracker/user/${user.id}`;
        fetch(url, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user)
            }).catch(
            (error) => { this.setState({error: error})}
        );
    }

    putUser(user) {
        const url = `/api/v1.0/tracker/user/${user.id}`;
        fetch(url, {
            method: 'PUT',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user)
        }).catch(
            (error) => { this.setState({error: error})}
        );
    }

    setAddUser = () => {
                console.log("called setAddUser");
        this.setState({
            editMode : userEditMode.addUser
        })
    }

    setEditUser = () => {
                console.log("called setEditUser")
        this.setState( {
            editMode : userEditMode.editUser
        })
    }

    finishAddUser = (e) => {
                console.log("called finishAddUser")
        const id = e.target.elements.userId.value // from elements property
        var user = { id: id }
        this.postUser(user)
        this.setState({
            users: this.state.users.concat([user]),
            editMode : userEditMode.none
            },
            () => sessionStorage.setItem("users", JSON.stringify(this.state.users))
        );
    }

    // todo...
    // delete our user
    deleteUser = (id, e) => {
                console.log("called finishUpdateUser")
        e.stopPropagation();
        // todo: ajax call
        this.setState({
            users: this.state.users.filter(user => user.id !== id),
        },
            () => sessionStorage.setItem("users", JSON.stringify(this.state.users)));
    }

    finishUpdateUser = (updatedUser) => {
        this.putUser(updatedUser)
        this.setState({
            users: this.state.users.filter(user => user.id !== id).concat([updatedUser]),
            editMode : userEditMode.none
        },
            () => sessionStorage.setItem("users", JSON.stringify(this.state.users)));
    }
}