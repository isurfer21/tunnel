# tunnel
It establishes a bridge between system and web for performing task on system via web.


### For end-user
Download the app from [releases](https://github.com/isurfer21/tunnel/releases) section of the [repository](https://github.com/isurfer21/tunnel).

#### Default or local
Accessible on host system at default port, i.e., 9999; without credentials so it uses default credentials

```bash
$ tunnel
```

#### Authenticative
Accessible at default host & port, i.e., localhost:9999 but APIs will become accessible only using credentials

```bash
$ tunnel -i admin -c 123456
```

#### Custom host & port
Accessible on all systems over network at custom port

```bash
$ tunnel -u 192.168.0.100 -p 8080
```


### For developer

#### Compile/Run
To directly compile and run the source code, 

```bash
$ cd tunnel
$ go run tunnel.go 
```

it will start the tunnel

#### Build
To build the executable from source code,

```bash
$ cd tunnel
$ go build tunnel.go
```

it will create the tunnel executable file for current OS

#### Cross-build
To build the cross-platform executable from source code,

on **macOS**,

```bash
$ cd tunnel
$ sh build.sh
```

on **Windows**,

```bash
$ cd tunnel
$ build.bat
```

it will generate the tunnel executable files for all supported OS



#### For tester
The APIs can be tested using python3 on command-line or using an standalone html page at browser that is being served via default embedded static web server in the app.

##### Python
```bash
$ tunnel
$ python3 ./test/test.py 
```

##### HTML+JS
```bash
$ tunnel
$ open http://127.0.0.1:9999/test/demo.html
```


#### Using API 
The API Key can be generated using `tunnel.genApiKey(username, password)` that can be kept separate, so as to avoid disclosing the credentials.

##### Example: Basic

```js
(function main() {
    console.log('main');
    let tunnel = new Tunnel();
    let failure = (options, status, error) => {
        console.log('main.failure', status, error);
    };
    let apiKey = tunnel.genApiKey('admin', '123456');
    tunnel.session(apiKey,
        (token) => {
            console.log('main.session.success', token);
            tunnel.terminal('ls', token,
                (result) => {
                    console.log('main.terminal.success', result);
                }, failure
            );
        }, failure
    );
}());
```


### References

1. [isurfer21/Suxm](https://github.com/isurfer21/Suxm)
2. [blueimp/JavaScript-MD5](https://github.com/blueimp/JavaScript-MD5)