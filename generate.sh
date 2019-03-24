#!/bin/bash

protoc pkg/pb/calculator.proto --go_out=plugins=grpc:.