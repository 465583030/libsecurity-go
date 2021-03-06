package entityRestful

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	defs "github.com/ibm-security-innovation/libsecurity-go/defs"
	ent "github.com/ibm-security-innovation/libsecurity-go/entity"
	acl "github.com/ibm-security-innovation/libsecurity-go/acl"
	cr "github.com/ibm-security-innovation/libsecurity-go/restful/common-restful"
	"github.com/ibm-security-innovation/libsecurity-go/restful/libsecurity-restful"
)

const (
	umPrefix          = "/entity"
	usersPath         = "/users"
	groupsPath        = "/groups"
	resourcesPath     = "/resources"
	permissionsPath     = "/permissions"
	groupIDToken      = "groups"
	groupIDParam      = "group-name"
	groupIDComment    = "identifier of the group"
	userIDToken       = "users"
	userIDParam       = "user-name"
	resourceIDParam   = "resource-name"
	userIDComment     = "identifier of the user"
	resourceIDComment = "identifier of the resource"
	permissionIDToken       = "permissions"
	permissionIDParam       = "permission"
	permissionIDComment = "permission"

	originToken = "Origin"
)

var (
	enServicePath       string //  = cr.ServicePathPrefix + "/um"
	usersServicePath    string // = enServicePath + usersPath
	groupServicePath    string // = enServicePath + groupsPath
	resourceServicePath string // = enServicePath + ResourcePath
	permissionServicePath string // = enServicePath + PermissionPath
)

// EnRestful : Entity structure
type EnRestful struct {
	st *libsecurityRestful.LibsecurityRestful
}

func init() {
	enServicePath = cr.ServicePathPrefix + umPrefix
	usersServicePath = enServicePath + usersPath
	groupServicePath = enServicePath + groupsPath
	resourceServicePath = enServicePath + resourcesPath
	permissionServicePath = enServicePath + permissionsPath

	initCommandToPath()
}

// NewEnRestful : return a pointer to the EnRestful structure
func NewEnRestful() *EnRestful {
	return &EnRestful{}
}

// SetData : initialize the entity structure
func (en *EnRestful) SetData(stR *libsecurityRestful.LibsecurityRestful) {
	en.st = stR
}

func (en EnRestful) getGroupURLPath(request *restful.Request, name string) cr.URL {
	return cr.URL{URL: fmt.Sprintf("%v%v/%v", enServicePath, groupsPath, name)}
}

func (en EnRestful) getUserURLPath(request *restful.Request, name string) cr.URL {
	return cr.URL{URL: fmt.Sprintf("%v%v/%v", enServicePath, usersPath, name)}
}

func (en EnRestful) getResourceURLPath(request *restful.Request, name string) cr.URL {
	return cr.URL{URL: fmt.Sprintf("%v%v/%v", enServicePath, resourcesPath, name)}
}

func (en EnRestful) getPermissionURLPath(request *restful.Request, name string) cr.URL {
	return cr.URL{URL: fmt.Sprintf("%v%v/%v", enServicePath, permissionsPath, name)}
}

func (en EnRestful) setError(response *restful.Response, httpStatusCode int, err error) {
	data, _ := json.Marshal(cr.Error{Code: httpStatusCode, Message: fmt.Sprintf("%v", err)})
	response.WriteErrorString(httpStatusCode, string(data))
}

func (en *EnRestful) restCreateGroup(request *restful.Request, response *restful.Response) {
	groupID := request.PathParameter(groupIDParam)
	err := en.st.UsersList.AddGroup(groupID)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, en.getGroupURLPath(request, groupID))
}

func (en EnRestful) getAllGroups() []string {
	var gList []string
	for name := range en.st.UsersList.Groups {
		gList = append(gList, name)
	}
	return gList
}

func (en EnRestful) restRemoveAllGroups(request *restful.Request, response *restful.Response) {
	gList := en.getAllGroups()
	for _, name := range gList {
		en.st.UsersList.RemoveGroup(name)
	}
	response.WriteHeader(http.StatusNoContent)
}

func (en EnRestful) restGetGroup(request *restful.Request, response *restful.Response) {
	groupID := request.PathParameter(groupIDParam)
	group, exist := en.st.UsersList.Groups[groupID]
	if exist == false {
		en.setError(response, http.StatusNotFound, fmt.Errorf("Group "+groupID+" could not be found."))
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, group)
}

func (en EnRestful) restGetAllGroups(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, en.st.UsersList.Groups)
}

func (en *EnRestful) restRemoveGroup(request *restful.Request, response *restful.Response) {
	groupID := request.PathParameter(groupIDParam)
	err := en.st.UsersList.RemoveGroup(groupID)
	if err != nil {
		en.setError(response, http.StatusNotFound, err)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func (en *EnRestful) restAddUserToGroup(request *restful.Request, response *restful.Response) {
	groupID := request.PathParameter(groupIDParam)
	userID := request.PathParameter(userIDParam)
	err := en.st.UsersList.AddUserToGroup(groupID, userID)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, en.getGroupURLPath(request, fmt.Sprintf("%v/%v/%v/%v", groupIDToken, groupID, userIDToken, userID)))
}

func (en *EnRestful) restRemoveUserFromGroup(request *restful.Request, response *restful.Response) {
	groupID := request.PathParameter(groupIDParam)
	userID := request.PathParameter(userIDParam)
	err := en.st.UsersList.RemoveUserFromGroup(groupID, userID)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (en *EnRestful) restCreateUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(userIDParam)
	err := en.st.UsersList.AddUser(id)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, en.getUserURLPath(request, id))
}

func (en EnRestful) restGetEntityManager(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, en.st.UsersList)
}

func (en EnRestful) restGetUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(userIDParam)
	user, exist := en.st.UsersList.Users[id]
	if exist == false {
		en.setError(response, http.StatusNotFound, fmt.Errorf("User "+id+" could not be found."))
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, user)
}

func (en EnRestful) restGetAllUsers(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, en.st.UsersList.Users)
}

func (en *EnRestful) restRemoveAllUsers(request *restful.Request, response *restful.Response) {
	for name := range en.st.UsersList.Users {
		if name == defs.RootUserName || name == defs.AclAllEntryName {
			continue
		}
		en.st.UsersList.RemoveUser(name)
	}
	response.WriteHeader(http.StatusNoContent)
}

func (en *EnRestful) restRemoveUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(userIDParam)
	err := en.st.UsersList.RemoveUser(id)
	if err != nil {
		en.setError(response, http.StatusNotFound, err)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func (en *EnRestful) restCreateResource(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(resourceIDParam)
	err := en.st.UsersList.AddResource(id)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	newAcl := acl.NewACL()
	err = en.st.UsersList.AddPropertyToEntity(id, defs.AclPropertyName, newAcl)
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, en.getResourceURLPath(request, id))
}

func (en EnRestful) restGetResource(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(resourceIDParam)
	user, exist := en.st.UsersList.Resources[id]
	if exist == false {
		en.setError(response, http.StatusNotFound, fmt.Errorf("Resource "+id+" could not be found."))
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, user)
}

func (en EnRestful) restGetAllResources(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, en.st.UsersList.Resources)
}

func (en *EnRestful) restRemoveAllResources(request *restful.Request, response *restful.Response) {
	for name := range en.st.UsersList.Resources {
		en.st.UsersList.RemoveResource(name)
	}
	response.WriteHeader(http.StatusNoContent)
}

func (en *EnRestful) restRemoveResource(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(resourceIDParam)
	err := en.st.UsersList.RemoveResource(id)
	if err != nil {
		en.setError(response, http.StatusNotFound, err)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func (en *EnRestful) restCreatePermission(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(permissionIDParam)
	err := en.st.UsersList.AddPermission(ent.Permission(id))
	if err != nil {
		en.setError(response, http.StatusPreconditionFailed, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusCreated, en.getPermissionURLPath(request, id))
}

func (en EnRestful) restGetAllPermissions(request *restful.Request, response *restful.Response) {
	response.WriteHeaderAndEntity(http.StatusOK, en.st.UsersList.Permissions)
}

func (en *EnRestful) restRemoveAllPermissions(request *restful.Request, response *restful.Response) {
	for name := range en.st.UsersList.Permissions {
		en.st.UsersList.RemovePermission(name)
	}
	response.WriteHeader(http.StatusNoContent)
}

func (en *EnRestful) restRemovePermission(request *restful.Request, response *restful.Response) {
	id := request.PathParameter(permissionIDParam)
	err := en.st.UsersList.RemovePermission(ent.Permission(id))
	if err != nil {
		en.setError(response, http.StatusNotFound, err)
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}
