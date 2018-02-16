import click

from .server import server as rpc_server
from .client import main as rpc_client


commands = click.Group('commands', help="RPC service")
commands.add_command(rpc_server, "server")
commands.add_command(rpc_client, "client")

if __name__ == "__main__":
    commands.main()
