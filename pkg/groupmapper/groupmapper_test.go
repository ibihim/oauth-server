package groupmapper

import (
	"context"
	"fmt"
	"testing"
	"time"

	userv1 "github.com/openshift/api/user/v1"
	fakeuserclient "github.com/openshift/client-go/user/clientset/versioned/fake"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/watch"
)

const testIDPName = "test-idp"

func TestUserGroupsMapper_removeUserFromGroupWithRetries(t *testing.T) {
	const testGroupName = "test-group"

	tests := []struct {
		name          string
		username      string
		group         *userv1.Group
		expectedGroup *userv1.Group
		expectEvent   bool
		wantErr       bool
		wantDeletion  bool
	}{
		{
			name:          "nonexistent group",
			username:      "user1",
			group:         nil,
			expectedGroup: nil,
		},
		{
			name:          "user not in target group",
			username:      "user1",
			group:         createGroupWithUsers(testGroupName, "user2", "user3", "user4"),
			expectedGroup: createGroupWithUsers(testGroupName, "user2", "user3", "user4"),
		},
		{
			name:          "first user gets removed from the group",
			username:      "user1",
			group:         createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
			expectedGroup: createGroupWithUsers(testGroupName, "user2", "user3", "user4"),
			expectEvent:   true,
		},
		{
			name:          "mid user gets removed from the group",
			username:      "user2",
			group:         createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
			expectedGroup: createGroupWithUsers(testGroupName, "user1", "user3", "user4"),
			expectEvent:   true,
		},
		{
			name:          "last user gets removed from the group",
			username:      "user4",
			group:         createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
			expectedGroup: createGroupWithUsers(testGroupName, "user1", "user2", "user3"),
			expectEvent:   true,
		},
		{
			name:          "last user on group w/o generated annotation -> empty users",
			username:      "user1",
			group:         removeGeneratedKeyFromGroup(createGroupWithUsers(testGroupName, "user1")),
			expectedGroup: removeGeneratedKeyFromGroup(createGroupWithUsers(testGroupName)),
			expectEvent:   true,
		},
		{
			name:          "last user on group -> delete group",
			username:      "user1",
			group:         createGroupWithUsers(testGroupName, "user1"),
			expectedGroup: createGroupWithUsers(testGroupName, "user1"), // this is the obj being deleted
			wantDeletion:  true,
			expectEvent:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := []runtime.Object{}
			if tt.group != nil {
				groups = append(groups, tt.group)
			}
			fakeUserClient := fakeuserclient.NewSimpleClientset(groups...)
			testCtx := context.Background()
			groupWatcher, err := fakeUserClient.UserV1().Groups().Watch(testCtx, metav1.ListOptions{})
			require.NoError(t, err)
			defer groupWatcher.Stop()

			m := &UserGroupsMapper{
				groupsClient: fakeUserClient.UserV1().Groups(),
			}

			finished, failed := make(chan bool), make(chan string)
			timedCtx, timedCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)

			expectedEventType := watch.Modified
			if tt.wantDeletion {
				expectedEventType = watch.Deleted
			}
			go watchForGroupEvents(groupWatcher, tt.expectedGroup, tt.expectEvent, expectedEventType, failed, finished, timedCtx)

			go func() {
				if err := m.removeUserFromGroupWithRetries(testIDPName, tt.username, testGroupName); (err != nil) != tt.wantErr {
					t.Errorf("UserGroupsMapper.removeUserFromGroupWithRetries() error = %v, wantErr %v", err, tt.wantErr)
				}

				// give the watch some time
				time.Sleep(1 * time.Second)
				timedCtxCancel()
			}()

			select {
			case <-finished:
			case errMsg := <-failed:
				t.Fatal(errMsg)
			}

		})
	}
}

func TestUserGroupsMapper_addUserToGroupWithRetries(t *testing.T) {
	const testGroupName = "test-group"

	tests := []struct {
		name          string
		username      string
		group         *userv1.Group
		expectedGroup *userv1.Group
		expectEvent   bool
		wantErr       bool
	}{
		{
			name:          "nonexistent group",
			username:      "user1",
			group:         nil,
			expectedGroup: createGroupWithUsers(testGroupName, "user1"),
			expectEvent:   true,
		},
		{
			name:          "user already in group",
			username:      "user3",
			group:         createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
			expectedGroup: createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
		},
		// {
		// 	name:          "user in group w/o the synced annotation", // TODO: check the enhancement on the expected behavior
		// 	username:      "user2",
		// 	group:         removeSyncedKeyFromGroup(createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"), testIDPName),
		// 	expectedGroup: createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4"),
		// 	expectEvent:   true,
		// },
		{
			name:          "user missing in the group",
			username:      "user99",
			group:         createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4", "user5", "user6"),
			expectedGroup: createGroupWithUsers(testGroupName, "user1", "user2", "user3", "user4", "user5", "user6", "user99"),
			expectEvent:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := []runtime.Object{}
			if tt.group != nil {
				groups = append(groups, tt.group)
			}
			fakeUserClient := fakeuserclient.NewSimpleClientset(groups...)
			testCtx := context.Background()
			groupWatcher, err := fakeUserClient.UserV1().Groups().Watch(testCtx, metav1.ListOptions{})
			require.NoError(t, err)
			defer groupWatcher.Stop()

			m := &UserGroupsMapper{
				groupsClient: fakeUserClient.UserV1().Groups(),
			}

			finished, failed := make(chan bool), make(chan string)
			timedCtx, timedCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)

			expectedEventType := watch.Modified
			if tt.group == nil && tt.expectedGroup != nil {
				expectedEventType = watch.Added
			}
			go watchForGroupEvents(groupWatcher, tt.expectedGroup, tt.expectEvent, expectedEventType, failed, finished, timedCtx)

			go func() {
				if err := m.addUserToGroupWithRetries(testIDPName, tt.username, testGroupName); (err != nil) != tt.wantErr {
					t.Errorf("UserGroupsMapper.addUserToGroupWithRetries() error = %v, wantErr %v", err, tt.wantErr)
				}

				// give the watch some time
				time.Sleep(1 * time.Second)
				timedCtxCancel()
			}()

			select {
			case <-finished:
			case errMsg := <-failed:
				t.Fatal(errMsg)
			}
		})
	}
}

func createGroupWithUsers(groupname string, users ...string) *userv1.Group {
	return &userv1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: groupname,
			Annotations: map[string]string{
				fmt.Sprintf(groupSyncedKeyFmt, testIDPName): "synced",
				groupGeneratedKey: "true",
			},
		},
		Users: users,
	}
}

func removeGeneratedKeyFromGroup(g *userv1.Group) *userv1.Group {
	delete(g.Annotations, groupGeneratedKey)
	return g
}

func removeSyncedKeyFromGroup(g *userv1.Group, idpName string) *userv1.Group {
	delete(g.Annotations, fmt.Sprintf(groupSyncedKeyFmt, idpName))
	return g
}

func watchForGroupEvents(
	groupWatcher watch.Interface,
	expectedGroup *userv1.Group,
	expectEvent bool,
	expectedEventType watch.EventType,
	failChan chan<- string,
	finishChan chan<- bool,
	timeOutCtx context.Context,
) {
	userChan := groupWatcher.ResultChan()
	eventCount := 0
	for {
		select {
		case groupEvent := <-userChan:
			if !expectEvent {
				failChan <- fmt.Sprintf("unexpected event: %v", groupEvent)
				return
			}

			eventCount++
			group, ok := groupEvent.Object.(*userv1.Group)
			if !ok {
				failChan <- "the retrieved object was not a group"
				return
			}
			if !equality.Semantic.DeepEqual(expectedGroup, group) {
				failChan <- fmt.Sprintf("the expected group is different from the actual: %s", diff.ObjectDiff(expectedGroup, group))
				return
			}
			if expectedEventType != groupEvent.Type {
				failChan <- fmt.Sprintf("expected event of type %s, got %s", expectedEventType, groupEvent.Type)
				return
			}
		case <-timeOutCtx.Done():
			if eventCount == 0 && expectEvent {
				failChan <- "timed out"
			} else {
				finishChan <- true
			}
			return
		}
	}
}
