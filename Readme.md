# Export Github Repo Stats

List repos for an organisation along with their last updated date as CSV or JSON. More feature may
come or not.

## Usage

- Create a personal access token on github
- Build:

  ```sh
  dep ensure
  go build
  ```

- Run:

  ```sh
  GITHUB_TOKEN=<my_token>
  ./gh-stats -o csv -org myorg
  ```

- All options:

  ```sh
  gh-stats -h
  ```
