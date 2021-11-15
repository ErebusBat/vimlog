# VimLog

## Install

```
go install github.com/erebusbat/vimlog
```

### What is VimLog?

VimLog (or `vimlog`) is a tool to help opening daily logs along with frequently used files.

I use vim along with [Obsidian](https://obsidian.md/) to manage my life.  Each day I create a daily log file to track what I do, notes, thoughts, etc.  I have a tmux pane dedicated to editing this log; however I want to have "todays" log as well as "yesterdays" log.  I also want to be able to open some specific pages (i.e. long running project notes).

I have my tmux setup to execute:

```
vimlog -- y 0 work/project1.md ref.md
```

This will open vim with:

  - `y` - Yesterdays log file
  - `0` - Todays log file
  - Specific files:
    - `work/project1.md` - project notes
    - `ref.md` - reference file

Then each day I can `xa` hit up and enter and will be presented with the "correct" files.

### Date Calculations

The main feature of `vimlog` is to calculate date offsets in order to open specific log files.

You can specify any numeric value and it will be calculated and replaced before vim is launched.

Assuming today is Mon 16-Oct-2006:

  - `vimlog -- 0 -1 -2` should be the same as executing `vim 2006-10-16.md 2006-10-15.md 2006-10-14.md`
  - `y` (yesterday) is special and really means:
    - -1, excepting weekends.
    - So if you use `y` on a monday then it should evaluate to the previous friday (13-Oct-2006 in our example above)
    - However on a Tue-Friday it should mean the same as `-1`

### Configuration

`vimlog` will look in the current working directory (assumed to be the root of your vimwiki / obsidian vault / notebook) for a `.vimlog.yaml` file.  You can also use `~/.config/.vimlog.yaml` if you wish.

Here is an example config file:

```yaml
DateBasePath: journal
Editor: /usr/local/bin/nvim
NoEdit: false
```

The `DateBasePath` option is the most interesting as it will be prepended to any date calculate files.

With the config above, assuming today is Mon 16-Oct-2006:

  - `vimlog -- 0 y -1` should be the same as executing `/usr/local/bin/nvim journal/2006-10-16.md journal/2006-10-13.md journal/2006-10-14.md`

The `NoEdit` option is useful when you are getting used to vimlog and want to check out what would be executed.  You can also prefix your command like so to temporarily turn it on: `VIMLOG_NOEDIT=1 vimlog -- ...`

If you would like to see the current configuration you can run `vimlog config print`

#### Editor Options

You can also pass options to your editor. For example if you are using vim and want to open all files as readonly then you can put the following in your config file:

```
EditorOptions:
  - "-R"
```

Which will cause the command executed to be something like: `vim -R 2006-10-16.md`

A useful setting is:

```yaml
EditorOptions:
  - "-c"
  - ":n"
```

With a command of `vimlog -- y 0` (Notice that yesterday is first, followed by today)

This will place you on todays log with the alternate buffer setup as yesterdays log.
