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

        (!!payload) ? xhr.send(this.toURLQuery(payload)): xhr.send();
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

export function tnl(argv, args) {
    if ((argv.h || argv.help) && !argv._.length)  {
        return `
Establishes a bridge between system and web for performing task on system via web

Syntax:
 tnl [cmd] [options] 

 where,
  [cmd]   system commands

Options:
 -h --help            show help options
 -u --user [=admin]   insert username
 -p --pass [=*****]   insert password

Examples:
 tnl -u=admin -p=123456
 tnl ls
        `;
    } else if (('u' in argv || 'user' in argv) && ('p' in argv || 'pass' in argv) && !argv._.length) {
        let u = 'u' in argv ? argv.u : 'user' in argv ? argv.user : 'admin';
        let p = 'p' in argv ? argv.p : 'pass' in argv ? argv.pass : '123456';
        tunnel.login(u, p);
        Shell.exit();
    } else if (!!argv._.length) {
        tunnel.authenticate((result) => {
            tunnel.terminal(args.join(' '), (result) => {
                result = JSON.parse(result);
                if (!!result.err) {
                    Shell.write(result.err);
                    Shell.exit();
                } else if (!!result.out) {
                    Shell.write(result.out);
                    Shell.exit();
                }
            }, failure);
        }, failure);
    } else {
        return 'Missing command/options parameters';
    }
}

const tunnel = new Tunnel();
const failure = (options, status, error) => {
    Shell.write(`${status}: ${error}`);
    Shell.exit();
};