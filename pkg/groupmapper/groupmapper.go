package groupmapper

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	kuser "k8s.io/apiserver/pkg/authentication/user"

	userv1 "github.com/openshift/api/user/v1"
	userclient "github.com/openshift/client-go/user/clientset/versioned/typed/user/v1"
	userinformer "github.com/openshift/client-go/user/informers/externalversions/user/v1"
	usercache "github.com/openshift/library-go/pkg/oauth/usercache"

	authapi "github.com/openshift/oauth-server/pkg/api"
)

const (
	groupGeneratedKey = "oauth.openshift.io/generated"
	groupSyncedKeyFmt = "oauth.openshift.io/idp.%s"
)

var _ authapi.UserIdentityMapper = &UserGroupsMapper{}

var _ kuser.Info = &UserInfoGroupsWrapper{}

// UserInfoGroupsWrapper wraps a UserInfo object in order to add extra groups
// retrieved from the identity providers
type UserInfoGroupsWrapper struct {
	userInfo         kuser.Info
	additionalGroups sets.String
}

func (w *UserInfoGroupsWrapper) GetName() string {
	return w.userInfo.GetName()
}

func (w *UserInfoGroupsWrapper) GetUID() string {
	return w.userInfo.GetUID()
}

func (w *UserInfoGroupsWrapper) GetExtra() map[string][]string {
	return w.userInfo.GetExtra()
}

func (w *UserInfoGroupsWrapper) GetGroups() []string {
	groups := w.additionalGroups.Union(sets.NewString(w.userInfo.GetGroups()...))
	return groups.List()
}

// UserGroupsMapper wraps a UserIdentityMapper with a struct that's capable to
// create the groups for a given user based on the provided UserIdentityInfo
type UserGroupsMapper struct {
	delegatedUserMapper authapi.UserIdentityMapper
	groupsClient        userclient.GroupInterface
	groupsCache         *usercache.GroupCache
	groupsSynced        func() bool
}

func NewUserGroupsMapper(delegate authapi.UserIdentityMapper, groupsClient userclient.GroupInterface, groupInformer userinformer.GroupInformer) *UserGroupsMapper {
	return &UserGroupsMapper{
		delegatedUserMapper: delegate,
		groupsClient:        groupsClient,
		groupsCache:         usercache.NewGroupCache(groupInformer),
		groupsSynced:        groupInformer.Informer().HasSynced,
	}
}

func (m *UserGroupsMapper) UserFor(identityInfo authapi.UserIdentityInfo) (kuser.Info, error) {
	userInfo, err := m.delegatedUserMapper.UserFor(identityInfo)
	if err != nil {
		return userInfo, err
	}

	identityGroups := sets.NewString(identityInfo.GetProviderGroups()...)
	if err := m.processGroups(identityInfo.GetProviderName(), identityInfo.GetProviderPreferredUserName(), identityGroups); err != nil {
		return nil, err
	}

	return &UserInfoGroupsWrapper{
		userInfo:         userInfo,
		additionalGroups: identityGroups,
	}, nil
}

func (m *UserGroupsMapper) processGroups(idpName, username string, groups sets.String) error {
	err := wait.PollImmediate(1*time.Second, 5*time.Second, func() (bool, error) {
		return m.groupsSynced(), nil
	})
	if err != nil {
		return err
	}

	cachedGroups, err := m.groupsCache.GroupsFor(username)
	if err != nil {
		return err
	}

	removeGroups, addGroups := groupsDiff(cachedGroups, groups)
	for _, g := range removeGroups {
		if err := m.removeUserFromGroupWithRetries(idpName, username, g); err != nil {
			return err
		}
	}

	for _, g := range addGroups {
		if err := m.addUserToGroupWithRetries(idpName, username, g); err != nil {
			return err
		}
	}

	return nil
}

// FIXME: work better with annotations
// FIXME: add retries
// FIXME: use idpName
func (m *UserGroupsMapper) removeUserFromGroupWithRetries(idpName, username, group string) error {
	updatedGroup, err := m.groupsClient.Get(context.TODO(), group, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if len(updatedGroup.Users) == 0 {
		return nil
	}

	if len(updatedGroup.Users) == 1 && updatedGroup.Users[0] == username && updatedGroup.Annotations[groupGeneratedKey] == "true" {
		return m.groupsClient.Delete(context.TODO(), group, metav1.DeleteOptions{})
	}

	// find the user and remove it from the slice
	userIdx := -1
	for i, groupUser := range updatedGroup.Users {
		if groupUser == username {
			userIdx = i
			break
		}
	}

	var newUsers []string
	switch userIdx {
	case -1:
		return nil
	case 0:
		newUsers = updatedGroup.Users[1:]
	default:
		newUsers = append(updatedGroup.Users[0:userIdx], updatedGroup.Users[userIdx+1:]...)
	}

	updatedGroup.Users = newUsers

	_, err = m.groupsClient.Update(context.TODO(), updatedGroup, metav1.UpdateOptions{})
	return err
}

// FIXME: work better with annotations
// FIXME: add retries
func (m *UserGroupsMapper) addUserToGroupWithRetries(idpName, username, group string) error {
	updatedGroup, err := m.groupsClient.Get(context.TODO(), group, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = m.groupsClient.Create(context.TODO(),
			&userv1.Group{
				ObjectMeta: metav1.ObjectMeta{
					Name: group,
					Annotations: map[string]string{
						fmt.Sprintf(groupSyncedKeyFmt, idpName): "synced",
						groupGeneratedKey:                       "true",
					},
				},
				Users: []string{username},
			},
			metav1.CreateOptions{},
		)
		return err
	}
	if err != nil {
		return err
	}

	for _, u := range updatedGroup.Users {
		if u == username {
			return nil
		}
	}

	updatedGroup.Users = append(updatedGroup.Users, username)

	_, err = m.groupsClient.Update(context.TODO(), updatedGroup, metav1.UpdateOptions{})
	return err
}

func groupsDiff(existing []*userv1.Group, required sets.String) (toRemove, toAdd []string) {
	existingNames := sets.NewString()
	for _, g := range existing {
		existingNames.Insert(g.Name)
	}

	return existingNames.Difference(required).UnsortedList(), required.Difference(existingNames).UnsortedList()
}
