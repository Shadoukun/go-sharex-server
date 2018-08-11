# go-sharex-server
A small server to use for for file and text upload and file sharing with ShareX


## Installation

### Download
 1. Go download the newest release from [Releases](https://github.com/Shadoukun/go-sharex-server/releases)
 2. Copy `config.yml.example` to `config.yml` and edit it.
 3. Run the executable from the command line
 
 ### Compile From Source

1. Clone this repo. \
```git clone https://github.com/Shadoukun/go-sharex-server.git```


2. Inside the folder, use `go get` to download all the dependencies. \
```go get -d ./...```

3. run `go build`. which will create the executable in the current directory. \
```go build .```

3. Copy `config.yml.example` to `config.yml` and edit it.

4. Run the executable from the command line

## Configuration Options
The server uses a configuration file. `config.yml`

```yaml
host: "127.0.0.1"
port: "80"

# Permissions
admin:
    password: "test"
    maxupload: 20
    filetypes: [.png, .jpg, .gif, .mp4, .mp3, .jpeg, .tiff,
                .bmp, .ico, .psd, .eps, .raw, .cr2, .nef,
                .sr2, .orf, .svg, .wav, .webm, .aac, .flac,
                .ogg, .wma, .m4a, .gifv]
default:
    password: "test2"
    maxupload: 10
    filetypes: [.png, .jpg, .gif, .mp4, .mp3, .jpeg, .tiff,
                .bmp, .ico, .psd, .eps, .raw, .cr2, .nef,
                .sr2, .orf, .svg, .wav, .webm, .aac, .flac,
                .ogg, .wma, .m4a, .gifv]
```

The top contains currently supported server options:

* **Host** - The domain name or IP address the server should run on. `Default: 127.0.0.1`
* **Port** - The port the server should run on. `Default: 80`

Followed by permission groups. There is a default group, and an admin group. \
Each group contains the following options:

* **Password** - The password to allow uploading for this group.
* **MaxUpload** - The maximum upload size in megabytes.
* **Filetypes** - a list of the file types this group is allowed to upload.

\
You can also use environment variables to set configuration options. \
If any variable is set it will override the corresponding option in the configuration file.

* **SHAREX_HOST** 
* **SHAREX_PORT**

* **SHAREX_PASS**
* **SHAREX_MAXUPLOAD**
* **SHAREX_FILETYPES** 

* **SHAREX_PASS_ADMIN**
* **SHAREX_MAXUPLOAD_ADMIN**
* **SHAREX_FILETYPES_ADMIN** 

**Note:** At the moment, If you use environment variables to set permissions and want to use both groups, you must explicitly set both in the environment. Otherwise it will default to one group.

For example:

If you set **SHAREX_PASS_ADMIN** but do not set **SHAREX_PASS**, then **SHAREX_PASS_ADMIN** will be used for all uploads, and vice versa.


## Configure ShareX

1. Open up **Destination Settings** \
\
![](https://i.imgur.com/t5HVRla.png)

2. Go to **Custom Uploaders** \
\
![](https://i.imgur.com/C2ZxDSF.png)

3. On the **Custom Uploaders** menu, do the following: \
\
![](https://i.imgur.com/e9CkFET.png)
    - Hit **Add**.
    -  Fill in the fields on the left.
        * **Name** - The display name for your uploader.
        * **Destination Type** - *File Uploader* for files and *Text Uploader* for text/pastes
        * **Request Type** - `POST`
        * **Request URL** - The upload URL. the URLs will be:
            ```
            File: http://<your host>/file/upload
            Text: http://<your host>/paste/upload
            ```
        * **File Form Name** - This should be set to `file`
    - Under *Arguments*: create a `pass` argument with the password for one of the upload groups you want to use. \
      and hit *Add*

4. In the box to the right, do the following: \
\
![](https://i.imgur.com/HgeB4nS.png)
      - **Response Type** - Set to `Response Text`
      - **URL** - `$json:Result.URL$`

 5. Hit the File Uploader/Text Uploader *Test* button, and hopefully it spits out a test URL.
