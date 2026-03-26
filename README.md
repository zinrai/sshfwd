# sshfwd

SSH port forwarding wrapper that simplifies specifying multiple local and remote port forwards.

Instead of writing:

```
ssh -L 8080:localhost:8080 -L 3000:localhost:3000 -R 5432:localhost:5432 -p 2222 user@host
```

You can write:

```
sshfwd -local 8080,3000 -remote 5432 -- -p 2222 user@host
```

## Usage

```
sshfwd -local <ports> -remote <ports> -- [ssh options] [user@]hostname
```

- `-local` : Comma-separated list of ports for local forwarding (`-L port:localhost:port`)
- `-remote` : Comma-separated list of ports for remote forwarding (`-R port:localhost:port`)
- `--` : Separator between sshfwd options and ssh options
- Everything after `--` is passed directly to `ssh`

At least one of `-local` or `-remote` must be specified.

For advanced forwarding such as `-L 8080:10.0.0.1:80`, pass it directly as an ssh option after `--`.

## Examples

Local forwarding:

```
sshfwd -local 8080,3000 -- user@host
```

Remote forwarding:

```
sshfwd -remote 5432,6379 -- user@host
```

Both local and remote forwarding with ssh options:

```
sshfwd -local 8080 -remote 5432 -- -p 2222 -i ~/.ssh/id_ed25519 user@host
```

## License

This project is licensed under the [MIT License](./LICENSE).
