package api

import (
	"context"
	"encoding/base64"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/user"

	oauthapi "github.com/openshift/api/oauth/v1"
)

const (
	// IdentityDisplayNameKey is the key for an optional display name in an identity's Extra map
	IdentityDisplayNameKey = "name"
	// IdentityEmailKey is the key for an optional email address in an identity's Extra map
	IdentityEmailKey = "email"
	// IdentityPreferredUsernameKey is the key for an optional preferred username in an identity's Extra map.
	// This is useful when the immutable providerUserName is different than the login used to authenticate
	// If present, this extra value is used as the preferred username
	IdentityPreferredUsernameKey = "preferred_username"
)

// UserIdentityInfo contains information about an identity.  Identities are distinct from users.  An authentication server of
// some kind (like oauth for example) describes an identity.  Our system controls the users mapped to this identity.
type UserIdentityInfo interface {
	// GetIdentityName returns the name of this identity. It must be equal to GetProviderName() + ":" + GetProviderUserName()
	GetIdentityName() string
	// GetProviderName returns the name of the provider of this identity.
	GetProviderName() string
	// GetProviderUserName uniquely identifies this particular identity for this provider.  It is NOT guaranteed to be unique across providers
	GetProviderUserName() string
	// GetProviderGroups returns the names of the groups for this identity
	GetProviderGroups() []string
	// GetExtra is a map to allow providers to add additional fields that they understand
	GetExtra() map[string]string
	// GetProviderPreferredUserName is a shortcut to retrieve the preferred username from the Extra map
	GetProviderPreferredUserName() string
}

// UserIdentityMapper maps UserIdentities into user.Info objects to allow different user abstractions within auth code.
type UserIdentityMapper interface {
	// UserFor takes an identity, ignores the passed identity.Provider, forces the provider value to some other value and then creates the mapping.
	// It returns the corresponding user.Info
	UserFor(identityInfo UserIdentityInfo) (user.Info, error)
}

type GroupIdentityMapper interface {
	GroupsFor(identityInfo UserIdentityInfo) (user.Info, error)
}

type Client interface {
	GetId() string
	GetSecret() string
	GetRedirectUri() string
	GetUserData() interface{}
}

type Grant struct {
	Client      Client
	Scope       string
	Expiration  int64
	RedirectURI string
}

type DefaultUserIdentityInfo struct {
	ProviderName     string
	ProviderUserName string
	ProviderGroups   []string
	Extra            map[string]string
}

// NewDefaultUserIdentityInfo returns a DefaultUserIdentityInfo with a non-nil Extra component
func NewDefaultUserIdentityInfo(providerName, providerUserName string) *DefaultUserIdentityInfo {
	userName := urlEncodeIfNecessary(providerUserName)
	return &DefaultUserIdentityInfo{
		ProviderName:     providerName,
		ProviderUserName: userName,
		Extra:            map[string]string{},
	}
}

func (i *DefaultUserIdentityInfo) GetIdentityName() string {
	return i.ProviderName + ":" + i.ProviderUserName
}

func (i *DefaultUserIdentityInfo) GetProviderName() string {
	return i.ProviderName
}

func (i *DefaultUserIdentityInfo) GetProviderUserName() string {
	return i.ProviderUserName
}

func (i *DefaultUserIdentityInfo) GetProviderPreferredUserName() string {
	if preferredUsername := i.Extra[IdentityPreferredUsernameKey]; len(preferredUsername) > 0 {
		return preferredUsername
	}
	return i.ProviderUserName
}

func (i *DefaultUserIdentityInfo) GetProviderGroups() []string {
	return i.ProviderGroups
}

func (i *DefaultUserIdentityInfo) GetExtra() map[string]string {
	return i.Extra
}

func urlEncodeIfNecessary(s string) string {
	if strings.ContainsAny(s, ":/") {
		return "b64:" + base64.RawStdEncoding.EncodeToString([]byte(s))
	}
	return s
}

// ProviderInfo represents display information for an oauth identity provider.  This is used by the
// selection provider template to render links to login using different identity providers.
type ProviderInfo struct {
	// Name is unique and corresponds to the name of the identity provider in the oauth configuration
	Name string
	// URL to login using this identity provider
	URL string
}

// OAuthClientGetter exposes a way to get a specific client.  This is useful for other registries to get scope limitations
// on particular clients.   This interface will make its easier to write a future cache on it
type OAuthClientGetter interface {
	Get(ctx context.Context, name string, options metav1.GetOptions) (*oauthapi.OAuthClient, error)
}
