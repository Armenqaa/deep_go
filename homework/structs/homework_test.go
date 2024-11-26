package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	nameMaxSize = 42
	endOfString byte = '\000'

	// shift and masks for manaHp
	manaShiftTopBits = 2
	manaShiftLowBits = 6
	manaMaskLowBits  = 0b11000000

	healthShiftTopBits = 8
	healthMaskTopBits  = 0b00000011
	healthMaskLowBits  = 0b11111111

	// shift and masks for respectStrength
	respectShift = 4
	respectMask  = 0b11110000
	strengthMask = 0b00001111

	// shift and masks for lvlExp
	levelShift     = 4
	levelMask      = 0b11110000
	experienceMask = 0b00001111

	// shift and masks for typeHouseGunFamily
	hasHouseBit  = 2
	hasGunBit    = 1
	hasFamilyBit = 0
	typeShift    = 3
	typeMask     = 0b00011000
)

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i, v := range name {
			person.name[i] = byte(v)
		}

		if len(name) < nameMaxSize {
			person.name[len(name)] = endOfString
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
		person.manaHp[0] = byte(mana >> manaShiftTopBits)
		person.manaHp[1] |= byte((mana << manaShiftLowBits) & manaMaskLowBits)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHp[1] |= byte(health >> healthShiftTopBits)
		person.manaHp[2] = byte(health & 0xFF)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength |= byte(respect << respectShift)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectStrength |= byte(strength)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.lvlExp |= byte(level << levelShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.lvlExp |= byte(experience)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 1 << hasHouseBit
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 1 << hasGunBit
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= 1 << hasFamilyBit
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typeHouseGunFamily |= byte(personType) << typeShift
	}
}

type GamePerson struct {
	x, y, z            int32
	gold               uint32
	manaHp             [3]byte // 10 bits for both numbers. 0s byte with 2 top bits of 1st byte is Mana, 2 low bits of 1st byte with 2nd byte is Health.
	respectStrength    byte    // 4 top bits for respect, 4 low bits for strength
	lvlExp             byte    // 4 top bits for level, 4 low bits for experience
	typeHouseGunFamily byte    // first low bit is family, next is gun, next is house, next 2 bites is type
	name               [nameMaxSize]byte
}

func NewGamePerson(options ...Option) GamePerson {
	p := GamePerson{}
	for _, o := range options {
		o(&p)
	}

	return p
}

func (p *GamePerson) Name() string {
	var currLength int
	for _, v := range p.name {
		if v == endOfString {
			break
		}
		currLength++
	}
	return unsafe.String(&p.name[0], currLength)
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
	return int(p.manaHp[0])<<manaShiftTopBits + int(p.manaHp[1]&manaMaskLowBits>>manaShiftLowBits)
}

func (p *GamePerson) Health() int {
	return int(p.manaHp[1])&healthMaskTopBits<<healthShiftTopBits + int(p.manaHp[2])
}

func (p *GamePerson) Respect() int {
	return int(p.respectStrength & respectMask >> respectShift)
}

func (p *GamePerson) Strength() int {
	return int(p.respectStrength & strengthMask)
}

func (p *GamePerson) Level() int {
	return int(p.lvlExp & levelMask >> levelShift)
}

func (p *GamePerson) Experience() int {
	return int(p.lvlExp & experienceMask)
}

func (p *GamePerson) HasHouse() bool {
	return p.typeHouseGunFamily&(1 << hasHouseBit) != 0
}

func (p *GamePerson) HasGun() bool {
	return p.typeHouseGunFamily&(1 << hasGunBit) != 0
}

func (p *GamePerson) HasFamily() bool {
	return p.typeHouseGunFamily&(1 << hasFamilyBit) != 0
}

func (p *GamePerson) Type() int {
	return int(p.typeHouseGunFamily & typeMask >> typeShift)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa"
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
