# `Terraform` 简介

## 什么是 Terraform 
[Terraform](https://www.terraform.io/) 是一个开源的IT基础设施编排管理工具，Terraform支持使用配置文件描述单个应用或整个数据中心。

通过Terraform您可以轻松的创建、管理、删除资源，并对其进行版本控制。目前几乎所有的主流云服务商都支持Terraform。

### HCL
Terraform是通过HashiCorp Configuration Language来编写代码的，HCL是声明式的。

管理人员用HCL来描述整个基础架构应该是什么样的，然后通过cli命令去创建所描述的资源，管理人无需了解具体的实施步骤和细节，这就是声明式的好处，只关心自己想要什么。

## [Terraform 状态管理](https://www.terraform.io/docs/language/state/index.html)

Terraform初始化以后，会生成一个[状态文件](https://www.terraform.io/docs/language/state/index.html) ，该状态文件记录了最近一次操作的时间、各资源的相关属性、各变量的当前值、状态文件的版本、等等。下一次再操作的时候，terraform首先会把当前状态文件与当前运行的状态进行一次更新(diff前一次与当前期望的差异)，找出是否后有被删除或者更改了的资源，然后再根据.tf文件，决定那些资源需要删除、更新、创建。操作完成后，会重新生成一个状态文件。

### [Terraform Backends](https://www.terraform.io/docs/language/state/backends.html)
Terraform 会根据每次变更后生成一个[状态文件](https://www.terraform.io/docs/language/state/index.html) ，状态文件记录了每次期望变更的记录，所以该文件是比较重要的，默认状态文件会存在当前tf文件下面，多人操作的时候也需要读取该文件来diff他们之间的差异，[团队协作](https://www.terraform.io/guides/core-workflow.html#working-as-a-team) 的时候会存在问题，比如两边同时操作就没法去判断谁先谁后，这时候会需要锁，当有人使用的时候状态锁定不允许其他人操作，这时候我们就希望有一个集中管理状态文件的地方。

[Terraform Backends](https://www.terraform.io/docs/language/settings/backends/configuration.html) 配置后会把状态文件进行存储到远端，这样无论每次执行diff的时候都会从backend读取相应的状态文件(默认是local)


## [Terraform 核心工作流程](https://www.terraform.io/guides/core-workflow.html)

-[x] [Write](https://www.terraform.io/guides/core-workflow.html#write) - 编写基础架构代码，去声明你想要什么 
-[x] [Plan](https://www.terraform.io/guides/core-workflow.html#plan)  - 预览本次变化与上次变化的差异
-[x] [Apply](https://www.terraform.io/guides/core-workflow.html#apply) - 去执行代码，满足write的声明期望

## Terraform的优势

- 基础架构即代码（Infrastructure as Code）
  
  基础设施可以通过HCL来声明你想要的资源来保存为代码，这些基础设施可以进行版本化。可以多人使用也可以重复使用
  
- 变更预览
  
  通过plan可以对比上次运行状态和期望状态的差异
  
- 资源图
  
  Terraform建立了一个所有资源的图，并行创建和修改任何非依赖性资源。从而使得Terraform可以尽可能高效地构建基础设施，操作人员可以深入了解基础设施中的依赖性。
  
- 变更自动化
  
  复杂的变更集可以应用于您的基础设施，而只需最少的人工干预。有了前面提到的执行计划和资源图，您就可以准确地知道Terraform将改变什么，以及改变的顺序，从而避免了许多可能的人为错误。

- 自动化管理基础结构
  
  Terraform能够创建配置文件的模板，以可重复、可预测的方式定义、预配和配置ECS资源，减少因人为因素导致的部署和管理错误。能够多次部署同一模板，创建相同的开发、测试和生产环境。

- 自动化管理基础结构
  
  Terraform能够创建配置文件的模板，以可重复、可预测的方式定义、预配和配置ECS资源，减少因人为因素导致的部署和管理错误。能够多次部署同一模板，创建相同的开发、测试和生产环境。

- 将基础结构部署到多个云
  
  Terraform适用于多云方案，将类似的基础结构部署到阿里云、其他云提供商或者本地数据中心。开发人员能够使用相同的工具和相似的配置文件同时管理不同云提供商的资源。

## Terraform 缺点
**个人认为**：当terraform使用apply之后，会创建对应的资源，满足你的期望，问题在于当人为对你生成好的资源做操作进行一些改变的时候，terraform是无法感知到这点的，因为他无法loop你当前期望的状态是否发生改变，只能等下次apply的时候去重新触发。

------------------------------------
**关注不迷路，不定期分享devops相关技术文章**

![](../public/img/11.png) ![](../public/img/12.jpg)