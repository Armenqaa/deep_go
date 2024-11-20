package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}

type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type SomeSingletonService struct {
	// not need to impelement
	NotEmptyStruct bool
}

type Constructor interface {
	Construct() any
}

type BasicConstructor struct {
	f	func() any
}

func (bc BasicConstructor) Construct() any {
	return bc.f()
}

type SingletonConstructor struct {
	f 	 func() any
	called	 bool
	instance any
}

func (sc *SingletonConstructor) Construct() any {
	if !sc.called {
		sc.called = true
		sc.instance = sc.f()
	}
	return sc.instance
}

type Container struct {
	constructors map[string]Constructor
}

func NewContainer() *Container {
	return &Container{
		constructors: map[string]Constructor{},
	}
}

func (c *Container) RegisterType(name string, constructor func() any) {
	c.constructors[name] = BasicConstructor{f: constructor}
}

func (c *Container) RegisterSingletonType(name string, constructor func() any) {
	c.constructors[name] = &SingletonConstructor{f: constructor}
}

func (c *Container) Resolve(name string) (any, error) {
	constructor, ok := c.constructors[name]
	if !ok {
		return nil, fmt.Errorf("constructor not found")
	}

	return constructor.Construct(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	container.RegisterSingletonType("SomeSingleton", func() interface{} {
		return &SomeSingletonService{}
	})
	
	singleton1, err := container.Resolve("SomeSingleton")
	assert.NoError(t, err)
	singleton2, err := container.Resolve("SomeSingleton")
	assert.NoError(t, err)
	
	s1 := singleton1.(*SomeSingletonService)
	s2 := singleton2.(*SomeSingletonService)
	assert.True(t, s1 == s2)
	
}
