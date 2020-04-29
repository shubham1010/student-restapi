# This is a students restapi built using golang as backend and mariadb as database.

## APIs

### http://localhost:8080/students => GET
### http://localhost:8080/student/id => GET
### http://localhost:8080/student/id => POST
### http://localhost:8080/student/id => UPDATE
### http://localhost:8080/student/id => DELETE


## POSTMAN INSTRUCTIONS:
### FOR INSERTION OF DATA SEND JSON FORMAT AS BELOW => POST METHOD

{"id": "18143", "name": "Satish Sangam","age": "22", "dept": "Comp Scie", "subject": ["Formal Method-II","DevOps"]}

### FOR UPDATION OF DATA SEND JSON FORMAT AS BELOW===> UPDATE METHOD
{"id": "existing id", "name": "newname","age": "23", "dept": "updated value", "subject": ["Formal Method-II","DevOps","newsubject"]}


### HOW TO RUN PROJECT
##### 1) You just clone or download the project into you host machine.
##### 2) Change the mode of run.sh file by using command ===> chmod +x run.sh
##### 3) Run command ./run.sh
##### 4) When you get one prompt then just run ./startup.sh
