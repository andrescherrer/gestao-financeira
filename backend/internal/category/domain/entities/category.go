package entities

import (
	"errors"
	"time"

	"gestao-financeira/backend/internal/category/domain/valueobjects"
	identityvalueobjects "gestao-financeira/backend/internal/identity/domain/valueobjects"
	"gestao-financeira/backend/internal/shared/domain/events"
)

// Category represents a category aggregate root in the Category context.
type Category struct {
	id          valueobjects.CategoryID
	userID      identityvalueobjects.UserID
	name        valueobjects.CategoryName
	description string
	createdAt   time.Time
	updatedAt   time.Time
	isActive    bool

	// Domain events
	events []events.DomainEvent
}

// NewCategory creates a new Category aggregate.
func NewCategory(
	userID identityvalueobjects.UserID,
	name valueobjects.CategoryName,
	description string,
) (*Category, error) {
	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if name.IsEmpty() {
		return nil, errors.New("category name cannot be empty")
	}

	// Validate description length
	if len(description) > 500 {
		return nil, errors.New("category description is too long (max 500 characters)")
	}

	now := time.Now()

	category := &Category{
		id:          valueobjects.GenerateCategoryID(),
		userID:      userID,
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
		isActive:    true,
		events:      []events.DomainEvent{},
	}

	// Add domain event
	category.addEvent(events.NewBaseDomainEvent(
		"CategoryCreated",
		category.id.Value(),
		"Category",
	))

	return category, nil
}

// CategoryFromPersistence reconstructs a Category aggregate from persisted data.
// This method does not trigger domain events, as it's used for loading existing data.
func CategoryFromPersistence(
	id valueobjects.CategoryID,
	userID identityvalueobjects.UserID,
	name valueobjects.CategoryName,
	description string,
	createdAt time.Time,
	updatedAt time.Time,
	isActive bool,
) (*Category, error) {
	if id.IsEmpty() {
		return nil, errors.New("category ID cannot be empty")
	}

	if userID.IsEmpty() {
		return nil, errors.New("user ID cannot be empty")
	}

	if name.IsEmpty() {
		return nil, errors.New("category name cannot be empty")
	}

	return &Category{
		id:          id,
		userID:      userID,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		isActive:    isActive,
		events:      []events.DomainEvent{},
	}, nil
}

// ID returns the category ID.
func (c *Category) ID() valueobjects.CategoryID {
	return c.id
}

// UserID returns the user ID.
func (c *Category) UserID() identityvalueobjects.UserID {
	return c.userID
}

// Name returns the category name.
func (c *Category) Name() valueobjects.CategoryName {
	return c.name
}

// Description returns the category description.
func (c *Category) Description() string {
	return c.description
}

// CreatedAt returns the creation timestamp.
func (c *Category) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the last update timestamp.
func (c *Category) UpdatedAt() time.Time {
	return c.updatedAt
}

// IsActive returns whether the category is active.
func (c *Category) IsActive() bool {
	return c.isActive
}

// UpdateName updates the category name.
func (c *Category) UpdateName(name valueobjects.CategoryName) error {
	if !c.isActive {
		return errors.New("cannot update name for inactive category")
	}

	if name.IsEmpty() {
		return errors.New("category name cannot be empty")
	}

	c.name = name
	c.updatedAt = time.Now()

	c.addEvent(events.NewBaseDomainEvent(
		"CategoryNameUpdated",
		c.id.Value(),
		"Category",
	))

	return nil
}

// UpdateDescription updates the category description.
func (c *Category) UpdateDescription(description string) error {
	if !c.isActive {
		return errors.New("cannot update description for inactive category")
	}

	if len(description) > 500 {
		return errors.New("category description is too long (max 500 characters)")
	}

	c.description = description
	c.updatedAt = time.Now()

	c.addEvent(events.NewBaseDomainEvent(
		"CategoryDescriptionUpdated",
		c.id.Value(),
		"Category",
	))

	return nil
}

// Deactivate deactivates the category.
func (c *Category) Deactivate() error {
	if !c.isActive {
		return errors.New("category is already inactive")
	}

	c.isActive = false
	c.updatedAt = time.Now()

	c.addEvent(events.NewBaseDomainEvent(
		"CategoryDeactivated",
		c.id.Value(),
		"Category",
	))

	return nil
}

// Activate activates the category.
func (c *Category) Activate() error {
	if c.isActive {
		return errors.New("category is already active")
	}

	c.isActive = true
	c.updatedAt = time.Now()

	c.addEvent(events.NewBaseDomainEvent(
		"CategoryActivated",
		c.id.Value(),
		"Category",
	))

	return nil
}

// GetEvents returns all domain events that occurred on this aggregate.
func (c *Category) GetEvents() []events.DomainEvent {
	return c.events
}

// ClearEvents clears all domain events from this aggregate.
func (c *Category) ClearEvents() {
	c.events = []events.DomainEvent{}
}

// addEvent adds a domain event to the aggregate.
func (c *Category) addEvent(event events.DomainEvent) {
	c.events = append(c.events, event)
}
