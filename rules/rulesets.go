package rules

import (
	"github.com/battlesnakeio/engine/controller/pb"
)

type Ruleset interface {
	UpdateSnakeLocations(game *pb.Game, frame *pb.GameFrame, moves []*SnakeUpdate)
	UpdateSnakeHealth(frame *pb.GameFrame)
	CheckForSnakesEating(frame *pb.GameFrame) []*pb.Point
	UpdateFood(game *pb.Game, gameFrame *pb.GameFrame, foodToRemove []*pb.Point) ([]*pb.Point, error)
	CheckForDeath(width, height int32, frame *pb.GameFrame) []DeathUpdate
}

type DefaultRuleset struct{}

func (rs *DefaultRuleset) UpdateSnakeLocations(game *pb.Game, frame *pb.GameFrame, moves []*SnakeUpdate) {
	UpdateSnakeLocations(game, frame, moves)
}
func (rs *DefaultRuleset) UpdateSnakeHealth(frame *pb.GameFrame) {
	UpdateSnakeHealth(frame)
}
func (rs *DefaultRuleset) CheckForSnakesEating(frame *pb.GameFrame) []*pb.Point {
	return CheckForSnakesEating(frame)
}
func (rs *DefaultRuleset) UpdateFood(game *pb.Game, gameFrame *pb.GameFrame, foodToRemove []*pb.Point) ([]*pb.Point, error) {
	return UpdateFood(game, gameFrame, foodToRemove)
}
func (rs *DefaultRuleset) CheckForDeath(width, height int32, frame *pb.GameFrame) []DeathUpdate {
	return CheckForDeath(width, height, frame)
}
