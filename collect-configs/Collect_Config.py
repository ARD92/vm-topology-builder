'''
__author__: Aravind Prabhakar

used to collect configs from devices.

Usage:

usage: Collect_Config.py [-h] [--action ACTION] [--dir DIR] [--user USER]
                         [--pwd PWD]

optional arguments:
  -h, --help       show this help message and exit
  --action ACTION
  --dir DIR
  --user USER
  --pwd PWD

Example:

1. Collect configs in set format and store it in a directory 
python Collect_Config.py --action save --dir SE-LAB-SET-CONFIG --user aprabh --pwd Juniper123 --type set

2. Collect configs as is in junos hierarchy and store it in a directory 
python Collect_Config.py --action save --dir SE-LAB-SET-CONFIG --user aprabh --pwd Juniper123

3. Load the config based on set or regular
python Collect_Config.py --action load --user aprabh --pwd Juniper123 --file add_more_vlans
'''

from jnpr.junos import Device
import argparse
import yaml
from jnpr.junos.exception import ConnectAuthError, ConnectRefusedError,ConnectError, ConnectTimeoutError, RpcError, CommitError,ConfigLoadError
from jnpr.junos import rpcmeta
from jnpr.junos.utils.config import Config
from jnpr.junos.factory.factory_loader import FactoryLoader

yaml_open = open('Collect_Config.yaml')
hostList = yaml.load(yaml_open)

parser = argparse.ArgumentParser()
parser.add_argument('--action', action='store', dest='action', type=str, default=None)
parser.add_argument('--dir', action='store', dest='dir', type=str, default=None)
parser.add_argument('--user', action='store', dest='user', type=str, default=None)
parser.add_argument('--pwd', action='store', dest='pwd', type=str, default=None)
parser.add_argument('--type', action='store', dest='type', type=str, default=None)
parser.add_argument('--file', action='store', dest='file', type=str, default=None)

args=parser.parse_args()

def gen_report(dev,hst):
    print(hst)
    with open('info_'+str(hst)+'.txt','w') as f:
        for op in map['opcommand']:
            f.write(20*"*"+" "+op+" "+20*"*")
            f.write('\n')
            f.write(dev.cli(op))
            #xml_rpc = dev.display_xml_rpc(op)
            #try:
            #    f.write(dev.execute(xml_rpc))
            #except RpcError:
            #    print("RPC Error")   
            f.write("\n")
            f.write("\n")
        f.close()

if args.action == 'save':
    for host in hostList['hosts']:
        dev = Device(host=host , user=args.user, password=args.pwd)
        dev.open()
        hostname = dev.facts['hostname']
        print('Collecting configuration for: ' , hostname)

        class Collect_Config():
            def __init__(self):
                config=str(self.get_conf(args.type))
                try:
                    if args.type == "set":
                        with open(args.dir +"/" + hostname + "_set", "w") as file:
                            file.write(config)
                            file.close()
                    else:
                        with open(args.dir +"/" + hostname, "w") as file:
                            file.write(config)
                            file.close()
                except:
                    print('Configs directory is missing, create it!')

            def get_conf(self,type):
                if type == "set":
                    return dev.cli("show configuration | display set | no-more")
                else:
                    return dev.cli("show configuration")


        run=Collect_Config()

elif args.action == 'load':
    for host in hostList['hosts']:
        dev = Device(host=host , user=args.user, password=args.pwd)
        dev.open()
        hostname = dev.facts['hostname']
        cu = Config(dev)
        print('loading configuration for: ' , hostname)

        try:
            with open(args.file, "r") as file:
                a = file.read()
                if args.type == "set":
                    cu.load(a, load="set")
                else:
                    cu.load(a, overwrite=True)
                cu.pdiff()
                if cu.commit_check():
                    cu.commit()
                else:
                    cu.rollback()
        except:
            print('Commit failed')

elif args.action =="report":
    dev = Device(host=hst, user=args.user, passwd=args.pwd,  normalize=True)
    open_dev(dev)
    gen_report(dev,hst)
