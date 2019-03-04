import React from 'react'
import EntryReport from './EntryReport'
import uuid from 'node-uuid'
import { entryEditMode } from './Datastructures'
import { AddJwtToken } from './HttpHelpers'

export default class Entries extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            editMode : entryEditMode.none,
            entries: [],
            loading : false, // caching - later
            selectedEntry : null
        };
        console.log("entries constructor")
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

    render() {
        console.log("Entries will render")

        const entries = this.state.entries;
        const error = this.state.error;
        const loading = this.state.loading;
        const selectedUser = this.props.selectedUser;

        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (selectedUser == null) {
            return <div>No user selected.</div>
        }
        else if (loading) {
            return <div>Loading...</div>;
        } else {
            console.log("printing %d entries", entries.length);

            const editMode = this.state.editMode;
            const selectedEntry = this.state.selectedEntry;

            return (
                <div>
                    <div>
                        <button onClick={this.reloadAll}>Reload Entries</button>
                    </div>
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
                                <div className="entry-header-column">Duration (min)</div>
                            </div>
                        </li>
                    </ul>
                    <ul>
                        {entries.map(entry => {
                            return (
                                <li key={entry.id}>
                                    <EntryReport
                                        entry={entry}
                                        deleteEntry={this.deleteEntry.bind(this, entry.id, entry.userid)}
                                        finishUpdateEntry={this.finishUpdateEntry.bind(this, entry.id)}
                                        selectEntry={this.selectEntry.bind(this, entry.id)}
                                        setEditEntry={this.setEditEntry.bind(this, entry.id)}
                                        cancelEditMode={this.cancelEditMode}
                                        selectedEntry={selectedEntry}
                                        editMode={editMode}
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
                <button type="button" onClick={this.cancelEditMode}>Cancel</button>
            </form>
        )
    }

    reloadAll = () => {
         this.loadEntries();
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
        else {
            this.loadEntries();
        }
    }

    loadEntries() {

        //if (!this.props.isAuthorized) { return; }

        const selectedUser = this.props.selectedUser;
        if (selectedUser == null) {
            return;
        }
        const url = `/api/v1.0/tracker/user/${selectedUser}`;
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
                        if (!res.ok) { throw Error(res.statusText) }
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
                            entries: result
                        },
                        () => {
                            const key = 'entries${selectedUser}'
                            sessionStorage.setItem(key, JSON.stringify(this.state.entries));
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

    sendPostEntry(entry) {
        const url = `/api/v1.0/tracker/user/${entry.userid}/entry/${entry.id}`;
        let options = {
            method: 'POST',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(entry)
        };
        AddJwtToken.call(options);

        fetch(url, options)
            .then(
                (res) =>  {
                    if (!res.ok) { throw Error(res.statusText) }
                    this.setState({
                            entries: this.state.entries.concat([entry]),
                            editMode : entryEditMode.none
                        },
                        () => {
                            const key = `entries${entry.userid}`;
                            sessionStorage.setItem(key, JSON.stringify(this.state.entries))
                        }
                    );
                },
                (error) => { this.setState({error:error})}
            )
            .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendPutEntry(entry) {
        const updatedEntry = entry;
        const url = `/api/v1.0/tracker/user/${entry.userid}/entry/${entry.id}`;
        let options = {
            method: 'PUT',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(entry)
        };
        AddJwtToken.call(options);

        fetch(url, options)
        .then((res) => {
                if (!res.ok) { throw Error(res.statusText) }
            this.setState({
                entries: this.state.entries.filter(entry => entry.id !== updatedEntry.id).concat([updatedEntry]),
                editMode : entryEditMode.none
            },
            () => {
                const key = `entries${entry.userid}`;
                sessionStorage.setItem(key, JSON.stringify(this.state.entries))
            }
            );
        },(error) => { this.setState({error:error})}
            )
        .catch(
            (error) => { this.setState({error: error})}
        );
    }

    sendDeleteEntry(id, userId) {
        const url = `/api/v1.0/tracker/user/${userId}/entry/${id}`;
        let options = {
            method: 'DELETE',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            }
        };
        AddJwtToken.call(options);

        fetch(url, options)
            .then((res) => {
                    if (!res.ok) { throw Error(res.statusText) }
                this.setState({
                        entries: this.state.entries.filter(entry => entry.id !== id),
                    },
                    () => {
                        const key = `entries${userId}`;
                        sessionStorage.setItem(key, JSON.stringify(this.state.entries))
                    });
            }, (error) => { this.setState({ error: error })}
            )
            .catch(
                (error) => { this.setState({error: error})}
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

    setEditEntry = (id) => {
        this.selectEntry(id);
        this.setState({
            editMode: entryEditMode.editEntry
        })
    }

    cancelEditMode  = (e)  => {
        //e.preventDefault();
        //e.stopPropagation();

        this.setState({
            editMode: entryEditMode.none
        })
    }

    finishAddEntry = (e) => {
        e.preventDefault();

        const selectedUser = this.props.selectedUser;
        if (selectedUser == null) {
           return;
        }

        const id = e.target.elements.entryId.value;
        const userId = selectedUser;
        const now = (new Date()).toJSON();
        var entry = {
            id: id,
            userid: userId,
            entryType: "1",
            startTime: now,
            length: 3600000000000
        }

        this.sendPostEntry(entry);
    }

    deleteEntry = (id, userId, e) => {
        this.sendDeleteEntry(id, userId);
    }

    finishUpdateEntry = (id, updatedEntry) => {
        this.sendPutEntry(updatedEntry);
    }
}