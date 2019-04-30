// @file tunnel.js
class Tunnel {
    hostUrl = 'http://127.0.0.1:9999';
    constructor(newUrl) {
        if (!!newUrl) {
            hostUrl = newUrl;
        }
    }
    request(service, payload, success, failure) {
        // console.log('Tunnel.request', service, payload);
        if (!!window.jQuery) {
            $.ajax({
                url: service,
                type: 'GET',
                dataType: 'jsonp',
                crossDomain: true,
                headers: {
                    'cache-control': 'no-cache',
                },
                data: payload,
                success: success,
                error: failure
            });
        } else {
            throw 'Error: jQuery library not found';
        }
    }
    session(success, failure) {
        this.request(this.hostUrl + '/session', {},
            (response, status, options) => {
                // console.log('Tunnel.session.success', response);
                if (response.status == 'failure') {
                    failure(options, response.status, response.response)
                } else {
                    success(JSON.parse(response.response).token);
                }
            }, failure
        );
    }
    terminal(command, token, success, failure) {
        this.request(this.hostUrl + '/terminal', {
                cmd: command,
                token: token
            },
            (response, status, options) => {
                // console.log('Tunnel.terminal.success', response);
                if (response.status == 'failure') {
                    failure(options, response.status, response.response)
                } else {
                    success(JSON.parse(response.response));
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
        console.log('main.failure', status, error);
    };
    tunnel.session((token) => {
        console.log('main.session.success', token);
        tunnel.terminal('ls', token,
            (result) => {
                console.log('main.terminal.success', result);
            }, failure
        );
    }, failure);
    tunnel.terminal('pwd', null,
        (result) => {
            console.log('main.terminal.success', result);
        }, failure
    );
}());
*/