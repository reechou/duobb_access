// 接入层:透传到后端系统

package duobb

import (
	"strings"

	"github.com/reechou/holmes"
)

type BackendService struct {
	ServiceName string
	Hosts       []string
	Type        string
	idx         int
}

type BackendServiceMethod struct {
	ServiceMethod string
	Uri           string
}

func (self *DuobbProcess) RegisterService(serviceName string, hosts []string, t string) {
	if serviceName == "" || len(hosts) == 0 {
		holmes.Error("register service argv error.")
		return
	}

	self.smMutex.Lock()
	defer self.smMutex.Unlock()

	if t == "" {
		t = BACKEND_GOJSONRPC
	}
	self.serviceMap[serviceName] = &BackendService{
		ServiceName: serviceName,
		Hosts:       hosts,
		Type:        t,
	}
	holmes.Info("duobb register backend service: %v", self.serviceMap[serviceName])
}

func (self *DuobbProcess) RegisterServiceMethod(serviceMethodName string, uri string) {
	if serviceMethodName == "" || uri == "" {
		holmes.Error("register service method argv error.")
		return
	}

	self.smMutex.Lock()
	defer self.smMutex.Unlock()

	self.serviceMethodMap[serviceMethodName] = &BackendServiceMethod{
		ServiceMethod: serviceMethodName,
		Uri:           uri,
	}
	holmes.Info("duobb register backend service method: %v", self.serviceMethodMap[serviceMethodName])
}

func (self *DuobbProcess) chooseHost(s *BackendService) string {
	self.smMutex.Lock()
	defer self.smMutex.Unlock()

	host := s.Hosts[s.idx]
	s.idx = (s.idx + 1) % len(s.Hosts)
	return host
}

func (self *DuobbProcess) chooseSMUri(method string) string {
	self.smMutex.Lock()
	defer self.smMutex.Unlock()

	serviceMethod := self.serviceMethodMap[method]
	if serviceMethod == nil {
		return ""
	}

	return serviceMethod.Uri
}

func (self *DuobbProcess) process(method string, decodeMsg []byte) ([]byte, error) {
	sm := strings.Split(method, ".")
	if len(sm) != 2 {
		holmes.Error("bad method[%s]", method)
		return nil, ErrorBadMethod
	}
	service := self.serviceMap[sm[0]]
	if service == nil {
		holmes.Error("unknown method[%s] of service[%v]", method, self.serviceMap)
		return nil, ErrorUnkownService
	}
	host := self.chooseHost(service)
	switch service.Type {
	case BACKEND_GOJSONRPC:
		request, err := JsonDecode(decodeMsg)
		if err != nil {
			holmes.Error("json decode[%s] error: %v", string(decodeMsg), err)
			return nil, err
		}
		result, err := self.jsonRpc.Call(host, method, request)
		if err != nil {
			holmes.Error("gojsonrpc call host[%s] method[%s] request[%v] error: %v", host, method, request, err)
			return nil, err
		}
		resultBytes, err := JsonEncode(result)
		if err != nil {
			holmes.Error("json encode error: %v", err)
			return nil, err
		}
		//holmes.Debug("jsonrpc result: %s", string(resultBytes))
		return resultBytes, nil
	case BACKEND_NORMAL_HTTP:
		serviceMethod := self.serviceMethodMap[method]
		if serviceMethod == nil {
			holmes.Error("unknown service method[%s] of service[%v]", method, self.serviceMethodMap)
			return nil, ErrorUnkownService
		}
		result, err := self.normalHttp.Call(host, serviceMethod.Uri, decodeMsg)
		if err != nil {
			holmes.Error("normalhttp call host[%s] method[%s] request[%s] error: %v", host, method, string(decodeMsg), err)
			return nil, err
		}
		return result, nil
	}

	return nil, nil
}
