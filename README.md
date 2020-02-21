[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/200106-uta-go/BAM-P2?include_prereleases)](https://github.com/200106-uta-go/BAM-P2/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/200106-uta-go/BAM-P2?style=flat-square)](https://goreportcard.com/report/github.com/200106-uta-go/BAM-P2)

# BAM
Brandon, Aaron, and Matt -- Revature-Project 2

BAM is an easy to use tool to create and manage a Kubernetes cluster on AWS. BAM gives user the ability to quickly spin up a kubernetes cluster on AWS from your host machine and then deploys a simple interface locally to control your cluster deployment. 

## Getting Started

To get started created a Kubernetes cluster on AWS with BAM, run the `startup` executable and follow the instructions to get your host machine set up to connect to AWS. The startup program will then prompt you for information about the cluster you want to create. Once all the user information has been recieved, BAM will deploy your cluster on AWS.

To control your cluster, launch the `Controller` application and use your browser to access http://localhost:4040/web to get to the cluster controller interface. 

## Using the installer

The installer is packaged with two binaries, one for dependency installation, and one for deploying the cluster to AWS. 

The only prerequisite to running the installer is to have an AWS account and a generated key pair so that KOPS can start ec2 and s3 instances for your cluster. 

If your machine does not have kops, kubectl, or aws cli installed, you need to run the `Installer` before a cluster can be created. Once the `Installer` has run, or if you already have all the necessary programs installed, you can run the `Deploy` executable to launch your kubernetes cluster on AWS.

## Using the Web Interface

The BAM web interface allows the user to wasily control a cluster from their host machine. Once your cluster has been created, you need to add a deployment to start running container in your kubernetes pods.

The web interface provides two ways to create a deployment. Id you have a .yaml formatted deployment ready to go, you can paste your deployment into the text box that appears when you click the `add deployment` button. Services must be created using a deployment using .yaml formatted text in the apply deployment window. 

![Image of add deployment button](https://i.ibb.co/5x6XJzm/controller-addpod.png)

The other way you can start a deployment is to add a new pod using a container image name from [DockerHub](https://hub.docker.com/). Make sure to include the user if present and the version you want to deploy.

![Image of add pod button](https://i.ibb.co/8gRV7qM/controller-deployment.png)

Once your deployment is running, you can increase and decrease the pod replicas for each deployment using the buttons next to the pods header.

![Image of scale buttons](https://i.ibb.co/Yy2kRPb/scale.png)

You can view the logs for the pod and delete the pod from your deployment using the buttons on the bottom of each pod's details.



## Development

To get started with a development environment for this project, you only need a go installation and the go dependencies for each application. These can be retrieved using:
```
go get -d
```
You will find multiple different applications in the command package, most of which were designed to be used inside containers running on a kubernetes cluster. The makefile has scripts to create docker images and deploy each app into a docker container. The apps can also be built into a go binary for easy testing during development. 

## Built With

* [KOPS](https://github.com/kubernetes/kops) - The cluster creation tool used
* [Kubernetes](https://kubernetes.io/docs/home/) - Cluster management and networking