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

        // todo: translate startTime and duration into something typed
        const duration = row.length / (60 * 1000000000)

        if ((editMode == entryEditMode.editEntry) && (row.id == selectedEntry)) {
            return this.renderEdit(row, duration)
        }
        else {
            const entryClassName = row.id == selectedEntry ? "entry-select" : "entry";
            const editEntryClassName = editMode != entryEditMode.none ? "edit-entry-hidden" : "edit-entry"
            const deleteEntryClassName = editMode != entryEditMode.none ? "delete-entry-hidden" : "delete-entry"

            return this.renderNormal(row, duration, entryClassName, editEntryClassName, deleteEntryClassName)
        }
    }

    renderNormal = (row, duration, entryClassName, editEntryClassName, deleteEntryClassName) => {
        return (
            <div className={entryClassName} onClick={this.props.selectEntry}>
                <div id="entry-id">{row.id}</div>
                <div id="entry-userId">{row.userID}</div>
                <div id="entry-entryType">{row.entryType}</div>
                <div id="entry-startTime">{row.startTime}</div>
                <div id="entry-length">{duration}</div>
                <button className={editEntryClassName} onClick={this.props.setEditEntry}>Edit</button>
                <button className={deleteEntryClassName} onClick={this.props.deleteEntry}>X</button>
            </div>
        )
    }

    renderEdit = (row, duration) => {
        return (
            <form className="entry" onSubmit={this.submitEntry.bind(this, row.id, row.userID)}>
                <div id="entry-id">{row.id}</div>
                <div id="entry-userId">{row.userID}</div>
                <div>
                    <input type="text" name="entryType" defaultValue={row.entryType}/>
                </div>
                <div>
                    <input type="datetime-local" name="startTime" defaultValue={row.startTime}/>
                </div>
                <div>
                    <input type="text" name="length" defaultValue={duration}/>
                </div>
                <button type="submit" className="submit-entry">Submit</button>
                <button type="button" className="cancel-entry" onClick={this.props.cancelEditMode}>Cancel</button>
            </form>
        )
    }

    submitEntry = (entryId, userId, e) => {

        e.preventDefault();

        // todo: translate these into strings
        var entryType = e.target.elements.entryType.value;
        var startTime = e.target.elements.startTime.value;
        var entryLength = e.target.elements.length.value;

        entryLength = entryLength * 60 * 1000000000;

        const entry = {
            id: entryId,
            userID : userId,
            entryType : entryType,
            startTime : startTime, // translate
            length: entryLength
        }

        this.props.finishUpdateEntry(entry)
    }
}