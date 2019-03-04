import auth0 from 'auth0-js';
import * as authConstants from './AuthConstants'

export default class Auth {
    home = window.location.origin;

    auth0 = new auth0.WebAuth({
        domain: authConstants.DOMAIN,
        clientID: authConstants.CLIENTID,
        redirectUri: window.location.origin,
        responseType: 'token id_token',
        scope: 'openid',
        audience : authConstants.AUDIENCE
    });

    login() {
        this.auth0.authorize();
    }

    logout() {
        localStorage.removeItem("access_token");
        localStorage.removeItem("id_token");
        localStorage.removeItem("profile");

        this.auth0.logout({
            returnTo: this.home,
            clientID: authConstants.CLIENTID
        })
    }

    parseHash() {
        this.auth0.parseHash((err, authResult) => {
            if (err) {
                return console.log(err);
            }
            if(authResult !== null && authResult.accessToken !== null && authResult.idToken !== null){
                localStorage.setItem('access_token', authResult.accessToken);
                localStorage.setItem('id_token', authResult.idToken);
                localStorage.setItem('profile', JSON.stringify(authResult.idTokenPayload));
                window.location = window.location.href.substr(0, window.location.href.indexOf('#'))
            }
        });
    }
}