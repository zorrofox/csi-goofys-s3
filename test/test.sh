#!/usr/bin/env bash

export AWS_REGION=cn-northwest-1
export AWS_ACCESS_KEY_ID=<YOUR KEY>
export AWS_SECRET_ACCESS_KEY=<YOUR KEY>

go test github.com/zorrofox/csi-goofys-s3/pkg/s3 -v -args -v 5 -logtostderr true