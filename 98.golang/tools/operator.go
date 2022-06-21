package devops

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	initConfig "k8s/config"
)

// 获取GVR客户端
func GetGVRClient(gvk *schema.GroupVersionKind, namespace string) (dr dynamic.ResourceInterface, err error) {
	config, err := initConfig.GetRestConf()
	// discoveryclient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)

	// GVK GVR 映射
	mapperGVKGVR := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))

	resourceMampper, err := mapperGVKGVR.RESTMapping(gvk.GroupKind(), gvk.Version)

	// 动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)

	if resourceMampper.Scope.Name() == meta.RESTScopeNameNamespace {
		dr = dynamicClient.Resource(resourceMampper.Resource).Namespace(namespace)
	} else {
		dr = dynamicClient.Resource(resourceMampper.Resource)
	}

	return

}

// 用于k8s clientset 客户端无法识别的api接口，用动态客户端去实现
func OperatorCreate(operatorData []byte) (err error) {

	var (
		objCreate *unstructured.Unstructured
		gvk       *schema.GroupVersionKind
		dr        dynamic.ResourceInterface
	)
	obj := &unstructured.Unstructured{}
	_, gvk, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(operatorData, nil, obj)
	dr, err = GetGVRClient(gvk, obj.GetNamespace())
	objCreate, err = dr.Create(context.Background(), obj, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("错误信息", err)
		fmt.Println(objCreate)
	}
	return
}

// 用于k8s clientset 客户端无法识别的api接口，用动态客户端去实现
func OperatorDelete(operatotData []byte) (err error) {

	var (
		gvk *schema.GroupVersionKind
		dr  dynamic.ResourceInterface
	)
	obj := &unstructured.Unstructured{}
	_, gvk, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(operatotData, nil, obj)
	dr, err = GetGVRClient(gvk, obj.GetNamespace())
	err = dr.Delete(context.Background(), obj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		fmt.Println("错误信息", err)
	}
	return
}
