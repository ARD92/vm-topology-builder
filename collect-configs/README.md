# Collect Configuration
we can collect configs from the devices spun up using the topology builder. To get a copy of the configs. Edit the file `collect_Config.yaml` and add the host IP addresses of the nodes respectively. 
Then run the below . Ensure the directory has been created if using the `--dir` option.
```
python3 Collect_Config.py --action save --dir test --user root --pwd juniper123 --type set
```

you can also use action `load` with `--file` provided 

## Options
```
python3 Collect_Config.py --help
Collect_Config.py:41: YAMLLoadWarning: calling yaml.load() without Loader=... is deprecated, as the default Loader is unsafe. Please read https://msg.pyyaml.org/load for full details.
  hostList = yaml.load(yaml_open)
usage: Collect_Config.py [-h] [--action ACTION] [--dir DIR] [--user USER] [--pwd PWD] [--type TYPE] [--file FILE]

optional arguments:
  -h, --help       show this help message and exit
  --action ACTION
  --dir DIR
  --user USER
  --pwd PWD
  --type TYPE
  --file FILE
```
# Generate report 
you can generate report based on various operational commands using the below.
Ensure you add the commands you intend to run on the nodes by editing the file `report.py`

```
python3 report.py -user <username> -pwd <password>
```

## Options
```
python3 report.py --help
report.py:14: YAMLLoadWarning: calling yaml.load() without Loader=... is deprecated, as the default Loader is unsafe. Please read https://msg.pyyaml.org/load for full details.
  map = yaml.load(yaml_open)
usage: report.py [-h] -user USER -pwd PWD

optional arguments:
  -h, --help  show this help message and exit
  -user USER
  -pwd PWD
```
