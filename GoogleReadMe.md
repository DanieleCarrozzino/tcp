
# GOOGLE VIRTUAL MACHINE

## Create virtual machine
this command start a virtual machine over your logged project
gcloud compute instances create bb-virtual-machine 
--machine-type e2-micro 
--image-family debian-10 
--image-project debian-cloud

## Stop virtual machine
gcloud compute instances stop bb-virtual-machine

## Re start virtual machine
gcloud compute instances start bb-virtual-machine

## Get status
gcloud compute instances list

## Eneter through ssh
gcloud compute ssh bb-virtual-machine