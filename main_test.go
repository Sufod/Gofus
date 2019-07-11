package main

import (
	"testing"

	"gotest.tools/assert"
)

func TestCryptPassword(t *testing.T) {
	assert.Equal(t, cryptPassword("MonSUperp4ssword", "zzybokxyrtkpjvxmmoxbnwiynojxdbqn"), "#1-haj_hNL00YR-9a75Y34YU3ZXX8f_6ZX")
}
