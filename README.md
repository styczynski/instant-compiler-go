# latte-compiler/Boilerplate For Golang Projects

[![Build Status](https://travis-ci.com/styczynski/latte-compiler.svg?branch=master)](https://travis-ci.com/styczynski/latte-compiler)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=MartinHeinz_go-project-latte-compiler&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=MartinHeinz_go-project-latte-compiler)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ec7ebefe63609984cb5c/test_coverage)](https://codeclimate.com/github/styczynski/latte-compiler/test_coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/styczynski/latte-compiler)](https://goreportcard.com/report/github.com/styczynski/latte-compiler)

## Blog Posts - More Information About This Repo

You can find more information about this project/repository and how to use it in following blog posts:

- [Ultimate Setup for Your Next Golang Project](https://towardsdatascience.com/ultimate-setup-for-your-next-golang-project-1cc989ad2a96)
- [Setting up GitHub Package Registry with Docker and Golang](https://towardsdatascience.com/setting-up-github-package-registry-with-docker-and-golang-7a75a2533139)
- [Building RESTful APIs in Golang](https://towardsdatascience.com/building-restful-apis-in-golang-e3fe6e3f8f95)
- [Setting Up Swagger Docs for Golang API](https://towardsdatascience.com/setting-up-swagger-docs-for-golang-api-8d0442263641)

### Setting Up
- Replace All Occurrences of `styczynski/latte-compiler` with your username repository name
- Replace All Occurrences of `latte-compiler` with your desired image name


### Adding New Libraries/Dependencies
```bash
go mod vendor
```

### Using GitHub Registry

Create and Push:

```bash
docker login docker.pkg.github.com -u <USERNAME> -p <GITHUB_TOKEN>
docker build -t  docker.pkg.github.com/styczynski/latte-compiler/latte-compiler:latest .
# make container
docker push docker.pkg.github.com/styczynski/latte-compiler/latte-compiler:latest
# make push
```

Pull and Run:

```bash
docker pull docker.pkg.github.com/styczynski/latte-compiler/latte-compiler:latest
docker run docker.pkg.github.com/styczynski/latte-compiler/latte-compiler:latest
```


### Setup new SonarCloud Project

- On _SonarCloud_:
    - Click _Plus_ Sign in Upper Right Corner
    - _Analyze New Project_
    - Click _GitHub app configuration_ link
    - Configure SonarCloud
    - Select Repository and Save
    - Go Back to Analyze Project
    - Tick Newly Added Repository
    - Click Set Up
    - Click Configure with Travis
    - Copy the Command to Encrypt the Travis Token
    - Run `travis encrypt --com <TOKEN_YOU_COPPIED>`
    - Populate the `secure` Field in `.travis.yml` with outputted string
    - Follow steps to populate your `sonar-project.properties`
    - Push
- On Travis CI:
    - Set `DOCKER_USERNAME`
    - Set `DOCKER_PASSWORD` to Your GitHub Registry Token

### Setup CodeClimate
- Go to <https://codeclimate.com/github/repos/new>
- Add Repository
- Go to Test Coverage Tab
- Copy Test Reporter ID
- Go to Travis and Open Settings for Your Repository
- Add Environment Variable: name: `CC_TEST_REPORTER_ID`, value: _Copied from CodeClimate_
