# Double-Down Controller

Learning the concepts of a Kubernetes controller by creating a simple one myself!

## Description

A very simple controller which watches for the `jdocklabs.co.uk/double-down: 'true'` annotation on a `Deployment` object and doubles the number
of replicas in it, thereby doubling down.

We avoid continually doubling the number of pods by using the `jdocklabs.co.uk/doubled: 'true'` annotation.
If this is found, we do nothing.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

To test this out, you can run the following

```bash
make run

# In another terminal, create a deployment with 2 replicas
kubectl create deployment double-down-nginx --image nginx --replicas 2

# Verify the number of replicas
kubectl get deployments.apps double-down-nginx -o jsonpath='{.spec.replicas}{"\n"}'

# Annotate the Deployment
kubectl annotate deployments.apps double-down-nginx 'jdocklabs.co.uk/double-down=true'

# Verify the number of replicas again, now 4 is expected.
kubectl get deployments.apps double-down-nginx -o jsonpath='{.spec.replicas}{"\n"}'
```

You can also accomplish the same using a Kubernetes manifest, but this is shown step-by-step for understanding.


### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster.

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

