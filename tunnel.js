// @file tunnel.js
class Tunnel {
    hostUrl = 'http://127.0.0.1:9999';
    constructor(newUrl) {
        if (!!newUrl) {
            this.hostUrl = newUrl;
        }
    }
    toURLQuery(json) {
        return Object.keys(json).map((k) => String(encodeURIComponent(k) + '=' + encodeURIComponent(json[k]))).join('&');
    }
    login(username, password) {
        this.authToken = btoa(username + ':' + password);
    }
    request(service, payload, success, failure) {
        // console.log('Tunnel.request', service, payload);
        var xhr = new XMLHttpRequest();
        xhr.withCredentials = true;

        xhr.addEventListener('readystatechange', () => {
            if (xhr.readyState === 4 && xhr.status === 200) {
                success(JSON.parse(xhr.responseText));
            }
        });
        xhr.addEventListener('error', () => failure(xhr, xhr.statusText, xhr.responseText));

        xhr.open('POST', service);
        xhr.setRequestHeader('authorization', 'Basic ' + this.authToken);
        xhr.setRequestHeader('content-type', 'application/x-www-form-urlencoded');

        (!!payload) ? xhr.send(this.toURLQuery(payload)) : xhr.send();
    }
    authenticate(success, failure) {
        this.request(this.hostUrl + '/authenticate', null,
            (res) => {
                // console.log('Tunnel.authenticate.success', res);
                if (res.status == 'failure') {
                    failure(res, res.status, res.response)
                } else {
                    success(res.response);
                }
            }, failure
        );
    }
    terminal(command, success, failure) {
        this.request(this.hostUrl + '/terminal', {
                cmd: command
            },
            (res) => {
                // console.log('Tunnel.terminal.success', res);
                if (res.status == 'failure') {
                    failure(res, res.status, res.response)
                } else {
                    success(res.response);
                }
            }, failure
        );
    }
}

/*
// e.g.
(function main() {
    console.log('main');
    let tunnel = new Tunnel();
    let failure = (options, status, error) => {
        console.log('[1] main.failure', status, error);
    };
    tunnel.login('admin', '123456');
    tunnel.authenticate((result) => {
        console.log('[1] main.authenticate.success', result);
        tunnel.terminal('ls', (result) => {
            console.log('[2] main.terminal.success', result);
        }, failure);
    }, failure);
}());
*/