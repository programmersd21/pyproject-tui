# Contributing to pyproject-tui

## Dev setup
```sh
git clone https://github.com/programmersd21/pyproject-tui
cd pyproject-tui
go mod tidy
make build
```

## Workflow
1. Fork -> branch -> PR against `main`
2. Run `make check` before pushing (fmt + vet + lint + test)
3. Follow [Conventional Commits](https://www.conventionalcommits.org/)

## Commit types
- `feat:` new feature
- `fix:` bug fix
- `docs:` docs only
- `refactor:` no behavior change
- `test:` test changes
- `chore:` tooling/deps
