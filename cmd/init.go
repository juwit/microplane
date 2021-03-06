package cmd

import (
	"fmt"
	"log"

	"github.com/Clever/microplane/initialize"

	"github.com/spf13/cobra"
)

var repoProviderFlag string
var initFlagReposFile string

var initCmd = &cobra.Command{
	Use:   "init [query]",
	Short: "Initialize a microplane workflow",
	Long: `Initialize a microplane workflow.

There are two ways to init, either (1) from a file or (2) via search

## (1) Init from File

$ mp init -f repos.txt

where repos.txt has lines like:

	clever/repo2
	clever/repo2

## (2) Init via Search

### GitHub

Search targets repos based on a Github Code Search query.

For example:

$ mp init "org:Clever filename:circle.yml"

would target all Clever repos with a circle.yml file.

See https://help.github.com/articles/searching-code/ for more details about the search syntax on Github.

### GitLab

Search targets repos based on a GitLab search.

If you are using the *public* version of GitLab, search is done via the Global "projects" scope.
See https://docs.gitlab.com/ce/api/search.html#scope-projects for more information on the search syntax. For example

$ mp init "mp-test-1"

would target a specific repo called mp-test-1.

If you are using an *enterprise* GitLab instance, we assume you have an ElasticSearch setup.
See https://docs.gitlab.com/ee/user/search/advanced_search_syntax.html for more details about the search syntax on Gitlab.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && initFlagReposFile == "" {
			log.Fatal("to init via search, you must pass a search query. otherwise, specify a repos file with -f")
		}

		if len(args) == 1 && initFlagReposFile != "" {
			log.Fatal("to init via search, you must pass a search query. otherwise, specify a repos file with -f")
		}

		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		output, err := initialize.Initialize(initialize.Input{
			Query:         query,
			WorkDir:       workDir,
			Version:       cliVersion,
			RepoProvider:  repoProviderFlag,
			ReposFromFile: initFlagReposFile,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = writeJSON(output, outputPath("", "init"))
		if err != nil {
			log.Fatal(err)
		}

		for _, repo := range output.Repos {
			fmt.Println(repo.Name)
		}
	},
}
