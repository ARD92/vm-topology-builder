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
- if compute instances need to be referred, then use either "server" or "vm"
- the interface naming convention has to be
    ge-x.y.z , xe-x.y.z, et-x.y.z for the vMX, vSRX, vQFX
    ethx , ensx for vms/server representations
