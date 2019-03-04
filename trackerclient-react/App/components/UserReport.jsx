import React from "react";
import { userEditMode } from './Datastructures'

export default class UserReport extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        console.log("UserReport render")
        const row = this.props.user;
        const selectedUser = this.props.selectedUser;
        const editMode = this.props.editMode;

        const userClassName = row.id == selectedUser ? "user-select" : "user";
        const editUserClassName = editMode != userEditMode.none ? "edit-user-hidden" : "edit-user"
        const deleteUserClassName = editMode != userEditMode.none ? "delete-user-hidden" : "delete-user"

        if (editMode == userEditMode.editUser && (row.id == selectedUser)) {
            return this.renderEdit(row)
        }
        else {
            return this.renderNormal(row, userClassName, editUserClassName, deleteUserClassName)
        }
    }

    renderNormal = (row, userClassName, editUserClassName, deleteUserClassName) => {
        return (
            <div className={userClassName} onClick={this.props.selectUser}>
                <div id="user-id">{row.id}</div>
                <div id="user-name">dummy</div>
                <button className={editUserClassName} onClick={this.props.setEditUser}>Edit</button>
                <button className={deleteUserClassName} onClick={this.props.deleteUser}>X</button>
            </div>
        )
    }

    renderEdit = (row) => {
        return (
        <form className="user" onSubmit={this.submitUser.bind(this, row.id)}>
            <div id="user-id">{row.id}</div>
            <div id="user-name">
                <input type="text" name="userName" defaultValue="dummy"/>
            </div>
            <button type="submit" className="submit-user">Submit</button>
            <button type="button" onClick={this.props.cancelEditMode}>Cancel</button>
        </form>
        )
    }

    submitUser = (userId, e) => {

        e.preventDefault();

         // todo - add user names
        var newName = e.target.elements.userName.value

        const user = { id: userId }

        this.props.finishUpdateUser(user)
    }
}