## Sober

Configuration Center for Applications.

## Modules

Server:
- watch store changes and send them to agent

Agent:
- read configurations from data source and save them to outputs


Store:

- [x] Etcd
- [ ] ZooKeeper
- [ ] Consul

DataSource:

- [x] gRPC
- [x] .ini file

Output:

- [x] json file
- [x] stdout
- [ ] shared memory

Authentication：

- [x] token

Others:

- auto watch and reload
- easy golang library

Todo:

- [ ] traffic optimization
- [ ] match algorithm optimization
- [ ] gRPC authentication

