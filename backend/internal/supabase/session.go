package supabase

import (
	"fmt"
	"log/slog"
	"net/http"
)

// SupabaseValidateSession asks Supabase whether the access token's session is still live.
func (s *Supabase) SupabaseValidateSession(Client *http.Client, access_token string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v1/user", s.URL), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	req.Header.Set("apikey", s.AnonKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Failed to execute Request", "err", err)
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("session validation failed: status %d", res.StatusCode)
	}

	return nil
}
