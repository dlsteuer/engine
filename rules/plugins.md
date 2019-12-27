# Ruleset Plugin Architecture

The game engine provides the ability to override the default ruleset via a plugin based architecture
On startup the game engine will search `~/.battlesnake/rulesets` for available .so files to load.
It will pull out a symbol called ruleset and make this available under the name of the library without the .so extension.
For example if you create a ruleset plugin called `awesome.so` the engine will use that ruleset under the name `awesome`.
There is an optional parameter on create game called ruleset, which will try and use the named ruleset when running that game.

Ruleset plugins must expose a symbol called `Ruleset` that implements the following interface

```go
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
```