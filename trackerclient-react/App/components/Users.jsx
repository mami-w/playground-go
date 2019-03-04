import React from 'react'
import UserReport from './UserReport'
import uuid from 'node-uuid'
import { userEditMode } from './Datastructures'
import { AddJwtToken } from './HttpHelpers'

export default class Users extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            editMode : userEditMode.none,
            users: [],
            loading: true
        };
        console.log("users constructor")
    }

    componentDidMount() {
        this.loadUsers();
        console.log("Users did mount")

    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        console.log("Users did Update")
    }

    componentWillUnmount() {
        console.log("Users will unmount")
    }

    render() {

        console.log("Users render")
        const users = this.state.users;
        const editMode = this.state.editMode;

        const error = this.state.error;
        const loading = this.state.loading;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (loading) {
            return <div>Loading...</div>;
        } else {
            const selectedUser = this.props.selectedUser;

            return (
                <div>
                    <div>
                        <button onClick={this.reloadAll}>Reload Users</button>
                    </div>
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
                                    setEditUser={this.setEditUser.bind(this, user.id)}
                                    cancelEditMode={this.cancelEditMode}
                                    selectedUser={selectedUser}
                                    editMode={editMode}
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
                <button type="button" onClick={this.cancelEditMode}>Cancel</button>
            </form>
        )
    }

    reloadAll = () => {
         this.loadUsers();
    }

    loadUsers() {

        //if (!this.props.isAuthorized) { return; }

        const url = "/api/v1.0/tracker/user";

        let options = {};
        AddJwtToken.call(options);

        fetch(url, options)
            .then(
                (res) => {
                    if (res.status == 404) {
                        return [];
                    }
                    if (res.status == 401) {
                        this.props.setUnauthorized();
                        return [];
                    }
                    if (!res.ok) {
                        throw Error(res.statusText);
                    }
                    return res.json();
                },
                (error) => {
                    this.setState({loading: false, error: error})
                }
            )
            .then(
                (result) => {
                    this.setState({
                        loading: false,
                        users: result
                    });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
                (error) => {
                    this.setState({
                        loading: false,
                        error: error
                    });
                }
            )
            .catch(
                (error) => {
                    this.setState({
                        loading: true,
                        error: error
                    })
                }
            )
    }

    sendPostUser(user) {

        const url = `/api/v1.0/tracker/user/${user.id}`;
        let options = {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user)
        };
        AddJwtToken.call(options);

        fetch(url, options)
            .then((res) => {
                if (!res.ok) { throw Error(res.statusText) }
                this.setState({
                        users: this.state.users.concat([user]),
                        editMode : userEditMode.none
                    }                )
            }, (error) => {
                this.setState({
                    error:error
                })
            })
            .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendPutUser(user) {
        const id = user.id;
        const url = `/api/v1.0/tracker/user/${id}`;
        let options = {
            method: 'PUT',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user)
        };
        AddJwtToken.call(options);

        fetch(url, options)
            .then((res) => {
                if (!res.ok) { throw Error(res.statusText) }
                this.setState({
                        users: this.state.users.filter(user => user.id !== id).concat([updatedUser]),
                        editMode : userEditMode.none
                    })
            }, (error) => {
                this.setState({
                    error:error
                })
            })
            .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendDeleteUser(id) {
        const url = `/api/v1.0/tracker/user/${id}`;
        let options = {
            method: 'DELETE',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(id)
        };
        AddJwtToken.call(options);

        fetch(url, options)
            .then((res) => {
                if (!res.ok) { throw Error(res.statusText) }
                this.setState({
                        users: this.state.users.filter(user => user.id !== id),
                    })
            }, (error) => {
                this.setState({
                    error:error
                })
            })
            .catch(
                (error) => { this.setState({error: error})}
            );
    }

    setAddUser = () => {
        this.setState({
            editMode : userEditMode.addUser
        })
    }

    setEditUser = (id) => {
        this.props.selectUser(id);
        this.setState( {
            editMode : userEditMode.editUser
        })
    }

    cancelEditMode = (e) => {
        this.setState( {
            editMode : userEditMode.none
        })
    }

    finishAddUser = (e) => {
        e.preventDefault();

        const id = e.target.elements.userId.value // from elements property
        var user = { id: id }

        this.sendPostUser(user)
    }

    deleteUser = (id, e) => {
        this.sendDeleteUser(id)
    }

    finishUpdateUser = (updatedUser) => {
        this.sendPutUser(updatedUser)
    }
}