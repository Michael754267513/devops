/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package harbor

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/TimeBye/go-harbor"
	rest2 "github.com/TimeBye/go-harbor/pkg/rest"
	"github.com/mittwald/goharbor-client/v4/apiv2/model"
	legacymodel "github.com/mittwald/goharbor-client/v4/apiv2/model/legacy"
	k8s_model "k8s/model"
)

type Harbor struct {
	User     string          `json:"user"`
	Password string          `json:"password"`
	Host     string          `json:"host"`
	c        rest2.Interface `json:"c"`
}

func (h *Harbor) InitClient() (err error) {
	ncs, err := harbor.NewClientSet(h.Host, h.User, h.Password)
	h.c = ncs.V2.RESTClient()
	return
}

// 获取所有的policys

func (h *Harbor) GetReplicationPolicies() (policies []legacymodel.ReplicationPolicy, err error) {
	err = h.c.Get().Resource("replication").Name("policies").Do().Into(&policies)
	return policies, err
}

// 获取policy的 policyID 判断policy 是否存在

func (h *Harbor) GetPolicyID(replicationName string) (replicationpolicy legacymodel.ReplicationPolicy, policyID int64, isExist bool) {
	policies, err := h.GetReplicationPolicies()
	if err != nil {
		fmt.Println("获取服务策略失败：", err)
	}
	for _, v := range policies {
		if v.Name == replicationName {
			return v, v.ID, true
		}
	}
	return
}

// 获取模板测试

func (h *Harbor) GetTemplatePolicy() (template legacymodel.ReplicationPolicy, err error) {
	policies, err := h.GetReplicationPolicies()
	if err != nil {
		fmt.Println("获取服务策略失败：", err)
		return
	}
	for _, v := range policies {
		if v.Name == "template" {
			return v, err
		}
	}
	return template, fmt.Errorf("复制名称为template的模板复制策略找不到")

}

// 创建复制策略

func (h *Harbor) CreateReplicationPloicy(projectURI, imageName, tag string) (err error) {
	template, err := h.GetTemplatePolicy()
	if err != nil {
		fmt.Println("获取服务策略失败：", err)
		return
	}
	var filters []*legacymodel.ReplicationFilter

	// 基于名称过滤
	filters = append(filters, &legacymodel.ReplicationFilter{
		Type:  "name",
		Value: projectURI,
	})
	// 基于tag过滤
	filters = append(filters, &legacymodel.ReplicationFilter{
		Type:  "tag",
		Value: tag,
	})

	var create_rep_policy legacymodel.ReplicationPolicy
	create_rep_policy.Filters = filters
	create_rep_policy.Enabled = true
	create_rep_policy.DestRegistry = template.DestRegistry
	create_rep_policy.SrcRegistry = template.SrcRegistry
	create_rep_policy.Trigger = template.Trigger
	create_rep_policy.Deletion = false
	create_rep_policy.Override = true
	create_rep_policy.Name = imageName
	err = h.c.Post().Resource("replication/policies").Body(&create_rep_policy).Do().Error()
	return
}

// 更新policy 用于镜像同步到生产

func (h *Harbor) UpdateReplicationPloicy(projectURI, imageName, tag string) (err error) {
	replicationPolicy, policyid, isexit := h.GetPolicyID(imageName)

	// 不存在
	if !isexit {
		return
	}

	// 判断本次上线和上次同步是否一致

	for _, v := range replicationPolicy.Filters {
		if v.Type == "tag" {
			if v.Value == tag {
				return
			}
		}
	}

	var filters []*legacymodel.ReplicationFilter

	// 基于名称过滤
	filters = append(filters, &legacymodel.ReplicationFilter{
		Type:  "name",
		Value: projectURI,
	})
	// 更新tag过滤
	filters = append(filters, &legacymodel.ReplicationFilter{
		Type:  "tag",
		Value: tag,
	})

	replicationPolicy.Filters = filters

	err = h.c.Put().Resource(fmt.Sprint("replication/policies/", policyid)).Body(&replicationPolicy).Do().Error()
	return

}

// 启动复制

func (h *Harbor) StartReplicationPolicy(imageName string) (err error) {
	_, policyid, _ := h.GetPolicyID(imageName)
	var replicationPolicy model.ReplicationExecution
	replicationPolicy.PolicyID = policyid
	err = h.c.Post().Resource("replication/executions").Body(&replicationPolicy).Do().Error()
	return
}

// 创建或者更新任务
func (h *Harbor) CreateOrUpdate(projectURI, imageName, tag string) (err error) {
	var (
		//replicationPolicy  legacymodel.ReplicationPolicy
		//policyID int64
		isExist bool
	)
	_, _, isExist = h.GetPolicyID(imageName)
	// 存在更新任务
	if isExist {
		if err := h.UpdateReplicationPloicy(projectURI, imageName, tag); err != nil {
			return err
		}
	}
	// 不存在新建任务
	if !isExist {
		if err := h.CreateReplicationPloicy(projectURI, imageName, tag); err != nil {
			return err
		}

	}

	return

}

// replicationName 为包名

func (h *Harbor) PushImage(projectURI, imageName, tag string) (err error) {
	var (
		//replicationPolicy  legacymodel.ReplicationPolicy
		//policyID int64
		isExist bool
	)
	_, _, isExist = h.GetPolicyID(imageName)

	if isExist {
		if err := h.UpdateReplicationPloicy(projectURI, imageName, tag); err != nil {
			return err
		}
	}

	if !isExist {
		if err := h.CreateReplicationPloicy(projectURI, imageName, tag); err != nil {
			return err
		}

	}
	if err := h.StartReplicationPolicy(imageName); err != nil {
		return err
	}

	return
}

// 获取replication POlicy 的最新的task状态

func (h *Harbor) GetPolicyTaskStatus(imageName string) (err error, rpts model.ReplicationExecution) {

	_, id, isExist := h.GetPolicyID(imageName)
	if !isExist {
		fmt.Errorf("不存在该复制策略")
		return fmt.Errorf("不存在该复制策略"), rpts
	}
	data := []model.ReplicationExecution{}
	query := k8s_model.PolicyTaskQuery{
		PolicyID: id,
		//Page: 1,
		//PageSize: 1,
	}
	err = h.c.Get().Resource("replication/executions").Body(query).Do().Into(&data)
	rpts = data[0]
	//err ,res := h.GetPolicyTaskLog (rpts.ID)
	//fmt.Println(res)
	return
}

// 查看 replicationtask log
// TODO 获取日志 目前格式化存在问题
func (h *Harbor) GetPolicyTaskLog(taskID int64) (err error, data []interface{}) {
	uri := fmt.Sprint("/replication/executions/", taskID, "/tasks/", taskID, "/log")

	err = h.c.Get().Resource(uri).Do().Into(&data)

	return
}

// 忽略https校验
func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
