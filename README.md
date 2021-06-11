# vm-topology-builder
Topology builder to build lab topologies for multiple juniper vNFs 

Directory structure:
- templates: templates for each vnf based on dir. genxml will store the files in respective directory
- build: creates a directory based on topology names with respective images copied and is referenced in xml
- images: original image
- topologies: store all topologies yaml files which can be referenced
- vm-topo.go : main file

