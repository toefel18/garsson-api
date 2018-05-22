package auth

import (

"testing"

"github.com/stretchr/testify/assert"

)

func TestUserEntity_GetRoles(t *testing.T) {
    configuredRoles := "admin, cs,    dev ,"
    expectedRoles := []string{"admin", "cs", "dev"}

    actualRoles := UserEntity{"","", configuredRoles, ""}.GetRoles()
    assert.EqualValues(t, actualRoles, expectedRoles, "expectedRoles %v configuredRoles, but got %v", expectedRoles, actualRoles)
}
