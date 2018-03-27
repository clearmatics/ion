import click

from .client import commands as rpc_client_commands
from .server import server as rpc_server

commands = click.Group('commands', help="RPC service")
commands.add_command(rpc_server, "server")
commands.add_command(rpc_client_commands, "client")

if __name__ == "__main__":
    commands.main()
