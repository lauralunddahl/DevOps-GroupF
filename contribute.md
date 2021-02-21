# Description of Distributed Workflow


## Which repository setup will we use?
We will be using the mono repository model, such that we have a single unified source-code repository that is accessible to all team members. 
## Which branching model will we use?
We make use of the Git-flow branching model, where we will have a master branch for storing versions and have another branch dev where the development will be. When creating a feature or working on some specific task we will create a new branch from dev and when it’s done then merge it back to dev. Once we have a stable version we want to “publish” we will merge from develop to master.
## Which distributed development workflow will we use?
We will be using a centralized workflow where every group member will contribute to our shared repository.
## How do we expect contributions to look like?
We will be using the Private Small Team collaboration setup. We are five developers on the project, as it is only us who are allowed to contribute to the project, but it is open on GitHub for others to follow. 

Commit guidelines:
- Commit regularly 
- Limit subject to 50 characters
- Summarize in bullet points what has been done
- If it is a smaller fix/change e.g. fixing a typo then a single line is fine


## Who is responsible for integrating/reviewing contributions?
We will be using the merging topic branches staged integration merging workflow. In order to merge a topic branch into dev, the developer will have to create a pull-request, and then another developer will review the changes and be in charge of then accepting the pull-request. 

