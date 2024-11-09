#
Repo contains logic for the creation of the github action yaml files
#

###
The files in this repo include 

```
skeletons.go
```
list of all supported github actions on the github marketplace

```
watcher.go
```
watches code to automatically update the github action in realtime

```
workflows.go
```
defines all structs and functions to be used on the workflow structs

```
yaml_generator
```
builds the actual yaml file checks for if the file already exists... etc
###