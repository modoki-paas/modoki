// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "Modoki API": Application User Types
//
// Command:
// $ goagen
// --design=github.com/cs3238-tsuzu/modoki/design
// --out=$(GOPATH)/src/github.com/cs3238-tsuzu/modoki
// --version=v1.3.1

package app

import (
	"mime/multipart"
	"unicode/utf8"

	"github.com/goadesign/goa"
)

// containerConfig user type.
type containerConfig struct {
	DefaultShell *string `form:"defaultShell,omitempty" json:"defaultShell,omitempty" yaml:"defaultShell,omitempty" xml:"defaultShell,omitempty"`
}

// Publicize creates ContainerConfig from containerConfig
func (ut *containerConfig) Publicize() *ContainerConfig {
	var pub ContainerConfig
	if ut.DefaultShell != nil {
		pub.DefaultShell = ut.DefaultShell
	}
	return &pub
}

// ContainerConfig user type.
type ContainerConfig struct {
	DefaultShell *string `form:"defaultShell,omitempty" json:"defaultShell,omitempty" yaml:"defaultShell,omitempty" xml:"defaultShell,omitempty"`
}

// uploadPayload user type.
type uploadPayload struct {
	// Allow for a existing directory to be replaced by a file
	AllowOverwrite *bool `form:"allowOverwrite,omitempty" json:"allowOverwrite,omitempty" yaml:"allowOverwrite,omitempty" xml:"allowOverwrite,omitempty"`
	// Copy all uid/gid information
	CopyUIDGID *bool `form:"copyUIDGID,omitempty" json:"copyUIDGID,omitempty" yaml:"copyUIDGID,omitempty" xml:"copyUIDGID,omitempty"`
	// File tar archive
	Data *multipart.FileHeader `form:"data,omitempty" json:"data,omitempty" yaml:"data,omitempty" xml:"data,omitempty"`
	// Path in the container to save files
	Path *string `form:"path,omitempty" json:"path,omitempty" yaml:"path,omitempty" xml:"path,omitempty"`
}

// Finalize sets the default values for uploadPayload type instance.
func (ut *uploadPayload) Finalize() {
	var defaultAllowOverwrite = false
	if ut.AllowOverwrite == nil {
		ut.AllowOverwrite = &defaultAllowOverwrite
	}
	var defaultCopyUIDGID = false
	if ut.CopyUIDGID == nil {
		ut.CopyUIDGID = &defaultCopyUIDGID
	}
}

// Validate validates the uploadPayload type instance.
func (ut *uploadPayload) Validate() (err error) {
	if ut.Path == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "path"))
	}
	if ut.Data == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "data"))
	}
	if ut.CopyUIDGID == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "copyUIDGID"))
	}
	return
}

// Publicize creates UploadPayload from uploadPayload
func (ut *uploadPayload) Publicize() *UploadPayload {
	var pub UploadPayload
	if ut.AllowOverwrite != nil {
		pub.AllowOverwrite = *ut.AllowOverwrite
	}
	if ut.CopyUIDGID != nil {
		pub.CopyUIDGID = *ut.CopyUIDGID
	}
	if ut.Data != nil {
		pub.Data = ut.Data
	}
	if ut.Path != nil {
		pub.Path = *ut.Path
	}
	return &pub
}

// UploadPayload user type.
type UploadPayload struct {
	// Allow for a existing directory to be replaced by a file
	AllowOverwrite bool `form:"allowOverwrite" json:"allowOverwrite" yaml:"allowOverwrite" xml:"allowOverwrite"`
	// Copy all uid/gid information
	CopyUIDGID bool `form:"copyUIDGID" json:"copyUIDGID" yaml:"copyUIDGID" xml:"copyUIDGID"`
	// File tar archive
	Data *multipart.FileHeader `form:"data" json:"data" yaml:"data" xml:"data"`
	// Path in the container to save files
	Path string `form:"path" json:"path" yaml:"path" xml:"path"`
}

// Validate validates the UploadPayload type instance.
func (ut *UploadPayload) Validate() (err error) {
	if ut.Path == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "path"))
	}

	return
}

// userAuthorizedKey user type.
type userAuthorizedKey struct {
	Key   *string `form:"key,omitempty" json:"key,omitempty" yaml:"key,omitempty" xml:"key,omitempty"`
	Label *string `form:"label,omitempty" json:"label,omitempty" yaml:"label,omitempty" xml:"label,omitempty"`
}

// Validate validates the userAuthorizedKey type instance.
func (ut *userAuthorizedKey) Validate() (err error) {
	if ut.Key == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "key"))
	}
	if ut.Label == nil {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`request`, "label"))
	}
	if ut.Key != nil {
		if utf8.RuneCountInString(*ut.Key) > 2048 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`request.key`, *ut.Key, utf8.RuneCountInString(*ut.Key), 2048, false))
		}
	}
	if ut.Label != nil {
		if ok := goa.ValidatePattern(`^[a-zA-Z0-9_]+$`, *ut.Label); !ok {
			err = goa.MergeErrors(err, goa.InvalidPatternError(`request.label`, *ut.Label, `^[a-zA-Z0-9_]+$`))
		}
	}
	if ut.Label != nil {
		if utf8.RuneCountInString(*ut.Label) < 1 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`request.label`, *ut.Label, utf8.RuneCountInString(*ut.Label), 1, true))
		}
	}
	if ut.Label != nil {
		if utf8.RuneCountInString(*ut.Label) > 32 {
			err = goa.MergeErrors(err, goa.InvalidLengthError(`request.label`, *ut.Label, utf8.RuneCountInString(*ut.Label), 32, false))
		}
	}
	return
}

// Publicize creates UserAuthorizedKey from userAuthorizedKey
func (ut *userAuthorizedKey) Publicize() *UserAuthorizedKey {
	var pub UserAuthorizedKey
	if ut.Key != nil {
		pub.Key = *ut.Key
	}
	if ut.Label != nil {
		pub.Label = *ut.Label
	}
	return &pub
}

// UserAuthorizedKey user type.
type UserAuthorizedKey struct {
	Key   string `form:"key" json:"key" yaml:"key" xml:"key"`
	Label string `form:"label" json:"label" yaml:"label" xml:"label"`
}

// Validate validates the UserAuthorizedKey type instance.
func (ut *UserAuthorizedKey) Validate() (err error) {
	if ut.Key == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "key"))
	}
	if ut.Label == "" {
		err = goa.MergeErrors(err, goa.MissingAttributeError(`type`, "label"))
	}
	if utf8.RuneCountInString(ut.Key) > 2048 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`type.key`, ut.Key, utf8.RuneCountInString(ut.Key), 2048, false))
	}
	if ok := goa.ValidatePattern(`^[a-zA-Z0-9_]+$`, ut.Label); !ok {
		err = goa.MergeErrors(err, goa.InvalidPatternError(`type.label`, ut.Label, `^[a-zA-Z0-9_]+$`))
	}
	if utf8.RuneCountInString(ut.Label) < 1 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`type.label`, ut.Label, utf8.RuneCountInString(ut.Label), 1, true))
	}
	if utf8.RuneCountInString(ut.Label) > 32 {
		err = goa.MergeErrors(err, goa.InvalidLengthError(`type.label`, ut.Label, utf8.RuneCountInString(ut.Label), 32, false))
	}
	return
}

// userConfig user type.
type userConfig struct {
	AuthorizedKeys []*userAuthorizedKey `form:"authorizedKeys,omitempty" json:"authorizedKeys,omitempty" yaml:"authorizedKeys,omitempty" xml:"authorizedKeys,omitempty"`
	DefaultShell   *string              `form:"defaultShell,omitempty" json:"defaultShell,omitempty" yaml:"defaultShell,omitempty" xml:"defaultShell,omitempty"`
}

// Validate validates the userConfig type instance.
func (ut *userConfig) Validate() (err error) {
	for _, e := range ut.AuthorizedKeys {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// Publicize creates UserConfig from userConfig
func (ut *userConfig) Publicize() *UserConfig {
	var pub UserConfig
	if ut.AuthorizedKeys != nil {
		pub.AuthorizedKeys = make([]*UserAuthorizedKey, len(ut.AuthorizedKeys))
		for i2, elem2 := range ut.AuthorizedKeys {
			pub.AuthorizedKeys[i2] = elem2.Publicize()
		}
	}
	if ut.DefaultShell != nil {
		pub.DefaultShell = ut.DefaultShell
	}
	return &pub
}

// UserConfig user type.
type UserConfig struct {
	AuthorizedKeys []*UserAuthorizedKey `form:"authorizedKeys,omitempty" json:"authorizedKeys,omitempty" yaml:"authorizedKeys,omitempty" xml:"authorizedKeys,omitempty"`
	DefaultShell   *string              `form:"defaultShell,omitempty" json:"defaultShell,omitempty" yaml:"defaultShell,omitempty" xml:"defaultShell,omitempty"`
}

// Validate validates the UserConfig type instance.
func (ut *UserConfig) Validate() (err error) {
	for _, e := range ut.AuthorizedKeys {
		if e != nil {
			if err2 := e.Validate(); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}
