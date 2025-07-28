# machv

machv is an easy to use, quick virtual machine set-up tool.

## Name

The name "machv" comes from its speed (machv -> mach5) and that it manages virtual machines (machv -> vmach -> vmachine -> virtualmachine)

## Configuration
### The code for configuration is NOT YET IMPLEMENTED, skip this

You need to create a `~/.config/machv/cfg.toml` configuration file that follows the following format:

```toml
[Config]
iso_fetch_url = "someurl"
```

## ISOs
You need run machv with the `fetch` verb (`fetch` not yet implemented), OR, manually create your own `~/.config/machv/iso.toml` file that follows the following the format: 

```toml
[[Entries]]
Url="https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.11.0-amd64-netinst.iso"
```