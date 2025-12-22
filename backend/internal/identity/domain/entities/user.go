package entities

import (
	"errors"
	"time"

	"gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
)

// User represents a user aggregate root in the Identity context.
type User struct {
	id           valueobjects.UserID
	email        valueobjects.Email
	passwordHash valueobjects.PasswordHash
	name         valueobjects.UserName
	createdAt    time.Time
	updatedAt    time.Time
	isActive     bool

	// Domain events
	events []events.DomainEvent
}

// NewUser creates a new User aggregate.
func NewUser(
	email valueobjects.Email,
	passwordHash valueobjects.PasswordHash,
	name valueobjects.UserName,
) (*User, error) {
	if email.IsEmpty() {
		return nil, errors.New("email cannot be empty")
	}

	if passwordHash.IsEmpty() {
		return nil, errors.New("password hash cannot be empty")
	}

	if name.IsEmpty() {
		return nil, errors.New("name cannot be empty")
	}

	now := time.Now()

	user := &User{
		id:           valueobjects.GenerateUserID(),
		email:        email,
		passwordHash: passwordHash,
		name:         name,
		createdAt:    now,
		updatedAt:    now,
		isActive:     true,
		events:       []events.DomainEvent{},
	}

	// Add domain event
	user.addEvent(events.NewBaseDomainEvent(
		"UserRegistered",
		user.id.Value(),
		"User",
	))

	return user, nil
}

// ID returns the user ID.
func (u *User) ID() valueobjects.UserID {
	return u.id
}

// Email returns the user email.
func (u *User) Email() valueobjects.Email {
	return u.email
}

// Name returns the user name.
func (u *User) Name() valueobjects.UserName {
	return u.name
}

// IsActive returns whether the user is active.
func (u *User) IsActive() bool {
	return u.isActive
}

// CreatedAt returns the creation timestamp.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// ChangePassword changes the user's password.
// It verifies the old password before setting the new one.
func (u *User) ChangePassword(oldPassword, newPassword string) error {
	if !u.isActive {
		return errors.New("cannot change password for inactive user")
	}

	if !u.passwordHash.Verify(oldPassword) {
		return errors.New("invalid old password")
	}

	newHash, err := valueobjects.NewPasswordHashFromPlain(newPassword)
	if err != nil {
		return err
	}

	u.passwordHash = newHash
	u.updatedAt = time.Now()

	u.addEvent(events.NewBaseDomainEvent(
		"UserPasswordChanged",
		u.id.Value(),
		"User",
	))

	return nil
}

// UpdateName updates the user's name.
func (u *User) UpdateName(name valueobjects.UserName) error {
	if !u.isActive {
		return errors.New("cannot update name for inactive user")
	}

	if name.IsEmpty() {
		return errors.New("name cannot be empty")
	}

	u.name = name
	u.updatedAt = time.Now()

	u.addEvent(events.NewBaseDomainEvent(
		"UserNameUpdated",
		u.id.Value(),
		"User",
	))

	return nil
}

// UpdateEmail updates the user's email.
func (u *User) UpdateEmail(email valueobjects.Email) error {
	if !u.isActive {
		return errors.New("cannot update email for inactive user")
	}

	if email.IsEmpty() {
		return errors.New("email cannot be empty")
	}

	u.email = email
	u.updatedAt = time.Now()

	u.addEvent(events.NewBaseDomainEvent(
		"UserEmailUpdated",
		u.id.Value(),
		"User",
	))

	return nil
}

// Deactivate deactivates the user.
func (u *User) Deactivate() error {
	if !u.isActive {
		return errors.New("user is already inactive")
	}

	u.isActive = false
	u.updatedAt = time.Now()

	u.addEvent(events.NewBaseDomainEvent(
		"UserDeactivated",
		u.id.Value(),
		"User",
	))

	return nil
}

// Activate activates the user.
func (u *User) Activate() error {
	if u.isActive {
		return errors.New("user is already active")
	}

	u.isActive = true
	u.updatedAt = time.Now()

	u.addEvent(events.NewBaseDomainEvent(
		"UserActivated",
		u.id.Value(),
		"User",
	))

	return nil
}

// VerifyPassword verifies if the provided password matches the user's password hash.
func (u *User) VerifyPassword(password string) bool {
	return u.passwordHash.Verify(password)
}

// GetEvents returns all domain events that occurred on this aggregate.
func (u *User) GetEvents() []events.DomainEvent {
	return u.events
}

// ClearEvents clears all domain events from this aggregate.
func (u *User) ClearEvents() {
	u.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (u *User) addEvent(event events.DomainEvent) {
	u.events = append(u.events, event)
}
