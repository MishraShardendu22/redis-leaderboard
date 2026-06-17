package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	redisclient "github.com/redis/go-redis/v9"

	"redis-leaderboard/internal/model"
)

const (
	leaderboardKey = "leaderboard"
	rankHistoryKey = "leaderboard:last_ranks"
	stepPoints     = 10
)

var seedPlayers = []model.Player{
	{Name: "Alice", Score: 120},
	{Name: "Bob", Score: 90},
	{Name: "Charlie", Score: 70},
	{Name: "David", Score: 40},
}

type Service struct {
	client *redisclient.Client
}

func NewService(client *redisclient.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Client() *redisclient.Client {
	return s.client
}

func (s *Service) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Service) EnsureSeedData() error {
	return s.ensureSeedData(context.Background())
}

func (s *Service) Leaderboard() ([]model.Player, error) {
	players, _, err := s.leaderboardSnapshot(context.Background())
	return players, err
}

func (s *Service) LeaderboardSnapshot(ctx context.Context) ([]model.Player, map[string]int64, error) {
	return s.leaderboardSnapshot(ctx)
}

func (s *Service) IncreaseScore(name string, points int) error {
	return s.changeScore(context.Background(), name, points)
}

func (s *Service) DecreaseScore(name string, points int) error {
	return s.changeScore(context.Background(), name, -points)
}

func (s *Service) Rank(name string) (int64, error) {
	rank, err := s.client.ZRevRank(context.Background(), leaderboardKey, name).Result()
	if err != nil {
		if errors.Is(err, redisclient.Nil) {
			return 0, fmt.Errorf("player %q not found", name)
		}
		return 0, err
	}
	return rank + 1, nil
}

func (s *Service) ensureSeedData(ctx context.Context) error {
	count, err := s.client.ZCard(ctx, leaderboardKey).Result()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	entries := make([]redisclient.Z, 0, len(seedPlayers))
	for _, player := range seedPlayers {
		entries = append(entries, redisclient.Z{Score: float64(player.Score), Member: player.Name})
	}

	if err := s.client.ZAdd(ctx, leaderboardKey, entries...).Err(); err != nil {
		return err
	}

	return s.storeCurrentRanks(ctx)
}

func (s *Service) changeScore(ctx context.Context, name string, delta int) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("player name is required")
	}

	if _, err := s.client.ZIncrBy(ctx, leaderboardKey, float64(delta), name).Result(); err != nil {
		return err
	}

	return nil
}

func (s *Service) leaderboardSnapshot(ctx context.Context) ([]model.Player, map[string]int64, error) {
	if err := s.ensureSeedData(ctx); err != nil {
		return nil, nil, err
	}

	previousRanks, err := s.client.HGetAll(ctx, rankHistoryKey).Result()
	if err != nil {
		return nil, nil, err
	}

	entries, err := s.client.ZRevRangeWithScores(ctx, leaderboardKey, 0, -1).Result()
	if err != nil {
		return nil, nil, err
	}

	players := make([]model.Player, 0, len(entries))
	movements := make(map[string]int64, len(entries))
	currentRanks := make(map[string]interface{}, len(entries))

	for index, entry := range entries {
		name := fmt.Sprint(entry.Member)
		rank := int64(index + 1)
		players = append(players, model.Player{Name: name, Score: int(entry.Score), Rank: int(rank)})
		if previousRank, ok := previousRanks[name]; ok {
			parsed, parseErr := strconv.ParseInt(previousRank, 10, 64)
			if parseErr == nil {
				movements[name] = parsed - rank
			}
		}
		currentRanks[name] = rank
	}

	if len(currentRanks) > 0 {
		if err := s.client.HSet(ctx, rankHistoryKey, currentRanks).Err(); err != nil {
			return nil, nil, err
		}
	}

	return players, movements, nil
}

func (s *Service) storeCurrentRanks(ctx context.Context) error {
	entries, err := s.client.ZRevRangeWithScores(ctx, leaderboardKey, 0, -1).Result()
	if err != nil {
		return err
	}

	currentRanks := make(map[string]interface{}, len(entries))
	for index, entry := range entries {
		currentRanks[fmt.Sprint(entry.Member)] = index + 1
	}

	if len(currentRanks) == 0 {
		return nil
	}

	return s.client.HSet(ctx, rankHistoryKey, currentRanks).Err()
}
