/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package issuer

import (
	"encoding/json"
	"fmt"

	"github.com/trustbloc/edge-core/pkg/storage"
)

const (
	keyPattern       = "%s_%s"
	profileKeyPrefix = "profile"

	storeName = "issuer"
)

// Profile db operation.
type Profile struct {
	store storage.Store
}

// ProfileData struct for profile.
type ProfileData struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	CallbackURL string `json:"callbackURL"`
}

// New returns new issuer profile instance.
func New(provider storage.Provider) (*Profile, error) {
	err := provider.CreateStore(storeName)
	if err != nil && err != storage.ErrDuplicateStore {
		return nil, err
	}

	store, err := provider.OpenStore(storeName)
	if err != nil {
		return nil, err
	}

	return &Profile{store: store}, nil
}

// SaveProfile saves the profile data.
func (c *Profile) SaveProfile(data *ProfileData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("issuer profile save - marshalling error: %s", err.Error())
	}

	return c.store.Put(getDBKey(data.ID), bytes)
}

// GetProfile retrieves the profile data based on id.
func (c *Profile) GetProfile(id string) (*ProfileData, error) {
	bytes, err := c.store.Get(getDBKey(id))
	if err != nil {
		return nil, err
	}

	response := &ProfileData{}

	err = json.Unmarshal(bytes, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getDBKey(id string) string {
	return fmt.Sprintf(keyPattern, profileKeyPrefix, id)
}
