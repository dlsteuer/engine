package rules

import (
	"github.com/battlesnakeio/engine/controller/pb"
)

type Ruleset interface {
	// UpdateSnakeLocations this function will update the snakes in the frame based upon the moves
	// returned by each snake
	UpdateSnakeLocations(game *pb.Game, frame *pb.GameFrame, moves []*SnakeUpdate)
	// UpdateSnakeHealth this function is where the snake health is updated each game tick
	UpdateSnakeHealth(frame *pb.GameFrame)
	// CheckForSnakesEating this function checks to see if a snake has eaten food, and returns a list of
	// points that is passed to UpdateFood for later processing, the standard ruleset uses these eaten food points
	// as a list of food to be removed from the game board
	CheckForSnakesEating(frame *pb.GameFrame) []*pb.Point
	// UpdateFood handles any processing of food that was eaten by snakes
	UpdateFood(game *pb.Game, gameFrame *pb.GameFrame, foodEaten []*pb.Point) ([]*pb.Point, error)
	// CheckForDeath checks to see if a snake has died, and returns any deaths as a DeathUpdate
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
