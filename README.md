# snarfle

Go-based program to work with files created by sniffle.

This does two things:

- With any IPv4 address in the log line that isn't on the 10.0.0.0/8 network, it performs a `dig +short -x <ip>` lookup and adds the domain name found as an extra element to the line.
- With local IPs, it looks in the entries in the `localhosts.cfg` file and sees if there's an entry matching the IP. If so, it replaces the IP with the name of that entry.

## COMMAND LINE OPTIONS

- **-d \<dir>** - The directory containing the input file (realtive or absolute). Default is set in `snarfle.cfg`. Out of the box, this is the current working directory.
- **-f \<filename>** - Filename of the input file. REQUIRED.
- **-o \<fmt>** - Output format: `log`, `txt` or `csv`. Default is set in `snarfle.cfg`. Output files are saved to the current working directory.

## CONFIG FILES

You need these config files. The program looks for them in the following locations (and in this order of priority): `/etc/`, an `etc/` subdir in the current working directory and the current working directory itself.

- **localhosts.cfg** - a key/value list of devices on the local network, in the format `<name>: <ip>`. It might be useful to have this in `/etc/` on *nix systems.
- **snarfle.cfg** - a key/value list of basic config settings, in the format `<name>: <value>`.
