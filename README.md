[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/200106-uta-go/BAM-P2?include_prereleases)](https://github.com/200106-uta-go/BAM-P2/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/200106-uta-go/BAM-P2?style=flat-square)](https://goreportcard.com/report/github.com/200106-uta-go/BAM-P2)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/200106-uta-go/BAM-P2)


# BAM
Brandon, Aaron, and Matt -- Revature-Project 2

BAM is an easy to use tool to create and manage a Kubernetes cluster on AWS. BAM gives user the ability to quickly spin up a kubernetes cluster on AWS from your host machine and then deploys a simple interface locally to control your cluster deployment. 

## Getting Started

To get started created a Kubernetes cluster on AWS with BAM, run the `startup` executable and follow the instructions to get your host machine set up to connect to AWS. The startup program will then prompt you for information about the cluster you want to create. Once all the user information has been recieved, BAM will deploy your cluster on AWS.

To control your cluster, launch the `Controller` application and use your browser to access http://localhost:4040/web to get to the cluster controller interface. 

## Using the installer

How to use the installer/launcher

## Using the Web Interface

The BAM web interface allows the user to wasily control a cluster from their host machine. Once your cluster has been created, you need to add a deployment to start running container in your kubernetes pods.

The web interface provides two ways to create a deployment. Id you have a .yaml formatted deployment ready to go, you can paste your deployment into the text box that appears when you click the `add deployment` button.

![Image of add deployment button](https://i.ibb.co/5x6XJzm/controller-addpod.png)

The other way you can start a deployment is to add a new pod using a container image name from [DockerHub](https://hub.docker.com/). Make sure to include the user if present and the version you want to deploy.

![Image of add pod button](https://i.ibb.co/CnWgL2x/controller-deployment.png)

Once your deployment is running, you can increase and decrease the pod replicas for each deployment using the buttons next to the pods header.

![Image of scale buttons](https://i.ibb.co/Yy2kRPb/scale.png)

Services must be created using a deployment using .yaml formatted text in the apply deployment window. 

## Development

Talk about how to set up the development environment

## Built With

* [KOPS](https://github.com/kubernetes/kops) - The cluster creation tool used
* [Kubernetes](https://kubernetes.io/docs/home/) - Cluster management and networking
