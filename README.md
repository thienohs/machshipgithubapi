# Configuration
## Environment variables
PORT: server port (default: 8777)

GITHUB_API_URL: API URL of github (default: https://api.github.com)

GITHUB_API_USER: Github API User path (default: users)

# Examples
## HTTP GET query
```
curl -L \
  "http://localhost:8777/retrieveUsers?usernames=machship,google,apache,kubernetes"
```

Demo:
```
curl -L \
  "https://machship.gevelation.com/retrieveUsers?usernames=machship,google,apache,kubernetes"
```

## GraphQL query
### Playground: 
http(s)://[host]:[port]/graphql/playground

Demo: https://machship.gevelation.com/graphql/playground
### Query: 
http(s)://[host]:[port]/graphql/query

Demo: https://machship.gevelation.com/graphql/query
### Examples query:
```
query retrieveUsers {
  retrieveUsers(usernames: ["machship", "apache", "google", "kubernetes"]) {
    users {
      name,
      login,
      company,
      followers,
      public_repos,
      avg_followers_per_public_repo,
    },
    errors {
      message
    },
  }
}
```
