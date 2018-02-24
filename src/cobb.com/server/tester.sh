#! /bin/bash
cd test
go test
cd ../events/controllers
go test
cd ../services
go test
