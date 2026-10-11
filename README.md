# testing-factory

Testing Fiocco Factory

## Usage

```
weather <city> [--units metric|imperial] [--forecast]
```

Units default to metric (`°C`, `km/h`). Use `--units`/`-u` to switch:

```
weather Rome
weather --units imperial Rome
weather -u imperial New York
```

Add `--forecast`/`-f` to also print a 3-day forecast after the current conditions:

```
weather --forecast Rome
weather -f Rome
```
