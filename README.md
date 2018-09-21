# eventmanager

## Set Up

1. Install [dep](https://golang.github.io/dep/docs/installation.html):
    ```sh
    brew install dep  # OS X example, see link for more
    ```

2. `cd` into the eventmanager directory:
    ```sh
    cd $(go env GOPATH)/src/github.com/megaspacelab/eventmanager
    ```

3. Install dependencies and set up git hooks:
    ```sh
    ./init-dev.sh
    ```
