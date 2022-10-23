# Go FTPd

* A fully functional FTP `client` and `server` module written in Golang.
* Currently supported commands are as follows:
  - `cd`: Changes directory in remote context
  - `pwd`: Prints current remote directory
  - `ls`: lists all Files/directories in current remote directory
  - `cat`: display contents of a particular file
  - `del`: delete a remote file
  - `get`: download a remote file
  - `put`: upload a local file

## Usage:

- Clone the repository with:
```bash
git clone https://github.com/The-Debarghya/goftpd
```
- Move to the `server` directory and compile with:
`go build .`
- Add username and encrypted-password hashes(AES-256) along with the key in `creds.json` file.
- Make two directories named `ftp` and `ftpfiles` in the directory above `server` or `client`.(These 2 directories actually work as work directories for server and client respectively.)
- Start the server first.
- Next switch to another console tab and move to `client` directory and compile with:
`go build .`
- Now run the client.
- The above mentioned steps works fine for `localhost` connections. To setup server to accept connections over same network, change the **HOST** constant to the required IP address in `client.go` file:
```go
const (
	HOST = "localhost"
	PORT = "9999"
	BUFFSIZE = 5120
)
```
