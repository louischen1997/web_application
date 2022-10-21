+ ## Import and modules
1. go.mod,go.sum
2. go gin framework
3. go gorm -> for database

+ ## Directory
1. config
   - database contrl function
   - connect to database MYSQL
   
2. src
   - router control
   - group to v1/account
   - call GET,PUT,POST etc...

3. pojo
   - create User structure
   - functions of getting or posting data from DB
   
4. call
   - api function control
   - implement HandleBA(basic Auth)function 
   - implement handle Bcrypt functions

+ ## Main 

1. GET healthz api
2. call function connect to DB
3. run on loaclhost3000
4. go run main.go
5. test all apis on Postman

+ ## Github actions
1. PR on go.yml
2. set branch protection on this action workflow
 
## Updated README.md 
