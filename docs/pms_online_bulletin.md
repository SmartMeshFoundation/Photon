# Photon Monitoring Service Online bulletin

We have launched the Spectrum main network and test network PMS services, developers can use PMS

## Spectrum Test network

The test network PMS has been launched, and the corresponding service configuration information is as follows:


Names|Description
--|--
ContractAddress|0xa2150A4647908ab8D0135F1c4BFBB723495e8d12
PMS'IP|transport01.smartmesh.cn
Port|7004
DelegatedChargeNode|0xaed9188842c05e07bf5abdde2fb400432ae49d28
DelegatedChargeToken|0x048257d9F5e671412E46f2Ff4B5F7AFDb7059A86



## Spectrum Main network

Spectrum main network PMS has been launched, and the corresponding service configuration information is as follows:

Names|Description
--|--
ContractAddress|0xa2150A4647908ab8D0135F1c4BFBB723495e8d12
PMS'IP|transport01.smartmesh.cn
Port|7003
DelegatedChargeNode|0xa94399b93da31e25ab5612de8c64556694d5f2fd
DelegatedChargeToken|0x6fdb6b4deb71c4D9AFbA4350e2e9D6CfD534F1cb


## How to use it？

If you want to use PMS on the main network or test network, you can modify it based on the above information.

1. The contract address of the node, ensuring that the contract address of the node is the same as the contract address of the PMS
2. The agent charges the token to ensure that the node has a proxy toll at the proxy charging node 
3. Modify the corresponding URL parameter when calling the interface.For example: `GET http://transport01.smartmesh.cn:7004/tx/<delegater_address>/<channel_address>`


If you still don't know how to use PMS, please go to the official website [tutorial](./sm_service.md).
