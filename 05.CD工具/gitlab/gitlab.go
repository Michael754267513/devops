/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package gitlab-mod

import (
	"bytes"

	"github.com/xanzy/go-gitlab"
)

type GitlabDevOps struct {
	Token 		string 			`json:"token"` 	// gitlab token
	Url   		string  		`json:"url"`	// gitlab url
	Client      *gitlab.Client	`json:"client"` // gitlab client

}
// 初始化连接
func (g *GitlabDevOps) Connect() (err error) {
	g.Client, err = gitlab.NewClient(
		g.Token,
		gitlab.WithBaseURL(g.Url),
	)
	return  err
}

// 查看所有分支
func (g *GitlabDevOps) ListBranch(pid interface{}) (branchList []*gitlab.Branch,err error) {
	branchList, _,err = g.Client.Branches.ListBranches(pid,&gitlab.ListBranchesOptions{})
	if err != nil {
		return
	}
	return

}

// 获取提交的唯一id号  sha
func (g *GitlabDevOps) GetBranchCommit(pid interface{},branch string) (commitid string) {
	signBranch,_,err := g.Client.Branches.GetBranch(pid,branch)
	if err !=nil {}
	commitid = signBranch.Commit.ID
	return
}

// 创建merge
func (g *GitlabDevOps) CreateMerge(pid interface{},sourceBranch string,targetBranch string,title string,removeSourceBranch bool) (mergeRequest *gitlab.MergeRequest,err error) {
	mergeRequest,_,err = g.Client.MergeRequests.CreateMergeRequest(pid,&gitlab.CreateMergeRequestOptions{
		Title:              &title, // merge 标签
		Description:        &title, // 描述
		SourceBranch:       &sourceBranch,
		TargetBranch:       &targetBranch,
		Labels:             nil,
		AssigneeID:         nil,
		AssigneeIDs:        nil,
		TargetProjectID:    nil,
		MilestoneID:        nil,
		RemoveSourceBranch: &removeSourceBranch,
		Squash:             nil,
		AllowCollaboration: nil,
	})
	if err != nil {}

	return
}

// 获取合并信息
func (g *GitlabDevOps) GetMerge(pid interface{},mergeID int) (mergeRequest *gitlab.MergeRequest,err error)  {
	mergeRequest,_,err =g.Client.MergeRequests.GetMergeRequest(pid,mergeID,&gitlab.GetMergeRequestsOptions{})
	if err != nil {}
	return
}

// 运行合并
func (g *GitlabDevOps) AcceptMerge(pid interface{},mergeID int) (mergeRequest *gitlab.MergeRequest,err error) {
	mergeRequest,_,err = g.Client.MergeRequests.AcceptMergeRequest(pid,mergeID,&gitlab.AcceptMergeRequestOptions{})
	if err!=nil{}
	return
}

// 获取单个pipeline
func (g *GitlabDevOps) GetPipeline(pid interface{},pipelineID int)(pipeline *gitlab.Pipeline,err error) {
	pipeline,_,err = g.Client.Pipelines.GetPipeline(pid,pipelineID)
	if err != nil {}
	return
}

// 获取commit merge 后生成的pipeline 列表
func (g *GitlabDevOps) ListPipeline(pid interface{},commitid *string) (pipelineInfo []*gitlab.PipelineInfo,err error) {
	 pipelineInfo,_,err = g.Client.Pipelines.ListProjectPipelines(pid,&gitlab.ListProjectPipelinesOptions{SHA: commitid })
	 if err != nil {}
	return
}

// 获取pipeline的jobs列表
func (g *GitlabDevOps) GetPipeLineJobs(pid interface{},pipelineID int)(jobs []*gitlab.Job,err error){
	jobs ,_,err = g.Client.Jobs.ListPipelineJobs(pid,pipelineID,&gitlab.ListJobsOptions{})
	if err !=nil {}
	return

}

// 获取job详情
func (g *GitlabDevOps) GetJob(pid interface{},jobID int) (job *gitlab.Job,err error) {
	job,_,err  =g.Client.Jobs.GetJob(pid,jobID)
	if err !=nil {}
	return
}

// 获取pipeline --> job --> 运行日志
func (g *GitlabDevOps) GetJobMessage(pid interface{},jobID int)(file *bytes.Reader,err error)  {
	file,_,err = g.Client.Jobs.GetTraceFile(pid,jobID)
	if err !=nil {}
	return
}

 
