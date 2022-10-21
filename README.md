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
 
+ ## Create group and users
1. create group csye6225-ta for TA
2. create DEVGroup for me
3. add ReadOnly policy,Administrator policy 

+ ## Configure
1. aws configure
   - configure 3 users(me)
   - dev, demo , root for profile names
   - export AWS_PROFILE=demo

+ ## GitHub Organization  

1. create repo infrastructure
2. fork to local
3. add README.md

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
 


 
