# Static analysis tools

### We have added the following two tools to our github such that both are run when creating a pull request.
- Sonarqube
- Code Climate

### We have two analysis tools embedded in our CI pipeline:
We run 'go vet' in our workflow in github actions, we have chosen this since it is a pretty standard analysis tool for checking Golang source code.
The command should help detect any suspicious or abnormal in the application. 

Another command we use is 'golint'. This one checks for code style, it should enforce the some commen coding conventions by Go developers. 
