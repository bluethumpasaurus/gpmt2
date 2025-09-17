
# gpmt2  [![Build Status](https://travis-ci.org/pivotal-gss/gpmt2.svg?branch=master)](https://travis-ci.org/pivotal-gss/gpmt2)  [![Go Report Card](https://goreportcard.com/badge/github.com/pivotal-gss/gpmt2)](https://goreportcard.com/report/github.com/pivotal-gss/gpmt2)
Greenplum Magic Tool 2 - open source rebuild of GPMT

**The code does not work at present**

## Goals

- The ultimate goal of the project is to make an OSS version of the tool
- Initial development will focus on replicating functionality of the original GPMT tool, in the following order of tools (subject to adjustment):
  1. gp_log_collector
  2. analyze_session
  3. gpstatscheck
  
  
- Other goals include:
  - Compatibility with gpdb 6x and 7x
  - Additional features
  

## Contributing

Please ensure all submitted code is formatted with gofmt. 

This repository uses [dep](https://github.com/golang/dep) so please ensure you run `dep ensure --update` if you add imports. 

Of course, github issue requests are the best way to start a request for new features or bugs.
