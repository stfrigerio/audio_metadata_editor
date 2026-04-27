# Install

The audio-meta-edit Go binary is wired into the system via a wrapper
script kept in `LinuxConfig/scripts` and a symlink in `~/.local/bin`.
This matches the existing pattern used for `wall`, `claude-switcher`,
`nav`, etc.

## Layout

```
~/Github/audio_metadata_editor/go/audio-meta-edit                      # compiled binary
~/Github/LinuxConfig/scripts/audio-meta-edit/audio-meta-edit.sh        # wrapper
~/.local/bin/audio-meta-edit       → .../LinuxConfig/scripts/audio-meta-edit/audio-meta-edit.sh
```

`~/.local/bin` is already on `$PATH`, so `audio-meta-edit` launches
from anywhere.

## One-time setup

```bash
cd ~/Github/audio_metadata_editor/go
make build

ln -sf ~/Github/LinuxConfig/scripts/audio-meta-edit/audio-meta-edit.sh \
       ~/.local/bin/audio-meta-edit
```

The wrapper script lives in the LinuxConfig repo and is checked in
there. If it's missing, recreate it from the wall pattern.

## Rebuild flow

After editing Go sources:

```bash
cd ~/Github/audio_metadata_editor/go
make build
```

The wrapper execs the binary by absolute path, so there is no
"re-install" step — the next `audio-meta-edit` invocation runs the
freshly built binary.

## Uninstall

```bash
rm ~/.local/bin/audio-meta-edit
rm -rf ~/Github/LinuxConfig/scripts/audio-meta-edit
```
