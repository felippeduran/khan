// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"encoding/json"

	"github.com/topfreegames/khan/util"

	"gopkg.in/gorp.v1"
)

// Game identifies uniquely one game
type Game struct {
	ID                            int                    `db:"id"`
	PublicID                      string                 `db:"public_id"`
	Name                          string                 `db:"name"`
	MinMembershipLevel            int                    `db:"min_membership_level"`
	MaxMembershipLevel            int                    `db:"max_membership_level"`
	MinLevelToAcceptApplication   int                    `db:"min_level_to_accept_application"`
	MinLevelToCreateInvitation    int                    `db:"min_level_to_create_invitation"`
	MinLevelToRemoveMember        int                    `db:"min_level_to_remove_member"`
	MinLevelOffsetToRemoveMember  int                    `db:"min_level_offset_to_remove_member"`
	MinLevelOffsetToPromoteMember int                    `db:"min_level_offset_to_promote_member"`
	MinLevelOffsetToDemoteMember  int                    `db:"min_level_offset_to_demote_member"`
	MaxMembers                    int                    `db:"max_members"`
	MaxClansPerPlayer             int                    `db:"max_clans_per_player"`
	MembershipLevels              map[string]interface{} `db:"membership_levels"`
	Metadata                      map[string]interface{} `db:"metadata"`
	CreatedAt                     int64                  `db:"created_at"`
	UpdatedAt                     int64                  `db:"updated_at"`
	CooldownAfterDeny             int                    `db:"cooldown_after_deny"`
	CooldownAfterDelete           int                    `db:"cooldown_after_delete"`
}

// PreInsert populates fields before inserting a new game
func (g *Game) PreInsert(s gorp.SqlExecutor) error {
	// Handle JSON fields
	sortedLevels := util.SortLevels(g.MembershipLevels)
	g.MinMembershipLevel = sortedLevels[0].Value
	g.MaxMembershipLevel = sortedLevels[len(sortedLevels)-1].Value
	g.CreatedAt = util.NowMilli()
	g.UpdatedAt = g.CreatedAt
	return nil
}

// PreUpdate populates fields before updating a game
func (g *Game) PreUpdate(s gorp.SqlExecutor) error {
	sortedLevels := util.SortLevels(g.MembershipLevels)
	g.MinMembershipLevel = sortedLevels[0].Value
	g.MaxMembershipLevel = sortedLevels[len(sortedLevels)-1].Value
	g.UpdatedAt = util.NowMilli()
	return nil
}

// GetGameByID returns a game by id
func GetGameByID(db DB, id int) (*Game, error) {
	obj, err := db.Get(Game{}, id)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, &ModelNotFoundError{"Game", id}
	}

	game := obj.(*Game)
	return game, nil
}

// GetGameByPublicID returns a game by their public id
func GetGameByPublicID(db DB, publicID string) (*Game, error) {
	var games []*Game
	_, err := db.Select(&games, "SELECT * FROM games WHERE public_id=$1", publicID)
	if err != nil {
		return nil, err
	}
	if games == nil || len(games) < 1 {
		return nil, &ModelNotFoundError{"Game", publicID}
	}
	return games[0], nil
}

// GetAllGames returns all games in the DB
func GetAllGames(db DB) ([]*Game, error) {
	var games []*Game
	_, err := db.Select(&games, "SELECT * FROM games")
	if err != nil {
		return nil, err
	}
	return games, nil
}

// CreateGame creates a new game
func CreateGame(
	db DB,
	publicID, name string,
	levels, metadata map[string]interface{},
	minLevelAccept, minLevelCreate, minLevelRemove,
	minOffsetRemove, minOffsetPromote, minOffsetDemote,
	maxMembers, maxClans, cooldownAfterDeny, cooldownAfterDelete int,
	upsert bool,
) (*Game, error) {
	var game *Game

	if upsert {
		levelsJSON, err := json.Marshal(levels)
		if err != nil {
			return nil, err
		}

		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}

		sortedLevels := util.SortLevels(levels)
		minMembershipLevel := sortedLevels[0].Value
		maxMembershipLevel := sortedLevels[len(sortedLevels)-1].Value

		_, err = db.Exec(`
			INSERT INTO games(
				public_id,
				name,
				min_level_to_accept_application,
				min_level_to_create_invitation,
				min_level_to_remove_member,
				min_level_offset_to_remove_member,
				min_level_offset_to_promote_member,
				min_level_offset_to_demote_member,
				max_members,
				max_clans_per_player,
				membership_levels,
				metadata,
				cooldown_after_delete,
				cooldown_after_deny,
				min_membership_level,
				max_membership_level,
				created_at,
				updated_at
			)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $17)
			ON CONFLICT (public_id)
			DO UPDATE set
				name=$2,
				min_level_to_accept_application=$3,
				min_level_to_create_invitation=$4,
				min_level_to_remove_member=$5,
				min_level_offset_to_remove_member=$6,
				min_level_offset_to_promote_member=$7,
				min_level_offset_to_demote_member=$8,
				max_members=$9,
				max_clans_per_player=$10,
				membership_levels=$11,
				metadata=$12,
				cooldown_after_delete=$13,
				cooldown_after_deny=$14,
				min_membership_level=$15,
				max_membership_level=$16,
				updated_at=$17
			WHERE games.public_id=$1`,
			publicID,            // $1
			name,                // $2
			minLevelAccept,      // $3
			minLevelCreate,      // $4
			minLevelRemove,      // $5
			minOffsetRemove,     // $6
			minOffsetPromote,    // $7
			minOffsetDemote,     // $8
			maxMembers,          // $9
			maxClans,            // $10
			levelsJSON,          // $11
			metadataJSON,        // $12
			cooldownAfterDelete, // $13
			cooldownAfterDeny,   // $14
			minMembershipLevel,  // $15
			maxMembershipLevel,  // $16
			util.NowMilli(),     // $17
		)
		if err != nil {
			return nil, err
		}
		return GetGameByPublicID(db, publicID)
	}

	game = &Game{
		PublicID: publicID,
		Name:     name,
		MinLevelToAcceptApplication:   minLevelAccept,
		MinLevelToCreateInvitation:    minLevelCreate,
		MinLevelToRemoveMember:        minLevelRemove,
		MinLevelOffsetToRemoveMember:  minOffsetRemove,
		MinLevelOffsetToPromoteMember: minOffsetPromote,
		MinLevelOffsetToDemoteMember:  minOffsetDemote,
		MaxMembers:                    maxMembers,
		MaxClansPerPlayer:             maxClans,
		MembershipLevels:              levels,
		Metadata:                      metadata,
		CooldownAfterDelete:           cooldownAfterDelete,
		CooldownAfterDeny:             cooldownAfterDeny,
	}
	err := db.Insert(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// UpdateGame updates an existing game
func UpdateGame(db DB, publicID, name string, levels, metadata map[string]interface{},
	minLevelAccept, minLevelCreate, minLevelRemove, minOffsetRemove, minOffsetPromote, minOffsetDemote, maxMembers, maxClans, cooldownAfterDeny, cooldownAfterDelete int,
) (*Game, error) {
	return CreateGame(
		db, publicID, name, levels, metadata, minLevelAccept,
		minLevelCreate, minLevelRemove, minOffsetRemove, minOffsetPromote,
		minOffsetDemote, maxMembers, maxClans, cooldownAfterDeny, cooldownAfterDelete, true,
	)
}
