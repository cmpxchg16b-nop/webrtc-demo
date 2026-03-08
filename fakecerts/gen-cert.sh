#!/bin/bash

cfssl genkey cert.json | cfssljson -bare cert
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem cert.json | cfssljson -bare cert
