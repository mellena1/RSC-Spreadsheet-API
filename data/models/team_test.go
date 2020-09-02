package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Team_String(t *testing.T) {
	team := Team{
		Franchise: "The Bear Den",
		Name:      "Care Bears",
		Tier:      "Master",
	}

	require.Equal(t, "[The Bear Den] Care Bears (Master)", team.String())
}
