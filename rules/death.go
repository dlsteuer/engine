package rules

import "github.com/battlesnakeio/engine/controller/pb"

type DeathUpdate struct {
	Snake *pb.Snake
	Death *pb.Death
}

// CheckForDeath looks through the snakes with the updated coords and checks to see if any have died
// possible death options are starvation (health has reached 0), wall collision, snake body collision
// snake head collision (other snake is same size or greater)
func CheckForDeath(width, height int32, frame *pb.GameFrame) []DeathUpdate {
	updates := []DeathUpdate{}
	for _, s := range frame.AliveSnakes() {
		if deathByHealth(s.Health) {
			updates = append(updates, DeathUpdate{
				Snake: s,
				Death: &pb.Death{
					Turn:  frame.Turn,
					Cause: DeathCauseStarvation,
				},
			})
			continue
		}
		head := s.Head()
		if head == nil {
			continue
		}
		if deathByOutOfBounds(head, width, height) {
			updates = append(updates, DeathUpdate{
				Snake: s,
				Death: &pb.Death{
					Turn:  frame.Turn,
					Cause: DeathCauseWallCollision,
				},
			})
			continue
		}

		for _, other := range frame.AliveSnakes() {
			if deathByHeadCollision(s, other) {
				updates = append(updates, DeathUpdate{
					Snake: s,
					Death: &pb.Death{
						Turn:  frame.Turn,
						Cause: DeathCauseHeadToHeadCollision,
					},
				})
			}

			for i, b := range other.Body {
				if i == 0 {
					continue
				}

				if deathByBodyCollision(s.Head(), b) {
					var cause string
					if s.ID == other.ID {
						cause = DeathCauseSnakeSelfCollision
					} else {
						cause = DeathCauseSnakeCollision
					}

					updates = append(updates, DeathUpdate{
						Snake: s,
						Death: &pb.Death{
							Turn:  frame.Turn,
							Cause: cause,
						},
					})
					break
				}
			}
		}
	}
	return updates
}

func deathByHealth(health int32) bool {
	return health <= 0
}

func deathByBodyCollision(head, body *pb.Point) bool {
	return head.Equal(body)
}

func deathByOutOfBounds(head *pb.Point, width, height int32) bool {
	return (head.X < 0) || (head.X >= width) || (head.Y < 0) || (head.Y >= height)
}

func deathByHeadCollision(snake, other *pb.Snake) bool {
	return (other.ID != snake.ID) && (snake.Head().Equal(other.Head())) && (len(snake.Body) <= len(other.Body))
}
