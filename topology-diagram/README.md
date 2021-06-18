# Generating Topology diagram using Drawthe.net

## Usage
```
python generate_topology_diagram.py -t <topology.yaml>
```

### Example
```
python generate_topology_diagram.py -t ipfabric-test.yaml
```
The drawit yaml file is saved as drawit_<topology_name>.yaml

## WIP
- the x and y axis have to be manually edited
- if the connection doesnt exist for the other direction, then the diagram isnt rendered. 
