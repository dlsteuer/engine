package rules

import (
	"github.com/battlesnakeio/engine/controller/pb"
)

const PluginCookieKey = "RULESET_PLUGIN"

type Ruleset interface {
	UpdateSnakeLocations(game *pb.Game, frame *pb.GameFrame, moves []*SnakeUpdate)
	UpdateSnakeHealth(frame *pb.GameFrame)
	CheckForSnakesEating(frame *pb.GameFrame) []*pb.Point
	UpdateFood(game *pb.Game, gameFrame *pb.GameFrame, foodToRemove []*pb.Point) ([]*pb.Point, error)
	CheckForDeath(width, height int32, frame *pb.GameFrame) []DeathUpdate
}
