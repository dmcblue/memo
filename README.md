# memo

## Application

`memo` is a small application that records a 'Memo,' a small note with a title, for the user to read later.  The user can tag memos to organize and search for them. Basically, gives quick, short reminders for the user.

Features:
- Ability to add content in a command or via the user's text editor, via the `$EDITOR` var.
- Terminal width aware formatting with the option for no formatting, making results pipable.
- Tagging memo's in order to group them.
- Basic search.

Inspired by the [hot keys popup in AwesomeWM](https://awesomewm.org/apidoc/popups_and_bars/awful.hotkeys_popup.widget.html) but for the terminal. See also the [macOS Stickies app](https://support.apple.com/guide/stickies/welcome/mac) or [Sticky Notes on Windows](https://apps.microsoft.com/detail/9nblggh4qghw?hl=en-us&gl=US). 

Tested on Linux and macOS.

A listing of memos will look like:
```shell
$ memo ls
HASH        TITLE                 CONTENT                          TAGS  
1031f355    Uncommit last set     git reset HEAD~                  git   
            of changes                                                   
28533b79    Check for updates     please softwareupdate --list     macos 
41778df4    Kill process using    kill $(lsof -t -i:3001)          system
            port                                                         
46237b3d    Calculator            bc # quit to quit                tools 
4804368e    Disk usage by file    du -m                                  
            in 1mb increments                                            
5b2d807a    Convert MOV to MP4    ffmpeg -i {in-video}.mov         tools 
                                  -vcodec h264 -acodec aac               
                                  {out-video}.mp4                        
78acb3f2    Vim: Search and       s/foo/bar/g                      vim   
            Replace                                                      
7e1190e3    Run pre-commit        pre-commit run --all-files       git            
926779cd    Run updates           please softwareupdate            macos 
                                  --install --all                        
9cae1706    List all services     docker compose config            docker
            in docker compose     --services                             
            group                                                        
a1d5a722    Run N times           for i in {1..10}; do             linux,
                                  my_command; done                 bash  
a4d65baf    Disk space usage      df -h                            system                                                                  system  
```

### Setup / Installation

Binaries can be found with the latest release. You can build it yourself in the [Development/Build](#build) section.

Upon first use, `memo` creates a config file in the [User Configuration Directory](https://pkg.go.dev/os#UserConfigDir) called `memo.conf`. This config file contains a JSON with one proprery, `SavesDir`. The value for this is the directory where information for the memos will be saved. The default value for this directory is in a folder `memo` also located in the [User Configuration Directory](https://pkg.go.dev/os#UserConfigDir).

### Usage

Below are some basic usages but do not represent all functionality.

Use `memo --help` to see more.

#### Add Memo

```shell
$ memo add "Name of my memo" "Content of my memo"
# The command returns a short Hash to identify the memo for programmatic interactions
03b7d9d8
# or
$ memo add "Name of my memo"
# which then opens a text editor for the user to write out the memo's contents
```

#### Viewing Memos

```shell
$ memo ls
HASH        TITLE                 CONTENT                          TAGS  
1031f355    Uncommit last set     git reset HEAD~                  git   
            of changes                                                   
28533b79    Check for updates     please softwareupdate --list     macos 
# Limit results by tag
$ memo ls --tag git 
HASH        TITLE                 CONTENT                          TAGS  
1031f355    Uncommit last set     git reset HEAD~                  git   
            of changes  
# Show a single memo by hash
$ memo show 1031f355 
HASH        TITLE                 CONTENT                          TAGS  
1031f355    Uncommit last set     git reset HEAD~                  git   
            of changes  
# Show by title
# Remove fancy formatting for pipability
$ memo show --no-format "Uncommit last set of changes"
1031f355	Uncommit last set of changes	git reset HEAD~	git                                          
```

#### Tag a Memo

```shell
$ memo tag 1031f355 my_tag                                       
```

#### Search

```shell
# Get all memos with the word 'commit' in them.
$ memo search commit
HASH        TITLE                           CONTENT                       TAGS                                                                      
1031f355    Uncommit last set of changes    git reset HEAD~               git                                                                       
7e1190e3    Run pre-commit                  pre-commit run --all-files    git 
```

#### Full Options

You can see all available commands with:
```shell
memo --help
```

## Development

### Setup

Install the correct version of [Golang](https://go.dev/) using [g](https://github.com/voidint/g).

```shell
g install $(cat .gorc)
```

### Build

Then build the project:
```shell
go build
```
