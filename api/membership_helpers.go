// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/topfreegames/khan/log"
	"github.com/topfreegames/khan/models"
	"github.com/uber-go/zap"
)

type membershipOptionalParams struct {
	Message string
}

func getMembershipOptionalParameters(app *App, c echo.Context) (*membershipOptionalParams, error) {
	data, err := GetRequestBody(c)
	if err != nil {
		return nil, err
	}

	var jsonPayload map[string]interface{}
	err = json.Unmarshal([]byte(data), &jsonPayload)
	if err != nil {
		return nil, err
	}

	var message string
	if val, ok := jsonPayload["message"]; ok {
		message = val.(string)
	} else {
		message = ""
	}

	return &membershipOptionalParams{
		Message: message,
	}, nil
}

func dispatchMembershipHookByPublicID(app *App, db models.DB, hookType int, gameID, clanID, playerID, requestorID, membershipLevel string) error {
	clan, err := models.GetClanByPublicID(db, gameID, clanID)
	if err != nil {
		return err
	}

	player, err := models.GetPlayerByPublicID(db, gameID, playerID)
	if err != nil {
		return err
	}

	requestor := player
	if requestorID != playerID {
		requestor, err = models.GetPlayerByPublicID(db, gameID, requestorID)
		if err != nil {
			return err
		}
	}

	return dispatchMembershipHook(app, db, hookType, gameID, clan, player, requestor, "", membershipLevel)
}

func dispatchMembershipHookByID(app *App, db models.DB, hookType int, gameID string, clanID, playerID, requestorID int, message, membershipLevel string) error {
	clan, err := models.GetClanByID(db, clanID)
	if err != nil {
		return err
	}

	player, err := models.GetPlayerByID(db, playerID)
	if err != nil {
		return err
	}

	requestor := player
	if requestorID != playerID {
		requestor, err = models.GetPlayerByID(db, requestorID)
		if err != nil {
			return err
		}
	}

	return dispatchMembershipHook(app, db, hookType, gameID, clan, player, requestor, message, membershipLevel)
}

func dispatchApproveDenyMembershipHookByID(app *App, db models.DB, hookType int, gameID string, clanID, playerID, requestorID, creatorID int, message, membershipLevel string) error {
	clan, err := models.GetClanByID(db, clanID)
	if err != nil {
		return err
	}

	player, err := models.GetPlayerByID(db, playerID)
	if err != nil {
		return err
	}

	requestor := player
	if requestorID != playerID {
		requestor, err = models.GetPlayerByID(db, requestorID)
		if err != nil {
			return err
		}
	}

	creator := player
	if creatorID != playerID {
		creator = requestor
		if creatorID != requestorID {
			creator, err = models.GetPlayerByID(db, creatorID)
			if err != nil {
				return err
			}
		}
	}

	return dispatchApproveDenyMembershipHook(app, db, hookType, gameID, clan, player, requestor, creator, message, membershipLevel)
}

func dispatchMembershipHook(app *App, db models.DB, hookType int, gameID string, clan *models.Clan, player *models.Player, requestor *models.Player, message, membershipLevel string) error {
	clanJSON := clan.Serialize()
	delete(clanJSON, "gameID")

	playerJSON := player.Serialize()
	playerJSON["membershipLevel"] = membershipLevel
	delete(playerJSON, "gameID")

	requestorJSON := requestor.Serialize()
	delete(requestorJSON, "gameID")

	result := map[string]interface{}{
		"gameID":    gameID,
		"clan":      clanJSON,
		"player":    playerJSON,
		"requestor": requestorJSON,
	}

	if message != "" {
		result["message"] = message
	}
	app.DispatchHooks(gameID, hookType, result)

	return nil
}

func dispatchApproveDenyMembershipHook(app *App, db models.DB, hookType int, gameID string, clan *models.Clan, player *models.Player, requestor *models.Player, creator *models.Player, message, playerMembershipLevel string) error {
	clanJSON := clan.Serialize()
	delete(clanJSON, "gameID")

	playerJSON := player.Serialize()
	playerJSON["membershipLevel"] = playerMembershipLevel
	delete(playerJSON, "gameID")

	requestorJSON := requestor.Serialize()
	delete(requestorJSON, "gameID")

	creatorJSON := creator.Serialize()
	delete(creatorJSON, "gameID")

	result := map[string]interface{}{
		"gameID":    gameID,
		"clan":      clanJSON,
		"player":    playerJSON,
		"requestor": requestorJSON,
		"creator":   creatorJSON,
	}

	if message != "" {
		result["message"] = message
	}
	app.DispatchHooks(gameID, hookType, result)

	return nil
}

func getPayloadAndGame(app *App, c echo.Context, l zap.Logger) (*BasePayloadWithRequestorAndPlayerPublicIDs, *models.Game, int, error) {
	gameID := c.Param("gameID")

	var payload BasePayloadWithRequestorAndPlayerPublicIDs
	if err := LoadJSONPayload(&payload, c, l.With(zap.String("gameID", gameID))); err != nil {
		return nil, nil, 400, err
	}

	game, err := app.GetGame(gameID)
	if err != nil {
		log.W(l, "Could not find game.")
		return nil, nil, 404, err
	}

	return &payload, game, 200, nil
}
