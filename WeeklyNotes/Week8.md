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
- GORM: MIT 
- mysql: Mozilla public license 2.0
- grafana: Apache License 2.0
- prometheus: Apache License 2.0
- ELK (kibana + elasticsearch): Server Side Public License (SSPL) and the Elastic License (https://www.zdnet.com/article/elastic-changes-open-source-license-to-monetize-cloud-service-use/)
- filebeat: Elastic License
- crypto: Boost Software License (?) (https://github.com/golang/crypto/blob/master/LICENSE)
- tawesoft dialog: The Unlicense (https://github.com/tawesoft/go/blob/master/dialog/LICENSE.txt)

Due to ELK and firebeat having the SSPL license, which is based on the AGPL3 license (the strongest copyleft license), we must use the same license and therefore cannot choose to use e.g. the MIT license.
