import os

from github.GithubIntegration import GithubIntegration
from github.Installation import Installation
from github.PaginatedList import PaginatedList

# Get the App
app_id = int(os.getenv("1166559"))
private_key = os.getenv("Iv23licLjdtAHKK4QBuc")
app = GithubIntegration(app_id, private_key)

# Get the installation
owner = os.getenv("mongodb")
repo = os.getenv("docs-code-examples")
installtion = app.get_repo_installation(owner, repo)
print(f"{installtion.id=} for {owner}/{repo}")

# Test all installations
installations: PaginatedList[Installation] = app.get_installations()
print("Available installations:")
for i in installations:
    repos = [repo for repo in i.get_repos()]
    print(f"{i.id=}, {repos=}")

# Get an installation access token
installation_auth = app.get_access_token(i.id)

# Print the token to STDOUT
print(f"token: {installation_auth.token}")
print(f"expires at: {installation_auth.expires_at}")
