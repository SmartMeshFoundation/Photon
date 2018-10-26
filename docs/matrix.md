# matrix-regservice
Matrix-regservice is used as a third-party application service for Matrix. The main function is to limit the Matrix HomeServer used by SmartRaiden to only accept valid Raiden user registration.
That is, users registered on this Matrix server can only be SmartRaiden nodes. These nodes have a uniform name and can guarantee that these users must have the corresponding private key to register.


## Matrix HomeServer Deployment
Matrix installation and deployment reference[matrix](https://github.com/matrix-org/synapse).

### Modify homeserver.yaml

1. `enable_registration` is changed to False, ensuring that users cannot register through the normal interface, and can only register users through third-party application service.
2. `search_all_users` is changed to True to ensure that users can be retrieved
3. `expire_access_token` is changed to True to ensure that the user will automatically log out to prevent third-party replay attacks.
4. `port` 8008 is modified to 8007
5. Remove the `webclient` under port 8007 and disable login via web.
6. `app_service_config_files` is modified to [ registration.yaml]
7. `trusted_third_party_id_servers` are all deleted
## matrix-regservice install and deployment

### install matrix-regservice
```bash
go get github.com/SmartMeshFoundation/matrix-regservice
```
then copy `matrix-regservice` to $PATH

### Generate configuration file
Switch to the matrix working directory (homeserver.yaml directory)
```bash
matrix-regservice --matrixdomain yourdomain.com genconfig
```
registration.yaml and run.sh are generated in the matrix working directory.
#### a registration.yaml sample
```yaml
id: Q7PM2E53RE-transport02.smartmesh.cn
hs_token: RNI4CGEDTKC4WJTB4RZWRK4NOKA7M4PREUW6F2GZ
as_token: LODE52N2CKVXMOURUMAWLEEEXMWB4DKIKRI246XD
url: http://127.0.0.1:8009/regapp/1
sender_localpart: transport02.smartmesh.cn
protocols:
- regapp.transport02.smartmesh.cn
namespaces:
  users:
  - exclusive: false
    regex: '@.*'
  aliases: []
  rooms: []
```
run.sh
```bash
matrix-regservice --astoken LODE52N2CKVXMOURUMAWLEEEXMWB4DKIKRI246XD --hstoken RNI4CGEDTKC4WJTB4RZWRK4NOKA7M4PREUW6F2GZ --matrixurl http://127.0.0.1:8008/_matrix/client/api/v1/createUser --host 127.0.0.1 --port 8009 --datapath .matrix --matrixdomain transport02.smartmesh.cn --verbosity 5
```

### nginx reverse-proxy configuration
the following is a example configuration file for nginx.
```conf
  server {
    listen 8008;
    server_name localhost;

  location /_matrix {
    proxy_pass http://127.0.0.1:8007;
    proxy_max_temp_file_size 0;
    proxy_connect_timeout 30;
    }

  location /regapp/1 {
    proxy_pass http://127.0.0.1:8009;
    proxy_max_temp_file_size 0;
    proxy_connect_timeout 30;
    }
  }
```