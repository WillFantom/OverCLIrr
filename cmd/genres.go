package cmd

import (
	"fmt"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/willfantom/goverseerr"
	"github.com/willfantom/overclirr/ui"
)

var genreCmd = &cobra.Command{
	Use:       "genres [tv/movie]",
	Aliases:   []string{"genre"},
	Short:     "Get a list of all tv/movie genres available",
	ValidArgs: []string{"tv", "movie"},
	Args:      cobra.ExactValidArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		var genreList []*goverseerr.Genre
		var err error
		switch args[0] {
		case "tv":
			genreList, err = overseerr.TVGenres()
		case "movie":
			genreList, err = overseerr.MovieGenres()
		default:
			logrus.WithField("givenArg", args[0]).
				Panicln("Invalid genre argument got past argument validator")
		}
		if err != nil {
			ui.Fatal("Could not get genre list from Overseerr instance", err)
		}
		logrus.WithField("genreCount", len(genreList)).Debug("collected genre list")
		var tableValues = [][]string{
			{"ID", "Name"},
		}
		for _, genre := range genreList {
			tableValues = append(tableValues, []string{fmt.Sprintf("%d", genre.ID), genre.Name})
		}
		ui.Table(tableValues)
	},
}

var genreSearchCmd = &cobra.Command{
	Use:   "genre-search [search_term]",
	Short: "Check if a genre exists",
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		setOverseer(overseerrProfileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		tvGenreList, movieGenreList := searchGenre(args[0])
		var tableValues = [][]string{
			{"Type", "ID", "Name"},
		}
		for _, genre := range tvGenreList {
			if fuzzy.MatchNormalizedFold(args[0], genre.Name) {
				tableValues = append(tableValues, []string{"TV", fmt.Sprintf("%d", genre.ID), genre.Name})
			}
		}
		for _, genre := range movieGenreList {
			if fuzzy.MatchNormalizedFold(args[0], genre.Name) {
				tableValues = append(tableValues, []string{"Movie", fmt.Sprintf("%d", genre.ID), genre.Name})
			}
		}
		ui.Table(tableValues)
	},
}

func searchGenre(searchTerm string) ([]*goverseerr.Genre, []*goverseerr.Genre) {
	logrus.WithField("searchTerm", searchTerm).Traceln("searching for genre id")
	tvGenreList, err := overseerr.TVGenres()
	if err != nil {
		ui.Fatal("Could not get Genre list from Overseerr instance", err)
	}
	movieGenreList, err := overseerr.MovieGenres()
	if err != nil {
		ui.Fatal("Could not get Genre list from Overseerr instance", err)
	}
	logrus.WithFields(logrus.Fields{
		"tv-genres":    len(tvGenreList),
		"movie-genres": len(movieGenreList),
	}).Debug("collected genre lists")
	return tvGenreList, movieGenreList
}

func init() {
	RootCmd.AddCommand(genreCmd)
	RootCmd.AddCommand(genreSearchCmd)
}
