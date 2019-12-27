package main

import (
	"github.com/battlesnakeio/engine/controller/pb"
	"github.com/battlesnakeio/engine/rules"
)

type StandardRuleset struct{}

func (sr *StandardRuleset) UpdateSnakeLocations(game *pb.Game, frame *pb.GameFrame, moves []*rules.SnakeUpdate) {
	rules.UpdateSnakeLocations(game, frame, moves)
}
func (sr *StandardRuleset) UpdateSnakeHealth(frame *pb.GameFrame) {
	rules.UpdateSnakeHealth(frame)
}
func (sr *StandardRuleset) CheckForSnakesEating(frame *pb.GameFrame) []*pb.Point {
	return rules.CheckForSnakesEating(frame)
}
func (sr *StandardRuleset) UpdateFood(game *pb.Game, gameFrame *pb.GameFrame, foodToRemove []*pb.Point) ([]*pb.Point, error) {
	return rules.UpdateFood(game, gameFrame, foodToRemove)
}
func (sr *StandardRuleset) CheckForDeath(width, height int32, frame *pb.GameFrame) []rules.DeathUpdate {
	return rules.CheckForDeath(width, height, frame)
}

var Ruleset = StandardRuleset{}
