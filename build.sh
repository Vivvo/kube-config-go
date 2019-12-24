#!/bin/bash

ps aux | grep kube-config-go | grep -v grep
rm rcc* kube-config-go
qtrcc 
go build .

