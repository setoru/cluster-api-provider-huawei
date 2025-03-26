# Cluster API Provider HuaweiCloud Roadmap

> 本路线图在不断演进，过程中会不定期更新调整.

## v0.0.x (v1alpha1)

- [x] [项目基本框架搭建]()
- [x] [支持基于 tilt 工具的本地开发工作流](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/issues/15)
- [x] [InfraCluster 整体控制器协调逻辑框架实现](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/issues/27)
- [x] [InfraMachine 整体控制器协调逻辑框架实现](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/44)
- [x] [支持 VPC 及 Subnet 服务创建及销毁](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/29)
- [x] [创建 InfraMachine 所需的 K8S 节点最小可用系统磁盘镜像](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/issues/34)
- [x] [支持 SecurityGroups 服务创建及销毁](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/40)
- [x] [支持 Elastic LB 服务的创建及销毁](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/43)
- [x] [支持 NAT Gateways 服务的创建及销毁](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/42)
- [x] [支持 ECS 及 Elastic IP 服务的创建及销毁](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/44)
- [x] [InfraCluster/InfraMachine 与基于 kubeadm 的 CAPI 内建资源协调打通](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/47)
- [x] [支持 MachineDeployment 创建工作节点](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/60)
- [x] [添加项目发版自动化 CI 支持](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/pull/54)
- [ ] [协调: 提供社区用户访问华为 ECS 公共镜像的方案]()
- [ ] [支持通过模版自动化部署 CNI 组件]()
- [ ] [协调: `kubernetes-sigs/cloud-provider-huaweicloud` 项目支持 v1.32.0 版本 k8s]()
- [ ] [支持通过模版自动化部署 cloud-provider-huaweicloud]()
- [ ] [集群控制平面节点的扩缩容功能可用]()
- [ ] [集群工作节点的扩缩容功能可用]()
- [ ] [集群负载均衡器功能可用]()
