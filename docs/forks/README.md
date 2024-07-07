## forks and pulling tls lib updates into the tls_fork library
the document explains the workflow of how to get latest tls package updates into the tls_fork repo.

initially, I cloned the go standard lib into github.com/jplaui/go, then I followed the descriptions of git subtree as explained in the pdf inside the forks folder. once I had my tls_fork repo connected to the go standard tls upstream, I cloned the tls_fork repo into the other repo where I provide updates to the tls folder. here, Im always working on master.
the idea is that the tls_fork repo is cloned into proxy and client at the same time and that updates are maintained in the same repo. 
to update the tls_fork repo and merge new modifications to the tls_fork, you must go back to the jplaui/go repo and switch branch to the one that has the upstream connected to the original go repo and tls standard lib. then from there pull, and merge the changes into the branch I used for the tls_fork repo, from there, you can jump into master and merge the changes into the lastest master branch that you are working on which maintains all latest changes. this process is described in the wegpage download pdf which has the green checkmark. 

## git submodule
to add the tls_fork repo into the client, use `git submodule add -b <branch> <remote_url> <destination_folder>`.
next, commit and push to the branch that you are working in with: `git commit -m "Added the submodule to the project."` and `git push origin kdc`
- works: `git submodule add -b master git@github.com:jplaui/tls_fork.git data/kdc/client/tls_fork`

