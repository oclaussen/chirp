# Chirp

Chirp is a simple tool that exposes your system keyboard over networ or unix socket connections.

## Usage

```bash
$ chirp --help
Access system clipboard over network

Usage:
  chirp [command]

Available Commands:
  copy        Send a copy request to the clipboard
  help        Help about any command
  paste       Send a paste request to the clipboard
  server      Start in server mode and wait for incoming clipboard requests
  service     Manage the chirp server as a system service

Flags:
  -a, --address string         address to bind to
  -h, --help                   help for chirp
      --tls-ca-file string     tls certificate authority file
      --tls-cert-file string   tls certificate file
      --tls-key-file string    tls key file

Use "chirp [command] --help" for more information about a command.
```

First, run the chirp server somewhere with a clipboard (e.g. a command like `xclip` or `pbcopy` must be available). For example, run a chirp server and listen to a unix socket:

```bash
$ chirp server --address unix:///var/run/chirp.sock
```

Then you can copy and paste from anywhere like this:

```bash
$ echo "Hello World!" | chirp copy --address unix:///var/run/chirp.sock
$ chirp paste --address unix:///var/run/chirp.sock
Hello World!
```

## Running as a LaunchDaemon (OSX)

Chirp supports running the server as a LaunchDaemon. Simply run the following to
set up the service files and start it:

```bash
$ chirp service install <flags>
$ chirp service start
```

## Current limitations

 * Chirp does not maintain any configuration or state. You have to pass exactly the same arguments to `copy` and `paste` commands that you used to run the server.
 * Works on Linux and OSX. Windows is not (yet) supported because I simply don't have a test setup.
 * Running as system service works only out of the box on OSX. Running the install commands with systemd will create the correct services, but not start them as the correct user. I.e. it will not access the desired clipboard.

## License & Authors

```text
Copyright 2021 Ole Claussen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
