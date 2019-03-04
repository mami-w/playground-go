import React from 'react'

export function AddJwtToken() {
    if (localStorage.getItem('access_token')) {
        if (this.hasOwnProperty('headers')) {
            this.headers.Authorization = 'Bearer ' + localStorage.getItem('access_token');
        }
        else {
            this.headers = {
                'Authorization': 'Bearer ' + localStorage.getItem('access_token')
            }
        }
    }
}

