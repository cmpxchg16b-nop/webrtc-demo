## Notes

### About the dependencies

This project relies on [hraban/opus](https://github.com/hraban/opus) golang library and libopus to compile, install the dependencies before trying to build the project:

On Debian, Ubuntu, ...:

```sh
sudo apt-get install pkg-config libopus-dev libopusfile-dev
```

On macOS:

```sh
brew install pkg-config opus opusfile
```
