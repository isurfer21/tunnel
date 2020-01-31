# tunnel
It establishes a bridge between system and web for performing task on system via web.


## Setup
Please follow the below instruction to install the `tunnel` on target system.

#### Install latest version
With *Shell*:

```bash
curl -fsSL https://isurfer21.github.io/tunnel/install.sh | sh
```

With *PowerShell*:

```powershell
iwr https://isurfer21.github.io/tunnel/install.ps1 -useb | iex
```

#### Install specific version
With *Shell*:

```bash
curl -fsSL https://isurfer21.github.io/tunnel/install.sh | sh -s v2.0.0
```

With *PowerShell*:

```powershell
iwr https://isurfer21.github.io/tunnel/install.ps1 -useb -outf install.ps1; .\install.ps1 v2.0.0
```

#### Manual setup

1. Download the app from [releases](https://github.com/isurfer21/tunnel/releases) section of the [repository](https://github.com/isurfer21/tunnel).
2. Access it from anywhere by adding the app path in the system path. 


## Usage
Below mentioned switches can be shuffled based on requirement, so explore it via help

    tunnel -h

#### Default or local
Accessible on host system at default port, i.e., 9999; without credentials so it uses default credentials

    tunnel

#### Authenticative
Accessible at default host & port, i.e., localhost:9999 but APIs will become accessible only using credentials

    tunnel -i admin -c 123456

#### Custom host & port
Accessible on all systems over network at custom port

    tunnel -u 192.168.0.100 -p 8080

#### Serve cross domain request
Serve cross-site incoming request from host system only

    tunnel -x

Serve cross-site request from any systems

    tunnel -x -u 192.168.0.100 -p 8080

#### Custom host & port
Serve along with static site at `./` location

    tunnel -d

Serve along with static site at `./app` location

    tunnel -d=./app


## API 
The API Key can be generated using `tunnel.genApiKey(username, password)` that can be kept separate, so as to avoid disclosing the credentials.

#### Example: Basic

```js
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
```


## Development

#### Compile/Run
To directly compile and run the source code, 

    cd tunnel
    go run tunnel.go 

it will start the tunnel

#### Build
To build the executable from source code,

    cd tunnel
    go build tunnel.go

it will create the tunnel executable file for current OS

#### Cross-build
To build the cross-platform executable from source code,

on **macOS**,

    cd tunnel
    sh build.sh

on **Windows**,

    cd tunnel
    build.bat

it will generate the tunnel executable files for all supported OS


## Testing
The APIs can be tested using python3 on command-line or using an standalone html page at browser that is being served via default embedded static web server in the app.

#### Python

    tunnel
    python3 ./test/test.py 

#### HTML+JS

    tunnel
    open http://127.0.0.1:9999/test/demo.html


## References

1. [isurfer21/Suxm](https://github.com/isurfer21/Suxm)