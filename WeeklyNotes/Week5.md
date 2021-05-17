# Week 5
 Here are our considerations to how our group adhere to the "Three Ways" characterizing DevOps

### Flow
To *make our work visible* we are using a board on github for tasks, in order to visualize progress. 

In regards to *Limit work in process (WIP)* we do not have an agreement for how many tasks must be in the work in progress bar. But generally when someone is taking long on a task other members of the group try to help to ensure that the task can be completed. 

In regards to *Reduce batch sizes* we do not have continuous deployment at the moment, with trying to deploy smaller parts and see how they work. At the end of the week everything made during the week is gathered for the deployment. But the plan is to figure out something with continuous deployment using Travis. UPDATE: We now have continous deployment with GitHub Actions. 

In regards to *Reduce the number of handoffs* we are in the process of automating things but have not managed to do much yet. But since we are only one team, and things don’t have to go between teams it has not really been much of a problem, since every member of the team has the same skillset. 

In regards to *Continually identify and elevate our constraints*, we are still missing a lot of automated steps in regards to deployments and test setup. UPDATE: Some static analysis tools has been setup and CI/CD with GitHub Actions.

*Eliminate hardship and waste in the value stream*, is not something we have yet to consider. But we have run into problems with partially done work, waiting, defects which should definitely be something to try and lessen.


### Feedback

In regards to *See problems as they occur*, our feedback loop for development is fairly quick. Before someone pushes their code they test it by manually testing the specific functionality, and after the test simulator is run to see if any error occurs when talking to the api. The problem is that (for now) all code testing steps are done manually which adds to the flow of our deployment. Manual testing can also impact the feedback loop negatively, if a manual test is forgotten, it might be some time before the error is discovered, eg. api not return correct status codes.  


When it comes to *Swarm and solve problems to build new knowledge*, when a problem is discovered it is often delegated to the person that knows most about the given part of the project, if that person can’t solve it, we choose a subgroup or even the entire group to help with the problem. No matter the case the problem is explained and the solution is presented to the entire group such that everyone can benefit from the new knowledge. 

In regards to *Keep pushing quality closer to the source*, the issue of relying on other teams or on other people to perform og approve a request is not a problem for a group our size, e.g since we don’t need approval from CEO before deploying. 


### Continual Learning and Experimentation

When it comes to *Enabling Organizational Learning and safety culture*, If an error is discovered in the system the group as whole tries to find the source of the issue. When it has been discovered the person/s that are most familiar with that certain part will work on the error. If they don’t manage to fix the error then the group will come together and work on it together. The responsibility is shared throughout the value stream

Regarding *Institutionalizing the improvement of daily work*, when we have a bug from last week that needs to be fixed, we reserve some time in the beginning of the new week to try to finish the outstanding issues instead of trying to work around them. Then we are able to fix them while they are still small and easy to fix.

In regards to *Transform Local Discoveries into global improvements*, for now we have not run into any issues such as they could be redundant errors. But usually for the tasks where everyone needs to know what is going on, we go through them together and work on them as a team so each member is aware of how it should work.

In regards to *Inject Resilience Patterns into our daily work*, we haven’t taken any precautions for now regarding if for example the database would be down.

