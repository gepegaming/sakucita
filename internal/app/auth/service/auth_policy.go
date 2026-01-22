package service

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	loginAttemptKey = "auth:login:attempt"
	loginBanKey     = "auth:login:ban"

	maxFreeAttempts = 3
	attemptTTL      = 1 * time.Hour
	maxBanDuration  = 1 * time.Hour
)

func (s *service) CheckLoginBan(ctx context.Context, id string) (time.Duration, error) {
	// Ambil sisa TTL dari ban key
	// Kalau key tidak ada → ttl < 0
	ttl, err := s.rdb.TTL(ctx, banKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			// key ga ada berarti baru pertama kali coba login
			return 0, nil
		}
		s.log.Error().Err(err).Msg("failed to check login ban ttl")
		return 0, err
	}

	// kalo ttl lebih dari 0 maka itu ke ban
	if ttl > 0 {
		return ttl, nil
	}

	// kalo ttl < 0 maka ga ke ban
	return 0, nil
}

func (s *service) OnLoginFail(ctx context.Context, id string) (time.Duration, error) {
	attemptKey := attemptKey(id)

	// txPipeline, supaya di eksekusi barengan jadi sekali kirim aja ke redis
	pipe := s.rdb.TxPipeline()

	// tambah 1 ke attempt counter
	incr := pipe.Incr(ctx, attemptKey)

	// pastikan counter punya ttl agar otomatis di reset/hilang
	pipe.Expire(ctx, attemptKey, attemptTTL)

	// exec semua command yang diatas
	if _, err := pipe.Exec(ctx); err != nil {
		s.log.Err(err).Msg("failed to increment login attempt")
		return 0, err
	}

	// ambil hasil incr 9jumlah gagal sekarang
	attemps := incr.Val()

	// kalo masih dibawah free attempt, return 0 berarrti belom ke ban aja
	if attemps <= maxFreeAttempts {
		return 0, nil
	}

	// durasi exponential backoff
	// attempts = 4 → 5 detik
	// attempts = 5 → 10 detik
	// attempts = 6 → 20 detik
	exp := attemps - maxFreeAttempts
	delay := time.Duration(5*(1<<exp)) * time.Second

	// set tll agar tidak lebih dari max ban duration
	if delay > maxBanDuration {
		delay = maxBanDuration
	}

	// set login ban, value nya 1 cuma buat tanda aja bebas itumah agar ada value doang.
	//  delay nya itu dari perhitungan delay di atas
	// jadi set ke redis dengan key auth:login:ban:user@email.com dan value nya 1 dengan ttl delay
	if err := s.rdb.SetEx(ctx, banKey(id), "1", delay).Err(); err != nil {
		s.log.Error().Err(err).Msg("failed to set login ban")
		return 0, err
	}

	return delay, nil
}

func (s *service) OnLoginSuccess(ctx context.Context, id string) error {
	// kalo berhasil login apus counter attempt + ban
	if _, err := s.rdb.Del(ctx, attemptKey(id), banKey(id)).Result(); err != nil {
		s.log.Error().Err(err).Msg("failed to reset login attempt")
		return err
	}
	return nil
}

// helper
func attemptKey(id string) string {
	return fmt.Sprintf("%s:%s", loginAttemptKey, id)
}

func banKey(id string) string {
	return fmt.Sprintf("%s:%s", loginBanKey, id)
}
