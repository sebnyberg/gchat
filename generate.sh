#!/bin/bash

protoc pkg/pb/gchat.proto --go_out=plugins=grpc:.