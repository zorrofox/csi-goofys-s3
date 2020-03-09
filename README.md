# CSI for S3 Using Goofys Mounter

This is a Container Storage Interface ([CSI](https://github.com/container-storage-interface/spec/blob/master/spec.md)) for S3 (or S3 compatible) storage and using [Goofys](https://github.com/kahing/goofys) Mounter. This can dynamically allocate buckets and mount them via Goofys FUSE into any container.

## Status

This is still very experimental and should not be used in any production environment. Unexpected data loss could occur depending on what mounter and S3 storage backend is being used.

## Kubernetes on EKS installation

### Requirements

* AWS EKS 1.14+
* Kubernetes 1.13+ (CSI v1.0.0 compatibility)
* Kubernetes has to allow privileged containers
* Docker daemon must allow shared mounts (systemd flag `MountFlags=shared`)

### 1. Create 3 EKS Pod Roles for Kubernetes Service Account
These 3 command will create 3 IAM roles with AmazonS3FullAccess access policy and associate a new Kubernetes Service Account.
```bash
eksctl create iamserviceaccount \
    --name csi-attacher-sa \
    --namespace kube-system \
    --cluster <your_cluster_name> \
    --attach-policy-arn arn:aws-cn:iam::aws:policy/AmazonS3FullAccess \
    --approve \
    --override-existing-serviceaccounts

eksctl create iamserviceaccount \
    --name csi-s3-sa \
    --namespace kube-system \
    --cluster <your_cluster_name> \
    --attach-policy-arn arn:aws-cn:iam::aws:policy/AmazonS3FullAccess \
    --approve \
    --override-existing-serviceaccounts

eksctl create iamserviceaccount \
    --name csi-provisioner-sa \
    --namespace kube-system \
    --cluster <your_cluster_name> \
    --attach-policy-arn arn:aws-cn:iam::aws:policy/AmazonS3FullAccess \
    --approve \
    --override-existing-serviceaccounts
```

The region can be empty if you are using some other S3 compatible storage.

### 2. Deploy the driver

```bash
cd deploy/kubernetes
kubectl create -f provisioner.yaml
kubectl create -f attacher.yaml
kubectl create -f csi-s3.yaml
```

### 3. Create the storage class

```bash
kubectl create -f storageclass.yaml
```

### 4. Test the S3 driver

1. Create a pvc using the new storage class:

```bash
kubectl create -f pvc.yaml
```

2. Check if the PVC has been bound:

```bash
$ kubectl get pvc csi-s3-pvc
NAME         STATUS    VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
csi-s3-pvc   Bound     pvc-c5d4634f-8507-11e8-9f33-0e243832354b   5Gi        RWO            csi-s3         9s
```

3. Create a test pod which mounts your volume:

```bash
kubectl create -f busybox_pod_1.yaml
kubectl create -f busybox_pod_2.yaml
```

If the pod can start, everything should be working.

4. Test the mount

```bash
$ kubectl exec -ti busybox-1 sh
$ mount | grep fuse
pvc-1a959cb4-603f-11ea-af00-021e45d100ce:csi-fs on /data type fuse (rw,nosuid,nodev,relatime,user_id=0,group_id=0,default_permissions,allow_other)
$ touch /data/hello_world
```

If something does not work as expected, check the troubleshooting section below.

## Additional configuration

### Mounter

As S3 is not a real file system there are some limitations to consider here. Depending on what mounter you are using, you will have different levels of POSIX compability. Also depending on what S3 storage backend you are using there are not always [consistency guarantees](https://github.com/gaul/are-we-consistent-yet#observed-consistency).

The driver can be configured to use one of these mounters to mount buckets:

* [goofys](https://github.com/kahing/goofys)

The mounter can be set as a parameter in the storage class. You can also create multiple storage classes for each mounter if you like.

All mounters have different strengths and weaknesses depending on your use case. Here are some characteristics which should help you choose a mounter:

#### goofys

* Weak POSIX compatibility
* Performance first
* Files can be viewed normally with any S3 client
* Does not support appends or random writes

## Troubleshooting

### Issues while creating PVC

Check the logs of the provisioner:

```bash
kubectl logs -l app=csi-provisioner-s3 -c csi-s3
```

### Issues creating containers

1. Ensure feature gate `MountPropagation` is not set to `false`
2. Check the logs of the s3-driver:

```bash
kubectl logs -l app=csi-s3 -c csi-s3
```

## Development

This project can be built like any other go application.

```bash
go get -u github.com/zorrofox/csi-s3
```

### Build executable

```bash
make build
```

### Tests

Currently the driver is tested by the [CSI Sanity Tester](https://github.com/kubernetes-csi/csi-test/tree/master/pkg/sanity). As end-to-end tests require S3 storage and a mounter like s3fs, this is best done in a docker container. A Dockerfile and the test script are in the `test` directory. The easiest way to run the tests is to just use the make command:

```bash
make test
```
