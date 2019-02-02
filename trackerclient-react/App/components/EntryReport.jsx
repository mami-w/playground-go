import React from "react";
import { entryEditMode } from './Datastructures'

export default class EntryReport extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const row = this.props.entry;
        const selectedEntry = this.props.selectedEntry;
        const editMode = this.props.editMode;

        const entryClassName = row.id == selectedEntry ? "entry-select" : "entry";
        const editEntryClassName = editMode != entryEditMode.none ? "edit-entry-hidden" : "edit-entry"
        const deleteEntryClassName = editMode != entryEditMode.none ? "delete-entry-hidden" : "delete-entry"

        if ((editMode == entryEditMode.editEntry) && (row.id == selectedEntry)) {
            return this.renderEdit(row)
        }
        else {
            return this.renderNormal(row, entryClassName, editEntryClassName, deleteEntryClassName)
        }

    }


    renderNormal = (row, entryClassName, editEntryClassName, deleteEntryClassName) => {
        return (
            <div className={entryClassName} onClick={this.props.selectEntry}>
            <div id="entry-id">{row.id}</div>
    <div id="entry-userId">{row.userId}</div>
    <div id="entry-entryType">{row.entryType}</div>
    <div id="entry-startTime">{row.startTime}</div>
    <div id="entry-length">{row.length}</div>
                <button className={editEntryClassName} onClick={this.props.setEditEntry}>Edit</button>
                <button className={deleteEntryClassName} onClick={this.props.deleteEntry}>X</button>
    </div>
        )
    }


    renderEdit = (row) => {
        return (
            <form className="entry" onSubmit={this.submitEntry.bind(this, row.id)}>
                <div id="entry-id">{row.id}</div>
                <div id="entry-userId">{row.userId}</div>
                <div>
                    <input type="text" name="entryType" defaultValue="1"/>
                </div>
                <div>
                    <input type="text" name="startTime" defaultValue="1"/>
                </div>
                <div>
                    <input type="text" name="length" defaultValue="1"/>
                </div>
                <button type="submit" className="submit-user">Submit</button>
            </form>
        )
    }

    submitEntry = (entryId, e) => {

        // todo - add user names
        var entryType = e.target.elements.entryType.value;
        var startTime = e.target.elements.startTime.value;
        var entryLength = e.target.elements.length.value;

        const entry = {
            id: entryId,
            userId : "dummy", // todo
            entryType : entryType,
            startTime : startTime, // translate
            length: entryLength // translate
        }

        this.props.finishUpdateUser(entry)
    }
}