
# Risk Identification
When setting up an application that is open to everyone on the world wide web it is very important to set up some defences against persons that would be willing to attack your system. Do know what defences to setup, it’s important to make an analysis on what vulnerabilities could be in your application. In this report we will be identifying the parts of the application that could be targets to attacks, identify threat sources and construct risk scenarios. Then we will be prioritizing the risk scenarios based on their likelihood and impact.

## Identifiy assets (e.g. web application)

When considering what assets could be vulnerable to attacks almost all of the parts come in mind. Following list contains the parts that are considered as possible victims of an attacker.

-   Minitwit web application
-   minitwit API
-   MySQL Database
-   Travis(ci/cd)
-   Kibana
-   Prometheus
-   Grafana

## Identify threat sources

Following is a list of possible threats sources that could be used to attack the application.

- SQL injection
-  Cross-site scripting
- DDOS Attack
- Bad Authentication


## Construct risk scenarios 
We summed up possible risk scenarios for our application they are following

1.  An attacker could attempt to do some cross site scripting through the form when a message is being posted to the application to try to retrieve cookies for some other user.
    
2.  An attacker could attempt to perform a SQL injection when registering a new user for example or logging into the application. He could do that to be able to login as some other user for example or try to delete some data from tables in the database.
    
3.  Attacker could crack the password to our droplet and ssh to the machine and delete everything or install malicious software
    
4.  An attacker could DDOS (Denial of service) our web application, making it impossible for users to view the site. The attacker could also DDOS our api such the it would stop sending request to the simulator.
    
5.  An attacker could crack the password to our database and delete everything or access information. Since our GitHub is public and all other information about the database is therefore public (username, database name, where it is hosted etc.) it would only requiring bruteforcing the password.


## Risk Analysis
<img src="https://github.com/lauralunddahl/DevOps-GroupF/blob/dev/documents/WeeklyNotes/risk_analysis.PNG"/>

### Determine likelihood
The likelihood of most of these attacks are very low. How every it could happen that some malicious web crawler would find our website and attempt to Cross-site-script our website or do an SQL injection attack

- Scenario 1: Possible
- Scenario 2: Likely
- Scenario 3: Unlikely
- Scenario 4: Possible
- Scenario 5: Possible
### Determine impact
 - Scenario 1: Moderate
 - Scenario 2: Significant
 - Scenario 3: Significant
 - Scenario 4: Moderate
 - Scenario 5: Significant

### Use a Risk Matrix to prioritize risk of scenarios

 - Scenario 1: Risk(Possible, Moderate) = 6
 - Scenario 2: Risk(Likely, Significant) = 9
 - Scenario 3: Risk(Unlikely, Significant) = 7
 - Scenario 4: Risk(Possible, Moderate) = 6
 - Scenario 5: Risk(Possible, Significant) = 8

Based on the ratings we get from the risk assessment matrix the prioritization of the possible attacks is following

Scenario 2, Scenario 5, Scenario 3, Scenario 4 and Scenario 1

### Discuss what are you going to do about each of the scenarios

1.  **Scenario 2**: The go package we are currently using in our program already sanitizes the input that is being used to query the database so as it is now, the input is already parameterized which is a good defence against SQL injections.
    
2.  **Scenario 5**: Change the password of the database so it is long and complicated making it difficult to brute force. Check if it is possible to add two-factor authentication.
    
3.  **Scenario 3**: Change the password to the droplet so it is long and complicated, making it difficult to brute force.
    
4.  **Scenario 4**: By adding some scaling to our setup, if the main system is down due to DDos attack we will be able to start the backup system so we have the application up and running. Other than that, there is not much you can do when it comes to this type of attack.
    
5.  **Scenario 1**: Go has some packages that can be used to defend against XXS attacks. We will look into adding the functionality provided by those packages to prevent this type of attack.

##