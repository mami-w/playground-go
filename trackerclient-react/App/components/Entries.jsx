import React from 'react'
import EntryReport from './EntryReport'
import uuid from 'node-uuid'
import { entryEditMode } from './Datastructures'

export default class Entries extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            editMode : entryEditMode.none,
            entries: [],
            needLoadEntries : true, // todo: needs major rethinking
            selectedEntry : null
        };
    }

    componentDidMount() {
        console.log("Entries did mount")
        this.checkReload();
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        console.log("Entries did update")
        if (prevProps.selectedUser != this.props.selectedUser) {
            this.checkReload();
        }
    }

    componentWillUnmount() {
        console.log("Entries will unmount")
    }

    render() {
        console.log("Entries will render")

        const entries = this.state.entries;

        const error = this.state.error;
        const needLoadEntries = this.state.needLoadEntries;
        const selectedUser = this.props.selectedUser;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (selectedUser == null) {
            return <div>No user selected...</div>
        }
        else if (needLoadEntries) {
            return <div>Loading...</div>;
        } else {
            console.log("printing %d entries", entries.length);
            const editMode = this.state.editMode;

            return (
                <div>
                    {(editMode == entryEditMode.addEntry) ? this.renderNewEntry() : null}
                    {(editMode == entryEditMode.none && (selectedUser != null)) ? <button onClick={this.setAddEntry}>New Entry</button> : null}

                    <div id="entries-caption">Entries</div>
                    <ul>
                        <li>
                            <div className="entry-header">
                                <div className="entry-header-column">Id</div>
                                <div className="entry-header-column">UserId</div>
                                <div className="entry-header-column">Entry Type</div>
                                <div className="entry-header-column">Start Time</div>
                                <div className="entry-header-column">Duration</div>
                            </div>
                        </li>
                    </ul>
                    <ul>
                        {entries.map(entry => {
                            return (
                                <li key={entry.id}>
                                    <EntryReport
                                        entry={entry}
                                        deleteEntry={this.deleteEntry.bind(this, entry.id)}
                                        finishUpdateEntry={this.finishUpdateEntry.bind(this, entry.id)}
                                        selectEntry={this.selectEntry.bind(this, entry.id)}
                                        setEditEntry={this.setEditEntry}
                                        selectedEntry={this.state.selectedEntry}
                                        editMode={this.state.editMode}
                                    />
                                </li>)
                        })}
                    </ul>
                </div>
            )
        }
    }

    renderNewEntry () {
        const newEntryId = uuid.v4()
        return (
            <form className="inputForm" onSubmit={this.finishAddEntry}>
                <input type="text" name="entryId" defaultValue={newEntryId} />
                <button type="submit">Create Entry</button>
            </form>
        )
    }

    checkReload() {
        const selectedUser = this.props.selectedUser;

        if (selectedUser == null) {
            this.setState({entries:[]})
            return;
        }

        const key = 'entries${selectedUser}';
        const entries = sessionStorage.getItem(key)
        if (entries != null) {
            this.setState({ entries : JSON.parse(entries) })
        }

        // double check this...
        const needLoadEntries = sessionStorage.getItem("needLoadEntries")
        if (needLoadEntries == null || needLoadEntries == true) {
            this.loadEntries();
        }
        else {
            this.setState({needLoadEntries:false})
        }
    }

    loadEntries() {
        const selectedUser = this.props.selectedUser;
        if (selectedUser == null) {
            return;
        }
        //const url = `/api/v1.0/tracker/user/${selectedUser}`; // todo: create correct url
        const url = "/api/v1.0/tracker/user/1"
        fetch(url)
            .then(
                res => res.json(),
                (error) => {
                    this.setState({needLoadEntries: false, error: error})
                }
            )
            .then(
                (result) => {
                    this.setState({
                            needLoadEntries: false,
                            entries: result
                        },
                        () => {
                            sessionStorage.setItem("needLoadEntries", false);
                            const key = 'entries${selectedUser}'
                            sessionStorage.setItem(key, JSON.stringify(this.state.entries));
                        });
                },
                // Note: it's important to handle errors here
                // instead of a catch() block so that we don't swallow
                // exceptions from actual bugs in components.
                (error) => {
                    this.setState({
                        needLoadEntries: false,
                        error: error
                    });
                }
            )
            .catch(
                (error) => {
                    this.setState({
                        needLoadEntries: true,
                        error: error
                    })
                }
            )
    }

    sendPostEntry(entry) {
        const url = `/api/v1.0/tracker/user/${entry.userId}/entries/${entry.id}`;
        fetch(url, {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(entry)
        })
            .then(
                () =>  {
                    this.setState({
                            entries: this.state.entries.concat([entry]),
                            editMode : entryEditMode.none
                        },
                        () => {
                            const key = `entries${entry.userId}`;
                            sessionStorage.setItem(key, JSON.stringify(this.state.entries))
                        }
                    );
                }
            )
            .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendPutEntry(entry) {
        const url = `/api/v1.0/tracker/user/${entry.userId}/entries/${entry.id}`;
        fetch(url, {
            method: 'PUT',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(entry)
        })
        .then(() => {
            this.setState({
                entries: this.state.entries.filter(entry => entry.id !== updatedEntry.id).concat([updatedEntry]),
                editMode : entryEditMode.none
            },
            () => {
                const key = `entries${entry.userId}`;
                sessionStorage.setItem(key, JSON.stringify(this.state.entries))
            }
            );
        })
        .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendDeleteEntry(id) {
        // todo: ajax missing
        this.setState({
                entries: this.state.entries.filter(entry => entry.id !== id),
            },
            () => {
                const key = `entries${entry.userId}`;
                sessionStorage.setItem(key, JSON.stringify(this.state.entries))
            }
        );
    }

    selectEntry = (id) => {
         this.setState({
            selectedEntry:id
        })
    }

    setAddEntry = () => {
        this.setState({
            editMode : entryEditMode.addEntry
        })
    }

    setEditEntry = () => {
        this.setState({
            editMode: entryEditMode.editEntry
        })
    }

    // add a user to the users
    finishAddEntry = (e) => {
        const selectedUser = this.props.selectedUser;
        if (selectedUser == null) {
           return;
        }

        const id = e.target.elements.entryId.value;

        var entry = {
            id: id,
            userId: selectedUser,
            // todo - default values; type, now, 1h
            entryType: "1"
        }

        sendPostEntry(entry);
    }

    deleteEntry = (id, e) => {
        e.stopPropagation();
        sendDeleteEntry(id);
    }

    finishUpdateEntry = (updatedEntry) => {
        sendPutEntry(updatedEntry);
    }
}