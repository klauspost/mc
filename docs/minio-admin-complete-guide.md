# MinIO Admin Complete Guide [![Slack](https://slack.min.io/slack?type=svg)](https://slack.min.io)

MinIO Client (mc) provides `admin` sub-command to perform administrative tasks on your MinIO deployments.

```
service  restart and stop all MinIO servers
update   update all MinIO servers
info     display MinIO server information
user     manage users
group    manage groups
policy   manage policies defined in the MinIO server
config   manage configuration file
heal     heal disks, buckets and objects on MinIO server
profile  generate profile data for debugging purposes
top      provide top like statistics for MinIO
trace    show http trace for minio server
console  show console logs for MinIO server
```

## 1.  Download MinIO Client
### Docker Stable
```
docker pull minio/mc
docker run minio/mc admin info server play
```

### Docker Edge
```
docker pull minio/mc:edge
docker run minio/mc:edge admin info server play
```

### Homebrew (macOS)
Install mc packages using [Homebrew](http://brew.sh/)

```
brew install minio/stable/mc
mc --help
```

### Binary Download (GNU/Linux)
| Platform | Architecture | URL |
| ---------- | -------- |------|
|GNU/Linux|64-bit Intel|https://dl.min.io/client/mc/release/linux-amd64/mc |
||64-bit PPC|https://dl.min.io/client/mc/release/linux-ppc64le/mc |

```
chmod +x mc
./mc --help
```

### Binary Download (Microsoft Windows)
| Platform | Architecture | URL |
| ---------- | -------- |------|
|Microsoft Windows|64-bit Intel|https://dl.min.io/client/mc/release/windows-amd64/mc.exe |

```
mc.exe --help
```

### Install from Source
Source installation is intended only for developers and advanced users. `mc update` command does not support update notifications for source based installations. Please download official releases from https://min.io/download/#minio-client.

If you do not have a working Golang environment, please follow [How to install Golang](https://golang.org/doc/install).

```
go get -d github.com/minio/mc
cd ${GOPATH}/src/github.com/minio/mc
make
```

## 2. Run MinIO Client

### GNU/Linux

```
chmod +x mc
./mc --help
```

### macOS

```
chmod 755 mc
./mc --help
```

### Microsoft Windows

```
mc.exe --help
```

## 3. Add a MinIO Storage Service
MinIO server displays URL, access and secret keys.

#### Usage

```
mc config host add <ALIAS> <YOUR-MINIO-ENDPOINT> <YOUR-ACCESS-KEY> <YOUR-SECRET-KEY>
```

Alias is simply a short name to your MinIO service. MinIO end-point, access and secret keys are supplied by your MinIO service. Admin API uses "S3v4" signature and cannot be changed.

```
mc config host add minio http://192.168.1.51:9000 BKIKJAA5BMMU2RHO6IBB V7f1CwQqAcwo80UEIJEjc5gVQUSSx5ohQ9GSrr12
```

## 4. Test Your Setup

*Example:*

Get MinIO server information for the configured alias `minio`

```
mc admin info server minio

●  192.168.1.51:9000
   Uptime : online since 1 day ago
  Version : 2018-05-28T04:31:38Z
   Region :
 SQS ARNs : <none>
    Stats : Incoming 82GiB, Outgoing 28GiB
  Storage : Used 7.4GiB

  CPU        min        avg      max
  current    0.12%      0.14%    0.17%
  historic   0.03%      1.27%    412.18%

  MEM        usage
  current    602 MiB
  historic   448 MiB
```

## 5. Everyday Use
You may add shell aliases for info, healing.

```
alias minfo='mc admin info'
alias mheal='mc admin heal'
```

## 6. Global Options

### Option [--debug]
Debug option enables debug output to console.

*Example: Display verbose debug output for `info` command.*

```
mc admin --debug info server minio
mc: <DEBUG> GET /minio/admin/v1/info HTTP/1.1
Host: 192.168.1.51:9000
User-Agent: MinIO (linux; amd64) madmin-go/0.0.1 mc/2018-05-23T23:43:34Z
Authorization: AWS4-HMAC-SHA256 Credential=**REDACTED**/20180530/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=**REDACTED**
X-Amz-Content-Sha256: UNSIGNED-PAYLOAD
X-Amz-Date: 20180530T001808Z
Accept-Encoding: gzip

mc: <DEBUG> HTTP/1.1 200 OK
Transfer-Encoding: chunked
Accept-Ranges: bytes
Content-Security-Policy: block-all-mixed-content
Content-Type: application/json
Date: Wed, 30 May 2018 00:18:08 GMT
Server: MinIO/DEVELOPMENT.2018-05-28T04-31-38Z (linux; amd64)
Vary: Origin
X-Amz-Request-Id: 1533440573A63034
X-Xss-Protection: "1; mode=block"

mc: <DEBUG> Response Time:  140.70112ms

●  192.168.1.51:9000
   Uptime : online since 1 day ago
  Version : 2018-05-28T04:31:38Z
   Region :
 SQS ARNs : <none>
    Stats : Incoming 82GiB, Outgoing 28GiB
  Storage : Used 7.4GiB

  CPU        min        avg      max
  current    0.12%      0.14%    0.17%
  historic   0.03%      1.27%    412.18%

  MEM        usage
  current    602 MiB
  historic   448 MiB
```

### Option [--json]
JSON option enables parseable output in [JSON lines](http://jsonlines.org/) format.

*Example: MinIO server information.*

```
mc admin --json info server minio
{
  "status": "success",
  "service": "on",
  "address": "192.168.1.51:9000",
  "error": "",
  "storage": {
    "used": 7979370172,
    "backend": {
      "backendType": "FS"
    }
  },
  "network": {
    "transferred": 90473434722,
    "received": 30662519192
  },
  "server": {
    "uptime": 157467244813288,
    "version": "2018-05-28T04:31:38Z",
    "commitID": "7d8c5ffb13334f4aec20a35bd2575bd7c740fb7a",
    "region": "",
    "sqsARN": []
  }
  "cpu": {
  "addr": "play.minio.io:9000",
  "load": [
   {
    "avg": 0.11,
    "max": 0.13,
    "min": 0.09
   }
  ],
  "historicLoad": [
   {
    "avg": 1.270331023187104,
    "max": 412.18000000000006,
    "min": 0.03
   }
  ]
 },
 "mem": {
  "addr": "play.minio.io:9000",
  "usage": [
   {
    "mem": 631800056
   }
  ],
  "historicUsage": [
   {
    "mem": 469437357
   }
  ]
 }
}
```

### Option [--no-color]
This option disables the color theme. It is useful for dumb terminals.

### Option [--quiet]
Quiet option suppress chatty console output.

### Option [--config-dir]
Use this option to set a custom config path.

### Option [ --insecure]
Skip SSL certificate verification.

## 7. Commands

|                                                                        |
|:-----------------------------------------------------------------------|
| [**service** - restart and stop all MinIO servers](#service)           |
| [**update** - updates all MinIO servers](#update)                      |
| [**info** - display MinIO server information](#info)                   |
| [**user** - manage users](#user)                                       |
| [**group** - manage groups](#group)                                    |
| [**policy** - manage canned policies](#policy)                         |
| [**config** - manage server configuration file](#config)               |
| [**heal** - heal disks, buckets and objects on MinIO server](#heal)    |
| [**profile** - generate profile data for debugging purposes](#profile) |
| [**top** - provide top like statistics for MinIO](#top)                |
| [**trace** - show http trace for MinIO server](#trace)                 |
| [**console** - show console logs for MinIO server](#console)           |


<a name="update"></a>
### Command `update` - updates all MinIO servers
`update` command provides a way to update all MinIO servers in a cluster.

> NOTE:
> - An alias pointing to a distributed setup this command will automatically update all MinIO servers in the cluster.
> - `update` is an disruptive operations for your MinIO service, any on-going API operations will be forcibly canceled. So, it should be used only when you are planning MinIO upgrades for your deployment.

*Example: Update all MinIO servers.*
```
mc admin update play
Server `play` updated successfully from RELEASE.2019-08-14T20-49-49Z to RELEASE.2019-08-21T19-59-10Z
```

<a name="service"></a>
### Command `service` - restart and stop all MinIO servers
`service` command provides a way to restart and stop all MinIO servers.

> NOTE:
> - An alias pointing to a distributed setup this command will automatically execute the same actions across all servers.
> - `restart` and `stop` sub-commands are disruptive operations for your MinIO service, any on-going API operations will be forcibly canceled. So, it should be used only under administrative circumstances. Please use it with caution.

```
NAME:
  mc admin service - restart and stop all MinIO servers

FLAGS:
  --help, -h                       show help

COMMANDS:
  restart  restart all MinIO servers
  stop     stop all MinIO servers
```

*Example: Restart all MinIO servers.*
```
mc admin service restart play
Restarted `play` successfully.
```

<a name="info"></a>
### Command `info` - Display MinIO server information
`info` command displays server information of one or many MinIO servers (under distributed cluster)

```
NAME:
  mc admin info - get MinIO server information

FLAGS:
  --help, -h                       show help
```

*Example: Display MinIO server information.*

```
mc admin info play
●  play.min.io
   Uptime : online since 1 day ago
  Version : 2018-05-28T04:31:38Z
   Region :
 SQS ARNs : <none>
    Stats : Incoming 82GiB, Outgoing 28GiB
  Storage : Used 8.2GiB
```

<a name="policy"></a>
### Command `policy` - Manage canned policies
`policy` command to add, remove, list policies on MinIO server.

```
NAME:
  mc admin policy - manage policies

FLAGS:
  --help, -h                       show help

COMMANDS:
  add      add new policy
  remove   remove policy
  list     list all policies
  info     show info on a policy
  set      set IAM policy on a user or group
```

*Example: Add a new policy 'newpolicy' on MinIO, with policy from /tmp/newpolicy.json.*

```
mc admin policy add myminio/ newpolicy /tmp/newpolicy.json
```

*Example: Remove policy 'newpolicy' on MinIO.*

```
mc admin policy remove myminio/ newpolicy
```

*Example: List all policies on MinIO.*

```
mc admin policy list --json myminio/
{"status":"success","policy":"newpolicy"}
```

*Example: Show info on a policy*

```
mc admin policy info myminio/ writeonly
```

*Example: Set the policy on a user or group*

```
mc admin policy set myminio writeonly user=someuser
mc admin policy set myminio writeonly group=somegroup
```

<a name="user"></a>
### Command `user` - Manage users
`user` command to add, remove, enable, disable, list users on MinIO server.

```
NAME:
  mc admin user - manage users

FLAGS:
  --help, -h                       show help

COMMANDS:
  add      add new user
  disable  disable user
  enable   enable user
  remove   remove user
  list     list all users
  info     display info of a user
```

*Example: Add a new user 'newuser' on MinIO.*

```
mc admin user add myminio/ newuser newuser123
```

*Example: Disable a user 'newuser' on MinIO.*

```
mc admin user disable myminio/ newuser
```

*Example: Enable a user 'newuser' on MinIO.*

```
mc admin user enable myminio/ newuser
```

*Example: Remove user 'newuser' on MinIO.*

```
mc admin user remove myminio/ newuser
```

*Example: List all users on MinIO.*

```
mc admin user list --json myminio/
{"status":"success","accessKey":"newuser","userStatus":"enabled"}
```

*Example: Display info of a user*

```
mc admin user info myminio someuser
```

<a name="group"></a>
### Command `group` - Manage groups
`group` command to add, remove, info, list, enable, disable groups on MinIO server.

```
NAME:
  mc admin group - manage groups

USAGE:
  mc admin group COMMAND [COMMAND FLAGS | -h] [ARGUMENTS...]

COMMANDS:
  add      add users to a new or existing group
  remove   remove group or members from a group
  info     display group info
  list     display list of groups
  enable   Enable a group
  disable  Disable a group
```

*Example: Add a pair of users to a group 'somegroup' on MinIO.*

Group is created if it does not exist.

```
mc admin group add myminio somegroup someuser1 someuser2
```

*Example: Remove a pair of users from a group 'somegroup' on MinIO.*

```
mc admin group remove myminio somegroup someuser1 someuser2
```

*Example: Remove a group 'somegroup' on MinIO.*

Only works if the given group is empty.

```
mc admin group remove myminio somegroup
```

*Example: Get info on a group 'somegroup' on MinIO.*

```
mc admin group info myminio somegroup
```

*Example: List all groups on MinIO.*

```
mc admin group list myminio
```

*Example: Enable a group 'somegroup' on MinIO.*

```
mc admin group enable myminio somegroup
```

*Example: Disable a group 'somegroup' on MinIO.*

```
mc admin group disable myminio somegroup
```

<a name="config"></a>
### Command `config` - Manage server configuration
`config` command to manage MinIO server configuration.

```
NAME:
  mc admin config - manage configuration file

USAGE:
  mc admin config COMMAND [COMMAND FLAGS | -h] [ARGUMENTS...]

COMMANDS:
  get     get config of a MinIO server/cluster.
  set     set new config file to a MinIO server/cluster.

FLAGS:
  --help, -h                       Show help.
```

*Example: Get server configuration of a MinIO server/cluster.*

```
mc admin config get myminio > /tmp/my-serverconfig
```

*Example: Set server configuration of a MinIO server/cluster.*

```
mc admin config set myminio < /tmp/my-serverconfig
```

<a name="heal"></a>
### Command `heal` - Heal disks, buckets and objects on MinIO server
`heal` command heals disks, missing buckets, objects on MinIO server. NOTE: This command is only applicable for MinIO erasure coded setup (standalone and distributed).

```
NAME:
  mc admin heal - heal disks, buckets and objects on MinIO server

FLAGS:
  --scan value                     select the healing scan mode (normal/deep) (default: "normal")
  --recursive, -r                  heal recursively
  --dry-run, -n                    only inspect data, but do not mutate
  --force-start, -f                force start a new heal sequence
  --force-stop, -s                 force stop a running heal sequence
  --remove                         remove dangling objects in heal sequence
  --help, -h                       show help
```

*Example: Heal MinIO cluster after replacing a fresh disk, recursively heal all buckets and objects, where 'myminio' is the MinIO server alias.*

```
mc admin heal -r myminio
```

*Example: Heal MinIO cluster on a specific bucket recursively, where 'myminio' is the MinIO server alias.*

```
mc admin heal -r myminio/mybucket
```

*Example: Heal MinIO cluster on a specific object prefix recursively, where 'myminio' is the MinIO server alias.*

```
mc admin heal -r myminio/mybucket/myobjectprefix
```

<a name="profile"></a>
### Command `profile` - generate profile data for debugging purposes

```
NAME:
  mc admin profile - generate profile data for debugging purposes

COMMANDS:
  start  start recording profile data
  stop   stop and download profile data
```

Start CPU profiling
```
mc admin profile start --type cpu myminio/
```

<a name="top"></a>
### Command `top` - provide top like statistics for MinIO
NOTE: This command is only applicable for a distributed MinIO setup. It is not supported on single node and gateway deployments.

```
NAME:
  mc admin top - provide top like statistics for MinIO

COMMANDS:
  locks  Get a list of the 10 oldest locks on a MinIO cluster.
```

*Example: Get a list of the 10 oldest locks on a distributed MinIO cluster, where 'myminio' is the MinIO cluster alias.*

```
mc admin top locks myminio
```

<a name="trace"></a>
### Command `trace` - Show http trace for MinIO server
`trace` command displays server http trace of one or all MinIO servers (under distributed cluster)

```sh
NAME:
  mc admin trace - show http trace for MinIO server

FLAGS:
  --verbose, -v                 print verbose trace
  --all, -a                     trace all traffic (including internode traffic between MinIO servers)
  --errors, -e                  trace failed requests only
  --help, -h                    show help
```

*Example: Display MinIO server http trace.*

```sh
mc admin trace myminio
172.16.238.1 [REQUEST (objectAPIHandlers).ListBucketsHandler-fm] [154828542.525557] [2019-01-23 23:17:05 +0000]
172.16.238.1 GET /
172.16.238.1 Host: 172.16.238.3:9000
172.16.238.1 X-Amz-Date: 20190123T231705Z
172.16.238.1 Authorization: AWS4-HMAC-SHA256 Credential=minio/20190123/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=8385097f264efaf1b71a9b56514b8166bb0a03af8552f83e2658f877776c46b3
172.16.238.1 User-Agent: MinIO (linux; amd64) minio-go/v6.0.8 mc/2019-01-23T23:15:38Z
172.16.238.1 X-Amz-Content-Sha256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
172.16.238.1
172.16.238.1 <BODY>
172.16.238.1 [RESPONSE] [154828542.525557] [2019-01-23 23:17:05 +0000]
172.16.238.1 200 OK
172.16.238.1 X-Amz-Request-Id: 157C9D641F42E547
172.16.238.1 X-Minio-Deployment-Id: 5f20fd91-6880-455f-a26d-07804b6821ca
172.16.238.1 X-Xss-Protection: 1; mode=block
172.16.238.1 Accept-Ranges: bytes
172.16.238.1 Content-Security-Policy: block-all-mixed-content
172.16.238.1 Content-Type: application/xml
172.16.238.1 Server: MinIO/RELEASE.2019-09-05T23-24-38Z
172.16.238.1 Vary: Origin
...
```

<a name="console"></a>
### Command `console` - show console logs for MinIO server
`console` command displays server logs of one or all MinIO servers (under distributed cluster)

```sh
NAME:
  mc admin console - show console logs for MinIO server

FLAGS:
  --limit value, -l value       show last n log entries (default: 10)
  --help, -h                    show help
```

*Example: Display MinIO server http trace.*

```sh
mc admin console myminio

 API: SYSTEM(bucket=images)
 Time: 22:48:06 PDT 09/05/2019
 DeploymentID: 6faeded5-5cf3-4133-8a37-07c5d500207c
 RequestID: <none>
 RemoteHost: <none>
 UserAgent: <none>
 Error: ARN 'arn:minio:sqs:us-east-1:1:webhook' not found
        4: cmd/notification.go:1189:cmd.readNotificationConfig()
        3: cmd/notification.go:780:cmd.(*NotificationSys).refresh()
        2: cmd/notification.go:815:cmd.(*NotificationSys).Init()
        1: cmd/server-main.go:375:cmd.serverMain()
```
