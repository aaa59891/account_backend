# README #

This README would normally document whatever steps are necessary to get your application up and running.

### How do I get set up? ###
 - use **go get**
    
   >$ go get github.com/aaa59891/account_backend/src
    
 - use **git clone & govendor**
   
   >$ git clone https://ChongChen@github.com/aaa59891/account_backend/src.git  
   >$ cd go_empty_gin  
   >$ govendor sync
   
### Configuration setting: ###

>Using environment variable: **GO_ENV** to set it
  
>This repo has three env:  
>1. **GO_ENV=prod**  
>2. **GO_ENV=dev**  
>3. **GO_ENV=test**  

>you can set all share config in **/configs/default.toml**,

>then set the different fields in **/configs/prod.toml** OR **/configs/dev.toml** OR **/configs/test.toml** 
  
### Database ###

>This repo uses **MySQL**, version: **at least 5.7**
  
>Can set MySQL configuration in **any env** you want
  
>database setting example:
>>[database]
>>>driveName = "mysql"  
>>>host = "localhost"  
>>>port = 5400  
>>>account = "account"  
>>>password = "password"  
>>>databaseName = "dbName" 

>**Need to create new init script**
>**There are SQL init scripts in /sqlScripts/initScript, you can create/insert into your MySQL**


>**Or comment routers/routes.go line 19 to disable authorization**

### How do I run this repo? ###

Install auto-reload tool first:
```
$ go get github.com/codegangsta/gin
```

Three run commands in different env:

**before testing, have to run all the scripts in the sqlScripts/initScripts and test_init.sql in your testing database **

* Testing: 
```
$ GO_ENV=test go test -v -cover ./controllers/
```
* Stage env:
```
$ GO_ENV=dev gin -i go run main.go
```
* Prod  env:
```
$ GO_ENV=prod gin -i go run main.go
```

### Docker Script Example ###
>MySQL
```
$ docker container run --name mysql_mosi -d -e MYSQL_USER=account -e MYSQL_PASSWORD=password -e MYSQL_DATABASE=dbName -p 5400:3306 --restart always -v /your/path/to/store:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=rootPassword mysql
```

### Build ###
Please check build.sh is executable, if not run the command:
```
$ chmod +x build.sh
```
Set the GOOS and GOARCH env depend on what os you want to deploy  
GOOS and GOARCH reference:
>https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

Ubuntu example:
```
$ GOOS=linux GOARCH=amd64 ./build.sh
```

Windows example:
```
$ GOOS=windows GOARCH=386 ./build.sh
```

Then you can find your deployment files in export/*ProjectName*/
Just move this folder to the machine and run the binary file!