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
1. PR check on go.yml
2. set branch protection on this action workflow
 
+ ## Create group and users
1. create group csye6225-ta 
2. create DEVGroup 
3. add ReadOnly policy,Administrator policy 

+ ## Configure
1. aws configure
   - configure 3 users
   - dev, demo , root for profile names
   - export AWS_PROFILE=demo

+ ## GitHub Organization  

1. create repo infrastructure
2. fork to local


+ ## AWS CLI Networking setup
1. export AWS_PROFILE=dev
2. aws configure list
3. aws ec2 create-vpc --cidr-block "10.0.0.0/16" --no-amazon-provided-ipv6-cidr-block --instance-tenancy default
4. aws ec2 create-subnet --availability-zone us-west-2c --vpc-id vpc-0e4a47c0e5bc09654 --cidr-block 10.0.16.0/20
5. aws ec2 create-internet-gateway
6. aws ec2 attach-internet-gateway --vpc-id vpc-0e4a47c0e5bc09654 --internet-gateway-id igw-087d869e093dcb388
7. aws ec2 create-route-table --vpc-id vpc-0e4a47c0e5bc09654
8. aws ec2 associate-route-table --route-table-id rtb-05df52a0a9963e38f --subnet-id subnet-01eb9cba2d6bf4c57 
9. aws ec2 create-route --route-table-id rtb-05df52a0a9963e38f --destination-cidr-block 0.0.0.0/0 --gateway-id igw-087d869e093dcb388
 
+ ## Infrastructure as Code with CloudFormation
1. csye6225-infra.yml 
2. replace hardcode to variables
3. aws cloudformation create-stack --stack-name myvpc --template-body file://csye6225-infra.yml --parameters ParameterKey=VpcCidrBlock,ParameterValue="10.1.1.0/24" ParameterKey=Subnet1ZN, ParameterValue="us-east-1a" --region us-west-2 

+ ## Build packer
1. with us-west-2 ami ID
2. use default vpc
3. file provisioner after building artifacts
4. shell provisioner to install Golang, mysql (initial),systemd
5. add ami_users to share ami to Demo account
6. initialize db with password, new user
7. add privilege to new user
8. add structure to auto migrate db table

+ ## Systemd
1. add app.service and cp to right location
2. sudo systemctl daemon-reload
3. sudo systemctl enable app.service
4. sudo systemctl start app.service
5. GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o systemd-test  main.go


+ ## Create resources with cloud formation
1. private subnets
2. private route table
3. private route
4. new security group for RDS instance -> database
5. s3 bucket
    - use stackid for bucket name (random) 
    - use deafault encryption
    - set lifecycle policy
6. rds parameter group -> mysql8.0
7. rds instance
    - set username password as parameters
    - hostname(db endpoint) pass to env file with user data
8. user data
   -  pass hostname, username, password , bucket name to env file in ec2 instance

9. IAM policy & IAM ROLE
    - attach S3 put and delete policy
   
     

+ ## Webapp
1. write document apis
2. new table for docs in mysql DB, storing docs metadata
3. able to upload and delete files to S3 bucket 
4. S3 default crendentials passing by aws ec2 (setting IAM role)


+ ## DNS setup
1. Get a domain form Namecheap louisdomain6225.me

+ ## Create public hostzone
1. confugure name servers on root account
2. configure name servers also in root account 


+ ## cloudformation
1. add route S3 recordset
2. set A record for prod subdomain
3. resource record for ec2 public IP

 
+ ## SES
1. request production access
2. Easy DKIM

+ ## AMI update
1. install cloudwatch Agent with scripts in packer file 
2. attach CloudWatchServerAgent in cloudformation

+ ## User Stories
1. create log file 
2. generating logs whenever api calls
3. including endpoint, request type, success or not, time&date
4. create counter metrics whenever api calls with statsd library golang support
5. copy cloudwatch-config into EC2 instance
6. run config scripts in user data 


+ ## Set Cloud Formation
1. SNS topic
2. lambda function with go(via S3 bucket)
3. lambda role with permission policy
4. dynamodb, Attribute: EMAIL_KEY,TOKEN,TTL: TTLATT

+ ## Lambda Function Set Up
1. lambda handler with golang
2. new package using golang 
3. triggered by SNS topic(with permissions)
4. sending EMAIL to verify user account
5. including email address, token, domain, sendgridkey

+ ## WebApp api
1. verify link: get request getting items form dynamodb to confirm token validation
2. the token expired after 300 seconds
3. set user account new variable: account verify
4. block all auth apis if not verified
5. create approproate logs in cloudwatch
 


 
