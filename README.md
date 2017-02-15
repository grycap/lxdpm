# lxdpm
LXD Platform Manager

### Requisites
* The key and crt of the client for the host you want to control
* The server crt from the host you want to control


### Instructions to get it running
1. Install go and set the $GOPATH.
2. Git clone the project in $GOPATH/src
3. Execute 'go get' in the folder to get all the dependencies.
4. Execute 'go build -o lxdpm'
5. Execute './lxdpm args' specifiying the arguments. You can do './lxdpm --' to get help.

The server then will start, on the specified port. Access the several routes to get responses from the host.
You should also get the certs (client/server) from the host by now, to pass them to the client, and set the source client as trusted(using the lxc client in the host or by curl) in order to be able to use methods that require authorization.
