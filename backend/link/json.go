package link

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// rawJSONInput represents the expected structure of JSON input for creating a new link.
type rawJSONInput struct {
	Target      string  `json:"target"`              // Required: destination URL
	SlugLength  int     `json:"slugLength"`          // Required: length of the generated slug
	SlugCharset string  `json:"slugCharset"`         // Required: allowed characters in the slug
	ExpiresAt   string  `json:"expiresAt,omitempty"` // Required: RFC3339 absolute expiration
	Password    *string `json:"password,omitempty"`  // Optional: password to protect the link
	MaxHits     *int    `json:"maxHits,omitempty"`   // Optional: max allowed hits
	ValidFrom   *string `json:"validFrom,omitempty"` // Optional: RFC3339 start time for link validity
}

// FromJSON reads, validates, and converts JSON input into a Validated Link.
//
// It expects JSON with the following fields:
//   - Required: target, slugLength, slugCharset
//   - Required expiration: either expiresAt (RFC3339) or expiresIn (days)
//   - Optional: password, maxHits, validFrom (RFC3339)
//
// Returns a validated link or an error.
func FromJSON(r io.Reader, now time.Time) (*Validated, error) {
	var input rawJSONInput
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	missing := make([]string, 0, 4)
	if input.Target == "" {
		missing = append(missing, "target")
	}
	if input.SlugCharset == "" {
		missing = append(missing, "slugCharset")
	}
	if input.ExpiresAt == "" {
		missing = append(missing, "expiresAt")
	}
	if len(missing) != 0 {
		return nil, errors.New("missing one or more required fields: " + strings.Join(missing, ", "))
	}

	expiresAt, err := time.Parse(time.RFC3339, input.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expiresAt: %w", err)
	}

	var validFrom *time.Time
	if input.ValidFrom != nil && *input.ValidFrom != "" {
		vf, err := time.Parse(time.RFC3339, *input.ValidFrom)
		if err != nil {
			return nil, fmt.Errorf("invalid validFrom: %w", err)
		}
		validFrom = &vf
	}

	maxHits := input.MaxHits

	adminExpiresAt := expiresAt.Add(24 * time.Hour)

	link := &Link{
		ID:             primitive.NewObjectID(),
		Slug:           "",
		AdminToken:     "",
		Target:         input.Target,
		PasswordHash:   nil,
		MaxHits:        maxHits,
		ValidFrom:      validFrom,
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      expiresAt,
		AdminExpiresAt: adminExpiresAt,
		HitCount:       0,
		SchemaVersion:  schemaVersion,
	}

	validated, err := Validate(link, now)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	err = validated.SetPasswordHash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("error generating the hash: %w", err)
	}

	slugLen := input.SlugLength
	err = validated.SetSlug(slugLen, strings.ToLower(input.SlugCharset))
	if err != nil {
		return nil, fmt.Errorf("error generating slug: %w", err)
	}

	err = validated.SetAdminToken()
	if err != nil {
		return nil, fmt.Errorf("error generating admin token: %w", err)
	}

	return validated, nil
}

// rawJSONPatch represents incoming PATCH JSON data.
// Each field is a double pointer to distinguish:
//   - Omitted fields: nil
//   - Explicit nulls: *T == nil
//   - Set values: **T
type rawJSONPatch struct {
	ExpiresAt **time.Time `json:"expiresAt"`
	MaxHits   **int       `json:"maxHits"`
	ValidFrom **time.Time `json:"validFrom"`
	Password  **string    `json:"password"`
}

// PatchFromJSON applies partial JSON updates to a Link.
//
// Accepts a JSON payload with any combination of:
//   - expiresAt: null (remove) or timestamp (update)
//   - maxHits: null (remove) or integer (update)
//   - validFrom: null (remove) or timestamp (update)
//   - password: null (remove) or string (update)
//
// Fields not provided in the JSON will not be changed.
// Returns a validated patch or an error.
func PatchFromJSON(data []byte, original *Link) (*ValidatedPatch, error) {
	var raw rawJSONPatch
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	now := time.Now()
	patch := NewPatchLink(now)

	if raw.ExpiresAt != nil {
		if *raw.ExpiresAt == nil {
			return nil, errors.New("expiresAt cannot be null")
		}
		patch.ExpiresAt = *raw.ExpiresAt

		adminExpiresAt := (*raw.ExpiresAt).Add(24 * time.Hour)
		patch.AdminExpiresAt = &adminExpiresAt
	}

	if raw.MaxHits != nil {
		if *raw.MaxHits == nil {
			patch.MaxHits.Remove = true
		} else {
			patch.MaxHits.Value = *raw.MaxHits
		}
	}

	if raw.ValidFrom != nil {
		if *raw.ValidFrom == nil {
			patch.ValidFrom.Remove = true
		} else {
			patch.ValidFrom.Value = *raw.ValidFrom
		}
	}

	validated, err := ValidatePatch(original, patch, now)
	if err != nil {
		return nil, err
	}

	if raw.Password != nil {
		if *raw.Password == nil {
			patch.PasswordHash.Remove = true
		} else {
			err = validated.SetPasswordHash(*raw.Password)
			if err != nil {
				return nil, ErrHashingPassword
			}
		}
	}

	return validated, nil
}
