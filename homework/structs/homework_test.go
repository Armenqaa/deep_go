package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, v := range name {
			person.name[i] = byte(v)
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x, person.y, person.z = int32(x), int32(y), int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHp[0] = byte(mana >> 2)
		person.manaHp[1] |= byte((mana << 6) & 0b11000000)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHp[1] |= byte(health >> 8)
		person.manaHp[2] = byte(health & 0b0011111111)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength |= byte(respect << 4)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength |= byte(strength)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.lvlExp |= byte(level << 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.lvlExp |= byte(experience)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 0b00000100
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 0b00000010
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 0b00000001
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= byte(personType) << 3
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z            int32
	gold               uint32
	manaHp             [3]byte
	respectStrength    byte
	lvlExp             byte
	typeHouseGunFamily byte
	name               [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}
	for _, o := range options {
		o(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaHp[0])<<2 + int(p.manaHp[1]&0b11000000>>6)
}

func (p *GamePerson) Health() int {
	return int(p.manaHp[1])&0b00000011<<8 + int(p.manaHp[2])
}

func (p *GamePerson) Respect() int {
	return int(p.respectStrength & 0b11110000 >> 4)
}

func (p *GamePerson) Strength() int {
	return int(p.respectStrength & 0b00001111)
}

func (p *GamePerson) Level() int {
	return int(p.lvlExp & 0b11110000 >> 4)
}

func (p *GamePerson) Experience() int {
	return int(p.lvlExp & 0b00001111)
}

func (p *GamePerson) HasHouse() bool {
	return p.typeHouseGunFamily&0b00000100>>2 == 1
}

func (p *GamePerson) HasGun() bool {
	return p.typeHouseGunFamily&0b00000010>>1 == 1
}

func (p *GamePerson) HasFamily() bool {
	return p.typeHouseGunFamily&0b00000001 == 1
}

func (p *GamePerson) Type() int {
	return int(p.typeHouseGunFamily & 0b00011000 >> 3)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
