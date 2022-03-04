# WP-Changer
Simple command line tool to change your wallpaper without repeating the same one twice (unless you want to)

## Dependencies
WP-Changer uses `feh` as its backend so you will need to install it or it will now work.

## Usage
First you will need to create a folder to store the sqlite file which indexes the wallpapers:
```bash
mkdir $HOME/.config/wallpaper-go/
```
First usage will create a sqlite file to store the order of the wallpapers, so it will take
some time indexing, next usages will go instantly.

For changing wallpaper to the next one:
```bash
$ wallpaper -d /path/to/wallpapers
```

For changing wallpaper to previous one:
```bash
$ wallpaper -d /path/to/wallpapers -p
```

## Caveats
In order to add new wallpapers you will need to remove the `database.db` and rerun the program

## TODO
- I've noticed that the wallpaper's indexing is not in order so... sorting them before indexing
should be a good idea.
- I plan to add a random option.

