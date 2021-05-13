## Dependency diagram

<img src="https://github.com/lauralunddahl/DevOps-GroupF/blob/dev/WeeklyNotes/Dependencies.png"/>

http://www.webgraphviz.com/

```
digraph G {
  "travis.yml" -> "docker-compose.yml"
  "Dockerfile" -> "main.go"
  "docker-compose.yml" -> "grafana 4.5.2"
  "docker-compose.yml" -> "prometheus"
  "docker-compose.yml" -> "Dockerfile"
  "docker-compose.yml" -> "filebeat 7.2.0"
  "docker-compose.yml" -> "kibana 7.2.0"
  "filebeat 7.2.0" -> "elasticsearch 7.2.0"
  "kibana 7.2.0" -> "elasticsearch 7.2.0"
  "main.go" -> "api/"
  "main.go" -> "css/"
  "main.go" -> "metrics/"
  "main.go" -> "minitwit/"
  "minitwit/" -> "templates/"
  "templates/" -> "css/"
  "api/" -> "db/"
  "api/" -> "helper/"
  "api/" -> "dto/"
  "minitwit/" -> "db/"
  "minitwit/" -> "helper/"
  "minitwit/" -> "dto/"
  "db/" -> "GORM 2.0"
  "GORM 2.0" -> "mysql 1.5.0"
  "metrics/" -> "prometheus"
  "main.go" -> "Gorilla mux 1.8.0"
  "minitwit/" -> "Gorilla mux 1.8.0"
  "api/" -> "Gorilla mux 1.8.0"
  "Dockerfile" -> "go.mod"
  "Dockerfile" -> "go.sum"
  "main.go" -> "go compiler"
  "Gorilla mux 1.8.0" -> "go compiler"
  "go compiler" -> "Ubuntu"
}
```
## Licenses
- Gorilla mux: BSD 3-Clause 
- GORM: MIT License
- GORM MySQL driver: MIT License
- mysql: Mozilla public license 2.0
- grafana: Apache License 2.0
- prometheus: Apache License 2.0
- ELK (kibana + elasticsearch): Server Side Public License (SSPL) and the Elastic License (dual-licensed) (https://www.zdnet.com/article/elastic-changes-open-source-license-to-monetize-cloud-service-use/
- filebeat: Server Side Public License (SSPL) & Elastic License (dual-licensed)
- crypto: BSD 3-Clause (https://github.com/golang/crypto/blob/master/LICENSE)
- tawesoft go/dialog: The Unlicense (https://github.com/tawesoft/go/blob/master/dialog/LICENSE.txt)
- logrus: MIT License
- gorilla/session: BSD 3-Clause
- godotenv: MIT License
- goupsutil: BSD License

Before choosing a license for the project, the advantages and disadvantages of the different licenses was discussed. An advantage of using a permissive license such as the MIT license is that since it is free and without intimidating legal constraints, more people may be inclined to try it and some of these people will eventually contribute to it. A trade-off is however that organizations can take the work and commercialize it. Considering the copyleft GPL license we found that this fitted best with our vision. The GPL license is beneficial for individual contributors and often create a community among them. By using the GPL license we still allow users to run, modify and share the software for free, but with the restriction that their improved versions must also stay free. This may discourage some larger organizations to use it, but since our software is a social media platform, we wish to create a community where we encourage cooperation and allow users to contribute freely while avoiding larger organizations turning it into proprietary software.

The GPL license is compatible with the licenses of the used dependencies. 

The EFK stack has recently changed from using Elastic license and Apache 2.0 license to using a dual-license under Elastic Search and Server Side Public License (SSPL - which is based on the strongest copyleft GPL3). The dual-license however means that the user can choose which terms to work under if working directly with the source code. If not working with the source code (as in our case, just using the default versions) it uses the Elastic Search license as before (source: https://www.elastic.co/pricing/faq/licensing). 
