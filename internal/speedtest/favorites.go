package speedtest

import (
	"encoding/json"
	"os"
	"path/filepath"

	"auraspeed/internal/config"
)

// FavoriteServer represents a server marked as favorite by the user
type FavoriteServer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Sponsor  string `json:"sponsor"`
	Host     string `json:"host"`
	URL      string `json:"url"`
	AddedAt  string `json:"added_at"`
}

// favoritesFilePath returns the path to the favorites file
func favoritesFilePath() string {
	return filepath.Join(config.GetDataDir(), "favorites.json")
}

// LoadFavorites loads the list of favorite servers
func LoadFavorites() ([]FavoriteServer, error) {
	path := favoritesFilePath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []FavoriteServer{}, nil
		}
		return nil, err
	}

	var favorites []FavoriteServer
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, err
	}

	return favorites, nil
}

// SaveFavorites saves the list of favorite servers
func SaveFavorites(favorites []FavoriteServer) error {
	path := favoritesFilePath()

	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// AddFavorite adds a server to the favorites list
func AddFavorite(serverID, name, country, sponsor, host, url string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	// Check if already exists
	for _, f := range favorites {
		if f.ID == serverID {
			return nil // Already a favorite
		}
	}

	favorite := FavoriteServer{
		ID:       serverID,
		Name:     name,
		Country:  country,
		Sponsor:  sponsor,
		Host:     host,
		URL:      url,
		AddedAt:  "",
	}

	favorites = append(favorites, favorite)
	return SaveFavorites(favorites)
}

// RemoveFavorite removes a server from the favorites list
func RemoveFavorite(serverID string) error {
	favorites, err := LoadFavorites()
	if err != nil {
		return err
	}

	newFavorites := make([]FavoriteServer, 0)
	for _, f := range favorites {
		if f.ID != serverID {
			newFavorites = append(newFavorites, f)
		}
	}

	return SaveFavorites(newFavorites)
}

// IsFavorite checks if a server is in the favorites list
func IsFavorite(serverID string) bool {
	favorites, err := LoadFavorites()
	if err != nil {
		return false
	}

	for _, f := range favorites {
		if f.ID == serverID {
			return true
		}
	}

	return false
}

// GetFavorites returns the list of favorite server IDs
func GetFavorites() []string {
	favorites, err := LoadFavorites()
	if err != nil {
		return []string{}
	}

	ids := make([]string, len(favorites))
	for i, f := range favorites {
		ids[i] = f.ID
	}

	return ids
}